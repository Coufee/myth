package server_handler

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"myth/go-essential/log/logf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	pb "myth/go-example/proto"
	"myth/go-example/server/api"
	//"github.com/pkg/errors"
)

type Handler struct {
	helloClient pb.GreeterClient
}

func NewHandler() *Handler {
	result := &Handler{}
	hello, err := api.NewClient()
	if err != nil {
		panic(err)
	}
	result.helloClient = hello
	return result
}

func (handler *Handler) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Debug("SayHello")
	resp, err := handler.helloClient.SayHello(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (handler *Handler) StreamHello(stream pb.Greeter_StreamHelloServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("收到客户端通过context发出的终止信号")
			return ctx.Err()
		default:
			in, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			ret := &pb.HelloReply{Message: "Hello " + in.Name, Success: true}
			err = stream.Send(ret)
			if err != nil {
				return err
			}
		}
	}

	log.Debug("StreamHello success")
	return nil
}

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

func ExampleAuthFunc() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}
		tokenInfo, err := parseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
		newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
		return newCtx, nil
	}
}
