package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

func requestCreateToken(conn *grpc.ClientConn, email string) (string, error) {
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

func newOauthTokenRequest(conn *grpc.ClientConn, useExpired bool, email string) (*oauth2.Token, error) {
	token, err := accessToken(conn, useExpired, email)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func accessToken(conn *grpc.ClientConn, useExpired bool, email string) (string, error) {
	var token string
	if useExpired {
		by, err := ioutil.ReadFile("cmd/sanitytestclient/expiredtoken.data")
		token = string(by)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Expired JWT for testing\n", token)
		return token, nil
	}
	token, err := requestCreateToken(conn, email)
	if err != nil {
		return "", err
	}
	log.Println("New JWT that can be used in next requests for Authentication\n", token)
	return token, nil
}
