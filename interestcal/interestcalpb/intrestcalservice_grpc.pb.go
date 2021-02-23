// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package interestcalpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// InterestCalServiceClient is the client API for InterestCalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InterestCalServiceClient interface {
	CreateInterest(ctx context.Context, in *CreateInterestRequest, opts ...grpc.CallOption) (*CreateInterestResponse, error)
}

type interestCalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInterestCalServiceClient(cc grpc.ClientConnInterface) InterestCalServiceClient {
	return &interestCalServiceClient{cc}
}

func (c *interestCalServiceClient) CreateInterest(ctx context.Context, in *CreateInterestRequest, opts ...grpc.CallOption) (*CreateInterestResponse, error) {
	out := new(CreateInterestResponse)
	err := c.cc.Invoke(ctx, "/interestcal.InterestCalService/CreateInterest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InterestCalServiceServer is the server API for InterestCalService service.
// All implementations must embed UnimplementedInterestCalServiceServer
// for forward compatibility
type InterestCalServiceServer interface {
	CreateInterest(context.Context, *CreateInterestRequest) (*CreateInterestResponse, error)
	mustEmbedUnimplementedInterestCalServiceServer()
}

// UnimplementedInterestCalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInterestCalServiceServer struct {
}

func (UnimplementedInterestCalServiceServer) CreateInterest(context.Context, *CreateInterestRequest) (*CreateInterestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInterest not implemented")
}
func (UnimplementedInterestCalServiceServer) mustEmbedUnimplementedInterestCalServiceServer() {}

// UnsafeInterestCalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InterestCalServiceServer will
// result in compilation errors.
type UnsafeInterestCalServiceServer interface {
	mustEmbedUnimplementedInterestCalServiceServer()
}

func RegisterInterestCalServiceServer(s grpc.ServiceRegistrar, srv InterestCalServiceServer) {
	s.RegisterService(&InterestCalService_ServiceDesc, srv)
}

func _InterestCalService_CreateInterest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInterestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterestCalServiceServer).CreateInterest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interestcal.InterestCalService/CreateInterest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterestCalServiceServer).CreateInterest(ctx, req.(*CreateInterestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InterestCalService_ServiceDesc is the grpc.ServiceDesc for InterestCalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InterestCalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "interestcal.InterestCalService",
	HandlerType: (*InterestCalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateInterest",
			Handler:    _InterestCalService_CreateInterest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interestcal/interestcalpb/intrestcalservice.proto",
}
