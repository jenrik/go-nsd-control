package main

import (
	"flag"
	"fmt"
	"log"
	"nsd/pkg/client"
)

const defaultSocket = "/var/run/nsd.sock"

func main() {
	flag.Parse()
	posArgs := flag.Args()

	if len(posArgs) < 1 {
		fmt.Println("Usage: nsd-control <cmd>")
		return
	}

	c, err := client.NewUNIXSocketClient(defaultSocket)
	if err != nil {
		log.Fatal(err)
	}
	defer func(c client.Client) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	switch posArgs[0] {
	case "stop":
		if err := c.Stop(); err != nil {
			log.Fatal(err)
		}
	}
}
