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
	fmt.Println("jmd")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}
	s := grpc.NewServer()
	interestcalpb.RegisterInterestCalServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error is %#v", err)
	}
}
