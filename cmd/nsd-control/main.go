package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"nsd/pkg/client"
	"os"
	"strings"
)

const defaultSocket = "/var/run/nsd.sock"

type SimpleAddr string

func (s SimpleAddr) Network() string {
	return "tcp"
}

func (s SimpleAddr) String() string {
	if !strings.Contains(string(s), ":") {
		return string(s) + ":8952"
	}
	return string(s)
}

func main() {
	connUrl := flag.String("i", defaultSocket, "server address and port, or socket path")
	caPath := flag.String("ca", "", "Server CA certificate path")
	clientCertPath := flag.String("client-cert", "", "Client certificate path")
	clientKeyPath := flag.String("client-key", "", "Client private key path")
	flag.Parse()
	posArgs := flag.Args()

	if len(posArgs) < 1 {
		_, _ = os.Stderr.WriteString("Usage: nsd-control <cmd>\n\n")
		flag.PrintDefaults()
		return
	}

	var c *client.Client
	if _, err := os.Stat(*connUrl); err == nil {
		c, err = client.NewUNIXSocketClient(*connUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// CA that signed the server certificate
		if caPath == nil || *caPath == "" {
			log.Fatal("missing -ca")
			return
		}
		caCert, err := os.ReadFile(*caPath)
		if err != nil {
			log.Fatal(err)
			return
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)

		// Client certificate
		if clientCertPath == nil || *clientCertPath == "" {
			log.Fatal("missing -client-cert")
			return
		}
		if clientKeyPath == nil || *clientKeyPath == "" {
			log.Fatal("missing -client-key")
			return
		}
		clientCert, err := tls.LoadX509KeyPair(*clientCertPath, *clientKeyPath)
		if err != nil {
			log.Fatal(err)
			return
		}

		c, err = client.NewSimpleTLSClient(SimpleAddr(*connUrl), caPool, clientCert)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	defer mustClose(c)

	doCommand(c, posArgs[0])
}

func mustClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}

func doCommand(c *client.Client, cmd string) {
	switch cmd {
	case "stop":
		if err := c.Stop(); err != nil {
			log.Fatal(err)
		}
	case "status":
		if lines, err := c.Status(); err != nil {
			log.Fatal(err)
		} else {
			for _, l := range lines {
				fmt.Printf("%s\n", l)
			}
		}
	}
}
