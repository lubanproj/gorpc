package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lubanproj/gorpc"

	ghttp "github.com/lubanproj/gorpc/http"
)

func init() {
	ghttp.HandleFunc("GET","/hello", sayHello)
}


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithProtocol("http"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithTimeout(time.Millisecond * 2000),
	}
	s := gorpc.NewServer(opts ...)
	s.ServeHttp()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Write([]byte("world"))
}
