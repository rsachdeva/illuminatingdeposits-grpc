// Calculations

package interestcal

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"go.opencensus.io/trace"
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// CalculateDelta calculations for all banks
// computeBanksDelta uses concurrency
func (svc ServiceServer) CalculateDelta(ctx context.Context, cireq *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	defer timeTrack(time.Now(), "CalculateDelta timed with withConcurrency for more I/O processing")
	var bks []*interestcalpb.Bank
	var delta float64
	var err error
	// withConcurrency := true
	// if withConcurrency {
	// 	bks, delta, err = computeBanksDelta(cireq)
	// }
	// if !withConcurrency {
	// 	bks, delta, err = computeBanksDeltaSequentially(cireq)
	// }
	log.Println("Starting overall calculate delta")
	bks, delta, err = svc.computeBanksDelta(ctx, cireq)
	if err != nil {
		return &interestcalpb.CreateInterestResponse{}, err
	}
	ciresp := interestcalpb.CreateInterestResponse{
		Banks: bks,
		Delta: roundToNearest(delta),
	}
	return &ciresp, nil
}

// func computeBanksDeltaSequentially(cireq *interestcalpb.CreateInterestRequest) ([]*interestcalpb.Bank, float64, error) {
// 	var bks []*interestcalpb.Bank
// 	var delta float64
// 	for _, nb := range cireq.NewBanks {
// 		ds, bDelta, err := computeBankDeltaNoChannel(nb)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		bk := interestcalpb.Bank{
// 			Name:     nb.Name,
// 			Deposits: ds,
// 			Delta:    roundToNearest(bDelta),
// 		}
// 		bks = append(bks, &bk)
// 		delta = delta + bk.Delta
// 	}
// 	return bks, delta, nil
// }
//
// func computeBankDeltaNoChannel(nb *interestcalpb.NewBank) ([]*interestcalpb.Deposit, float64, error) {
// 	// time.Sleep(5 * time.Second)
// 	var ds []*interestcalpb.Deposit
// 	var bDelta float64
// 	for _, nd := range nb.NewDeposits {
// 		d := interestcalpb.Deposit{
// 			Account:     nd.Account,
// 			AccountType: nd.AccountType,
// 			Apy:         nd.Apy,
// 			Years:       nd.Years,
// 			Amount:      nd.Amount,
// 		}
// 		err := computeDepositDelta(&d)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		ds = append(ds, &d)
// 		bDelta = bDelta + d.Delta
// 	}
// 	return ds, bDelta, nil
// }

type BankResult struct {
	name   string
	ds     []*interestcalpb.Deposit
	bDelta float64
	err    error
}

func (svc ServiceServer) computeBanksDelta(ctx context.Context, cireq *interestcalpb.CreateInterestRequest) ([]*interestcalpb.Bank, float64, error) {
	var bks []*interestcalpb.Bank
	var delta float64
	bkCh := make(chan BankResult)
	// so as to know when to close channel
	var waitGroup sync.WaitGroup

	// for all banks just 1 go routine
	// go func() {
	// 	for _, nb := range cireq.NewBanks {
	// 		// ds, bDelta, err := computeBankDelta(nb)
	// 		// doing for each bank calculation concurrently
	// 		computeBankDelta(nb, bkCh)
	// 	}
	// 	close(bkCh)
	// }()

	// for each bank calculation with all deposists 1 goroutine
	log.Println("len(cireq.NewBanks) is", len(cireq.NewBanks))
	waitGroup.Add(len(cireq.NewBanks))
	for i, nb := range cireq.NewBanks {
		go func(i int, nb *interestcalpb.NewBank) {
			defer waitGroup.Done()
			svc.computeBankDelta(ctx, nb, bkCh)
		}(i, nb)
	}

	go func() {
		waitGroup.Wait()
		close(bkCh)
	}()

	for bkRes := range bkCh {
		if bkRes.err != nil {
			return nil, 0, bkRes.err
		}
		bk := interestcalpb.Bank{
			Name:     bkRes.name,
			Deposits: bkRes.ds,
			Delta:    roundToNearest(bkRes.bDelta),
		}
		log.Println("bk bank being added with delta calculation is", bk)
		bks = append(bks, &bk)
		delta = delta + bk.Delta
	}
	return bks, delta, nil
}

// return []*interestcalpb.Deposit, float64, error now in channel
func (svc ServiceServer) computeBankDelta(ctx context.Context, nb *interestcalpb.NewBank, bkCh chan<- BankResult) {
	//  I/O processing kafka broker and tracing (instead of time.Sleep(5 * time.Second) )
	_, span := trace.StartSpan(ctx, "interestcal.computeBankDelta")
	defer span.End()
	var ds []*interestcalpb.Deposit
	var bDelta float64
	for _, nd := range nb.NewDeposits {
		log.Printf("for nb.name %v the nd is %v", nb.Name, nd)
		d := interestcalpb.Deposit{
			Account:     nd.Account,
			AccountType: nd.AccountType,
			Apy:         nd.Apy,
			Years:       nd.Years,
			Amount:      nd.Amount,
		}
		err := computeDepositDelta(&d)
		// if err != nil {
		// 	return nil, 0, err
		// }
		// Sending err result
		fmt.Println("any err in compute deposit delta", err)
		if err != nil {
			fmt.Printf("error in goroutine running computeDepositDelta %v\n", err)
			bkCh <- BankResult{
				name:   "",
				ds:     nil,
				bDelta: 0,
				err:    err,
			}
			return
		}
		dc := DepositCalculation{
			BankName:    nb.Name,
			Account:     d.Account,
			AccountType: d.AccountType,
			Apy:         d.Apy,
			Years:       d.Years,
			Amount:      d.Amount,
			Delta:       d.Delta,
		}
		if svc.KafkaWriter != nil {
			log.Println("writing message to message broker..")
			svc.writeMessage(ctx, dc)
		}
		ds = append(ds, &d)
		bDelta = bDelta + d.Delta
		fmt.Println("inside loop ds is", ds)
		fmt.Println("inside loop bDelta is", bDelta)
	}
	fmt.Println("outside loop ds is", ds)
	fmt.Println("outside loop bDelta is", bDelta)
	// return ds, bDelta, nil
	bkRes := BankResult{
		name:   nb.Name,
		ds:     ds,
		bDelta: bDelta,
		err:    nil,
	}
	bkCh <- bkRes
}

// CalculateDelta calculates interest for 30 days for output/response Deposit
func computeDepositDelta(d *interestcalpb.Deposit) error {
	e := earned(d)
	e30Days, err := earned30days(e, d.Years)
	fmt.Println("err in computeDepositDelta is", err)
	if err != nil {
		fmt.Println("error in computeDepositDelta ", err)
		return errors.Wrapf(err, "calculation for Account: %s", d.Account)
	}
	d.Delta = roundToNearest(e30Days)
	log.Println("no error returned computeDepositDelta d.Delta is", d.Delta)
	return nil
}

func earned(d *interestcalpb.Deposit) float64 {
	fmt.Println("d.AccountType is", d.AccountType)
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
