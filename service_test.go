package gorpc

import (
	"testing"
	"time"
)

func TestServe(t *testing.T) {

	opts := &ServerOptions {
		network: "tcp",
		address: "127.0.0.1:8000",
		timeout: time.Millisecond * 1000,
	}
	s := &service{}
	go func() {
		s.Serve(opts)
	}()
	time.Sleep(time.Millisecond * 1000)
	s.Close()

}
