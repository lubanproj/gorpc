package transport

import (
	"context"
)

func (s *serverTransport) ListenAndServeUdp(ctx context.Context, opts ...ServerTransportOption) error {

	//lis, err := net.ListenPacket(s.opts.Network, s.opts.Address)
	//if err != nil {
	//	return err
	//}

	return nil

}
