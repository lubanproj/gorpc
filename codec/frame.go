package codec

import (
	"encoding/binary"
	"github.com/lubanproj/gorpc/codes"
	"io"
	"net"
)

type FrameHeader struct {
	Magic uint8  // 魔数
	Version uint8 // 版本号
	MsgType uint8  // 消息类型   0x0 普通请求  0x1 心跳包
	ReqType uint8  // 请求类型   0x0 一发一收  0x1 只发不收 0x2 客户端流式请求 0x3 服务端流式请求 0x4 双向流式请求
	CompressType uint8 // 是否压缩 0x0 不压缩 0x1 压缩
	StreamID uint16    // 流 ID
	Length uint32  	// 数据包总长度
	Reserved uint32  // 4 字节保留字段
}

func ReadFrame(conn net.Conn) ([]byte, error) {

	frameHeader := make([]byte, FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != FrameHeadLen || err != nil {
		return nil, err
	}

	// 校验魔数
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
