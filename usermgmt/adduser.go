package usermgmt

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func addUser(ctx context.Context, mdb *mongo.Database, cureq *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	n := cureq.NewUser
	fmt.Printf("new user from request is %+v", n)
	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := usermgmtpb.User{
		Id:           uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  ptypes.TimestampNow(),
		DateUpdated:  ptypes.TimestampNow(),
	}

	curesp := usermgmtpb.CreateUserResponse{
		User: &u,
	}

	return &curesp, nil
}
