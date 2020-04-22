package api

import (
	"context"
	"google.golang.org/grpc"
	"myth/go-essential/net/rpc/warden"
	pb "myth/go-example/proto"
)

// AppID unique app id for service discovery
const AppID = "account.service.member"

// NewClient new member grpc client
func NewClient(opts ...grpc.DialOption) (pb.GreeterClient, error) {
	client := warden.NewClient(opts...)
	conn, err := client.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}

	return pb.NewGreeterClient(conn), nil
}
