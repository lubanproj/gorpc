package consul

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	err := Init("localhost:8500")
	assert.Nil(t,err)
}
