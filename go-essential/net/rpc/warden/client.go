package warden

import (
	"myth/go-essential/base/rpc/client"
)

type Client struct {
	opts client.Options
}

//Init(...Option) error
//Options() Options
//NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
//NewRequest(service, endpoint string, req interface{}, reqOpts ...RequestOption) Request
//Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
//Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
//Publish(ctx context.Context, msg Message, opts ...PublishOption) error
//String() string

func NewClient() client.Client {
	client := &Client{}

	return client
}



func (s *Client) Options() client.Options {
	return s.opts
}

func (s *Client) Init(opts ...client.Option) error {
	return nil
}

func (s *Client) Handle() error {
	return nil
}

func (s *Client) NewHandler() {
	return
}

func (s *Client) NewSubscriber() {
	return
}

func (s *Client) Subscribe() {
	return
}

func (s *Client) Start() error {
	return nil
}

func (s *Client) Stop() error {
	return nil
}

func (s *Client) String() string {
	return "grpc"
}