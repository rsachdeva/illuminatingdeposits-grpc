package userauthn

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"go.mongodb.org/mongo-driver/mongo"
)

func generateToken(ctx context.Context, Mdb *mongo.Database, ctreq *userauthnpb.CreateTokenRequest) (*userauthnpb.CreateTokenResponse, error) {
	return nil, nil
}
