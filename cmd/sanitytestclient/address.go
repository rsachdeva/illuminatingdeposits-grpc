package main

import (
	"fmt"
	"os"
)

func svcAddress(tls bool) string {
	address := "localhost"
	port := "50052"
	if addr := os.Getenv("DEPOSITS_GRPC_SERVICE_ADDRESS"); addr != "" {
		address = addr
		if tls {
			port = "443"
		}
	}
	return fmt.Sprintf("%v:%v", address, port)
}
