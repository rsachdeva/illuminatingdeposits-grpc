package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func serveWithShutdown(ctx context.Context, s *grpc.Server, lis net.Listener, mt *mongo.Client) {
	// Start the service listening for requests.
	log.Println("Ready to Serve now")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error is %#v", err)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh
	fmt.Println("getting os interrupt")
	quitGracefully(ctx, s, lis, mt)
}

func quitGracefully(ctx context.Context, s *grpc.Server, lis net.Listener, mt *mongo.Client) {
	mongodbconn.DisconnectMongodb(ctx, mt)
	log.Println("Stopping the server...")
	s.Stop()
	log.Println("Closing the listener...")
	lis.Close()
	log.Println("End of program")
}
