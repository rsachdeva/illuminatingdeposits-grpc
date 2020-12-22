package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/readenv"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	address = "localhost:50052"
)

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

func requestCreateUser(conn *grpc.ClientConn) {
	// Set up a connection to the server.
	fmt.Println("starting requestCreateUser")

	fmt.Println("calling NewUserMgmtServiceClient(conn)")
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")

	req := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           "growth-a@drinnovations.us",
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

func requestCreateToken(conn *grpc.ClientConn) (string, error) {
	// Set up a connection to the server.
	fmt.Println("starting requestCreateToken")

	fmt.Println("calling NewUserMgmtServiceClient(conn)")
	uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
	fmt.Println("uAuthnSvcClient client created")

	req := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    "growth-a@drinnovations.us",
			Password: "kubernetes",
		},
	}

	uaresp, err := uAuthnSvcClient.CreateToken(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateUser service", err)
		return "", err
	}
	log.Printf("ciresp is %+v", uaresp)
	return uaresp.Verifieduser.Token, nil
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

func main() {
	tls := readenv.TlsEnabled()
	fmt.Println("tls is", tls)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if tls {
		opts = []grpc.DialOption{tlsOption()}
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	requestGetMongoDBHealth(conn)
	// requestCreateUser(conn)
	token, err := requestCreateToken(conn)
	if err != nil {
		log.Fatalf("could not get token; cannot proceed %v ", err)
	}
	log.Println("JWT that can be used in next requests for Authentication\n", token)
	// requestCreateInterest(conn)
}
