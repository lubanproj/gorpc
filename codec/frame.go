package codec

import (
	"encoding/binary"
	"github.com/lubanproj/gorpc/codes"
	"io"
	"net"
)

type FrameHeader struct {
	Magic uint16  // 魔数
	Version uint8 // 版本号
	Type uint8   // 帧类型
	Length uint32  // 数据包总长度
	HeaderLength uint32  // 数据包包头长度
	Reserved uint32  // 4 字节保留字段
}

func ReadFrame(conn net.Conn) ([]byte, error) {

	frameHeader := make([]byte, FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != FrameHeadLen || err != nil {
		return nil, err
	}

	// 校验魔数
	if magic := binary.BigEndian.Uint16(frameHeader[0:2]); magic != Magic {
		return nil, codes.ClientMsgError
	}


	length := binary.BigEndian.Uint32(frameHeader[4:8])

	data := make([]byte, length)
	if num, err := io.ReadFull(conn, data); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, data ...), nil
}
