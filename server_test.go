package gorpc

import (
	"testing"

	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"

	"github.com/stretchr/testify/assert"
)

func TestRegisterService(t *testing.T) {
	s := NewServer()
	err := s.RegisterService("helloworld", new(helloworld.Service))
	assert.Nil(t, err)
}
