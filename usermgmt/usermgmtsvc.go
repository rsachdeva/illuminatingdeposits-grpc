package usermgmt

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceServer struct {
	usermgmtpb.UnimplementedUserMgmtServiceServer
	Mdb *mongo.Database
}

func (svc ServiceServer) CreateUser(ctx context.Context, cureq *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	resp, err := addUser(ctx, svc.Mdb, cureq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
