package userauthn

import (
	"regexp"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestValidAuthHeaderExpired(t *testing.T) {
	err := valid(authHeaderExpired(), verify)
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("Token is expired"), err)
}

func TestValidAuthHeaderNotPresent(t *testing.T) {
	err := valid(authHeaderNotPresent(), verify)
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("no authorization header"), err)
}

func TestValidAuthNoEmailInClaims(t *testing.T) {
	err := valid(authHeaderSampleForMocked(), verifyWithNoEmailClaims)
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("invalid token without email"), err)
}

func verifyWithNoEmailClaims(_ string) (*customClaims, error) {
	claims := &customClaims{
		Email:          "",
		Roles:          []string{"TestRole"},
		StandardClaims: jwt.StandardClaims{},
	}
	return claims, nil
}

func authHeaderExpired() []string {
	return []string{
		"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imdyb3d0aEBkcmlubm92YXRpb25zLnVzIiwicm9sZXMiOlsiVVNFUiJdLCJleHAiOjE2MTAwNTA2NjAsImlzcyI6ImdpdGh1Yi5jb20vcnNhY2hkZXZhL2lsbHVtaW5hdGluZ2RlcG9zaXRzLXJlc3QifQ.Q6bOd3qO-2sJoZLm0unR0XXefg8FnRaYZtyoUollE20",
	}
}

func authHeaderNotPresent() []string {
	return nil
}

func authHeaderSampleForMocked() []string {
	return []string{
		"Bearer a.b.c",
	}
}
