// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.11.2
// source: db_hash.proto

package db_hash

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

const (
	HashDatabase_HSet_FullMethodName    = "/api_db_hash.HashDatabase/HSet"
	HashDatabase_HGet_FullMethodName    = "/api_db_hash.HashDatabase/HGet"
	HashDatabase_HGetAll_FullMethodName = "/api_db_hash.HashDatabase/HGetAll"
	HashDatabase_HDel_FullMethodName    = "/api_db_hash.HashDatabase/HDel"
	HashDatabase_HExists_FullMethodName = "/api_db_hash.HashDatabase/HExists"
	HashDatabase_HKeys_FullMethodName   = "/api_db_hash.HashDatabase/HKeys"
	HashDatabase_HVals_FullMethodName   = "/api_db_hash.HashDatabase/HVals"
	HashDatabase_HLen_FullMethodName    = "/api_db_hash.HashDatabase/HLen"
)

// HashDatabaseClient is the client API for HashDatabase service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HashDatabaseClient interface {
	HSet(ctx context.Context, in *HSetRequest, opts ...grpc.CallOption) (*HSetResponse, error)
	HGet(ctx context.Context, in *HGetRequest, opts ...grpc.CallOption) (*HGetResponse, error)
	HGetAll(ctx context.Context, in *HGetAllRequest, opts ...grpc.CallOption) (*HGetAllResponse, error)
	HDel(ctx context.Context, in *HDelRequest, opts ...grpc.CallOption) (*HDelResponse, error)
	HExists(ctx context.Context, in *HExistsRequest, opts ...grpc.CallOption) (*HExistsResponse, error)
	HKeys(ctx context.Context, in *HKeysRequest, opts ...grpc.CallOption) (*HKeysResponse, error)
	HVals(ctx context.Context, in *HValsRequest, opts ...grpc.CallOption) (*HValsResponse, error)
	HLen(ctx context.Context, in *HLenRequest, opts ...grpc.CallOption) (*HLenResponse, error)
}

type hashDatabaseClient struct {
	cc grpc.ClientConnInterface
}

func NewHashDatabaseClient(cc grpc.ClientConnInterface) HashDatabaseClient {
	return &hashDatabaseClient{cc}
}

