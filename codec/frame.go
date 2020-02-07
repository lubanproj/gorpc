package codec

import (
	"encoding/binary"
	"github.com/lubanproj/gorpc/codes"
	"io"
	"net"
)

type FrameHeader struct {
	Magic uint8
	Version uint8
	MsgType uint8  // msg type e.g. :   0x0: general req,  0x1: heartbeat
	ReqType uint8  // request type e.g. :   0x0: send and receive,   0x1: send but not receive,  0x2: client stream request, 0x3: server stream request, 0x4: bidirectional streaming request
	CompressType uint8 // compression or not :  0x0: not compression,  0x1: compression
	StreamID uint16    // stream ID
	Length uint32  	// total packet length
	Reserved uint32  // 4 bytes reserved
}

func ReadFrame(conn net.Conn) ([]byte, error) {

	frameHeader := make([]byte, FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != FrameHeadLen || err != nil {
		return nil, err
	}

	// validate magic
	if magic := uint8(frameHeader[0]); magic != Magic {
		return nil, codes.ClientMsgError
	}


	length := binary.BigEndian.Uint32(frameHeader[7:11])

	data := make([]byte, length)
	if num, err := io.ReadFull(conn, data); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, data ...), nil
}
