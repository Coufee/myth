package warden

import (
	"google.golang.org/grpc"
	"myth/go-essential/base/rpc/server"

	"context"
	"math"
	"net"
	"sync"
)

var (
	_abortIndex int8 = math.MaxInt8 / 2
)

type Server struct {
	server grpc.Server
	// marks the serve as started
	started bool
	// used for first registration
	registered bool

	opts     server.Options
	mutex    sync.RWMutex
	handlers []grpc.UnaryServerInterceptor
}

func NewServer() server.Server {
	srv := &Server{}

	srv.Use(srv.Handle())
	return srv
}

func (s *Server) Options() server.Options {
	return s.opts
}

func (s *Server) Init(opts ...server.Option) error {
	return nil
}

func (s *Server) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//var (
		//	cancel func()
		//	addr   string
		//)
		//s.mutex.RLock()
		//conf := s.opts
		//s.mutex.RUnlock()
		//// get derived timeout from grpc context,
		//// compare with the warden configured,
		//// and use the minimum one
		//timeout := time.Duration(conf.Timeout)
		//if dl, ok := ctx.Deadline(); ok {
		//	ctimeout := time.Until(dl)
		//	if ctimeout-time.Millisecond*20 > 0 {
		//		ctimeout = ctimeout - time.Millisecond*20
		//	}
		//	if timeout > ctimeout {
		//		timeout = ctimeout
		//	}
		//}
		//ctx, cancel = context.WithTimeout(ctx, timeout)
		//defer cancel()
		//
		//// get grpc metadata(trace & remote_ip & color)
		//var t trace.Trace
		//cmd := nmd.MD{}
		//if gmd, ok := metadata.FromIncomingContext(ctx); ok {
		//	t, _ = trace.Extract(trace.GRPCFormat, gmd)
		//	for key, vals := range gmd {
		//		if nmd.IsIncomingKey(key) {
		//			cmd[key] = vals[0]
		//		}
		//	}
		//}
		//if t == nil {
		//	t = trace.New(args.FullMethod)
		//} else {
		//	t.SetTitle(args.FullMethod)
		//}
		//
		//if pr, ok := peer.FromContext(ctx); ok {
		//	addr = pr.Addr.String()
		//	t.SetTag(trace.String(trace.TagAddress, addr))
		//}
		//defer t.Finish(&err)
		//
		//// use common meta data context instead of grpc context
		//ctx = nmd.NewContext(ctx, cmd)
		//ctx = trace.NewContext(ctx, t)
		//
		//resp, err = handler(ctx, req)
		//return resp, status.FromError(err).Err()
		return
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(ch)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		err := ctx.Err()
		return err
	case <-ch:
	}

	return nil
}

func (s *Server) String() string {
	return "grpc"
}

func (s *Server) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

// interceptor is a single interceptor out of a chain of many interceptors.
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
func (s *Server) interceptor(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	n := len(s.handlers)
	if n == 0 {
		return handler(ctx, req)
	}

	chain = func(ic context.Context, ir interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ic, ir)
		}
		i++
		return s.handlers[i](ic, ir, args, chain)
	}

	return s.handlers[0](ctx, req, args, chain)
}

// Use attachs a global inteceptor to the server.
// For example, this is the right place for a rate limiter or error management inteceptor.
func (s *Server) Use(handlers ...grpc.UnaryServerInterceptor) *Server {
	finalSize := len(s.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: server use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, finalSize)
	copy(mergedHandlers, s.handlers)
	copy(mergedHandlers[len(s.handlers):], handlers)
	s.handlers = mergedHandlers
	return s
}
