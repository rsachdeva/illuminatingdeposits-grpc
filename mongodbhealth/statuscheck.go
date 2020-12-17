package mongodbhealth

import (
	"context"
	"fmt"

	"github.com/rsachdeva/illuminatingdeposits-grpc/api/mongodbhealthpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func statusCheck(ctx context.Context, mct *mongo.Client) *mongodbhealthpb.GetMongoDBHealthResponse {

	//  Running this query forces a
	// round trip to the mongodb.
	status := "Ok"
	if err := mct.Ping(
		ctx,
		&readpref.ReadPref{},
	); err != nil {
		fmt.Printf("Going to set MongoDB status as not Ready for err %v\n", err)
		status = "MongoDB Not Ready"
	}

	h := mongodbhealthpb.Health{
		Status: status,
	}

	mdbresp := mongodbhealthpb.GetMongoDBHealthResponse{
		Health: &h,
	}

	return &mdbresp
}
