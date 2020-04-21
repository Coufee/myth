package warden

//import (
//	"context"
//	"github.com/pkg/errors"
//	log "github.com/sirupsen/logrus"
//	"google.golang.org/grpc"
//	"io"
//	"myth/go-essential/container/pool"
//	"os"
//)
//
//type poolConn struct {
//	addr string
//	opts []grpc.DialOption
//	conn *grpc.ClientConn
//}
//
//func (c *poolConn) Close() error {
//	c.conn.Close()
//	return nil
//}
//
//type ClientPool struct {
//	*pool.Slice
//	conf *pool.Config
//}
//
//func NewClientPool(conf *pool.Config) *ClientPool {
//	p := pool.NewSlice(conf)
//	p.New = func(ctx context.Context) (io.Closer, error) {
//		return &poolConn{conn: &grpc.ClientConn{}}, nil
//	}
//
//	clientPool := &ClientPool{}
//	clientPool.conf = conf
//	clientPool.Slice = p
//	return clientPool
//}
//
//func (p *ClientPool) Get(ctx context.Context, target string, opts ...grpc.DialOption) (*poolConn, error) {
//	if p.Slice == nil {
//		return nil, errors.New("pool slice is empty")
//	}
//	c, err := p.Slice.Get(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	pc, _ := c.(*poolConn)
//	pc.addr = target
//	pc.opts = opts
//
//	conn1, err := grpc.DialContext(ctx, target, opts...)
//	if err != nil {
//		err = errors.WithStack(err)
//		log.Error(os.Stderr, "warden: client dial %s error %v!", target, err)
//		return nil, err
//	}
//	pc.conn = conn1
//	return pc, nil
//}
//
//func (p *ClientPool) Put(ctx context.Context, conn *poolConn, target string, opts ...grpc.DialOption) error {
//	conn.addr = target
//	conn.opts = opts
//	err := p.Slice.Put(ctx, conn, false)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *ClientPool) IdleCount() int {
//	return p.conf.Idle
//}
//
//func (p *ClientPool) ActiveCount() int {
//	return p.conf.Active
//}
//
//func (p *ClientPool) Close() error {
//	p.Slice.Close()
//	return nil
//}
