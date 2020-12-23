package userauthn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	secretKey     = "kubernetessecret"
	tokenDuration = 1 * time.Minute
)

type customClaims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

func generateToken(ctx context.Context, mdb *mongo.Database, ctreq *userauthnpb.CreateTokenRequest) (*userauthnpb.CreateTokenResponse, error) {
	vyu := ctreq.VerifyUser
	coll := mdb.Collection("user")
	uFound, err := usermgmt.FindUserByEmail(ctx, coll, vyu.Email)
	log.Printf("user found byb email is %v", uFound)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("NotFound Error: User not found for email %v", vyu.Email))
	}
	log.Printf("we were actually able to find the user %v\n", uFound)
	pwMatchErr := passwordMatch(uFound.PasswordHash, vyu.Password)
	log.Printf("Password match Err is %v\n", pwMatchErr)
	if pwMatchErr != nil {
		return nil, pwMatchErr
	}

	claims := customClaims{
		Email: uFound.Email,
		Roles: uFound.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			Issuer:    "github.com/rsachdeva/illuminatingdeposits-grpc",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	fmt.Println("signedToken generated finally is", signedToken)
	uaresp := userauthnpb.CreateTokenResponse{
		Verifieduser: &userauthnpb.VerifiedUser{
			Token: signedToken,
		},
	}
	return &uaresp, nil
}
