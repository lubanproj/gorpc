package codec

import "github.com/golang/protobuf/proto"

type Serialization interface {
	Marshal(interface{})([]byte, error)
	Unmarshal([]byte, interface{}) error
}

const (
	Proto         =  "proto" 	// protobuf
	MsgPack 	  =  "msgpack"   // msgpack
	Json      	  =  "json" 	// json
)

var serializationMap = make(map[string]Serialization)

var DefaultSerialization = NewSerialization()

var NewSerialization = func () Serialization {
	return &pbSerialization{}
}

func init() {
	registerSerialization("proto", DefaultSerialization)
}

func registerSerialization(name string, serialization Serialization) {
	if serializationMap == nil {
		serializationMap = make(map[string]Serialization)
	}
	serializationMap[name] = serialization
}

func GetSerialization(name string) Serialization {
	if v, ok := serializationMap[name]; ok {
		return v
	}
	return DefaultSerialization
}

type pbSerialization struct {}

func (d *pbSerialization) Marshal(v interface{}) ([]byte, error) {
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
	buffer.SetBuf(nil)
	bufferPool.Put(buffer)

	return data, nil
}

func (d *pbSerialization) Unmarshal(data []byte, v interface{}) error {
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
