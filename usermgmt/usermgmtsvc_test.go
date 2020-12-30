// Adds test that starts a gRPC server and client tests the user mgmt service with RPC
package usermgmt_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbtestconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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
	mdb := mt.Database("depositsmongodb")
	log.Println("Registering gRPC proto UserMgmtService...")
	usermgmtpb.RegisterUserMgmtServiceServer(s, usermgmt.ServiceServer{
		Mdb: mdb,
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
func TestServiceServer_CreateUser(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	mt, pool, resource := mongodbtestconn.Connect(ctx, 1)
	address := "localhost:50055"
	initGRPCServerHTTP2(t, mt, address) // Starting a conventional gRPC server runs on HTTP2

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")
	req := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           "growth@drinnovations.us",
			Roles:           []string{"USER"},
			Password:        "kubernetes",
			PasswordConfirm: "kubernetes",
		},
	}
	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateUser service", err)
	}
	log.Printf("response %s", umresp.User)
	require.Equal(t, umresp.User.Email, "growth@drinnovations.us")

	t.Cleanup(func() {
		fmt.Println("Purging dockertest for mongodb always; unless commented out to examine any data")
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
	})
}
