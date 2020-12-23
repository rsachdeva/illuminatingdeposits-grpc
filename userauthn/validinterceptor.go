package userauthn

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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
	err := valid(md["authorization"])
	if err != nil {
		return nil, err
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
