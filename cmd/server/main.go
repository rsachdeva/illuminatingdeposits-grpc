package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting ServiceServer...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Registering InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})

	fmt.Println("Ready to serve now")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error is %#v", err)
	}
}
