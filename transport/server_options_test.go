package transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithServerAddress(t *testing.T) {
	var sto ServerTransportOptions
	fSto := WithServerAddress("127.0.0.1")
	fSto(&sto)
	assert.Equal(t, "127.0.0.1", sto.Address)
	fSto = WithServerAddress("")
	fSto(&sto)
	assert.Equal(t, "", sto.Address)
}

func TestWithServerNetwork(t *testing.T) {
	var sto ServerTransportOptions
	fSto := WithServerNetwork("tcp")
	fSto(&sto)
	assert.Equal(t, "tcp", sto.Network)
	fSto = WithServerNetwork("")
	fSto(&sto)
	assert.Equal(t, "", sto.Network)
}

func TestWithServerTimeout(t *testing.T) {
	var sto ServerTransportOptions
	fSto := WithServerTimeout(time.Second * time.Duration(2))
	fSto(&sto)
	assert.Equal(t, time.Second*time.Duration(2), sto.Timeout)
}

func TestWithHandler(t *testing.T) {

}

func TestWithSerialization(t *testing.T) {
	var sto ServerTransportOptions
	fSto := WithSerialization("test")
	fSto(&sto)
	assert.Equal(t, "test", sto.Serialization)
	fSto = WithSerialization("")
	fSto(&sto)
	assert.Equal(t, "", sto.Serialization)
}

func TestWithKeepAlivePeriod(t *testing.T) {
	var sto ServerTransportOptions
	fSto := WithKeepAlivePeriod(time.Second * time.Duration(2))
	fSto(&sto)
	assert.Equal(t, time.Second*time.Duration(2), sto.KeepAlivePeriod)
}
