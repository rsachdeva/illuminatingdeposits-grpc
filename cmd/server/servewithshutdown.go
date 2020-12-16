package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func serveWithShutdown(s *grpc.Server, lis net.Listener, mt *mongo.Client, ctx context.Context) {
	// Start the service listening for requests.
	fmt.Println("Ready to Serve now; Press Ctrl+C for Graceful shutdown")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error is %#v", err)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh
	quitGracefully(s, lis, mt, ctx)
}

func quitGracefully(s *grpc.Server, lis net.Listener, mt *mongo.Client, ctx context.Context) {
	disconnectMongodb(mt, ctx)
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener...")
	lis.Close()
	fmt.Println("End of program")
}
