package userauthn

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceServer struct {
	userauthnpb.UnimplementedUserAuthnServiceServer
	Mdb *mongo.Database
}

func (svc ServiceServer) CreateToken(ctx context.Context, ctreq *userauthnpb.CreateTokenRequest) (*userauthnpb.CreateTokenResponse, error) {
	resp, err := generateToken(ctx, svc.Mdb, ctreq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
