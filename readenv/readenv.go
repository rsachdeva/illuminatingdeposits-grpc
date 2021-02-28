package readenv

import (
	"log"
	"os"
	"strconv"
)

func TlsEnabled() bool {
	enabled := envEnabled("DEPOSITS_GRPC_SERVICE_TLS", true)
	log.Println("TlsEnabled from env is ", enabled)
	return enabled
}

func MessageBrokerLogEnabled() bool {
	enabled := envEnabled("DEPOSITS_GRPC_SERVICE_MESSAGE_BROKER_LOG", false)
	log.Println("MessageBrokerLogEnabled from env is ", enabled)
	return enabled
}

func envEnabled(name string, defaultValue bool) bool {
	if envValue, ok := os.LookupEnv(name); ok {
		enabled, err := strconv.ParseBool(envValue)
		if err != nil {
			log.Fatalf("name %v env error", name)
		}
		return enabled
	}
	return defaultValue
}
