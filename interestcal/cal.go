// Calculations

package interestcal

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
)

const (

	// Sa for aving type
	Sa = "Saving"
	// CD for cd type
	CD = "CD"
	// Ch gor checking type
	Ch = "Checking"
	// Br for Brokered type
	Br = "Brokered CD"
)

// CalculateDelta calculations for all banks
func CalculateDelta(cireq *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	bks, delta, err := computeBanksDelta(cireq)
	if err != nil {
		return &interestcalpb.CreateInterestResponse{}, err
	}
	ciresp := interestcalpb.CreateInterestResponse{
		Banks: bks,
		Delta: roundToNearest(delta),
	}
	return &ciresp, nil
}

func computeBanksDelta(cireq *interestcalpb.CreateInterestRequest) ([]*interestcalpb.Bank, float64, error) {
	var bks []*interestcalpb.Bank
	var delta float64
	for _, nb := range cireq.NewBanks {
		ds, bDelta, err := computeBankDelta(nb)
		if err != nil {
			return nil, 0, err
		}
		bk := interestcalpb.Bank{
			Name:     nb.Name,
			Deposits: ds,
			Delta:    roundToNearest(bDelta),
		}
		bks = append(bks, &bk)
		delta = delta + bk.Delta
	}
	return bks, delta, nil
}

func computeBankDelta(nb *interestcalpb.NewBank) ([]*interestcalpb.Deposit, float64, error) {
	var ds []*interestcalpb.Deposit
	var bDelta float64
	for _, nd := range nb.NewDeposits {
		d := interestcalpb.Deposit{
			Account:     nd.Account,
			AccountType: nd.AccountType,
			Apy:         nd.Apy,
			Years:       nd.Years,
			Amount:      nd.Amount,
		}
		err := computeDepositDelta(&d)
		if err != nil {
			return nil, 0, err
		}
		ds = append(ds, &d)
		bDelta = bDelta + d.Delta
	}
	return ds, bDelta, nil
}

// CalculateDelta calculates interest for 30 days for output/response Deposit
func computeDepositDelta(d *interestcalpb.Deposit) error {
	e := earned(d)
	e30Days, err := earned30days(e, d.Years)
	if err != nil {
		return errors.Wrapf(err, "calculation for Account: %s", d.Account)
	}
	d.Delta = roundToNearest(e30Days)
	return nil
}

func earned(d *interestcalpb.Deposit) float64 {
	switch d.AccountType {
	case Sa, CD:
		return compoundInterest(d.Apy, d.Years, d.Amount)
	case Br:
		return simpleInterest(d.Apy, d.Years, d.Amount)
	default:
		return 0.0
	}
}

func roundToNearest(n float64) float64 {
	return math.Round(n*100) / 100
}

func simpleInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / 100
	intEarned := amount * rateInDecimal * years
	return intEarned
}

func compoundInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / 100
	calInProcess := math.Pow(1+rateInDecimal, years)
	intEarned := amount*calInProcess - amount
	return intEarned
}

func earned30days(iEarned float64, years float64) (float64, error) {
	if years*365 < 30 {
		return 0, fmt.Errorf("NewDeposit period in years %v should not be less than 30 days", years)
	}
	i1Day := iEarned / (years * 365)
	i30 := i1Day * 30
	return math.Round(i30*100) / 100, nil
}
