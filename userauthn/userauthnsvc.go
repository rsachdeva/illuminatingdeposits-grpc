package userauthn

import (
	"context"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

type ServiceServer struct {
	userauthnpb.UnimplementedUserAuthnServiceServer
	Mdb *mongo.Database
}

func (svc ServiceServer) CreateToken(ctx context.Context, ctreq *userauthnpb.CreateTokenRequest) (*userauthnpb.CreateTokenResponse, error) {
	ctx, span := trace.StartSpan(ctx, "userauthn.svc.createtoken")
	defer span.End()

	resp, err := generateAccessToken(ctx, svc.Mdb, ctreq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
