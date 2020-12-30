package readenv

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func TlsEnabled() bool {
	var enabled bool
	var err error
	// change DEPOSITS_GRPC_SERVICE_TLS env variable in command line or editor
	enabled = true
	if tlsAllowed, ok := os.LookupEnv("DEPOSITS_GRPC_SERVICE_TLS"); ok {
		fmt.Println("tlsAllowed from env is", tlsAllowed)
		enabled, err = strconv.ParseBool(tlsAllowed)
		if err != nil {
			log.Fatal("tls DEPOSITS_GRPC_SERVICE_TLS reading env error")
		}
	}
	return enabled
}
