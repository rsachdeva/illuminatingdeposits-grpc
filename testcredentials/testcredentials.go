package testcredentials

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func ServerTlsOption(opts []grpc.ServerOption) []grpc.ServerOption {
	certFile := path("tls/servercrtto.pem")
	keyFile := path("tls/serverkeyto.pem")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("tls error failed loading certificates %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}
	return opts
}

// export GODEBUG=x509ignoreCN=0
func ClientTlsOption(t *testing.T) grpc.DialOption {
	certFile := path("tls/cacrtto.pem")
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("loading certificate error is %v", err)
	}
	opt := grpc.WithTransportCredentials(creds)
	return opt
}

func path(rel string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(currentFile)
	pa := filepath.Join(basepath, rel)
	log.Printf("pa for test is %v", pa)
	return pa
}
