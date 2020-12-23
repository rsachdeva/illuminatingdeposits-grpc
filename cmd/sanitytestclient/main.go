// Proivides sanity test client with all requests/ TLS/ JWT
// This is to help with quick check of overall system
// useful when doing refactoring as well
// uses unique email to allow new user creation and use access token for the newly created user
// replace alrerady persisted email if requestCreateUser is not desired in nonAccessTokenRequests
// otherwise user not found error will happen
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/readenv"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	address  = "localhost:50052"
	emailFmt = "growth-%v@drinnovations.us"
)

func main() {
	tls := readenv.TlsEnabled()
	fmt.Println("tls is", tls)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if tls {
		opts = []grpc.DialOption{tlsOption()}
	}
	for _, v := range opts {
		fmt.Printf("initial opts v type is %T and val is %v\n", v, v)
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	email := fmt.Sprintf(emailFmt, uuid.New().String())
	nonAccessTokenRequests(conn, email)
	// replace alrerady persisted email if requestCreateUser is not desired in nonAccessTokenRequests
	// otherwise user not found error will happen
	oaToken := newAccessTokenRequest(conn, email)
	accessTokenRequiredRequests(oaToken, opts)
}

func tlsOption() grpc.DialOption {
	certFile := "conf/tls/cacrtto.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("loading certificate error is %v", err)
	}
	opt := grpc.WithTransportCredentials(creds)
	return opt
}

func requestGetMongoDBHealth(conn *grpc.ClientConn) {
	fmt.Println("starting requestGetMongoDBHealth")
	fmt.Println("calling NewMongoDbHealthServiceClient(conn)")
	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)
	fmt.Println("mdbSvcClient client created")
	mdbresp, err := mdbSvcClient.GetMongoDBHealth(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Println("error calling MongoDBHealth service", err)
	}
	log.Printf("mdbresp is %+v", mdbresp)
}

func requestCreateUser(conn *grpc.ClientConn, email string) {
	// Set up a connection to the server.
	fmt.Println("starting requestCreateUser")

	fmt.Println("calling NewUserMgmtServiceClient(conn)")
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")

	req := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           email,
			Roles:           []string{"USER"},
			Password:        "kubernetes",
			PasswordConfirm: "kubernetes",
		},
	}

	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateUser service", err)
	}
	log.Printf("ciresp is %+v", umresp)

}

func requestCreateInterest(conn *grpc.ClientConn) {
	fmt.Println("starting requestCreateInterest")
	iCalSvcClient := interestcalpb.NewInterestCalServiceClient(conn)
	fmt.Printf("iCalSvcClient client created")

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
	// endpoint CreateInterest method in InterestCalculationService
	ciresp, err := iCalSvcClient.CreateInterest(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateInterest service", err)
	}
	log.Printf("\nciresp is %+v", ciresp)
}

func nonAccessTokenRequests(conn *grpc.ClientConn, email string) {
	requestGetMongoDBHealth(conn)
	requestCreateUser(conn, email)
}

func newAccessTokenRequest(conn *grpc.ClientConn, email string) *oauth2.Token {
	// see NewOauthTokenRequest to force to use expired token
	oaToken, err := newOauthTokenRequest(conn, false, email)
	if err != nil {
		log.Fatalf("could not get token for the verification of user; cannot proceed without token %v ", err)
	}
	log.Printf("New JWT that can be used in next requests for Authentication %#v\n", oaToken)
	return oaToken
}

func accessTokenRequiredRequests(oaToken *oauth2.Token, opts []grpc.DialOption) {
	perRPC := oauth.NewOauthAccess(oaToken)
	opts = append(opts, grpc.WithPerRPCCredentials(perRPC))
	for _, v := range opts {
		fmt.Printf("Opts v type is %T and val is %v\n", v, v)
	}
	connWithToken, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connWithToken.Close()
	requestCreateInterest(connWithToken)
}
