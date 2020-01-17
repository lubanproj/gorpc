package connpool

import (
	"context"
	"errors"
	"github.com/lubanproj/gorpc/codes"
	"net"
	"sync"
)

type Pool interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}

type pool struct {
	opts *Options
	conns *sync.Map
}

var poolMap = make(map[string]Pool)

func init() {
	registorPool("default", DefaultPool)
}

func registorPool(poolName string, pool Pool) {
	poolMap[poolName] = pool
}

func GetPool(poolName string) Pool {
	if v, ok := poolMap[poolName]; ok {
		return v
	}
	return DefaultPool
}

// TODO 暴露 ConnPool 属性
var DefaultPool = NewConnPool()

func NewConnPool(opt ...Option) *pool {
	// 默认值
	opts := &Options {
		initialCap: 5,
		maxCap: 1000,
	}
	m := &sync.Map{}

	p := &pool {
		conns : m,
		opts : opts,
	}
	for _, o := range opt {
		o(p.opts)
	}

	return p
}

func (p *pool) Get(ctx context.Context, network string, address string) (net.Conn, error) {

	if value, ok := p.conns.Load(address); ok {
		if cp, ok := value.(*channelPool); ok {
			conn, err := cp.Get(ctx)
			return cp.wrapConn(conn), err
		}
	}

	cp, err := p.NewChannelPool(ctx, network, address)
	if err != nil {
		return nil, codes.ConnectionPoolInitError
	}

	p.conns.Store(address, cp)

	return cp.Get(ctx)
}

type channelPool struct {
	net.Conn
	initialCap int
	maxCap int
	Dial func(context.Context) (net.Conn, error)
	conns chan net.Conn
	mu sync.Mutex
}


func (p *pool) NewChannelPool(ctx context.Context, network string, address string) (*channelPool, error){
	c := &channelPool {
		initialCap: p.opts.initialCap,
		maxCap: p.opts.maxCap,
		Dial : func(ctx context.Context) (net.Conn, error) {
			return net.Dial(network, address)
		},
		conns : make(chan net.Conn, p.opts.maxCap),
	}
	conn , err := c.Dial(ctx);
	if err != nil {
		c.Close()
		return nil, codes.ConnectionPoolInitError
	}
	c.conns <- conn
	return c, nil
}

func (c *channelPool) Get(ctx context.Context) (net.Conn, error) {
	if c.conns == nil {
		return nil, errors.New("connection closed")
	}
	select {
		case conn := <-c.conns :
			if conn == nil {
				return nil, errors.New("connection closed")
			}
			return c.wrapConn(conn), nil
		default:
			conn, err := c.Dial(ctx)
			if err != nil {
				return nil, codes.ClientNetworkError
			}
			return c.wrapConn(conn), nil
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.Dial = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) Put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection closed")
	}
	if c.conns == nil {
		conn.Close()
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
		case c.conns <- conn :
			return nil
	default:
		// 连接池满
		return conn.Close()
	}
}



