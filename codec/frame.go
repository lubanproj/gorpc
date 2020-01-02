package codec

import (
	"encoding/binary"
	"io"
	"net"
)

type FrameHeader struct {
	Magic uint16  // 魔数
	Version uint8 // 版本号
	Type uint8   // 帧类型
	Length uint32  // 消息总长度
	HeaderLength uint32  // 包头长度
	Reserved uint32  // 4 字节保留字段
}

func ReadFrame(conn net.Conn) ([]byte, error) {

	// TODO 帧头校验

	buffer := make([]byte, FrameHeadLen)
	if num ,err := io.ReadFull(conn, buffer); num != FrameHeadLen || err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(buffer[4:8])

	data := make([]byte, length)
	if num, err := io.ReadFull(conn, data); num != FrameHeadLen || err != nil {
		return nil, err
	}

	return append(buffer, data ...), nil
}