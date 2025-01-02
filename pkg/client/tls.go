package client

import (
	"crypto/tls"
	"crypto/x509"
	"net"
)

func NewSimpleTLSClient(addr net.Addr, serverCA *x509.CertPool, clientCert tls.Certificate) (*Client, error) {
	tlsConfig := &tls.Config{
		RootCAs:      serverCA,
		Certificates: []tls.Certificate{clientCert},
	}
	conn, err := tls.Dial(addr.Network(), addr.String(), tlsConfig)
	if err != nil {
		return nil, err
	}

	client := &Client{
		socket: conn,
	}
	if err := client.init(); err != nil {
		return nil, err
	}
	return client, err
}

//goland:noinspection GoUnusedExportedFunction
func NewTLSClient(addr net.TCPAddr, tlsConfig *tls.Config) (*Client, error) {
	conn, err := tls.Dial(addr.Network(), addr.String(), tlsConfig)
	if err != nil {
		return nil, err
	}

	client := &Client{
		socket: conn,
	}
	if err := client.init(); err != nil {
		return nil, err
	}
	return client, err
}
