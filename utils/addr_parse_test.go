package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseServicePath(t *testing.T) {
	_, _, err := ParseServicePath("Greeter.Hello")
	assert.NotNil(t,err)
}
