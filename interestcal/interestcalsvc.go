// Package interestcal provides interest calculation service for the server api
package interestcal

import (
	"context"
	"log"
	"strings"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/segmentio/kafka-go"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/metadata"
)

type ServiceServer struct {
	KafkaWriter *kafka.Writer
	AccessToken string
	interestcalpb.UnimplementedInterestCalServiceServer
}

func (svc ServiceServer) CreateInterest(ctx context.Context, cireq *interestcalpb.CreateInterestRequest) (*interestcalpb.CreateInterestResponse, error) {
	_, span := trace.StartSpan(ctx, "interestcal.svc.createinterest")
	defer span.End()

	if svc.KafkaWriter != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Fatalln("missing metadata from incoming context")
		}
		svc.AccessToken = strings.TrimPrefix(md["authorization"][0], "Bearer ")
		log.Println("AccessToken for Kafka Key is ", svc.AccessToken)
		log.Println("KafkaWriter is ", svc.KafkaWriter)
	}

	resp, err := svc.CalculateDelta(ctx, cireq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
