package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

func withoutTlsRequestCreateUser() {
	// Set up a connection to the server.
	fmt.Println("starting withoutTlsRequestCreateUser")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	fmt.Println("calling NewUserMgmtServiceClient(conn)")
	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")

	req := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-Admin",
			Email:           "growth@drinnovations.us",
			Roles:           []string{"Admin"},
			Password:        "kubernetes",
			PasswordConfirm: "kubernetes",
		},
	}

	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &req)
	if err != nil {
		log.Fatal("error calling CreateUser service", err)
	}
	log.Printf("ciresp is %+v", umresp)

}

func withoutTlsRequestCreateInterest() {
	fmt.Println("starting withoutTlsRequestCreateInterest")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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
		log.Fatal("error calling CreateInterest service", err)
	}
	log.Printf("\nciresp is %+v", ciresp)
}

func main() {
	withoutTlsRequestCreateUser()
	withoutTlsRequestCreateInterest()
}
