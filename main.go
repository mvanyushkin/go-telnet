package main

import (
	"flag"
	"fmt"
	"github.com/mvanyushkin/go-telnet/client"
	"time"
)

func main() {
	timeout := flag.String("timeout", "60s", "connection timeout")
	flag.Parse()
	duration, _ := time.ParseDuration(*timeout)
	args := flag.Args()
	ip := "127.0.0.1"
	port := "23"
	if len(args) > 0 {
		ip = args[0]
	}

	if len(args) > 1 {
		port = args[1]
	}

	client.RunClient(fmt.Sprintf("%v:%v", ip, port), duration)
}
