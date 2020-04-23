package warden

import (
	"github.com/pkg/errors"
	"myth/go-essential/log/logf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"myth/go-essential/base/rpc/server"
	"net"

	"context"
	"math"
	"sync"
)

var (
	//中间件终止索引
	_abortIndex int8 = math.MaxInt8 / 2
)

type Server struct {
	// marks the serve as started
	started bool
	// used for first registration
	registered bool
	// middleware
	UnaryHandlers  []grpc.UnaryServerInterceptor
	StreamHandlers []grpc.StreamServerInterceptor

	server *grpc.Server
	opts   server.Options
	mutex  sync.RWMutex
}

func NewServer() *Server {
	srv := &Server{}
	return srv
}

func (s *Server) Options() server.Options {
	return s.opts
}

func (s *Server) RpcServer() *grpc.Server {
	return s.server
}

func (s *Server) String() string {
	return "grpc"
}

func (s *Server) Init(opts ...server.Option) error {
	//if len(opts) == 0 {
	//	return errors.New("warden: rpc server options is empty")
	//}

	if s.server != nil {
		return errors.New("warden: grpc server already init")
	}

	if err := s.register(); err != nil {
		return err
	}

	return nil
}


func (s *Server) RegisterRpc(){
	s.server = grpc.NewServer(
		WithUnaryServerChain(s.UnaryHandlers...),
		WithStreamServerChain(s.StreamHandlers...),
	)
}

func (s *Server) register() error {
	return nil
}

func (s *Server) deRegister() error {
	return nil
}

//普通中间件
func (s *Server) UseUnary(handlers ...grpc.UnaryServerInterceptor) *Server {
	finalSize := len(s.UnaryHandlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: server use unary interceptor too many handlers")
	}

	mergedHandlers := make([]grpc.UnaryServerInterceptor, finalSize)
	copy(mergedHandlers, s.UnaryHandlers)
	copy(mergedHandlers[len(s.UnaryHandlers):], handlers)
	s.UnaryHandlers = mergedHandlers

	return s
}

//数据流中间件
func (s *Server) UseStream(handlers ...grpc.StreamServerInterceptor) *Server {
	finalSize := len(s.StreamHandlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: server use stream interceptor too many handlers")
	}

	mergedHandlers := make([]grpc.StreamServerInterceptor, finalSize)
	copy(mergedHandlers, s.StreamHandlers)
	copy(mergedHandlers[len(s.StreamHandlers):], handlers)
	s.StreamHandlers = mergedHandlers

	return s
}

func (s *Server) Start() error {
	if s.server == nil {
		return errors.New("warden: grpc server is null")
	}

	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *Server) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

// RunUnix create a unix listener and start goroutine for serving each incoming request.
// RunUnix will return a non-nil error unless Stop or GracefulStop is called.
func (s *Server) ServeUnix(file string) error {
	lis, err := net.Listen("unix", file)
	if err != nil {
		log.Error("warden: failed to listen unix: %v", err)
		err = errors.WithStack(err)
		return err
	}

	reflection.Register(s.server)
	return s.Serve(lis)
}

func (s *Server) defaultUnaryHandlers() {

}

func (s *Server) defaultStreamHandlers() {

}

func (s *Server) Stop(ctx context.Context) error {
	s.mutex.RLock()
	if !s.started {
		s.mutex.RUnlock()
		return nil
	}
	s.mutex.RUnlock()

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

	if err := s.deRegister(); err != nil {
		return err
	}

	return nil
}
