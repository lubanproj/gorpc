package selector

import (
	"math/rand"
	"time"
)

// Balancer defines a universal standard for load Balancer
type Balancer interface {
	Balance(string, []*Node) *Node
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
	RegisterBalancer(Random, DefaultBalancer)
	RegisterBalancer(RoundRobin, RRBalancer)
	RegisterBalancer(WeightedRoundRobin, WRRBalancer)
}

// RandomBalancer is adopted as the default load balancer
var DefaultBalancer = newRandomBalancer()
// A unique RoundRobinBalancer instance is used globally
var RRBalancer = newRoundRobinBalancer()
// A unique WeightedRoundRobinBalancer instance is used globally
var WRRBalancer = newWeightedRoundRobinBalancer()

// RegisterBalancer supports business custom registered Balancer
func RegisterBalancer(name string, balancer Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = balancer
}

// GetBalancer get a Balancer by a balancer name
func GetBalancer(name string) Balancer {
	if balancer, ok := balancerMap[name]; ok {
		return balancer
	}
	return DefaultBalancer
}

func newRandomBalancer() *randomBalancer {
	return &randomBalancer{}
}

type randomBalancer struct {

}

func (r *randomBalancer) Balance(serviceName string, nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}



