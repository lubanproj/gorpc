package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var NewTest = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{ServiceName: "test"},
	}
}

func TestGetClientTransport(t *testing.T) {
	var testClientTransport = NewTest()
	clientTransportMap["clinetTest"] = testClientTransport
	clientTranspot := GetClientTransport("clinetTest")
	assert.Equal(t, clientTranspot, testClientTransport)
	clientTranspot = GetClientTransport("test")
	assert.Equal(t, clientTranspot, DefaultClientTransport)
}
