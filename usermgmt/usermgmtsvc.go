package usermgmt

import "github.com/rsachdeva/illuminatingdeposits-grpc/api/usermgmtpb"

type ServiceServer struct {
	usermgmtpb.UnimplementedUserMgmtServiceServer
}
