package recovery

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// recovery is a server interceptor that recovers from any panics.
func ServerRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc server panic: %v\n%v\n%s\n", req, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				err = status.Errorf(codes.Unknown, "server middleware error")
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

// recovery return a client interceptor  that recovers from any panics.
func ClientRecovery() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc client panic: %v\n%v\n%v\n%s\n", req, reply, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				err = errors.New("client middleware error")
			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
