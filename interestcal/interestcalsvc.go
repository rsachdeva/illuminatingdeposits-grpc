// Package interestcal provides interest calculation service for the server api
package interestcal

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestvalue"
)


type ServiceServer struct {
	interestcalpb.UnimplementedInterestCalServiceServer
}

func (ServiceServer) CreateInterest(ctx context.Context, cireq *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	resp, err := interestvalue.CalculateDelta(cireq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
