package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

func requestCreateToken(conn *grpc.ClientConn) (string, error) {
	// Set up a connection to the server.
	fmt.Println("starting requestCreateToken")

	fmt.Println("calling NewUserMgmtServiceClient(conn)")
	uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
	fmt.Println("uAuthnSvcClient client created")

	req := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    "growth-a@drinnovations.us",
			Password: "kubernetes",
		},
	}

	uaresp, err := uAuthnSvcClient.CreateToken(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateUser service", err)
		return "", err
	}
	log.Printf("uaresp is %+v", uaresp)
	// return "haha", nil
	return uaresp.Verifieduser.Token, nil
}

func oAuthToken(conn *grpc.ClientConn) (*oauth2.Token, error) {
	token, err := requestCreateToken(conn)
	if err != nil {
		return nil, err
	}
	log.Println("JWT that can be used in next requests for Authentication\n", token)
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}
