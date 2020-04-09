package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init("localhost:8500")
	assert.Nil(t,err)
}
