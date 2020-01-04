package transport

import (
	"context"
	"github.com/lubanproj/gorpc/codec"
	"time"
)

type ServerTransportOptions struct{
	Address string // 地址，格式例如 ip://127.0.0.1：8080
	Network string  // 网络类型
	Timeout time.Duration  // 传输层请求超时时间，默认为 2 min
	Codec codec.Codec    // 解析数据帧和请求体
	Serialization codec.Serialization  // 序列化方式 json/proto
	Handler Handler
}

type Handler func (context.Context, []byte) ([]byte, error)

type ServerTransportOption func(*ServerTransportOptions)

func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Address = address
	}
}

func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Network = network
	}
}

func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Timeout = timeout
	}
}

func WithServerCodec(codec codec.Codec) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Codec = codec
	}
}

func WithServerSerialization(serialization codec.Serialization) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Serialization = serialization
	}
}