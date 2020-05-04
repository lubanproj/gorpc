package main

import (
	"net/http"
	"time"

	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/plugin/jaeger"
	"github.com/lubanproj/gorpc/testdata"

	_ "net/http/pprof"
)


func main() {

	pprof()

	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithSerializationType("msgpack"),
		gorpc.WithTimeout(time.Millisecond * 2000),
		gorpc.WithTracingSvrAddr("localhost:6831"),
		gorpc.WithTracingSpanName("helloworld.Greeter"),
		gorpc.WithPlugin(jaeger.Name),
	}
	s := gorpc.NewServer(opts ...)
	if err := s.RegisterService("helloworld.Greeter", new(testdata.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}

func pprof() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", http.DefaultServeMux)
	}()
}