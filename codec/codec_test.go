package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
