package userauthn

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func generateToken(ctx context.Context, mdb *mongo.Database, ctreq *userauthnpb.CreateTokenRequest) (*userauthnpb.CreateTokenResponse, error) {
	vyu := ctreq.VerifyUser
	coll := mdb.Collection("user")
	uFound, err := usermgmt.FindUserByEmail(ctx, coll, vyu.Email)
	log.Printf("user found byb email is %v", uFound)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("NotFound Error: User not found for email %v", vyu.Email))
	}
	log.Printf("we were actually able to find the user %v\n", uFound)
	uaresp := userauthnpb.CreateTokenResponse{
		Verifieduser: &userauthnpb.VerifiedUser{
			Token: "still-to-generate",
		},
	}
	return &uaresp, nil
}
