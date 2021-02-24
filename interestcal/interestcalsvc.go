// Package interestcal provides interest calculation service for the server api
package interestcal

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"go.opencensus.io/trace"
)

type ServiceServer struct {
	interestcalpb.UnimplementedInterestCalServiceServer
}

func (ServiceServer) CreateInterest(ctx context.Context, cireq *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	_, span := trace.StartSpan(ctx, "interestcal.svc.createinterest")
	defer span.End()

	resp, err := CalculateDelta(ctx, cireq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
