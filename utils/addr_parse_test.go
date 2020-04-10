package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServicePath(t *testing.T) {
	_, _, err := ParseServicePath("Greeter.Hello")
	assert.NotNil(t, err)

	seriveName, method, err := ParseServicePath("SerivceName/Method")
	assert.Equal(t, seriveName, "erivceName")
	assert.Equal(t, method, "Method")
	assert.Equal(t, err, nil)
}
