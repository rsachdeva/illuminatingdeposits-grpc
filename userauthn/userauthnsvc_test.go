package userauthn_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbtestconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
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
	mt, pool, resource := mongodbtestconn.Connect(ctx, 1)
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
	log.Println("Registering gRPC proto UserAuthenticationService...")
	userauthnpb.RegisterUserAuthnServiceServer(s, userauthn.ServiceServer{
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

func TestServiceServer_CreateToken(t *testing.T) {
	t.Parallel()

	address := "localhost:50056"
	initGRPCServerHTTP2(t, address) // Starting a conventional gRPC server runs on HTTP2
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	email := "growth@drinnovations.us"
	password := "kubernetes"
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")
	cureq := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           email,
			Roles:           []string{"USER"},
			Password:        password,
			PasswordConfirm: password,
		},
	}
	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &cureq)
	if err != nil {
		log.Println("error calling CreateUser service", err)
	}
	log.Printf("response %s", umresp.User)
	require.Equal(t, umresp.User.Email, "growth@drinnovations.us")

	uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
	fmt.Println("uAuthnSvcClient client created")

	ctreq := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    email,
			Password: password,
		},
	}

	uaresp, err := uAuthnSvcClient.CreateToken(context.Background(), &ctreq)
	require.Nil(t, err, "Error should be nil when creating token")
	require.NotNil(t, uaresp, "Response should not be nil")
	token := uaresp.VerifiedUser.AccessToken
	t.Logf("access token is %v", token)
	require.NotNil(t, token, "Access token should not be nil")
}

func TestServiceServer_CreateTokenNotAllowed(t *testing.T) {
	t.Parallel()

	address := "localhost:50057"
	initGRPCServerHTTP2(t, address) // Starting a conventional gRPC server runs on HTTP2
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	email := "growth@drinnovations.us"
	password := "kubernetes"
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")
	cureq := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           email,
			Roles:           []string{"USER"},
			Password:        password,
			PasswordConfirm: password,
		},
	}
	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &cureq)
	if err != nil {
		log.Println("error calling CreateUser service", err)
	}
	log.Printf("response %s", umresp.User)
	require.Equal(t, umresp.User.Email, "growth@drinnovations.us")

	uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
	fmt.Println("uAuthnSvcClient client created")

	ctreq := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    email,
			Password: "wrong",
		},
	}

	_, err = uAuthnSvcClient.CreateToken(context.Background(), &ctreq)
	require.NotNil(t, err, "Error should not be nil when creating token with incorrect password")
}
