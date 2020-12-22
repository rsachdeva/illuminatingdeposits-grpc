package usermgmt

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func addUser(ctx context.Context, mdb *mongo.Database, cureq *usermgmtpb.CreateUserRequest) (*usermgmtpb.CreateUserResponse, error) {
	n := cureq.NewUser
	fmt.Printf("new user from request is %+v", n)
	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}
	uuid := uuid.New().String()
	u := usermgmtpb.User{
		Uuid:         uuid,
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  ptypes.TimestampNow(),
		DateUpdated:  ptypes.TimestampNow(),
	}

	coll := mdb.Collection("user")

	res, err := coll.InsertOne(ctx, &u)
	log.Printf("user InsertOne err is %v", err)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Internal Error: Cannot convert to OID")
	}

	foundUsr, err := findUserByUuid(ctx, coll, &u, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error: After inserting user, verification failed to look up user")
	}

	if foundUsr.Uuid != uuid {
		return nil, status.Errorf(codes.Internal, "Internal Error: After inserting user, uuid of found user should match to its uuid when persisted")
	}

	// to omit password hash on client side; used only for persistence
	u.PasswordHash = nil

	curesp := usermgmtpb.CreateUserResponse{
		User: &u,
	}

	return &curesp, nil
}
