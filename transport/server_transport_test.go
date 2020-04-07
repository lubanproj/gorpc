package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var NewTestServerTransport = func() ServerTransport {
	return &serverTransport{
		opts : &ServerTransportOptions{Address:"127.0.0.1"},
	}
}

func TestGetServerTransport(t *testing.T) {
	var testServerTransport = NewTestServerTransport()
	serverTransportMap["serverTest"] = testServerTransport
	serverTransport := GetServerTransport("serverTest")
	assert.Equal(t, serverTransport, testServerTransport)
	serverTransport = GetServerTransport("test")
	assert.Equal(t, serverTransport, DefaultServerTransport)
}


