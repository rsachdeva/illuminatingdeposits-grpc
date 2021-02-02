package interestcal_test

import (
	"fmt"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/stretchr/testify/require"
)

func TestCalculateDelta(t *testing.T) {
	t.Parallel()

	cireq := interestcalpb.CreateInterestRequest{
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
			{
				Name: "CALM",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "2662",
						AccountType: "Brokered CD",
						Apy:         6,
						Years:       7,
						Amount:      12662,
					},
					{
						Account:     "2552",
						AccountType: "CD",
						Apy:         2.22,
						Years:       8,
						Amount:      12552,
					},
				},
			},
		},
	}

	ciresp, err := interestcal.CalculateDelta(&cireq)
	t.Logf("ciresp is %+v", ciresp)
	require.Nil(t, err)
	require.Equal(t, 4, len(ciresp.Banks), "overall number of banks must match")
	require.Equal(t, 423.94, ciresp.Delta, "overall delta for all deposists in all banks must match")
	// require.Equal(t, 23.46, ciresp.Banks[0].Deposits[2].Delta, "delta for a deposit in a bank must match")
	// require.Equal(t, 259.86, ciresp.Banks[0].Delta, "delta for a bank must match")
	for _, bkResp := range ciresp.Banks {
		switch bkResp.Name {
		case "CALM":
			fmt.Printf("bkResp is %+v\n", bkResp)
			require.Equal(t, 87.2, bkResp.Delta, "delta for CALM bank must match")
		case "HAPPIEST":
			fmt.Printf("bkResp is %+v\n", bkResp)
			require.Equal(t, 259.86, bkResp.Delta, "delta for HAPPIEST bank must match")
			require.Equal(t, 23.46, bkResp.Deposits[2].Delta, "delta for a deposit in HAPPIEST bank must match")
		case "ANGRY":
			fmt.Printf("bkResp is %+v\n", bkResp)
			require.Equal(t, 55.2, bkResp.Delta, "delta for ANGRY bank must match")
		case "NICE":
			fmt.Printf("bkResp is %+v\n", bkResp)
			require.Equal(t, 21.68, bkResp.Delta, "delta for NICE bank must match")
		default:
			panic("this should not have been used")
		}
	}
}
