package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCodec(t *testing.T) {
	RegisterCodec("testCodec", nil)

	codec := GetCodec("testCodec")
	assert.Equal(t, codec, nil)
}


func TestDefaultCodec_Decode(t *testing.T) {

}

func TestDefaultCodec_Encode(t *testing.T) {

}
