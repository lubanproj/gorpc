package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServicePath(t *testing.T) {
	_, _, err := ParseServicePath("Greeter.Hello")
	assert.NotNil(t, err)

	_, _, err = ParseServicePath("Greeter/Hello")
	assert.NotNil(t, err)

	seriveName, method, err := ParseServicePath("/Greeter/Hello")
	assert.Equal(t, seriveName, "Greeter")
	assert.Equal(t, method, "Hello")
	assert.Equal(t, err, nil)
}
