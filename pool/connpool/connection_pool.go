package connpool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

// Pool provides a pooling capability for connections, enabling connection reuse
type Pool interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}

type pool struct {
	opts *Options
	conns *sync.Map
}

var poolMap = make(map[string]Pool)
var oneByte = make([]byte, 1)

func init() {
	registorPool("default", DefaultPool)
}

func registorPool(poolName string, pool Pool) {
	poolMap[poolName] = pool
}

// GetPool get a Pool by a pool name
func GetPool(poolName string) Pool {
	if v, ok := poolMap[poolName]; ok {
		return v
	}
	return DefaultPool
}

// TODO expose the ConnPool options
var DefaultPool = NewConnPool()

func NewConnPool(opt ...Option) *pool {
	// default options
	opts := &Options {
		maxCap: 1000,
		idleTimeout: 1 * time.Minute,
		dialTimeout: 200 * time.Millisecond,
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
			return conn, err
		}
	}

	cp, err := p.NewChannelPool(ctx, network, address)
	if err != nil {
		return nil, err
	}

	p.conns.Store(address, cp)

	return cp.Get(ctx)
}

type channelPool struct {
	net.Conn
	initialCap int  // initial capacity
	maxCap int      // max capacity
	maxIdle int     // max idle conn number
	idleTimeout time.Duration  // idle timeout
	dialTimeout time.Duration  // dial timeout
	Dial func(context.Context) (net.Conn, error)
	conns chan *PoolConn
	mu sync.RWMutex
}


func (p *pool) NewChannelPool(ctx context.Context, network string, address string) (*channelPool, error){
	c := &channelPool {
		initialCap: p.opts.initialCap,
		maxCap: p.opts.maxCap,
		Dial : func(ctx context.Context) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := p.opts.dialTimeout
			if t , ok := ctx.Deadline(); ok {
				timeout = t.Sub(time.Now())
			}

			return net.DialTimeout(network, address, timeout)
		},
		conns : make(chan *PoolConn, p.opts.maxCap),
		idleTimeout: p.opts.idleTimeout,
		dialTimeout: p.opts.dialTimeout,
	}

	if p.opts.initialCap == 0 {
		// default initialCap is 1
		p.opts.initialCap = 1
	}

	for i := 0; i < p.opts.initialCap; i++ {
		conn , err := c.Dial(ctx);
		if err != nil {
			return nil, err
		}
		c.Put(c.wrapConn(conn))
	}

	c.RegisterChecker(3 * time.Second, c.Checker)
	return c, nil
}

func (c *channelPool) Get(ctx context.Context) (net.Conn, error) {
	if c.conns == nil {
		return nil, ErrConnClosed
	}
	select {
		case pc := <-c.conns :
			if pc == nil {
				return nil, ErrConnClosed
			}

			if pc.unusable {
				return nil, ErrConnClosed
			}

			return pc, nil
		default:
			conn, err := c.Dial(ctx)
			if err != nil {
				return nil, err
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
		conn.MarkUnusable()
		conn.Close()
	}
}

func (c *channelPool) Put(conn *PoolConn) error {
	if conn == nil {
		return errors.New("connection closed")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conns == nil {
		conn.MarkUnusable()
		conn.Close()
	}

	select {
	case c.conns <- conn :
		return nil
	default:
		// 连接池满
		return conn.Close()
	}
}

func (c *channelPool) RegisterChecker(internal time.Duration, checker func(conn *PoolConn) bool) {

	if internal <= 0 || checker == nil {
		return
	}

	go func() {

		for {

			time.Sleep(internal)

			length := len(c.conns)

			for i:=0; i < length; i++ {

				select {
				case pc := <- c.conns :

					if !checker(pc) {
						pc.MarkUnusable()
						pc.Close()
						break
					} else {
						c.Put(pc)
					}
				default:
					break
				}

			}
		}

	}()
}

func (c *channelPool) Checker (pc *PoolConn) bool {

	// check timeout
	if pc.t.Add(c.idleTimeout).Before(time.Now()) {
		return false
	}

	// check conn is alive or not
	if !isConnAlive(pc.Conn) {
		return false
	}

	return true
}

func isConnAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond))

	if n, err := conn.Read(oneByte); n > 0 || err == io.EOF {
		return false
	}

	conn.SetReadDeadline(time.Time{})
	return true
}




