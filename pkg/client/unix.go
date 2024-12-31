package client

import (
	"net"
)

func NewUNIXSocketClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
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
