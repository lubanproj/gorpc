package transport

import (
	"context"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/stream"
	"net"
	"time"
)


func (s *serverTransport) ListenAndServeUdp(ctx context.Context, opts ...ServerTransportOption) error {

	conn , err := net.ListenPacket(s.opts.Network, s.opts.Address)
	defer conn.Close()

	buffer := make([]byte, 65536)
	if err != nil {
		return err
	}

	var tempDelay time.Duration

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err();
		default:
		}

		num, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		req := buffer[:num]

		go func() {

			// build stream
			ctx, _ := stream.NewServerStream(ctx)

			if err := s.handleUdpConn(ctx, conn, addr, req); err != nil {
				log.Error("gorpc handle udp conn error, %v", err)
			}

		}()


	}


	return nil

}


func (s *serverTransport) handleUdpConn(ctx context.Context, conn net.PacketConn, addr net.Addr, req []byte) error {

	rsp , err := s.handle(ctx, req)
	if err != nil{
		return err
	}

	_, err = conn.WriteTo(rsp, addr)
	return err
}