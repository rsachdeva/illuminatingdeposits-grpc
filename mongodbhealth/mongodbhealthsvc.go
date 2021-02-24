package mongodbhealth

import (
	"context"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceServer struct {
	mongodbhealthpb.UnimplementedMongoDbHealthServiceServer
	Mct *mongo.Client
}

func (svc ServiceServer) Health(ctx context.Context, em *emptypb.Empty) (*mongodbhealthpb.HealthResponse, error) {
	log.Println("Health svc status check!")
	ctx, span := trace.StartSpan(ctx, "mongodb.svc.health")
	defer span.End()

	mdbresp := statusCheck(ctx, svc.Mct)
	return mdbresp, nil
}
