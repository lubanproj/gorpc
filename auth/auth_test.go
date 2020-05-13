package auth

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerTLSAuthFromFile(t *testing.T) {
	transportAuth, err := NewServerTLSAuthFromFile("../testdata/server.crt", "../testdata/server.key")
	assert.Nil(t, err, nil)
	fmt.Printf("server conf : %v \n", transportAuth.(*tlsAuth))
}

func TestNewClientTLSAuthFromFile(t *testing.T) {
	transportAuth, err := NewClientTLSAuthFromFile("../testdata/server.crt", "helloworld")
	assert.Nil(t, err, nil)
	fmt.Printf("client conf : %v \n", transportAuth.(*tlsAuth))
}

func TestClientHandshake(t *testing.T) {
	ch := make(chan int)
	go newServer(t, ch)
	<- ch
	conn , err := net.Dial("tcp","127.0.0.1:8002")
	assert.Nil(t, err)
	defer conn.Close()
	tAuth , err := NewClientTLSAuthFromFile("../testdata/server.crt","lubanstudio.cn")
	assert.Nil(t, err)
	var ctx = context.Background()
	wrapperConn, _ , err := tAuth.ClientHandshake(ctx, "lubanstudio.cn", conn)
	assert.Nil(t, err)

	data , err := wrapperConn.Write([]byte("hello\n"))
	assert.Nil(t, err)

	buf := make([]byte, 100)
	data, err = wrapperConn.Read(buf)
	assert.Nil(t, err)

	fmt.Println(string(buf[:data]))
}

func newServer(t *testing.T, ch chan int) {
	ln, err := net.Listen("tcp",":8002")
	assert.Nil(t, err)
	ch <- 1
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		assert.Nil(t, err)
		tlsAuth , err := NewServerTLSAuthFromFile("../testdata/server.crt","../testdata/server.key")
		assert.Nil(t, err)
		wrapperConn, _, err := tlsAuth.ServerHandshake(conn)
		assert.Nil(t, err)
		go handleConn(t,wrapperConn)
	}
}

func handleConn(t *testing.T, conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				assert.Nil(t, err)
			}
		}
		fmt.Println(msg)
		conn.Write([]byte("world\n"))
	}
}