package selector

import (
	"sync"
	"time"
)

type roundRobinBalancer struct {
	pickers *sync.Map
	duration time.Duration  // time duration to update again
}

type roundRobinPicker struct {
	length int       // service nodes length
	lastUpdateTime time.Time  // last update time
	duration time.Duration    // time duration to update again
	lastIndex int    // last accessed index
}

func (rp *roundRobinPicker) pick(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	// update picker after timeout
	if time.Now().Sub(rp.lastUpdateTime) > rp.duration ||
		len(nodes) != rp.length {
		rp.length = len(nodes)
		rp.lastUpdateTime = time.Now()
		rp.lastIndex = 0
	}

	if rp.lastIndex == len(nodes) - 1 {
		rp.lastIndex = 0
		return nodes[0]
	}

	rp.lastIndex += 1
	return nodes[rp.lastIndex]
}

func (r *roundRobinBalancer) Balance(serviceName string, nodes []*Node) *Node {

	var picker *roundRobinPicker

	if p, ok := r.pickers.Load(serviceName); !ok {
		picker = &roundRobinPicker{
			lastUpdateTime: time.Now(),
			duration : r.duration,
			length : len(nodes),
		}
	} else {
		picker = p.(*roundRobinPicker)
	}

	node := picker.pick(nodes)
	r.pickers.Store(serviceName,picker)
	return node
}

func newRoundRobinBalancer() *roundRobinBalancer {
	return &roundRobinBalancer{
		pickers : new(sync.Map),
		duration : 3 * time.Minute,
	}
}
