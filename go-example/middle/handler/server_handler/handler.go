package server_handler

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
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

func (handler *Handler) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	for i := 0; i < 3; i++ {
		in, err := ss.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		ret := &pb.HelloReply{Message: "Hello " + in.Name, Success: true}
		err = ss.Send(ret)
		if err != nil {
			return err
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
