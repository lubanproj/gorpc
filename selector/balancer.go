package selector

import (
	"math/rand"
	"time"
)

type Balancer interface {
	Balance([]*Node) *Node
}

var balancerMap = make(map[string]Balancer, 0)

const (
	Random = "random"
	RoundRobin = "roundRobin"
	WeightedRoundRobin = "weightedRoundRobin"
	ConsistentHash = "consistentHash"

	Custom = "custom"
)

func init() {
	RegisterBalancer(Random, &randomBalancer{})
}

var DefaultBalancer = &randomBalancer{}

func RegisterBalancer(name string, balancer Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = balancer
}

func GetBalancer(name string) Balancer {
	if balancer, ok := balancerMap[name]; ok {
		return balancer
	}
	return DefaultBalancer
}

type randomBalancer struct {

}

func (r *randomBalancer) Balance(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}
