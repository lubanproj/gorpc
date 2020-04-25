package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"strings"
	"syscall"

	"github.com/lubanproj/gorpc/codes"
)

// tlsAuth defines the implementation of TLS authentication
// and implements TransportAuth, PerRPCAuth, AuthInfo
type tlsAuth struct {
	config *tls.Config
	state tls.ConnectionState
}

// AuthType returns the protocol name
func (t *tlsAuth) AuthType() string {
	return "tls"
}

// NewClientTLSAuthFromFile instantiates client-side authentication information
// with certificates and service names
func NewClientTLSAuthFromFile(certFile, serverName string) (TransportAuth, error) {
	cert , err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		return nil, codes.ClientCertFailError
	}
	conf := &tls.Config {
		ServerName: serverName,
		RootCAs: cp,
	}
	return &tlsAuth{config : conf}, nil
}

// NewServerTLSAuthFromFile generates server-side authentication information
// with certificates and keys
func NewServerTLSAuthFromFile(certFile, keyFile string) (TransportAuth, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, codes.ClientCertFailError
	}
	conf := &tls.Config{
		Certificates:[]tls.Certificate{cert},
	}
	return &tlsAuth{config : conf}, nil
}

// ClientHandshake implements the client's handshake
func (t *tlsAuth) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (net.Conn, AuthInfo, error) {
	// 防止使用不同的 endpoints 时 ServerName 被污染
	cfg := cloneTLSConfig(t.config)
	if cfg.ServerName == "" {
		colonPos := strings.LastIndex(authority, ":")
		if colonPos == -1 {
			colonPos = len(authority)
		}
		cfg.ServerName = authority[:colonPos]
	}
	conn := tls.Client(rawConn, cfg)
	errChan := make(chan error, 1)

	go func() {
		errChan <- conn.Handshake()
	}()
	select {
	case err := <- errChan :
		if err != nil {
			return nil, nil, err
		}
	case <- ctx.Done() :
		return nil, nil, ctx.Err()
	}

	return WrapConn(rawConn,conn) , &tlsAuth{state : conn.ConnectionState()}, nil
}

// the ServerHandshake implements the server handshake
func (t *tlsAuth) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) {
	conn := tls.Server(rawConn, t.config)
	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	return WrapConn(rawConn,conn), &tlsAuth{state : conn.ConnectionState()}, nil
}

func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	if cfg == nil {
		return &tls.Config{}
	}

	return cfg.Clone()
}


// WrapConn synthesizes two conn's into one conn
func WrapConn(rawConn, newConn net.Conn) net.Conn {
	sysConn, ok := rawConn.(syscall.Conn)
	if !ok {
		return newConn
	}
	return &wrapperConn{
		Conn:    newConn,
		sysConn: sysConn,
	}
}

type sysConn = syscall.Conn

type wrapperConn struct {
	net.Conn
	// sysConn is a type alias of syscall.Conn. It's necessary because the name
	// `Conn` collides with `net.Conn`.
	sysConn
}


