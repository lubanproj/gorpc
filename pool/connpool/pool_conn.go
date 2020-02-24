package connpool

import (
	"net"
	"sync"
	"time"
)

type PoolConn struct {
	net.Conn
	c *channelPool
	unusable bool		// if unusable is true, the conn should be closed
	mu sync.RWMutex
	t time.Time  // connection idle time
	checked bool        // flags to be used by the checker
	dialTimeout time.Duration // connection timeout duration
}

// overwrite conn Close for connection reuse
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}

	// reset connection deadline
	p.Conn.SetReadDeadline(time.Now().Add(p.dialTimeout))

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
		t : time.Now(),
		dialTimeout: c.dialTimeout,
	}
	p.Conn = conn
	return p
}