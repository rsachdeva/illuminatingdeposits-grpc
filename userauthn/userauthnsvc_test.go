package userauthn_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testserver"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestServiceServer_AccessTokenCreation(t *testing.T) {
	tt := []struct {
		name               string
		createTokenRequest *userauthnpb.CreateTokenRequest
		authnTestFunc      func(ctresp *userauthnpb.CreateTokenResponse, err error)
	}{
		{
			name: "Allowed",
			createTokenRequest: &userauthnpb.CreateTokenRequest{
				VerifyUser: &userauthnpb.VerifyUser{
					Email:    "growth@drinnovations.us",
					Password: "kubernetes",
				},
			},
			authnTestFunc: func(ctresp *userauthnpb.CreateTokenResponse, err error) {
				require.Nil(t, err, "Error should be nil when creating accessToken")
				require.NotNil(t, ctresp, "Response should not be nil")
				accessToken := ctresp.VerifiedUser.AccessToken
				t.Logf("access accessToken is %v", accessToken)
				require.NotNil(t, accessToken, "Access accessToken should not be nil")
			},
		},
		{
			name: "NotAllowed",
			createTokenRequest: &userauthnpb.CreateTokenRequest{
				VerifyUser: &userauthnpb.VerifyUser{
					Email:    "growth@drinnovations.us",
					Password: "wrong",
				},
			},
			authnTestFunc: func(ctresp *userauthnpb.CreateTokenResponse, err error) {
				require.NotNil(t, err, "Error should not be nil when creating token with incorrect password")
			},
		},
	}
	for _, tc := range tt {
		tc := tc // capture range variable https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()
			cr := testserver.InitGrpcTLSWithBuffConn(ctx, t, true)
			opts := []grpc.DialOption{grpc.WithContextDialer(testserver.GetBufDialer(cr.Listener)), testcredentials.ClientTlsOption(t)}
			conn, err := grpc.DialContext(ctx, "localhost", opts...)
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()

			email := "growth@drinnovations.us"
			password := "kubernetes"
			uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
			fmt.Println("uMgmtSvcClient client created")
			cureq := usermgmtpb.CreateUserRequest{
				NewUser: &usermgmtpb.NewUser{
					Name:            "Rohit-Sachdeva-User",
					Email:           email,
					Roles:           []string{"USER"},
					Password:        password,
					PasswordConfirm: password,
				},
			}
			umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &cureq)
			if err != nil {
				log.Println("error calling CreateUser service", err)
			}
			log.Printf("response %s", umresp.User)

			uAuthnSvcClient := userauthnpb.NewUserAuthnServiceClient(conn)
			fmt.Println("uAuthnSvcClient client created")

			ctresp, err := uAuthnSvcClient.CreateToken(context.Background(), tc.createTokenRequest)
			tc.authnTestFunc(ctresp, err)

		})
	}
}
