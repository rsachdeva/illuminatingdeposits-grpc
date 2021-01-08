package userauthn

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type tokenVerifyFunc func(accessToken string) (*customClaims, error)

// Unary interceptor grpc.UnaryServerInterceptor function parameter for grpc.UnaryInterceptor function
func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("EnsureValidToken unary server info method", info.FullMethod)
	if info.FullMethod != "/interestcal.InterestCalService/CreateInterest" {
		log.Printf("Nothing needs to be authenticated for %v", info.FullMethod)
		return handler(ctx, req)
	}
	log.Printf("Authentication is needed for %v", info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	fmt.Print("md is\n", md)
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	err := valid(md["authorization"], verify)
	if err != nil {
		return nil, err
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// valid validates the authorization.
func valid(authorization []string, vryFunc tokenVerifyFunc) error {
	if len(authorization) < 1 {
		return status.Errorf(codes.Unauthenticated, "no authorization header")
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	fmt.Println("token extracted for verification is: ", token)
	claims, err := vryFunc(token)
	fmt.Printf("tkv.vryFunc(token) err is %v\n", err)
	fmt.Printf("tkv.vryFunc(token) claims is %v\n", claims)
	if err != nil {
		return err
	}
	email := claims.Email
	if len(email) < 1 {
		return status.Errorf(codes.Unauthenticated, "invalid token without email")
	}
	return nil
}
