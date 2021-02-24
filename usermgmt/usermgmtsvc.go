package usermgmt

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

type ServiceServer struct {
	usermgmtpb.UnimplementedUserMgmtServiceServer
	Mdb *mongo.Database
}

func (svc ServiceServer) CreateUser(ctx context.Context, cureq *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "usermgmt.svc.createuser")
	defer span.End()

	resp, err := addUser(ctx, svc.Mdb, cureq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
