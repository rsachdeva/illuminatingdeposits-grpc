package main

import (
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := interestcalpb.NewInterestCalServiceClient(conn)
	fmt.Printf("client created c is %f", c)
}
