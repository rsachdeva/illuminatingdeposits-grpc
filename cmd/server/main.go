package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

func main() {
	tls := true
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

	fmt.Println("tls option is ", tls)
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
	log.Println("Registering InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})

	serveWithShutdown(ctx, s, lis, mt)
}
