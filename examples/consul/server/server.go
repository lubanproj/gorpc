package main

import (
	"time"

	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/plugin/consul"
	"github.com/lubanproj/gorpc/testdata"
)


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithSerializationType("msgpack"),
		gorpc.WithTimeout(time.Millisecond * 2000),
		gorpc.WithSelectorSvrAddr("localhost:8500"),
		gorpc.WithPlugin(consul.Name),
	}
	s := gorpc.NewServer(opts ...)
	if err := s.RegisterService("helloworld.Greeter", new(testdata.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}