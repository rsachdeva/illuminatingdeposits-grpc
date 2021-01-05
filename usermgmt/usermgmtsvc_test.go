// Adds test that starts a gRPC server and client tests the user mgmt service with RPC
package usermgmt_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/testcredentials"
	"github.com/rsachdeva/illuminatingdeposits-grpc/testserver"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestServiceServer_CreateUser(t *testing.T) {
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

	uMgmtSvcClient := usermgmtpb.NewUserMgmtServiceClient(conn)
	fmt.Println("uMgmtSvcClient client created")
	req := usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User",
			Email:           "growth@drinnovations.us",
			Roles:           []string{"ADMIN"},
			Password:        "kubernetes",
			PasswordConfirm: "kubernetes",
		},
	}
	umresp, err := uMgmtSvcClient.CreateUser(context.Background(), &req)
	if err != nil {
		log.Println("error calling CreateUser service", err)
	}
	log.Printf("response %s", umresp.User)
	require.Equal(t, "growth@drinnovations.us", umresp.User.Email)

	req = usermgmtpb.CreateUserRequest{
		NewUser: &usermgmtpb.NewUser{
			Name:            "Rohit-Sachdeva-User2",
			Email:           "growth@drinnovations.us",
			Roles:           []string{"USER"},
			Password:        "kubernetes2",
			PasswordConfirm: "kubernetes2",
		},
	}
	_, err = uMgmtSvcClient.CreateUser(context.Background(), &req)
	fmt.Printf("Again persisting same email err is %v", err)
	require.NotNil(t, err, "should not create user accounts with duplicate email")
}
