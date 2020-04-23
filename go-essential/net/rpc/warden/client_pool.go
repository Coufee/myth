package warden

import (
	"context"
	"github.com/pkg/errors"
	"myth/go-essential/log/logf"
	"google.golang.org/grpc"
	"io"
	"myth/go-essential/container/pool"
	"os"
)

type rpcConn struct {
	*grpc.ClientConn
}

func (c *rpcConn) Close() error {
	c.ClientConn.Close()
	return nil
}

type ClientPool struct {
	*pool.Slice
	conf *pool.Config
}

func NewClientPool(conf *pool.Config) *ClientPool {
	p := pool.NewSlice(conf)
	p.New = func(ctx context.Context) (io.Closer, error) {
		return &rpcConn{ClientConn: &grpc.ClientConn{}}, nil
	}

	clientPool := &ClientPool{}
	clientPool.conf = conf
	clientPool.Slice = p
	return clientPool
}

func (p *ClientPool) Get(ctx context.Context, target string, opts ...grpc.DialOption) (*rpcConn, error) {
	if p.Slice == nil {
		return nil, errors.New("pool slice is empty")
	}
	c, err := p.Slice.Get(ctx)
	if err != nil {
		return nil, err
	}

	pc, ok := c.(*rpcConn)
	log.Debug(ok)

	conn1, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		err = errors.WithStack(err)
		log.Error(os.Stderr, "warden: client dial %s error %v!", target, err)
		return nil, err
	}
	pc.ClientConn = conn1
	return pc, nil
}

func (p *ClientPool) Put(ctx context.Context, conn *rpcConn) error {
	err := p.Slice.Put(ctx, conn, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientPool) IdleCount() int {
	return p.conf.Idle
}

func (p *ClientPool) ActiveCount() int {
	return p.conf.Active
}

func (p *ClientPool) Close() error {
	p.Slice.Close()
	return nil
}
