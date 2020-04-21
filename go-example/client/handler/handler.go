package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"myth/go-essential/net/rpc/warden"
	pb "myth/go-example/proto"
	//"github.com/pkg/errors"
)

type Handler struct {
	Client *warden.Client
}

func NewHandler(client *warden.Client) *Handler {
	result := &Handler{
		Client: client,
	}

	return result
}

func (handler *Handler) SayHello(name string) (*pb.HelloReply, error) {
	ctx := context.Background()
	conn, err := handler.Client.Dial(ctx, "127.0.0.1:8080")
	if err != nil {
		log.Error(context.Background(), "did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	sa := pb.NewGreeterClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "koala_trace_id", "888888888888888888888888888")
	r, err := sa.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Error(ctx, "could not greet: %v", err)
		return nil, err
	}

	log.Info(ctx, " Success: %v %v", r.Message, r.Success)
	return r, nil
}

func (handler *Handler) StreamHello(name string) (*pb.HelloReply, error) {
	ctx := context.Background()
	conn, err := handler.Client.Dial(ctx, "127.0.0.1:8080")
	if err != nil {
		log.Error(context.Background(), "did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	sa := pb.NewGreeterClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "koala_trace_id", "888888888888888888888888888")
	r, err := sa.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Error(ctx, "could not greet: %v", err)
		return nil, err
	}

	log.Info(ctx, " Success: %v %v", r.Message, r.Success)
	return r, nil
}

//中间件
func ExampleAuthFunc() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Debug("ni hao")
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
