package interestcal_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbconntest"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func initGRPCServerHTTP2(t *testing.T, address string) {
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}

	var opts []grpc.ServerOption
	opts = testcredentials.ServerTlsOption(opts)
	opts = append(opts, grpc.UnaryInterceptor(userauthn.EnsureValidToken))
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
	log.Println("Registering gRPC proto UserAuthenticationService...")
	userauthnpb.RegisterUserAuthnServiceServer(s, userauthn.ServiceServer{
		Mdb: mdb,
	})
	log.Println("Registering gRPC proto InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{})

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

func TestServiceServer_CreateInterest(t *testing.T) {
	t.Parallel()

	address := "localhost:50058"
	initGRPCServerHTTP2(t, address) // Starting a conventional gRPC server runs on HTTP2
	opts := []grpc.DialOption{testcredentials.ClientTlsOption(t)}
	conn, err := grpc.Dial(address, opts...)
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

	uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
	t.Log("uAuthnSvcClient client created")
	ctreq := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    email,
			Password: password,
		},
	}
	uaresp, err := uAuthnSvcClient.CreateToken(context.Background(), &ctreq)
	require.Nil(t, err, "Error should be nil when creating accessToken")
	accessToken := uaresp.VerifiedUser.AccessToken
	t.Logf("access accessToken is %v", accessToken)

	oaToken := oauth2.Token{
		AccessToken: accessToken,
		// https://stackoverflow.com/questions/34013299/web-api-authentication-basic-vs-bearer
		TokenType: "Bearer",
	}

	oAccess := oauth.NewOauthAccess(&oaToken)
	opts = append(opts, grpc.WithPerRPCCredentials(oAccess))
	for _, v := range opts {
		fmt.Printf("Opts v type is %T and val is %v\n", v, v)
	}
	connWithToken, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connWithToken.Close()

	req := interestcalpb.CreateInterestRequest{
		// &interestcalpb.NewBank is reduntant type
		// []*interestcalpb.NewBank{&interestcalpb.NewBank{ changed to []*interestcalpb.NewBank{{
		NewBanks: []*interestcalpb.NewBank{
			{
				Name: "HAPPIEST",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Checking",
						Apy:         0,
						Years:       1,
						Amount:      100,
					},
					{
						Account:     "1256",
						AccountType: "CD",
						Apy:         24,
						Years:       2,
						Amount:      10700,
					},
					{
						Account:     "1111",
						AccountType: "CD",
						Apy:         1.01,
						Years:       10,
						Amount:      27000,
					},
				},
			},
			{
				Name: "NICE",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         2.4,
						Years:       7,
						Amount:      10990,
					},
				},
			},
			{
				Name: "ANGRY",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         5,
						Years:       7,
						Amount:      10990,
					},
					{
						Account:     "9898",
						AccountType: "CD",
						Apy:         2.22,
						Years:       1,
						Amount:      5500,
					},
				},
			},
		},
	}
	iCalSvcClient := interestcalpb.NewInterestCalServiceClient(connWithToken)
	t.Log("iCalSvcClient client created")
	// endpoint CreateInterest method in InterestCalculationService
	ciresp, err := iCalSvcClient.CreateInterest(context.Background(), &req)
	if err != nil {
		t.Log("error calling CreateInterest service", err)
	}
	t.Logf("ciresp is %+v", ciresp)
	require.Equal(t, ciresp.Banks[0].Deposits[2].Delta, 23.46, "delta for a deposit in a bank must match")
	require.Equal(t, ciresp.Banks[0].Delta, 259.86, "delta for a bank must match")
	require.Equal(t, ciresp.Delta, 336.74, "overall delta for all deposists in all banks must match")
}
