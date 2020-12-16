package usermgmt

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
)

type ServiceServer struct {
	usermgmtpb.UnimplementedUserMgmtServiceServer
}

func (ServiceServer) CreateUser(context.Context, *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	u := usermgmtpb.User{
		Id:          "xjajs",
		Name:        "",
		Email:       "someone3-ajmd@drinnovations.us",
		Roles:       []string{"USER"},
		DateCreated: ptypes.TimestampNow(),
		DateUpdated: ptypes.TimestampNow(),
	}
	resp := usermgmtpb.CreateUserResponse{
		User: &u,
	}
	return &resp, nil
}
