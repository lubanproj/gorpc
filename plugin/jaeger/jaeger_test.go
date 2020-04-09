package jaeger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJaegerInit(t *testing.T) {
	tracingSvrAddr := "localhost:6831"
	_, err := Init(tracingSvrAddr)
	assert.Nil(t,err)
}
