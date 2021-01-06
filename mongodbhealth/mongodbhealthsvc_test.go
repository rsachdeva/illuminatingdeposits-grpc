// Adds tests that starts a gRPC server and client tests the mongodb health service with RPC
package mongodbhealth_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testserver"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Conventional test that starts a gRPC server and client test the service with RPC
func TestServiceServer_HealthOk(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	cr := testserver.InitGrpcTLSWithBuffConn(ctx, t, true)
	opts := []grpc.DialOption{grpc.WithContextDialer(testserver.GetBufDialer(cr.Listener)), testcredentials.ClientTlsOption(t), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, "localhost", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)
	mdbresp, err := mdbSvcClient.Health(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Could not check Mongodb status: %v", err)
	}
	log.Printf("response %s", mdbresp.Health)
	require.Equal(t, "MongoDb Ok", mdbresp.Health.Status)
}

// Conventional test that starts a gRPC server and client test the service with RPC
func TestServiceServer_HealthNotOk(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	cr := testserver.InitGrpcTLSWithBuffConn(ctx, t, true)
	opts := []grpc.DialOption{grpc.WithContextDialer(testserver.GetBufDialer(cr.Listener)), testcredentials.ClientTlsOption(t), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, "localhost", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	t.Log("Disconnecting server connection to Mongodb")
	err = cr.MongoClient.Disconnect(ctx)
	if err != nil {
		t.Fatalf("Could not disconnect mongodb: %v", err)
	}

	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)

	mdbresp, err := mdbSvcClient.Health(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Could not check Mongodb status: %v", err)
	}
	log.Printf("response %s", mdbresp.Health)
	require.Equal(t, "MongoDB Not Ready", mdbresp.Health.Status)
}
