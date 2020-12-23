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
			Email:    email,
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
	return uaresp.Verifieduser.AccessToken, nil
}

func newOauthTokenRequest(conn *grpc.ClientConn, useExpired bool) (*oauth2.Token, error) {
	token, err := accessToken(conn, useExpired)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func accessToken(conn *grpc.ClientConn, useExpired bool) (string, error) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imdyb3d0aC1hQGRyaW5ub3ZhdGlvbnMudXMiLCJyb2xlcyI6WyJVU0VSIl0sImV4cCI6MTYwODc0MzA4NywiaXNzIjoiZ2l0aHViLmNvbS9yc2FjaGRldmEvaWxsdW1pbmF0aW5nZGVwb3NpdHMtZ3JwYyJ9.KJJI-GU_kDqDXK_NOa-BC3eBfOvcOtIQ6Ho61kN7rMI"
	if useExpired {
		log.Println("Excpired JWT for testing\n", token)
		return token, nil
	}
	token, err := requestCreateToken(conn)
	if err != nil {
		return "", err
	}
	log.Println("New JWT that can be used in next requests for Authentication\n", token)
	return token, nil
}
