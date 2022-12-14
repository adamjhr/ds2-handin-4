// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: grpc/interface.proto

package critical

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

// CriticalClient is the client API for Critical service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CriticalClient interface {
	PassToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Empty, error)
	Election(ctx context.Context, in *Candidate, opts ...grpc.CallOption) (*Empty, error)
}

type criticalClient struct {
	cc grpc.ClientConnInterface
}

func NewCriticalClient(cc grpc.ClientConnInterface) CriticalClient {
	return &criticalClient{cc}
}

func (c *criticalClient) PassToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ping.Critical/passToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalClient) Election(ctx context.Context, in *Candidate, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ping.Critical/election", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CriticalServer is the server API for Critical service.
// All implementations must embed UnimplementedCriticalServer
// for forward compatibility
type CriticalServer interface {
	PassToken(context.Context, *Token) (*Empty, error)
	Election(context.Context, *Candidate) (*Empty, error)
	mustEmbedUnimplementedCriticalServer()
}

// UnimplementedCriticalServer must be embedded to have forward compatible implementations.
type UnimplementedCriticalServer struct {
}

func (UnimplementedCriticalServer) PassToken(context.Context, *Token) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PassToken not implemented")
}
func (UnimplementedCriticalServer) Election(context.Context, *Candidate) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Election not implemented")
}
func (UnimplementedCriticalServer) mustEmbedUnimplementedCriticalServer() {}

// UnsafeCriticalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CriticalServer will
// result in compilation errors.
type UnsafeCriticalServer interface {
	mustEmbedUnimplementedCriticalServer()
}

func RegisterCriticalServer(s grpc.ServiceRegistrar, srv CriticalServer) {
	s.RegisterService(&Critical_ServiceDesc, srv)
}

func _Critical_PassToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalServer).PassToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ping.Critical/passToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalServer).PassToken(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _Critical_Election_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Candidate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalServer).Election(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ping.Critical/election",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalServer).Election(ctx, req.(*Candidate))
	}
	return interceptor(ctx, in, info, handler)
}

// Critical_ServiceDesc is the grpc.ServiceDesc for Critical service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Critical_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ping.Critical",
	HandlerType: (*CriticalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "passToken",
			Handler:    _Critical_PassToken_Handler,
		},
		{
			MethodName: "election",
			Handler:    _Critical_Election_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/interface.proto",
}
