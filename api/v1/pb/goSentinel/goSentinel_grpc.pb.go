// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: goSentinel.proto

package goSentinel

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GoSentinelServiceClient is the client API for GoSentinelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoSentinelServiceClient interface {
	SendVerifcationNotification(ctx context.Context, in *SendApplicationNotificationRequest, opts ...grpc.CallOption) (*SendApplicationNotificationResponse, error)
	VerifyApplication(ctx context.Context, in *VerifyApplicationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetApplicationSecrets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetApplicationSecretResponse, error)
}

type goSentinelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGoSentinelServiceClient(cc grpc.ClientConnInterface) GoSentinelServiceClient {
	return &goSentinelServiceClient{cc}
}

func (c *goSentinelServiceClient) SendVerifcationNotification(ctx context.Context, in *SendApplicationNotificationRequest, opts ...grpc.CallOption) (*SendApplicationNotificationResponse, error) {
	out := new(SendApplicationNotificationResponse)
	err := c.cc.Invoke(ctx, "/goSentinel.goSentinelService/SendVerifcationNotification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goSentinelServiceClient) VerifyApplication(ctx context.Context, in *VerifyApplicationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/goSentinel.goSentinelService/VerifyApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goSentinelServiceClient) GetApplicationSecrets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetApplicationSecretResponse, error) {
	out := new(GetApplicationSecretResponse)
	err := c.cc.Invoke(ctx, "/goSentinel.goSentinelService/GetApplicationSecrets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoSentinelServiceServer is the server API for GoSentinelService service.
// All implementations should embed UnimplementedGoSentinelServiceServer
// for forward compatibility
type GoSentinelServiceServer interface {
	SendVerifcationNotification(context.Context, *SendApplicationNotificationRequest) (*SendApplicationNotificationResponse, error)
	VerifyApplication(context.Context, *VerifyApplicationRequest) (*emptypb.Empty, error)
	GetApplicationSecrets(context.Context, *emptypb.Empty) (*GetApplicationSecretResponse, error)
}

// UnimplementedGoSentinelServiceServer should be embedded to have forward compatible implementations.
type UnimplementedGoSentinelServiceServer struct {
}

func (UnimplementedGoSentinelServiceServer) SendVerifcationNotification(context.Context, *SendApplicationNotificationRequest) (*SendApplicationNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendVerifcationNotification not implemented")
}
func (UnimplementedGoSentinelServiceServer) VerifyApplication(context.Context, *VerifyApplicationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyApplication not implemented")
}
func (UnimplementedGoSentinelServiceServer) GetApplicationSecrets(context.Context, *emptypb.Empty) (*GetApplicationSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationSecrets not implemented")
}

// UnsafeGoSentinelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoSentinelServiceServer will
// result in compilation errors.
type UnsafeGoSentinelServiceServer interface {
	mustEmbedUnimplementedGoSentinelServiceServer()
}

func RegisterGoSentinelServiceServer(s grpc.ServiceRegistrar, srv GoSentinelServiceServer) {
	s.RegisterService(&GoSentinelService_ServiceDesc, srv)
}

func _GoSentinelService_SendVerifcationNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendApplicationNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoSentinelServiceServer).SendVerifcationNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/goSentinel.goSentinelService/SendVerifcationNotification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoSentinelServiceServer).SendVerifcationNotification(ctx, req.(*SendApplicationNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoSentinelService_VerifyApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoSentinelServiceServer).VerifyApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/goSentinel.goSentinelService/VerifyApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoSentinelServiceServer).VerifyApplication(ctx, req.(*VerifyApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoSentinelService_GetApplicationSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoSentinelServiceServer).GetApplicationSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/goSentinel.goSentinelService/GetApplicationSecrets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoSentinelServiceServer).GetApplicationSecrets(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// GoSentinelService_ServiceDesc is the grpc.ServiceDesc for GoSentinelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoSentinelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goSentinel.goSentinelService",
	HandlerType: (*GoSentinelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendVerifcationNotification",
			Handler:    _GoSentinelService_SendVerifcationNotification_Handler,
		},
		{
			MethodName: "VerifyApplication",
			Handler:    _GoSentinelService_VerifyApplication_Handler,
		},
		{
			MethodName: "GetApplicationSecrets",
			Handler:    _GoSentinelService_GetApplicationSecrets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goSentinel.proto",
}
