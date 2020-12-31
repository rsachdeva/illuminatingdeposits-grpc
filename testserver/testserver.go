package testserver

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbconntest"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize = 1024 * 1024
)

func InitGRPCServerBuffConn(ctx context.Context, t *testing.T) *bufconn.Listener {
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	listener := bufconn.Listen(bufSize)

	var opts []grpc.ServerOption
	opts = testcredentials.ServerTlsOption(opts)
	opts = append(opts, grpc.UnaryInterceptor(userauthn.EnsureValidToken))
	s := grpc.NewServer(opts...)

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
	log.Println("Registering gRPC proto UserAuthenticationService...")
	userauthnpb.RegisterUserAuthnServiceServer(s, userauthn.ServiceServer{
		Mdb: mdb,
	})
	log.Println("Registering gRPC proto InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})

	log.Println("Ready to Serve now")
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("error is %#v", err)
		}
	}()

	t.Cleanup(func() {
		log.Println("Stopping the server...")
		s.Stop()
		log.Println("Closing the listener...")
		listener.Close()
		fmt.Println("Purging dockertest for mongodb always; unless commented out to examine any data")
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
		log.Println("End of program")
	})

	return listener
}
