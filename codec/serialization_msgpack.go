package codec

import (
	"bytes"
	"errors"

	"github.com/vmihailenco/msgpack"
)

func init() {
	registerSerialization("msgpack", &MsgpackSerialization{})
}

type MsgpackSerialization struct {}


func (c *MsgpackSerialization) Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("marshal nil interface{}")
	}

	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

func (c *MsgpackSerialization) Unmarshal(data []byte, v interface{}) error {
	if data == nil || len(data) == 0 {
		return errors.New("unmarshal nil or empty bytes")
	}

	decoder := msgpack.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return err
}