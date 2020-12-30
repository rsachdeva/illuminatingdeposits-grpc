// Adds tests that starts a gRPC server and client tests the mongodb health service with RPC
package mongodbhealth_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbtestconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func initGRPCServerHTTP2(t *testing.T, mt *mongo.Client, address string) {
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	log.Println("Registering gRPC proto MongoDBHealthService...")
	mongodbhealthpb.RegisterMongoDbHealthServiceServer(s, mongodbhealth.ServiceServer{
		Mct: mt,
	})

	log.Println("Ready to Serve now")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error is %#v", err)
		}
	}()

	t.Cleanup(func() {
		log.Println("Stopping the server...")
		s.Stop()
		log.Println("Closing the listener...")
		lis.Close()
		log.Println("End of program")
	})
}

// Conventional test that starts a gRPC server and client test the service with RPC
func TestServiceServer_GetMongoDBHealthOk(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mt, pool, resource := mongodbtestconn.Connect(ctx, 1)
	address := "localhost:50053"
	initGRPCServerHTTP2(t, mt, address) // Starting a conventional gRPC server runs on HTTP2

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)
	mdbresp, err := mdbSvcClient.GetMongoDBHealth(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Could not check Mongodb status: %v", err)
	}
	log.Printf("response %s", mdbresp.Health)
	require.Equal(t, mdbresp.Health.Status, "MongoDb Ok")
	t.Cleanup(func() {
		fmt.Println("Purging dockertest for mongodb always; unless commented out to examine any data")
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
	})
}

// Conventional test that starts a gRPC server and client test the service with RPC
func TestServiceServer_GetMongoDBHealthNotOk(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	mt, pool, resource := mongodbtestconn.Connect(ctx, 1)
	address := "localhost:50054"
	initGRPCServerHTTP2(t, mt, address) // Starting a conventional gRPC server runs on HTTP2

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)
	err = pool.Purge(resource)
	if err != nil {
		t.Fatalf("Could not purge container: %v", err)
	}
	mdbresp, err := mdbSvcClient.GetMongoDBHealth(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Could not check Mongodb status: %v", err)
	}
	log.Printf("response %s", mdbresp.Health)
	require.Equal(t, mdbresp.Health.Status, "MongoDB Not Ready")
}
