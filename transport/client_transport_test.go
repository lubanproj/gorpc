package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var NewClientTransport = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{ServiceName: "test"},
	}
}

func TestGetClientTransport(t *testing.T) {
	var testClientTransport = NewClientTransport()
	clientTransportMap["clinetTest"] = testClientTransport
	clientTransport := GetClientTransport("clinetTest")
	assert.Equal(t, clientTransport, testClientTransport)
	clientTransport = GetClientTransport("test")
	assert.Equal(t, clientTransport, DefaultClientTransport)
}
