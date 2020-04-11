package selector

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mockconsulName = "mockconsul"
	nodeList       = []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
)

type MockConsul struct {
	serivePool map[string][]string
}

func InitMockConsul() *MockConsul {
	return &MockConsul{
		serivePool: map[string][]string{
			"Greeter": nodeList,
		},
	}
}

func (c *MockConsul) Select(serviceName string) (string, error) {

	rand.Seed(time.Now().UnixNano())
	nodeList, ok := c.serivePool[serviceName]
	if !ok {
		return "", fmt.Errorf("service not be registered in mockconsul!")
	}

	index := rand.Int() % len(nodeList)
	return nodeList[index], nil
}

func TestGetSelector(t *testing.T) {
	m := InitMockConsul()
	RegisterSelector(mockconsulName, m)

	selector := GetSelector(mockconsulName)
	node, err := selector.Select("Greeter")
	assert.Nil(t, err)
	assert.Contains(t, nodeList, node)
	_, err = selector.Select("")
	assert.NotNil(t, err)

	selector = GetSelector("")
	assert.Equal(t, selector, DefaultSelector)

}