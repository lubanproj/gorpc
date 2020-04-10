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

	seriveName, method, err := ParseServicePath("/serivceName/method")
	assert.Equal(t, seriveName, "serivceName")
	assert.Equal(t, method, "method")
	assert.Equal(t, err, nil)
}
