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
	nodes []*Node			// service nodes
	lastUpdateTime time.Time  // last update time
	duration time.Duration    // time duration to update again
	lastIndex int    // last accessed index
}

func (rr *roundRobinBalancer) updateServer(picker *roundRobinPicker, nodes []*Node) {

}

func (rp *roundRobinPicker) pick(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	// update picker after timeout
	if time.Now().Sub(rp.lastUpdateTime) > rp.duration ||
		len(nodes) != len(rp.nodes){
		rp.nodes = nodes
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
			nodes : nodes,
		}
		r.pickers.Store(serviceName,picker)
	} else {
		picker = p.(*roundRobinPicker)
	}

	r.updateServer(picker, nodes)
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
