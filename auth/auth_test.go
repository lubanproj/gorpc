package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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