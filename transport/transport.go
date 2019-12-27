// 网络通讯层，负责底层网络通讯，主要包括 tcp && udp 两种协议方式实现
package transport

import "context"

type ServerTransport interface {
	// 请求的监听和处理
	ListenAndServe(context.Context, ...ServerTransportOption) error
}

type ClientTransport interface {
	// 发送请求
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}


