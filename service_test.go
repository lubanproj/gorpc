package gorpc

import (
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.Serve(opts)
		wg.Done()
	}()
	wg.Wait()
	s.Close()
}
