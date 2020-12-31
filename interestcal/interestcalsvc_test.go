package interestcal_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testserver"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func TestServiceServer_CreateInterest(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	cr := testserver.InitGRPCServerBuffConn(ctx, t)
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
	t.Log("uAuthnSvcClient client created")
	ctreq := userauthnpb.CreateTokenRequest{
		VerifyUser: &userauthnpb.VerifyUser{
			Email:    email,
			Password: password,
		},
	}
	uaresp, err := uAuthnSvcClient.CreateToken(context.Background(), &ctreq)
	require.Nil(t, err, "Error should be nil when creating accessToken")
	accessToken := uaresp.VerifiedUser.AccessToken
	t.Logf("access accessToken is %v", accessToken)

	oaToken := oauth2.Token{
		AccessToken: accessToken,
		// https://stackoverflow.com/questions/34013299/web-api-authentication-basic-vs-bearer
		TokenType: "Bearer",
	}

	oAccess := oauth.NewOauthAccess(&oaToken)
	opts = append(opts, grpc.WithPerRPCCredentials(oAccess))
	for _, v := range opts {
		fmt.Printf("Opts v type is %T and val is %v\n", v, v)
	}
	connWithToken, err := grpc.DialContext(ctx, "localhost", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connWithToken.Close()

	req := interestcalpb.CreateInterestRequest{
		// &interestcalpb.NewBank is reduntant type
		// []*interestcalpb.NewBank{&interestcalpb.NewBank{ changed to []*interestcalpb.NewBank{{
		NewBanks: []*interestcalpb.NewBank{
			{
				Name: "HAPPIEST",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Checking",
						Apy:         0,
						Years:       1,
						Amount:      100,
					},
					{
						Account:     "1256",
						AccountType: "CD",
						Apy:         24,
						Years:       2,
						Amount:      10700,
					},
					{
						Account:     "1111",
						AccountType: "CD",
						Apy:         1.01,
						Years:       10,
						Amount:      27000,
					},
				},
			},
			{
				Name: "NICE",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         2.4,
						Years:       7,
						Amount:      10990,
					},
				},
			},
			{
				Name: "ANGRY",
				NewDeposits: []*interestcalpb.NewDeposit{
					{
						Account:     "1234",
						AccountType: "Brokered CD",
						Apy:         5,
						Years:       7,
						Amount:      10990,
					},
					{
						Account:     "9898",
						AccountType: "CD",
						Apy:         2.22,
						Years:       1,
						Amount:      5500,
					},
				},
			},
		},
	}
	iCalSvcClient := interestcalpb.NewInterestCalServiceClient(connWithToken)
	t.Log("iCalSvcClient client created")
	// endpoint CreateInterest method in InterestCalculationService
	ciresp, err := iCalSvcClient.CreateInterest(context.Background(), &req)
	if err != nil {
		t.Log("error calling CreateInterest service", err)
	}
	t.Logf("ciresp is %+v", ciresp)
	require.Equal(t, ciresp.Banks[0].Deposits[2].Delta, 23.46, "delta for a deposit in a bank must match")
	require.Equal(t, ciresp.Banks[0].Delta, 259.86, "delta for a bank must match")
	require.Equal(t, ciresp.Delta, 336.74, "overall delta for all deposists in all banks must match")
}
