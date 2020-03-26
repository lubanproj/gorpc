package jaeger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJaegerInit(t *testing.T) {
	tracingSvrAddr := "localhost:6831"
	_, err := Init(tracingSvrAddr)
	assert.Nil(t,err)
}
