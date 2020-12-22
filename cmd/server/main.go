package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/readenv"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"google.golang.org/grpc"
)

const (
	// https://stackoverflow.com/questions/64093550/grpc-server-not-working-in-docker-compose
	address = "0.0.0.0:50052"
)

func main() {
	tls := readenv.TlsEnabled()
	fmt.Println("tls is", tls)

	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mt := mongodbconn.ConnectMongoDB(ctx, 10)
	mdb := mt.Database("depositsmongodb")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}

	var opts []grpc.ServerOption
	if tls {
		opts = tlsOpts(opts)
	}
	s := grpc.NewServer(opts...)
	log.Println("Registering MongoDBHealthService...")
	mongodbhealthpb.RegisterMongoDbHealthServiceServer(s, mongodbhealth.ServiceServer{
		Mct: mongodbconn.ConnectMongoDB(ctx, 2),
	})
	log.Println("Registering UserMgmtService...")
	usermgmtpb.RegisterUserMgmtServiceServer(s, usermgmt.ServiceServer{
		Mdb: mdb,
	})
	log.Println("Registering UserAuthenticationService...")
	userauthnpb.RegisterUserAuthnServiceServer(s, userauthn.ServiceServer{
		Mdb: mdb,
	})
	log.Println("Registering InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})

	serveWithShutdown(ctx, s, lis, mt)
}
