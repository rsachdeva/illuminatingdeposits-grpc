package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	iCalSvcClient := interestcalpb.NewInterestCalServiceClient(conn)
	fmt.Printf("client created")

	req := interestcalpb.CreateInterestRequest{
		// &interestcalpb.NewBank is reduntant type
		// []*interestcalpb.NewBank{&interestcalpb.NewBank{ changed to []*interestcalpb.NewBank{{
		NewBanks: []*interestcalpb.NewBank{
			&interestcalpb.NewBank{
				Name: "HAPPIEST",
				NewDeposits: []*interestcalpb.NewDeposit{
					&interestcalpb.NewDeposit{
						Account:     "1234",
						AccountType: "Checking",
						Apy:         0,
						Years:       1,
						Amount:      100,
					},
					&interestcalpb.NewDeposit{
						Account:     "1256",
						AccountType: "CD",
						Apy:         24,
						Years:       2,
						Amount:      10700,
					},
					&interestcalpb.NewDeposit{
						Account:     "1111",
						AccountType: "CD",
						Apy:         1.01,
						Years:       10,
						Amount:      27000,
					},
				},
			},
			&interestcalpb.NewBank{
				Name: "NICE",
				NewDeposits: []*interestcalpb.NewDeposit{
					&interestcalpb.NewDeposit{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         2.4,
						Years:       7,
						Amount:      10990,
					},
				},
			},
			&interestcalpb.NewBank{
				Name: "ANGRY",
				NewDeposits: []*interestcalpb.NewDeposit{
					&interestcalpb.NewDeposit{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         5,
						Years:       7,
						Amount:      10990,
					},
					&interestcalpb.NewDeposit{
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
	log.Printf("ciresp is %+v", ciresp)
}
