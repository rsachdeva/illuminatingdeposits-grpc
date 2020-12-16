package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func serveWithShutdown(s *grpc.Server, lis net.Listener, mt *mongo.Client, ctx context.Context) {
	// Start the service listening for requests.
	log.Println("Ready to Serve now; Press Ctrl+C for Graceful shutdown")
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
	mongodbconn.DisconnectMongodb(mt, ctx)
	log.Println("Stopping the server...")
	s.Stop()
	log.Println("Closing the listener...")
	lis.Close()
	log.Println("End of program")
}
