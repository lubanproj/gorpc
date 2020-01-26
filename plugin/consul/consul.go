package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/lubanproj/gorpc/selector"
	"net/http"
)

type Consul struct {
	opts *selector.Options
	client *api.Client
	config *api.Config
	balancerName string  // 负载均衡模式，包括随机、轮询、加权轮询、一致性hash 等
}

type KVPair struct {

}

const Name = "consul"

func New(consulAddr string, opts ...selector.Option) (*Consul, error) {

	c := &Consul{
		opts : &selector.Options{},
	}

	for _, o := range opts {
		o(c.opts)
	}

	config := api.DefaultConfig()
	c.config = config

	config.HttpClient = http.DefaultClient
	config.Address = consulAddr
	config.Scheme = "http"

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c.client = client

	return nil, nil
}


func (c *Consul) Resolve(serviceName string) ([]*selector.Node, error) {

	pairs, _, err := c.client.KV().List(serviceName, nil)
	if err != nil {
		return nil, err
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("no services find in path : %s", serviceName)
	}
	var nodes []*selector.Node
	for _, pair := range pairs {
		nodes = append(nodes, &selector.Node {
			Key : pair.Key,
			Value : pair.Value,
		})
	}
	return nodes, nil
}

func (c *Consul) Select(serviceName string) (string, error) {

	nodes, err := c.Resolve(serviceName)

	if nodes == nil || len(nodes) == 0 || err != nil {
		return "", err
	}

	balancer := selector.GetBalancer(c.balancerName)
	node := balancer.Balance(nodes)

	if node == nil {
		return "", fmt.Errorf("no services find in %s", serviceName)
	}

	return "", nil
}

func (c *Consul) Start() {
	c.Register()
}

func (c *Consul) Register() {

}
