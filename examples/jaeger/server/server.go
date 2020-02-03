package main

import (
	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"
	"github.com/lubanproj/gorpc/plugin/consul"
	"time"
)


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithSerializationType("msgpack"),
		gorpc.WithTimeout(time.Millisecond * 2000000),
		gorpc.WithSelectorSvrAddr("localhost:8500"),
		gorpc.WithPlugin(consul.Name),
	}
	s := gorpc.NewServer(opts ...)
	if err := s.RegisterService("helloworld.Greeter", new(helloworld.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}