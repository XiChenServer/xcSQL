// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.11.2
// source: db_set.proto

package db_set

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
	SetDatabase_SAdd_FullMethodName      = "/api_db_set.SetDatabase/SAdd"
	SetDatabase_SRem_FullMethodName      = "/api_db_set.SetDatabase/SRem"
	SetDatabase_SMembers_FullMethodName  = "/api_db_set.SetDatabase/SMembers"
	SetDatabase_SIsMember_FullMethodName = "/api_db_set.SetDatabase/SIsMember"
	SetDatabase_SCard_FullMethodName     = "/api_db_set.SetDatabase/SCard"
)

// SetDatabaseClient is the client API for SetDatabase service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SetDatabaseClient interface {
	// SAdd 将一个或多个成员添加到集合中。
	SAdd(ctx context.Context, in *SAddRequest, opts ...grpc.CallOption) (*SAddResponse, error)
	// SRem 从集合中移除一个或多个成员。
	SRem(ctx context.Context, in *SRemRequest, opts ...grpc.CallOption) (*SRemResponse, error)
	// SMembers 获取集合中的所有成员。
	SMembers(ctx context.Context, in *SMembersRequest, opts ...grpc.CallOption) (*SMembersResponse, error)
	// SIsMember 判断一个成员是否是集合的成员。
	SIsMember(ctx context.Context, in *SIsMemberRequest, opts ...grpc.CallOption) (*SIsMemberResponse, error)
	// SCard 获取集合中成员的数量。
	SCard(ctx context.Context, in *SCardRequest, opts ...grpc.CallOption) (*SCardResponse, error)
}

type setDatabaseClient struct {
	cc grpc.ClientConnInterface
}

func NewSetDatabaseClient(cc grpc.ClientConnInterface) SetDatabaseClient {
	return &setDatabaseClient{cc}
}

func (c *setDatabaseClient) SAdd(ctx context.Context, in *SAddRequest, opts ...grpc.CallOption) (*SAddResponse, error) {
	out := new(SAddResponse)
	err := c.cc.Invoke(ctx, SetDatabase_SAdd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *setDatabaseClient) SRem(ctx context.Context, in *SRemRequest, opts ...grpc.CallOption) (*SRemResponse, error) {
	out := new(SRemResponse)
	err := c.cc.Invoke(ctx, SetDatabase_SRem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *setDatabaseClient) SMembers(ctx context.Context, in *SMembersRequest, opts ...grpc.CallOption) (*SMembersResponse, error) {
	out := new(SMembersResponse)
	err := c.cc.Invoke(ctx, SetDatabase_SMembers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *setDatabaseClient) SIsMember(ctx context.Context, in *SIsMemberRequest, opts ...grpc.CallOption) (*SIsMemberResponse, error) {
	out := new(SIsMemberResponse)
	err := c.cc.Invoke(ctx, SetDatabase_SIsMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *setDatabaseClient) SCard(ctx context.Context, in *SCardRequest, opts ...grpc.CallOption) (*SCardResponse, error) {
	out := new(SCardResponse)
	err := c.cc.Invoke(ctx, SetDatabase_SCard_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SetDatabaseServer is the server API for SetDatabase service.
// All implementations must embed UnimplementedSetDatabaseServer
// for forward compatibility
type SetDatabaseServer interface {
	// SAdd 将一个或多个成员添加到集合中。
	SAdd(context.Context, *SAddRequest) (*SAddResponse, error)
	// SRem 从集合中移除一个或多个成员。
	SRem(context.Context, *SRemRequest) (*SRemResponse, error)
	// SMembers 获取集合中的所有成员。
	SMembers(context.Context, *SMembersRequest) (*SMembersResponse, error)
	// SIsMember 判断一个成员是否是集合的成员。
	SIsMember(context.Context, *SIsMemberRequest) (*SIsMemberResponse, error)
	// SCard 获取集合中成员的数量。
	SCard(context.Context, *SCardRequest) (*SCardResponse, error)
	mustEmbedUnimplementedSetDatabaseServer()
}

// UnimplementedSetDatabaseServer must be embedded to have forward compatible implementations.
type UnimplementedSetDatabaseServer struct {
}

func (UnimplementedSetDatabaseServer) SAdd(context.Context, *SAddRequest) (*SAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SAdd not implemented")
}
func (UnimplementedSetDatabaseServer) SRem(context.Context, *SRemRequest) (*SRemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SRem not implemented")
}
func (UnimplementedSetDatabaseServer) SMembers(context.Context, *SMembersRequest) (*SMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SMembers not implemented")
}
func (UnimplementedSetDatabaseServer) SIsMember(context.Context, *SIsMemberRequest) (*SIsMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SIsMember not implemented")
}
func (UnimplementedSetDatabaseServer) SCard(context.Context, *SCardRequest) (*SCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SCard not implemented")
}
func (UnimplementedSetDatabaseServer) mustEmbedUnimplementedSetDatabaseServer() {}

// UnsafeSetDatabaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SetDatabaseServer will
// result in compilation errors.
type UnsafeSetDatabaseServer interface {
	mustEmbedUnimplementedSetDatabaseServer()
}

func RegisterSetDatabaseServer(s grpc.ServiceRegistrar, srv SetDatabaseServer) {
	s.RegisterService(&SetDatabase_ServiceDesc, srv)
}

func _SetDatabase_SAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetDatabaseServer).SAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetDatabase_SAdd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetDatabaseServer).SAdd(ctx, req.(*SAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SetDatabase_SRem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SRemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetDatabaseServer).SRem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetDatabase_SRem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetDatabaseServer).SRem(ctx, req.(*SRemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SetDatabase_SMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetDatabaseServer).SMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetDatabase_SMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetDatabaseServer).SMembers(ctx, req.(*SMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SetDatabase_SIsMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SIsMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetDatabaseServer).SIsMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetDatabase_SIsMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetDatabaseServer).SIsMember(ctx, req.(*SIsMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SetDatabase_SCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetDatabaseServer).SCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetDatabase_SCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetDatabaseServer).SCard(ctx, req.(*SCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SetDatabase_ServiceDesc is the grpc.ServiceDesc for SetDatabase service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SetDatabase_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api_db_set.SetDatabase",
	HandlerType: (*SetDatabaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SAdd",
			Handler:    _SetDatabase_SAdd_Handler,
		},
		{
			MethodName: "SRem",
			Handler:    _SetDatabase_SRem_Handler,
		},
		{
			MethodName: "SMembers",
			Handler:    _SetDatabase_SMembers_Handler,
		},
		{
			MethodName: "SIsMember",
			Handler:    _SetDatabase_SIsMember_Handler,
		},
		{
			MethodName: "SCard",
			Handler:    _SetDatabase_SCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "db_set.proto",
}
