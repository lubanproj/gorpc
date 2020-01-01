package codec

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc/codes"
	"golang.org/x/net/http2"
	"io"
	"math"
	"sync"
)

type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte, interface{}) error
}

const FrameHeadLen = 9

func GetCodec(name string) Codec {
	return codecMap[name]
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

	header := http2.FrameHeader{
		Length : uint32(len(data)),
	}

	head, err := DefaultSerialization.Marshal(header)
	if err != nil {
		return nil, codes.ServerDecodeError
	}
	return append(head, data ...), nil
}


func (c *defaultCodec) Decode(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	header, err := http2.ReadFrameHeader(reader)
	if err != nil {
		return codes.ServerDecodeError
	}
	msg := make([]byte, header.Length - FrameHeadLen)
	n, err := io.ReadFull(reader,msg)
	if err != nil || n != int(header.Length - FrameHeadLen) {
		return codes.ServerDecodeError
	}
	return DefaultSerialization.Unmarshal(msg[FrameHeadLen:header.Length], v)
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