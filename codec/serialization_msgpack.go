package codec

import (
	"bytes"
	"github.com/vmihailenco/msgpack"
)


func init() {
	registerSerialization("msgpack", &MsgpackSerialization{})
}

type MsgpackSerialization struct {}


func (c *MsgpackSerialization) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

func (c *MsgpackSerialization) Unmarshal(data []byte, v interface{}) error {
	decoder := msgpack.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return err
}