package main

import (
	"log"
	"net"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	ctx, mt := connectMongoDB()
	mdb := mt.Database("depositsmongodb")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}
	s := grpc.NewServer()
	log.Println("Registering InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})
	log.Println("Registering UserMgmtService...")
	usermgmtpb.RegisterUserMgmtServiceServer(s, usermgmt.ServiceServer{
		Mdb: mdb,
	})

	serveWithShutdown(s, lis, mt, ctx)
}
