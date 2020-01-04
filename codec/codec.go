package codec

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc/protocol"
	"math"
	"sync"
)

type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

const FrameHeadLen = 16
const Magic = 0x1111
const Version = 0

func GetCodec(name string) Codec {
	if v, ok := codecMap[name]; ok {
		return v
	}
	return DefaultCodec
}

var codecMap = make(map[string]Codec)

var DefaultCodec = NewCodec()

var NewCodec = 	func () Codec {
	return &defaultCodec{}
}

func init() {
	registerCodec("proto", DefaultCodec)
}

func registerCodec(name string, codec Codec) {
	if codecMap == nil {
		codecMap = make(map[string]Codec)
	}
	codecMap[name] = codec
}

func (c *defaultCodec) Encode(data []byte) ([]byte, error) {

	requestHeader := &protocol.Request{}
	requestHeader.Length = uint32(len(data))

	reqHeadBuf, err := proto.Marshal(requestHeader)
	if err != nil {
		return nil, err
	}

	totalLen := FrameHeadLen + len(reqHeadBuf) + len(data)
	buffer := bytes.NewBuffer(make([]byte, 0, totalLen))

	frame := FrameHeader{
		Magic : Magic,
		Version : Version,
		Type : 0x1,
		Length: uint32(totalLen),
		HeaderLength: uint32(len(reqHeadBuf)),
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Magic); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Type); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.HeaderLength); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Reserved); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, reqHeadBuf); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}


func (c *defaultCodec) Decode(data []byte) ([]byte,error) {

	totalLen := binary.BigEndian.Uint32(data[4:8])
	headerLen := binary.BigEndian.Uint32(data[8:12])

	return data[totalLen + headerLen :], nil
}

type defaultCodec struct{}

func upperLimit(val int) uint32 {
	if val > math.MaxInt32 {
		return uint32(math.MaxInt32)
	}
	return uint32(val)
}

var bufferPool = &sync.Pool{
	New : func() interface {} {
		return &cachedBuffer {
			Buffer : proto.Buffer{},
			lastMarshaledSize : 16,
		}
	},
}

type cachedBuffer struct {
	proto.Buffer
	lastMarshaledSize uint32
}