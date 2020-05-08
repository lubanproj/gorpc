package codec

import (
	"fmt"
	"testing"

	"github.com/lubanproj/gorpc/protocol"
	"github.com/stretchr/testify/assert"
)

func TestMsgpackSerializationMarshal(t *testing.T) {
	msgSer := &MsgpackSerialization{}
	data, err := msgSer.Marshal(nil)
	assert.Nil(t, err)
	fmt.Println(string(data), err)
	err = msgSer.Unmarshal(data, &protocol.Response{})
	assert.Nil(t,err)
	err = msgSer.Unmarshal(nil, &protocol.Response{})
	assert.NotNil(t,err)
	fmt.Println(err)
}
