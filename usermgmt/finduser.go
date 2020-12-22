package usermgmt

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FindUserByEmail(ctx context.Context, coll *mongo.Collection, email string) (*usermgmtpb.User, error) {
	log.Printf("\n findUserByEmail email is %v\n", email)
	result := coll.FindOne(ctx, bson.M{"email": email})
	log.Printf("result is %+v\n", result)
	usr := usermgmtpb.User{}
	if err := result.Decode(&usr); err != nil {
		return &usermgmtpb.User{}, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find user with specified Email: %v", err))
	}
	return &usr, nil
}

func findUserByUuid(ctx context.Context, coll *mongo.Collection, usr *usermgmtpb.User, uuid string) (*usermgmtpb.User, error) {
	log.Printf("\n uuid is %v\n", uuid)
	result := coll.FindOne(ctx, bson.M{"uuid": uuid})
	log.Printf("result is %+v\n", result)
	if err := result.Decode(usr); err != nil {
		return &usermgmtpb.User{}, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find user with specified uuid %v due to error: %v", uuid, err))
	}
	return usr, nil
}
