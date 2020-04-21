package warden

import (
	"context"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"myth/go-essential/base/rpc/client"
	"myth/go-essential/container/pool"
	"sync"
)

type Client struct {
	// TLS
	OpenTLS bool
	// the underlying connection is up
	NonBlock bool
	// middleware
	UnaryHandlers  []grpc.UnaryClientInterceptor
	StreamHandlers []grpc.StreamClientInterceptor

	dialOpts   []grpc.DialOption
	mutex      sync.RWMutex
	opts       client.Options
	clientPool *ClientPool
	//rpcPool    *RpcPool
}

func NewClient() client.Client {
	client := &Client{}
	conf := &pool.Config{
		Active: 3,
		Idle:   3,
		Wait:   true,
	}
	client.clientPool = NewClientPool(conf)
	//client.rpcPool = newPool(3, 3, 3, 3)

	return client
}

func (c *Client) Options() client.Options {
	return c.opts
}

func (c *Client) String() string {
	return "grpc"
}

func (c *Client) Init(opts ...client.Option) error {
	//if len(opts) == 0 {
	//	return errors.New("warden: rpc client options is empty")
	//}

	return nil
}

func (c *Client) UseOpts(opts ...grpc.DialOption) *Client {
	c.dialOpts = append(c.dialOpts, opts...)
	return c
}

func (c *Client) cloneOpts() []grpc.DialOption {
	dialOptions := make([]grpc.DialOption, len(c.dialOpts))
	copy(dialOptions, c.dialOpts)
	return dialOptions
}

func (c *Client) UseUnary(handlers ...grpc.UnaryClientInterceptor) *Client {
	finalSize := len(c.UnaryHandlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: client use unary interceptor too many handlers")
	}

	mergedHandlers := make([]grpc.UnaryClientInterceptor, finalSize)
	copy(mergedHandlers, c.UnaryHandlers)
	copy(mergedHandlers[len(c.UnaryHandlers):], handlers)
	c.UnaryHandlers = mergedHandlers
	return c
}

//数据流中间件
func (c *Client) UseStream(handlers ...grpc.StreamClientInterceptor) *Client {
	finalSize := len(c.StreamHandlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: client use stream interceptor too many handlers")
	}

	mergedHandlers := make([]grpc.StreamClientInterceptor, finalSize)
	copy(mergedHandlers, c.StreamHandlers)
	copy(mergedHandlers[len(c.StreamHandlers):], handlers)
	c.StreamHandlers = mergedHandlers

	return c
}

func (c *Client) Stop() error {
	//c.pool.Close()
	return nil
}

func (c *Client) dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	dialOptions := c.cloneOpts()
	dialOptions = append(dialOptions, opts...)
	if !c.NonBlock {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	//dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
	//	Time:                time.Duration(c.conf.KeepAliveInterval),
	//	Timeout:             time.Duration(c.conf.KeepAliveTimeout),
	//	PermitWithoutStream: !c.conf.KeepAliveWithoutStream,
	//}))

	dialOptions = append(
		dialOptions,
		WithUnaryClientChain(c.UnaryHandlers...),
		WithStreamClientChain(c.StreamHandlers...),
	)

	if c.clientPool == nil {
		return nil, errors.New("client pool is empty")
	}

	rpcConn, err := c.clientPool.Get(ctx, target, dialOptions...)
	if err != nil {
		return nil, err
	}
	defer func() {
		rpcConn.Close()
		c.clientPool.Put(ctx, rpcConn)
	}()

	conn = rpcConn.ClientConn

	//cc, err := c.rpcPool.getConn(ctx, target, dialOptions...)
	//if err != nil {
	//	log.Error(err)
	//}
	//conn = cc.ClientConn

	//conn, err = grpc.DialContext(ctx, target, dialOptions...)
	//if err != nil {
	//	err = errors.WithStack(err)
	//	log.Error(os.Stderr, "warden: client dial %s error %v!", target, err)
	//	return nil, err
	//}

	return
}

func (c *Client) Dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts, grpc.WithInsecure())
	return c.dial(ctx, target, opts...)
}

func (c *Client) DialTLS(ctx context.Context, target string, file string, name string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))
	return c.dial(ctx, target, opts...)
}
