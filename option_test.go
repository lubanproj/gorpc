package gorpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithAddress(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithAddress("127.0.0.1")
	fServerops(&serverops)
	assert.Equal(t, "127.0.0.1", serverops.address)
	fServerops = WithAddress("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.address)
}

func TestWithNetwork(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithNetwork("test")
	fServerops(&serverops)
	assert.Equal(t, "test", serverops.network)
	fServerops = WithNetwork("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.network)
}

func TestWithProtocol(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithProtocol("http")
	fServerops(&serverops)
	assert.Equal(t, "http", serverops.protocol)
	fServerops = WithProtocol("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.protocol)
}

func TestWithTimeout(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithTimeout(time.Second * time.Duration(2))
	fServerops(&serverops)
	assert.Equal(t, time.Second*time.Duration(2), serverops.timeout)
}

func TestWithSerializationType(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithSerializationType("test")
	fServerops(&serverops)
	assert.Equal(t, "test", serverops.serializationType)
	fServerops = WithSerializationType("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.serializationType)
}

func TestWithSelectorSvrAddr(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithSelectorSvrAddr("127.0.0.1")
	fServerops(&serverops)
	assert.Equal(t, "127.0.0.1", serverops.selectorSvrAddr)
	fServerops = WithSelectorSvrAddr("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.selectorSvrAddr)
}

func TestWithPlugin(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithPlugin("test")
	fServerops(&serverops)
	assert.Equal(t, []string{"test"}, serverops.pluginNames)
	serverops.pluginNames=[]string(nil)
	fServerops = WithPlugin("test","another_test")
	fServerops(&serverops)
	assert.Equal(t, []string{"test","another_test"}, serverops.pluginNames)
	serverops.pluginNames=[]string(nil)
	fServeropsNew := WithPlugin()
	fServeropsNew(&serverops)
	assert.Equal(t,[]string(nil), serverops.pluginNames)
}

func TestWithInterceptor(t *testing.T) {

}

func TestWithTracingSvrAddr(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithTracingSvrAddr("127.0.0.1")
	fServerops(&serverops)
	assert.Equal(t, "127.0.0.1", serverops.tracingSvrAddr)
	fServerops = WithTracingSvrAddr("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.tracingSvrAddr)
}

func TestWithTracingSpanName(t *testing.T) {
	var serverops ServerOptions
	fServerops := WithTracingSpanName("test")
	fServerops(&serverops)
	assert.Equal(t, "test", serverops.tracingSpanName)
	fServerops = WithTracingSpanName("")
	fServerops(&serverops)
	assert.Equal(t, "", serverops.tracingSpanName)
}
