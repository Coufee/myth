package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "myth/go-example/proto"
	//"github.com/pkg/errors"
)

type Handler struct {
}

func NewHandler() *Handler {
	result := &Handler{}
	return result
}

func (handler *Handler) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := &pb.HelloResponse{}
	resp.Reply = "aaa"
	log.Debug("success")

	return resp, nil
}
