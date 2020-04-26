package client

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {

	var wg sync.WaitGroup
	var ch = make(chan struct{})
	wg.Add(1)
	go func() {
		serverOpts := []gorpc.ServerOption{
			gorpc.WithAddress("localhost:8000"),
			gorpc.WithNetwork("tcp"),
			gorpc.WithSerializationType("msgpack"),
			gorpc.WithTimeout(time.Millisecond * 2000),
		}
		s := gorpc.NewServer(serverOpts ...)
		if err := s.RegisterService("helloworld.Greeter", new(helloworld.Service)); err != nil {
			panic(err)
		}
		wg.Done()

		go func() {
			s.Serve()
		}()

		<- ch
		s.Close()
	}()

	wg.Wait()

	opts := []Option {
		WithTarget("localhost:8000"),
		WithNetwork("tcp"),
		WithTimeout(2000 * time.Millisecond),
		WithSerializationType("msgpack"),
	}
	c := DefaultClient
	req := &helloworld.HelloRequest{
		Msg: "hello",
	}
	rsp := &helloworld.HelloReply{}

	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)

	close(ch)

	assert.Nil(t, err)
}

