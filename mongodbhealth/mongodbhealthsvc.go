package mongodbhealth

import (
	"context"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceServer struct {
	mongodbhealthpb.UnimplementedMongoDbHealthServiceServer
	Mct *mongo.Client
}

func (svc ServiceServer) Health(ctx context.Context, em *emptypb.Empty) (*mongodbhealthpb.HealthResponse, error) {
	log.Println("Health svc status check!")

	mdbresp := statusCheck(ctx, svc.Mct)
	return mdbresp, nil
}
