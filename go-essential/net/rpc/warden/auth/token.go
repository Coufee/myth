package auth

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"myth/go-essential/metadata"
)

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func (a *Authentication) Auth(ctx context.Context) error {
	var appid string
	var appkey string
	appid = metadata.String(ctx, "login")
	appkey = metadata.String(ctx, "password")

	if appid != a.User || appkey != a.Password {
		return errors.Errorf("invalid token")
	}

	return nil
}

func WithAuth() grpc.DialOption {
	auth := Authentication{User: "gopher", Password: "password"}
	return grpc.WithPerRPCCredentials(&auth)
}

type grpcServer struct {
	auth *Authentication
}

//func (p *grpcServer) SomeMethod(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
//	if err := p.auth.Auth(ctx); err != nil {
//		return nil, err
//	}
//	return &HelloReply{Message: "Hello " + in.Name}, nil
//}

func filter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("fileter:", info)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return handler(ctx, req)
}
