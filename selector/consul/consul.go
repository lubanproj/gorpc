package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/lubanproj/gorpc/selector"
)

type consulSelector struct {
	opts *Options
	client *api.Client
	config *api.Config
	balancer string  // 负载均衡模式，包括随机、轮询、加权轮询、一致性hash 等
}

type Options struct {
	queryOpts *api.QueryOptions
}

type Option func(*Options)

type KVPair struct {

}

func (c *consulSelector) Resolve() ([]*selector.Node, error) {
	pairs, _, err := c.client.KV().List(serviceName, c.opts.queryOpts)
	if err != nil {
		return nil, err
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("no services find in path : %s", serviceName)
	}
	var nodes []*selector.Node
	for _, pair := range pairs {
		nodes = append(nodes, &selector.Node {
			key : pair.Key,
			value : pair.Value,
		})
	}
	return nodes, nil
}

func (c *consulSelector) Select(serviceName string) (string, error) {

	nodes, err := c.Resolve()

	if nodes == nil || len(nodes) == 0 || err != nil {
		return "", err
	}

	node := c.balancer.Balance(nodes)

	return node.Value, nil
}

