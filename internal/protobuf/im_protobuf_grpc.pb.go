// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: im_protobuf.proto

package protobuf

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

// AccServerClient is the client API for AccServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccServerClient interface {
	// 发送消息
	SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRsp, error)
	// 给这台机器的房间内所有用户发送消息
	SendMsgAll(ctx context.Context, in *SendMsgAllReq, opts ...grpc.CallOption) (*SendMsgAllRsp, error)
	// 获取用户列表
	GetRoomUserList(ctx context.Context, in *GetRoomUserListReq, opts ...grpc.CallOption) (*GetRoomUserListRsp, error)
}

type accServerClient struct {
	cc grpc.ClientConnInterface
}

func NewAccServerClient(cc grpc.ClientConnInterface) AccServerClient {
	return &accServerClient{cc}
}

func (c *accServerClient) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRsp, error) {
	out := new(SendMsgRsp)
	err := c.cc.Invoke(ctx, "/protobuf.AccServer/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accServerClient) SendMsgAll(ctx context.Context, in *SendMsgAllReq, opts ...grpc.CallOption) (*SendMsgAllRsp, error) {
	out := new(SendMsgAllRsp)
	err := c.cc.Invoke(ctx, "/protobuf.AccServer/SendMsgAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accServerClient) GetRoomUserList(ctx context.Context, in *GetRoomUserListReq, opts ...grpc.CallOption) (*GetRoomUserListRsp, error) {
	out := new(GetRoomUserListRsp)
	err := c.cc.Invoke(ctx, "/protobuf.AccServer/GetRoomUserList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccServerServer is the server API for AccServer service.
// All implementations must embed UnimplementedAccServerServer
// for forward compatibility
type AccServerServer interface {
	// 发送消息
	SendMsg(context.Context, *SendMsgReq) (*SendMsgRsp, error)
	// 给这台机器的房间内所有用户发送消息
	SendMsgAll(context.Context, *SendMsgAllReq) (*SendMsgAllRsp, error)
	// 获取用户列表
	GetRoomUserList(context.Context, *GetRoomUserListReq) (*GetRoomUserListRsp, error)
	mustEmbedUnimplementedAccServerServer()
}

// UnimplementedAccServerServer must be embedded to have forward compatible implementations.
type UnimplementedAccServerServer struct {
}

func (UnimplementedAccServerServer) SendMsg(context.Context, *SendMsgReq) (*SendMsgRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (UnimplementedAccServerServer) SendMsgAll(context.Context, *SendMsgAllReq) (*SendMsgAllRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsgAll not implemented")
}
func (UnimplementedAccServerServer) GetRoomUserList(context.Context, *GetRoomUserListReq) (*GetRoomUserListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomUserList not implemented")
}
func (UnimplementedAccServerServer) mustEmbedUnimplementedAccServerServer() {}

// UnsafeAccServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccServerServer will
// result in compilation errors.
type UnsafeAccServerServer interface {
	mustEmbedUnimplementedAccServerServer()
}

func RegisterAccServerServer(s grpc.ServiceRegistrar, srv AccServerServer) {
	s.RegisterService(&AccServer_ServiceDesc, srv)
}

func _AccServer_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.AccServer/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).SendMsg(ctx, req.(*SendMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccServer_SendMsgAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).SendMsgAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.AccServer/SendMsgAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).SendMsgAll(ctx, req.(*SendMsgAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccServer_GetRoomUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomUserListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).GetRoomUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.AccServer/GetRoomUserList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).GetRoomUserList(ctx, req.(*GetRoomUserListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AccServer_ServiceDesc is the grpc.ServiceDesc for AccServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.AccServer",
	HandlerType: (*AccServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMsg",
			Handler:    _AccServer_SendMsg_Handler,
		},
		{
			MethodName: "SendMsgAll",
			Handler:    _AccServer_SendMsgAll_Handler,
		},
		{
			MethodName: "GetRoomUserList",
			Handler:    _AccServer_GetRoomUserList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "im_protobuf.proto",
}
