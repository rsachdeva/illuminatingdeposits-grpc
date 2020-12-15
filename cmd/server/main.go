package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"google.golang.org/grpc"
)

// server is used to implement ??.
type server struct {
	interestcalpb.UnimplementedInterestCalServiceServer
}


func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Registering InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, &server{})

	fmt.Println("Ready to serve now")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error is %#v", err)
	}
}
