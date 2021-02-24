package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func tlsOpts() []grpc.ServerOption {
	certFile := "conf/tls/servercrtto.pem"
	keyFile := "conf/tls/serverkeyto.pem"
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("tls error failed loading certificates %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	return opts
}
