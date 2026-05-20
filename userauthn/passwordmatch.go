package userauthn

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func passwordMatch(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return status.Errorf(codes.NotFound, "NotFound Error: User email/password combination not found")
	}
	return nil
}
