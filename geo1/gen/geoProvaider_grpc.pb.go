// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: geoProvaider.proto

package geo_provider

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

// GeoProviderGRPCClient is the client API for GeoProviderGRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GeoProviderGRPCClient interface {
	AddressSearchGRPC(ctx context.Context, in *RequestAddressSearch, opts ...grpc.CallOption) (*RespAddress, error)
	AddressGeoCodeGRPC(ctx context.Context, in *RequestAddressGeocode, opts ...grpc.CallOption) (*RespAddress, error)
}

type geoProviderGRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewGeoProviderGRPCClient(cc grpc.ClientConnInterface) GeoProviderGRPCClient {
	return &geoProviderGRPCClient{cc}
}

func (c *geoProviderGRPCClient) AddressSearchGRPC(ctx context.Context, in *RequestAddressSearch, opts ...grpc.CallOption) (*RespAddress, error) {
	out := new(RespAddress)
	err := c.cc.Invoke(ctx, "/geoprovider.GeoProviderGRPC/AddressSearchGRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoProviderGRPCClient) AddressGeoCodeGRPC(ctx context.Context, in *RequestAddressGeocode, opts ...grpc.CallOption) (*RespAddress, error) {
	out := new(RespAddress)
	err := c.cc.Invoke(ctx, "/geoprovider.GeoProviderGRPC/AddressGeoCodeGRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GeoProviderGRPCServer is the server API for GeoProviderGRPC service.
// All implementations must embed UnimplementedGeoProviderGRPCServer
// for forward compatibility
type GeoProviderGRPCServer interface {
	AddressSearchGRPC(context.Context, *RequestAddressSearch) (*RespAddress, error)
	AddressGeoCodeGRPC(context.Context, *RequestAddressGeocode) (*RespAddress, error)
	mustEmbedUnimplementedGeoProviderGRPCServer()
}

// UnimplementedGeoProviderGRPCServer must be embedded to have forward compatible implementations.
type UnimplementedGeoProviderGRPCServer struct {
}

func (UnimplementedGeoProviderGRPCServer) AddressSearchGRPC(context.Context, *RequestAddressSearch) (*RespAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddressSearchGRPC not implemented")
}
func (UnimplementedGeoProviderGRPCServer) AddressGeoCodeGRPC(context.Context, *RequestAddressGeocode) (*RespAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddressGeoCodeGRPC not implemented")
}
func (UnimplementedGeoProviderGRPCServer) mustEmbedUnimplementedGeoProviderGRPCServer() {}

// UnsafeGeoProviderGRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GeoProviderGRPCServer will
// result in compilation errors.
type UnsafeGeoProviderGRPCServer interface {
	mustEmbedUnimplementedGeoProviderGRPCServer()
}

func RegisterGeoProviderGRPCServer(s grpc.ServiceRegistrar, srv GeoProviderGRPCServer) {
	s.RegisterService(&GeoProviderGRPC_ServiceDesc, srv)
}

func _GeoProviderGRPC_AddressSearchGRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAddressSearch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoProviderGRPCServer).AddressSearchGRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geoprovider.GeoProviderGRPC/AddressSearchGRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoProviderGRPCServer).AddressSearchGRPC(ctx, req.(*RequestAddressSearch))
	}
	return interceptor(ctx, in, info, handler)
}

func _GeoProviderGRPC_AddressGeoCodeGRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAddressGeocode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoProviderGRPCServer).AddressGeoCodeGRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geoprovider.GeoProviderGRPC/AddressGeoCodeGRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoProviderGRPCServer).AddressGeoCodeGRPC(ctx, req.(*RequestAddressGeocode))
	}
	return interceptor(ctx, in, info, handler)
}

// GeoProviderGRPC_ServiceDesc is the grpc.ServiceDesc for GeoProviderGRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GeoProviderGRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "geoprovider.GeoProviderGRPC",
	HandlerType: (*GeoProviderGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddressSearchGRPC",
			Handler:    _GeoProviderGRPC_AddressSearchGRPC_Handler,
		},
		{
			MethodName: "AddressGeoCodeGRPC",
			Handler:    _GeoProviderGRPC_AddressGeoCodeGRPC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "geoProvaider.proto",
}
