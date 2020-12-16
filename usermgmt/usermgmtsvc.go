package usermgmt

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceServer struct {
	usermgmtpb.UnimplementedUserMgmtServiceServer
	Mdb *mongo.Database
}

// func (svc ServiceServer) CreateUser(ctx context.Context, cureq *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
// 	resp, err := addUser(ctx, svc.Mdb, cureq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }

func (ServiceServer) CreateUser(ctx context.Context, req *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	u := req.NewUser
	resp := usermgmtpb.CreateUserResponse{
		User: &usermgmtpb.User{
			Id:          "random-uuid-jmd",
			Name:        u.Name,
			Email:       u.Email,
			Roles:       u.Roles,
			DateCreated: ptypes.TimestampNow(),
			DateUpdated: ptypes.TimestampNow(),
		},
	}
	return &resp, nil
}
