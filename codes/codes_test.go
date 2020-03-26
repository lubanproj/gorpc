package codes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodes(t *testing.T) {
	err := NewFrameworkError(1001, "server timeout")
	assert.Equal(t, err.Type, FrameworkError)
	err = New(-1, "params error")
	assert.NotNil(t,err)
	assert.Equal(t, err.Type, BusinuessError)
}
