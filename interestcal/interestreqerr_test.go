package interestcal_test

import (
	"context"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/stretchr/testify/require"
)

func TestCalculateDeltaLessThan30Days(t *testing.T) {
	t.Parallel()

	svc := interestcal.ServiceServer{}
	cireq := interestcalpb.CreateInterestRequest{
		NewBanks: []*interestcalpb.NewBank{
			{
				Name: "HAPPIEST",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "28282",
						AccountType: "CD",
						Apy:         5,
						Years:       0.0001,
						Amount:      1000000,
					},
				},
			},
		},
	}
	_, err := svc.CreateInterest(context.Background(), &cireq)
	require.NotNil(t, err, "NewDeposit period in years 0.0001 should not be less than 30 days")
}

func TestCalculateDeltaNoDeposits(t *testing.T) {
	t.Parallel()

	svc := interestcal.ServiceServer{}
	cireq := interestcalpb.CreateInterestRequest{
		NewBanks: []*interestcalpb.NewBank{
			{
				Name:        "HAPPIEST",
				NewDeposits: []*interestcalpb.NewDeposit{},
			},
		},
	}

	ciresp, err := svc.CreateInterest(context.Background(), &cireq)
	require.Nil(t, err)
	require.Equal(t, "HAPPIEST", ciresp.Banks[0].Name)
	require.Equal(t, 0.0, ciresp.Banks[0].Delta)
}
