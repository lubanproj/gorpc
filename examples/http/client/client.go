package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	rsp, err := http.Get("http://127.0.0.1:8000/hello")
	if err != nil {
		fmt.Println(err)
	}

    defer rsp.Body.Close()
	
    rspbody, err := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(rspbody), err)
}
