package codec

import (
	"fmt"
	"testing"

	"github.com/lubanproj/gorpc/protocol"
	"github.com/stretchr/testify/assert"
)

func TestPbSerializationMarshal(t *testing.T) {
	pbSer := &pbSerialization{}
	data, err := pbSer.Marshal(nil)
	assert.NotNil(t, err)
	fmt.Println(string(data), err)
	err = pbSer.Unmarshal(data, &protocol.Response{})
	assert.NotNil(t, err)
	err = pbSer.Unmarshal(nil, &protocol.Response{})
	assert.NotNil(t, err)
	fmt.Println(err)
}
