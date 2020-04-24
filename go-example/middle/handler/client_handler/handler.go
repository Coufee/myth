package client_handler

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	log "myth/go-essential/log/logc"
	"myth/go-essential/net/rpc/warden"
	pb "myth/go-example/proto"
	"os"

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
		log.Errorc(context.Background(), "did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	sa := pb.NewGreeterClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "koala_trace_id", "888888888888888888888888888")
	r, err := sa.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Errorc(ctx, "could not greet: %v", err)
		return nil, err
	}

	log.Infoc(ctx, " Success: %v %v", r.Message, r.Success)
	return r, nil
}

func (handler *Handler) StreamHello(name string)  error {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		log.Error("conn fail: [%v]\n", err)
		return err
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	ctx := context.Background()
	stream, err := client.StreamHello(ctx)
	if err != nil {
		log.Error("create stream error: [%v]\n", err)
	}

	go func() {
		in := bufio.NewReader(os.Stdin)
		for {
			inStr, _ := in.ReadString('\n')
			if err := stream.Send(&pb.HelloRequest{Name : inStr}); err != nil {
				return
			}
		}
	}()

	for {
		out, err := stream.Recv()
		if err == io.EOF {
			log.Debug("recv end sign")
			break
		}

		if err != nil {
			log.Debug("recv data error :", err)
		}

		// 没有错误的情况下，打印来自服务端的消息
		log.Debug("client recv data : %s", out.Message)
	}

	return nil
}

//中间件
func ExampleAuthFunc(i *int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Debug("ni hao %v", *i)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
