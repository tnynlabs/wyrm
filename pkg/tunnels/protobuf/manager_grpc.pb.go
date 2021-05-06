// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TunnelManagerClient is the client API for TunnelManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TunnelManagerClient interface {
	RevokeDevice(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	InvokeDevice(ctx context.Context, in *InvokeRequest, opts ...grpc.CallOption) (*InvokeResponse, error)
}

type tunnelManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewTunnelManagerClient(cc grpc.ClientConnInterface) TunnelManagerClient {
	return &tunnelManagerClient{cc}
}

func (c *tunnelManagerClient) RevokeDevice(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/wyrm.tunnel.TunnelManager/RevokeDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tunnelManagerClient) InvokeDevice(ctx context.Context, in *InvokeRequest, opts ...grpc.CallOption) (*InvokeResponse, error) {
	out := new(InvokeResponse)
	err := c.cc.Invoke(ctx, "/wyrm.tunnel.TunnelManager/InvokeDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TunnelManagerServer is the server API for TunnelManager service.
// All implementations must embed UnimplementedTunnelManagerServer
// for forward compatibility
type TunnelManagerServer interface {
	RevokeDevice(context.Context, *RevokeRequest) (*emptypb.Empty, error)
	InvokeDevice(context.Context, *InvokeRequest) (*InvokeResponse, error)
	mustEmbedUnimplementedTunnelManagerServer()
}

// UnimplementedTunnelManagerServer must be embedded to have forward compatible implementations.
type UnimplementedTunnelManagerServer struct {
}

func (UnimplementedTunnelManagerServer) RevokeDevice(context.Context, *RevokeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeDevice not implemented")
}
func (UnimplementedTunnelManagerServer) InvokeDevice(context.Context, *InvokeRequest) (*InvokeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvokeDevice not implemented")
}
func (UnimplementedTunnelManagerServer) mustEmbedUnimplementedTunnelManagerServer() {}

// UnsafeTunnelManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TunnelManagerServer will
// result in compilation errors.
type UnsafeTunnelManagerServer interface {
	mustEmbedUnimplementedTunnelManagerServer()
}

func RegisterTunnelManagerServer(s grpc.ServiceRegistrar, srv TunnelManagerServer) {
	s.RegisterService(&_TunnelManager_serviceDesc, srv)
}

func _TunnelManager_RevokeDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TunnelManagerServer).RevokeDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wyrm.tunnel.TunnelManager/RevokeDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TunnelManagerServer).RevokeDevice(ctx, req.(*RevokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TunnelManager_InvokeDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TunnelManagerServer).InvokeDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wyrm.tunnel.TunnelManager/InvokeDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TunnelManagerServer).InvokeDevice(ctx, req.(*InvokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TunnelManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wyrm.tunnel.TunnelManager",
	HandlerType: (*TunnelManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RevokeDevice",
			Handler:    _TunnelManager_RevokeDevice_Handler,
		},
		{
			MethodName: "InvokeDevice",
			Handler:    _TunnelManager_InvokeDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tunnel_service/manager/manager.proto",
}
