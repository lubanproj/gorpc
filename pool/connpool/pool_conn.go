package connpool

import (
	"net"
	"sync"
)

type PoolConn struct {
	net.Conn
	c *channelPool
	unusable bool
	mu sync.RWMutex
}

// 覆盖 conn Close, 实现连接复用
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}

	return p.c.Put(p.Conn)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn {
		c : c,
	}
	p.Conn = conn
	return p
}