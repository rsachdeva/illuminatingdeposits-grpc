// Package interestcal provides interest calculation service for the server api
package interestcal

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/interestcalpb"
)


type ServiceServer struct {
	interestcalpb.UnimplementedInterestCalServiceServer
}

func (ServiceServer) CreateInterest(ctx context.Context, req *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	resp := interestcalpb.CreateInterestResponse{
		Banks: nil,
		Delta: 0,
	}
	return &resp, nil
}
