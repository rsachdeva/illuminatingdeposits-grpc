package usermgmt

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func findUserByUuid(ctx context.Context, coll *mongo.Collection, usr *usermgmtpb.User, uuid string) (*usermgmtpb.User, error) {
	log.Printf("\n uuid is %v\n", uuid)
	result := coll.FindOne(ctx, bson.M{"uuid": uuid})
	log.Printf("result is %+v\n", result)
	// As in json
	if err := result.Decode(usr); err != nil {
		return &usermgmtpb.User{}, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}
	return usr, nil
}

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

	curesp := usermgmtpb.CreateUserResponse{
		User: &u,
	}

	return &curesp, nil
}
