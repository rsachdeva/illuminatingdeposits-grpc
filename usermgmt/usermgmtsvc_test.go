// Adds test that starts a gRPC server and client tests the user mgmt service with RPC
package usermgmt_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbconntest"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func initGRPCServerHTTP2(t *testing.T, address string) {
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	mt, pool, resource := mongodbconntest.Connect(ctx, 1)
	mdb := mt.Database("depositsmongodb")
	coll := mdb.Collection("user")

	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}
	name, err := coll.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal("err in index creation with CreateOne is ", err)
	}
	log.Println("Index created with name", name)
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
		fmt.Println("Purging dockertest for mongodb always; unless commented out to examine any data")
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
		log.Println("End of program")
	})
}

func TestServiceServer_CreateUser(t *testing.T) {
	t.Parallel()

	address := "localhost:50055"
	initGRPCServerHTTP2(t, address) // Starting a conventional gRPC server runs on HTTP2

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

	req = usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User2",
			Email:           "growth@drinnovations.us",
			Roles:           []string{"USER"},
			Password:        "kubernetes2",
			PasswordConfirm: "kubernetes2",
		},
	}
	_, err = uMgmtSvcClient.CreateUser(context.Background(), &req)
	fmt.Printf("Again persisting same email err is %v", err)
	require.NotNil(t, err, "should not create user accounts with duplicate email")
}
