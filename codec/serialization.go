package codec

import "github.com/golang/protobuf/proto"

type Serialization interface {
	Marshal(interface{})([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var DefaultSerialization = New()

var New =  func() Serialization {
	return &defaultSerialization{}
}

type defaultSerialization struct {}

func (d *defaultSerialization) Marshal(v interface{}) ([]byte, error) {
	if pm, ok := v.(proto.Marshaler); ok {
		// 可以 marshal 自身，无需 buffer
		return pm.Marshal()
	}
	buffer := bufferPool.Get().(*cachedBuffer)
	protoMsg := v.(proto.Message)
	lastMarshaledSize := make([]byte, 0, buffer.lastMarshaledSize)
	buffer.SetBuf(lastMarshaledSize)
	buffer.Reset()

	if err := buffer.Marshal(protoMsg); err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	buffer.lastMarshaledSize = upperLimit(len(data))

	return data, nil
}

func (d *defaultSerialization) Unmarshal(data []byte, v interface{}) error {
	protoMsg := v.(proto.Message)
	protoMsg.Reset()

	if pu, ok := protoMsg.(proto.Unmarshaler); ok {
		// 可以 unmarshal 自身，无需 buffer
		return pu.Unmarshal(data)
	}

	buffer := bufferPool.Get().(*cachedBuffer)
	buffer.SetBuf(data)
	err := buffer.Unmarshal(protoMsg)
	buffer.SetBuf(nil)
	bufferPool.Put(buffer)
	return err
}
