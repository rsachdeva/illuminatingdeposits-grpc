package mongodbhealth

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceServer struct {
	mongodbhealthpb.UnimplementedMongoDbHealthServiceServer
	Mct *mongo.Client
}

func (svc ServiceServer) GetMongoDBHealth(ctx context.Context, em *emptypb.Empty) (*mongodbhealthpb.GetMongoDBHealthResponse, error) {
	mdbresp := statusCheck(ctx, svc.Mct)
	return mdbresp, nil
}
