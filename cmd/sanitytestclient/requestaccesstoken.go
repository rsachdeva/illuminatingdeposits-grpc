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

func requestAccessToken(conn *grpc.ClientConn, email string, useExpired bool) *oauth2.Token {
	// see NewOauthTokenRequest to force to use expired token
	oaToken, err := requestOauth(conn, useExpired, email)
	if err != nil {
		log.Fatalf("could not get token for the verification of user; cannot proceed without token %v ", err)
	}
	log.Printf("New JWT that can be used in next requests for Authentication %#v\n", oaToken)
	return oaToken
}

func requestOauth(conn *grpc.ClientConn, useExpired bool, email string) (*oauth2.Token, error) {
	token, err := accessToken(conn, useExpired, email)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: token,
		// https://stackoverflow.com/questions/34013299/web-api-authentication-basic-vs-bearer
		TokenType: "Bearer",
	}, nil
}

func accessToken(conn *grpc.ClientConn, useExpired bool, email string) (string, error) {
	var token string
	if useExpired {
		storedTk, err := ioutil.ReadFile("cmd/sanitytestclient/expiredtoken.data")
		token = string(storedTk)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Expired JWT for testing\n", token)
		return token, nil
	}
	token, err := createToken(conn, email)
	if err != nil {
		return "", err
	}
	log.Println("New JWT that can be used in next requests for Authentication\n", token)
	return token, nil
}

func createToken(conn *grpc.ClientConn, email string) (string, error) {
	// Set up a connection to the server.
	fmt.Println("executing createToken")

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
	return uaresp.VerifiedUser.AccessToken, nil
}
