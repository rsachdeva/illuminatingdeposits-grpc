package main

import (
	"fmt"
	"os"
)

func svcAddress(tls bool) string {
	address := "localhost"
	port := "50052"
	addr := os.Getenv("DEPOSITS_GRPC_SERVICE_ADDRESS")
	if !(addr == "" || addr == "localhost") {
		address = addr
		if tls {
			port = "443"
		}
	}
	return fmt.Sprintf("%v:%v", address, port)
}
