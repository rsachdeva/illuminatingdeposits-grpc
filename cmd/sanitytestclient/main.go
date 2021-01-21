// Provides sanity test client with all gRPC requests with TLS and when required JWT Authentication.
// This is to help with quick check of overall system.
// It is useful when doing refactoring as well.
// Uses unique email every time to allow new user creation and uses access token for the newly created user.
// Replace already persisted email in requestAccessToken, if requestCreateUser is not desired,
// otherwise user not found error will happen.
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

	conn, err := grpc.Dial(svcAddress(tls), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	email := fmt.Sprintf(emailFmt, uuid.New().String())
	nonAccessTokenRequests(conn, email)
	// replace alrerady persisted email if requestCreateUser is not desired in nonAccessTokenRequests
	// otherwise user not found error will happen
	oaToken := requestAccessToken(conn, email, false)
	accessTokenRequiredRequests(oaToken, tls, opts)
}

func nonAccessTokenRequests(conn *grpc.ClientConn, email string) {
	requestMongoDBHealth(conn)
	requestCreateUser(conn, email)
}

func accessTokenRequiredRequests(oaToken *oauth2.Token, tls bool, opts []grpc.DialOption) {
	oAccess := oauth.NewOauthAccess(oaToken)
	opts = append(opts, grpc.WithPerRPCCredentials(oAccess))
	for _, v := range opts {
		fmt.Printf("Opts v type is %T and val is %v\n", v, v)
	}
	connWithToken, err := grpc.Dial(svcAddress(tls), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connWithToken.Close()
	requestCreateInterest(connWithToken)
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

func requestMongoDBHealth(conn *grpc.ClientConn) {
	fmt.Println("executing requestMongoDBHealth")
	mdbSvcClient := mongodbhealthpb.NewMongoDbHealthServiceClient(conn)
	fmt.Println("mdbSvcClient client created")
	mdbresp, err := mdbSvcClient.Health(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Println("error calling MongoDBHealth service", err)
	}
	log.Printf("mdbresp is %+v", mdbresp)
}

func requestCreateUser(conn *grpc.ClientConn, email string) {
	// Set up a connection to the server.
	fmt.Println("executing requestCreateUser")

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

func requestCreateInterest(connWithToken *grpc.ClientConn) {
	fmt.Println("executing requestCreateInterest")
	iCalSvcClient := interestcalpb.NewInterestCalServiceClient(connWithToken)
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
