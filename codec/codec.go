package codec

import (
	"github.com/diubrother/gorpc/codes"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"io"
	"math"
	"sync"
)

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(io io.Reader, v interface{}) error
}

const FrameHeadLen = 9

func GetCodec(name string) Codec {
	return codecMap[name]
}

var codecMap map[string]Codec

func init() {
	register("proto",&defaultCodec{})
}

func register(name string, codec Codec) {
	codecMap[name] = codec
}

func (c *defaultCodec) Encode(v interface{}) ([]byte, error) {
	body, err := marshal(v)
	if err != nil {
		return nil, codes.ServerDecodeError
	}

	header := http2.FrameHeader{
		Length : uint32(len(body)),
	}
	head, err := marshal(header)
	if err != nil {
		return nil, codes.ServerDecodeError
	}
	return append(head, body ...), nil
}


func (c *defaultCodec) Decode(reader io.Reader, v interface{}) error {
	header, err := http2.ReadFrameHeader(reader)
	if err != nil {
		return codes.ServerDecodeError
	}
	msg := make([]byte, header.Length - FrameHeadLen)
	n, err := io.ReadFull(reader,msg)
	if err != nil || n != int(header.Length - FrameHeadLen) {
		return codes.ServerDecodeError
	}
	return unmarshal(msg[FrameHeadLen:header.Length], v)
}

type defaultCodec struct{}

func marshal(v interface{}) ([]byte, error) {
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

func unmarshal(data []byte, v interface{}) error {
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