func (c *hashDatabaseClient) HSet(ctx context.Context, in *HSetRequest, opts ...grpc.CallOption) (*HSetResponse, error) {
	out := new(HSetResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HGet(ctx context.Context, in *HGetRequest, opts ...grpc.CallOption) (*HGetResponse, error) {
	out := new(HGetResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HGetAll(ctx context.Context, in *HGetAllRequest, opts ...grpc.CallOption) (*HGetAllResponse, error) {
	out := new(HGetAllResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HGetAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HDel(ctx context.Context, in *HDelRequest, opts ...grpc.CallOption) (*HDelResponse, error) {
	out := new(HDelResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HDel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HExists(ctx context.Context, in *HExistsRequest, opts ...grpc.CallOption) (*HExistsResponse, error) {
	out := new(HExistsResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HExists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HKeys(ctx context.Context, in *HKeysRequest, opts ...grpc.CallOption) (*HKeysResponse, error) {
	out := new(HKeysResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HKeys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HVals(ctx context.Context, in *HValsRequest, opts ...grpc.CallOption) (*HValsResponse, error) {
	out := new(HValsResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HVals_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashDatabaseClient) HLen(ctx context.Context, in *HLenRequest, opts ...grpc.CallOption) (*HLenResponse, error) {
	out := new(HLenResponse)
	err := c.cc.Invoke(ctx, HashDatabase_HLen_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HashDatabaseServer is the server API for HashDatabase service.
// All implementations must embed UnimplementedHashDatabaseServer
// for forward compatibility
type HashDatabaseServer interface {
	HSet(context.Context, *HSetRequest) (*HSetResponse, error)
	HGet(context.Context, *HGetRequest) (*HGetResponse, error)
	HGetAll(context.Context, *HGetAllRequest) (*HGetAllResponse, error)
	HDel(context.Context, *HDelRequest) (*HDelResponse, error)
	HExists(context.Context, *HExistsRequest) (*HExistsResponse, error)
	HKeys(context.Context, *HKeysRequest) (*HKeysResponse, error)
	HVals(context.Context, *HValsRequest) (*HValsResponse, error)
	HLen(context.Context, *HLenRequest) (*HLenResponse, error)
	mustEmbedUnimplementedHashDatabaseServer()
}

// UnimplementedHashDatabaseServer must be embedded to have forward compatible implementations.
type UnimplementedHashDatabaseServer struct {
}

func (UnimplementedHashDatabaseServer) HSet(context.Context, *HSetRequest) (*HSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HSet not implemented")
}
func (UnimplementedHashDatabaseServer) HGet(context.Context, *HGetRequest) (*HGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HGet not implemented")
}
func (UnimplementedHashDatabaseServer) HGetAll(context.Context, *HGetAllRequest) (*HGetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HGetAll not implemented")
}
func (UnimplementedHashDatabaseServer) HDel(context.Context, *HDelRequest) (*HDelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HDel not implemented")
}
func (UnimplementedHashDatabaseServer) HExists(context.Context, *HExistsRequest) (*HExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HExists not implemented")
}
func (UnimplementedHashDatabaseServer) HKeys(context.Context, *HKeysRequest) (*HKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HKeys not implemented")
}
func (UnimplementedHashDatabaseServer) HVals(context.Context, *HValsRequest) (*HValsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HVals not implemented")
}
func (UnimplementedHashDatabaseServer) HLen(context.Context, *HLenRequest) (*HLenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HLen not implemented")
}
func (UnimplementedHashDatabaseServer) mustEmbedUnimplementedHashDatabaseServer() {}

// UnsafeHashDatabaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HashDatabaseServer will
// result in compilation errors.
type UnsafeHashDatabaseServer interface {
	mustEmbedUnimplementedHashDatabaseServer()
}

func RegisterHashDatabaseServer(s grpc.ServiceRegistrar, srv HashDatabaseServer) {
	s.RegisterService(&HashDatabase_ServiceDesc, srv)
}

func _HashDatabase_HSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HSet(ctx, req.(*HSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HGet(ctx, req.(*HGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HGetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HGetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HGetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HGetAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HGetAll(ctx, req.(*HGetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HDelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HDel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HDel(ctx, req.(*HDelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HExists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HExists(ctx, req.(*HExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HKeys(ctx, req.(*HKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HVals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HValsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HVals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HVals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HVals(ctx, req.(*HValsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashDatabase_HLen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HLenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashDatabaseServer).HLen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashDatabase_HLen_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashDatabaseServer).HLen(ctx, req.(*HLenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HashDatabase_ServiceDesc is the grpc.ServiceDesc for HashDatabase service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HashDatabase_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api_db_hash.HashDatabase",
	HandlerType: (*HashDatabaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HSet",
			Handler:    _HashDatabase_HSet_Handler,
		},
		{
			MethodName: "HGet",
			Handler:    _HashDatabase_HGet_Handler,
		},
		{
			MethodName: "HGetAll",
			Handler:    _HashDatabase_HGetAll_Handler,
		},
		{
			MethodName: "HDel",
			Handler:    _HashDatabase_HDel_Handler,
		},
		{
			MethodName: "HExists",
			Handler:    _HashDatabase_HExists_Handler,
		},
		{
			MethodName: "HKeys",
			Handler:    _HashDatabase_HKeys_Handler,
		},
		{
			MethodName: "HVals",
			Handler:    _HashDatabase_HVals_Handler,
		},
		{
			MethodName: "HLen",
			Handler:    _HashDatabase_HLen_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "db_hash.proto",
}