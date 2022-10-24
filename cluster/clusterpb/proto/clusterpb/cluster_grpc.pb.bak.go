//// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
//
package clusterpb

//
//import (
//	context "context"
//	"github.com/lonng/nano/cluster/clusterpb"
//	grpc "google.golang.org/grpc"
//	codes "google.golang.org/grpc/codes"
//	status "google.golang.org/grpc/status"
//)
//
//// This is a compile-time assertion to ensure that this generated file
//// is compatible with the grpc package it is being compiled against.
//// Requires gRPC-Go v1.32.0 or later.
//const _ = grpc.SupportPackageIsVersion7
//
//// MasterClient is the client API for Master service.
////
//// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//type MasterClient interface {
//	Register(ctx context.Context, in *clusterpb.RegisterRequest, opts ...grpc.CallOption) (*clusterpb.RegisterResponse, error)
//	Unregister(ctx context.Context, in *clusterpb.UnregisterRequest, opts ...grpc.CallOption) (*clusterpb.UnregisterResponse, error)
//}
//
//type masterClient struct {
//	cc grpc.ClientConnInterface
//}
//
//func NewMasterClient(cc grpc.ClientConnInterface) MasterClient {
//	return &masterClient{cc}
//}
//
//func (c *masterClient) Register(ctx context.Context, in *clusterpb.RegisterRequest, opts ...grpc.CallOption) (*clusterpb.RegisterResponse, error) {
//	out := new(clusterpb.RegisterResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Master/Register", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *masterClient) Unregister(ctx context.Context, in *clusterpb.UnregisterRequest, opts ...grpc.CallOption) (*clusterpb.UnregisterResponse, error) {
//	out := new(clusterpb.UnregisterResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Master/Unregister", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//// MasterServer is the server API for Master service.
//// All implementations should embed UnimplementedMasterServer
//// for forward compatibility
//type MasterServer interface {
//	Register(context.Context, *clusterpb.RegisterRequest) (*clusterpb.RegisterResponse, error)
//	Unregister(context.Context, *clusterpb.UnregisterRequest) (*clusterpb.UnregisterResponse, error)
//}
//
//// UnimplementedMasterServer should be embedded to have forward compatible implementations.
//type UnimplementedMasterServer struct {
//}
//
//func (UnimplementedMasterServer) Register(context.Context, *clusterpb.RegisterRequest) (*clusterpb.RegisterResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
//}
//func (UnimplementedMasterServer) Unregister(context.Context, *clusterpb.UnregisterRequest) (*clusterpb.UnregisterResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
//}
//
//// UnsafeMasterServer may be embedded to opt out of forward compatibility for this service.
//// Use of this interface is not recommended, as added methods to MasterServer will
//// result in compilation errors.
//type UnsafeMasterServer interface {
//	mustEmbedUnimplementedMasterServer()
//}
//
//func RegisterMasterServer(s grpc.ServiceRegistrar, srv MasterServer) {
//	s.RegisterService(&Master_ServiceDesc, srv)
//}
//
//func _Master_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.RegisterRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MasterServer).Register(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Master/Register",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MasterServer).Register(ctx, req.(*clusterpb.RegisterRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Master_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.UnregisterRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MasterServer).Unregister(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Master/Unregister",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MasterServer).Unregister(ctx, req.(*clusterpb.UnregisterRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//// Master_ServiceDesc is the grpc.ServiceDesc for Master service.
//// It's only intended for direct use with grpc.RegisterService,
//// and not to be introspected or modified (even as a copy)
//var Master_ServiceDesc = grpc.ServiceDesc{
//	ServiceName: "clusterpb.Master",
//	HandlerType: (*MasterServer)(nil),
//	Methods: []grpc.MethodDesc{
//		{
//			MethodName: "Register",
//			Handler:    _Master_Register_Handler,
//		},
//		{
//			MethodName: "Unregister",
//			Handler:    _Master_Unregister_Handler,
//		},
//	},
//	Streams:  []grpc.StreamDesc{},
//	Metadata: "cluster.proto",
//}
//
//// MemberClient is the client API for Member service.
////
//// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//type MemberClient interface {
//	HandleRequest(ctx context.Context, in *clusterpb.RequestMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error)
//	HandleNotify(ctx context.Context, in *clusterpb.NotifyMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error)
//	HandlePush(ctx context.Context, in *clusterpb.PushMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error)
//	HandleResponse(ctx context.Context, in *clusterpb.ResponseMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error)
//	NewMember(ctx context.Context, in *clusterpb.NewMemberRequest, opts ...grpc.CallOption) (*clusterpb.NewMemberResponse, error)
//	DelMember(ctx context.Context, in *clusterpb.DelMemberRequest, opts ...grpc.CallOption) (*clusterpb.DelMemberResponse, error)
//	SessionClosed(ctx context.Context, in *clusterpb.SessionClosedRequest, opts ...grpc.CallOption) (*clusterpb.SessionClosedResponse, error)
//	CloseSession(ctx context.Context, in *clusterpb.CloseSessionRequest, opts ...grpc.CallOption) (*clusterpb.CloseSessionResponse, error)
//}
//
//type memberClient struct {
//	cc grpc.ClientConnInterface
//}
//
//func NewMemberClient(cc grpc.ClientConnInterface) MemberClient {
//	return &memberClient{cc}
//}
//
//func (c *memberClient) HandleRequest(ctx context.Context, in *clusterpb.RequestMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error) {
//	out := new(clusterpb.MemberHandleResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/HandleRequest", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) HandleNotify(ctx context.Context, in *clusterpb.NotifyMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error) {
//	out := new(clusterpb.MemberHandleResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/HandleNotify", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) HandlePush(ctx context.Context, in *clusterpb.PushMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error) {
//	out := new(clusterpb.MemberHandleResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/HandlePush", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) HandleResponse(ctx context.Context, in *clusterpb.ResponseMessage, opts ...grpc.CallOption) (*clusterpb.MemberHandleResponse, error) {
//	out := new(clusterpb.MemberHandleResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/HandleResponse", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) NewMember(ctx context.Context, in *clusterpb.NewMemberRequest, opts ...grpc.CallOption) (*clusterpb.NewMemberResponse, error) {
//	out := new(clusterpb.NewMemberResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/NewMember", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) DelMember(ctx context.Context, in *clusterpb.DelMemberRequest, opts ...grpc.CallOption) (*clusterpb.DelMemberResponse, error) {
//	out := new(clusterpb.DelMemberResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/DelMember", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) SessionClosed(ctx context.Context, in *clusterpb.SessionClosedRequest, opts ...grpc.CallOption) (*clusterpb.SessionClosedResponse, error) {
//	out := new(clusterpb.SessionClosedResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/SessionClosed", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func (c *memberClient) CloseSession(ctx context.Context, in *clusterpb.CloseSessionRequest, opts ...grpc.CallOption) (*clusterpb.CloseSessionResponse, error) {
//	out := new(clusterpb.CloseSessionResponse)
//	err := c.cc.Invoke(ctx, "/clusterpb.Member/CloseSession", in, out, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//// MemberServer is the server API for Member service.
//// All implementations should embed UnimplementedMemberServer
//// for forward compatibility
//type MemberServer interface {
//	HandleRequest(context.Context, *clusterpb.RequestMessage) (*clusterpb.MemberHandleResponse, error)
//	HandleNotify(context.Context, *clusterpb.NotifyMessage) (*clusterpb.MemberHandleResponse, error)
//	HandlePush(context.Context, *clusterpb.PushMessage) (*clusterpb.MemberHandleResponse, error)
//	HandleResponse(context.Context, *clusterpb.ResponseMessage) (*clusterpb.MemberHandleResponse, error)
//	NewMember(context.Context, *clusterpb.NewMemberRequest) (*clusterpb.NewMemberResponse, error)
//	DelMember(context.Context, *clusterpb.DelMemberRequest) (*clusterpb.DelMemberResponse, error)
//	SessionClosed(context.Context, *clusterpb.SessionClosedRequest) (*clusterpb.SessionClosedResponse, error)
//	CloseSession(context.Context, *clusterpb.CloseSessionRequest) (*clusterpb.CloseSessionResponse, error)
//}
//
//// UnimplementedMemberServer should be embedded to have forward compatible implementations.
//type UnimplementedMemberServer struct {
//}
//
//func (UnimplementedMemberServer) HandleRequest(context.Context, *clusterpb.RequestMessage) (*clusterpb.MemberHandleResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method HandleRequest not implemented")
//}
//func (UnimplementedMemberServer) HandleNotify(context.Context, *clusterpb.NotifyMessage) (*clusterpb.MemberHandleResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method HandleNotify not implemented")
//}
//func (UnimplementedMemberServer) HandlePush(context.Context, *clusterpb.PushMessage) (*clusterpb.MemberHandleResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method HandlePush not implemented")
//}
//func (UnimplementedMemberServer) HandleResponse(context.Context, *clusterpb.ResponseMessage) (*clusterpb.MemberHandleResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method HandleResponse not implemented")
//}
//func (UnimplementedMemberServer) NewMember(context.Context, *clusterpb.NewMemberRequest) (*clusterpb.NewMemberResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method NewMember not implemented")
//}
//func (UnimplementedMemberServer) DelMember(context.Context, *clusterpb.DelMemberRequest) (*clusterpb.DelMemberResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method DelMember not implemented")
//}
//func (UnimplementedMemberServer) SessionClosed(context.Context, *clusterpb.SessionClosedRequest) (*clusterpb.SessionClosedResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method SessionClosed not implemented")
//}
//func (UnimplementedMemberServer) CloseSession(context.Context, *clusterpb.CloseSessionRequest) (*clusterpb.CloseSessionResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method CloseSession not implemented")
//}
//
//// UnsafeMemberServer may be embedded to opt out of forward compatibility for this service.
//// Use of this interface is not recommended, as added methods to MemberServer will
//// result in compilation errors.
//type UnsafeMemberServer interface {
//	mustEmbedUnimplementedMemberServer()
//}
//
//func RegisterMemberServer(s grpc.ServiceRegistrar, srv MemberServer) {
//	s.RegisterService(&Member_ServiceDesc, srv)
//}
//
//func _Member_HandleRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.RequestMessage)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).HandleRequest(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/HandleRequest",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).HandleRequest(ctx, req.(*clusterpb.RequestMessage))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_HandleNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.NotifyMessage)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).HandleNotify(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/HandleNotify",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).HandleNotify(ctx, req.(*clusterpb.NotifyMessage))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_HandlePush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.PushMessage)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).HandlePush(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/HandlePush",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).HandlePush(ctx, req.(*clusterpb.PushMessage))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_HandleResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.ResponseMessage)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).HandleResponse(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/HandleResponse",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).HandleResponse(ctx, req.(*clusterpb.ResponseMessage))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_NewMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.NewMemberRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).NewMember(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/NewMember",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).NewMember(ctx, req.(*clusterpb.NewMemberRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_DelMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.DelMemberRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).DelMember(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/DelMember",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).DelMember(ctx, req.(*clusterpb.DelMemberRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_SessionClosed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.SessionClosedRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).SessionClosed(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/SessionClosed",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).SessionClosed(ctx, req.(*clusterpb.SessionClosedRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//func _Member_CloseSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//	in := new(clusterpb.CloseSessionRequest)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//	if interceptor == nil {
//		return srv.(MemberServer).CloseSession(ctx, in)
//	}
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/clusterpb.Member/CloseSession",
//	}
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(MemberServer).CloseSession(ctx, req.(*clusterpb.CloseSessionRequest))
//	}
//	return interceptor(ctx, in, info, handler)
//}
//
//// Member_ServiceDesc is the grpc.ServiceDesc for Member service.
//// It's only intended for direct use with grpc.RegisterService,
//// and not to be introspected or modified (even as a copy)
//var Member_ServiceDesc = grpc.ServiceDesc{
//	ServiceName: "clusterpb.Member",
//	HandlerType: (*MemberServer)(nil),
//	Methods: []grpc.MethodDesc{
//		{
//			MethodName: "HandleRequest",
//			Handler:    _Member_HandleRequest_Handler,
//		},
//		{
//			MethodName: "HandleNotify",
//			Handler:    _Member_HandleNotify_Handler,
//		},
//		{
//			MethodName: "HandlePush",
//			Handler:    _Member_HandlePush_Handler,
//		},
//		{
//			MethodName: "HandleResponse",
//			Handler:    _Member_HandleResponse_Handler,
//		},
//		{
//			MethodName: "NewMember",
//			Handler:    _Member_NewMember_Handler,
//		},
//		{
//			MethodName: "DelMember",
//			Handler:    _Member_DelMember_Handler,
//		},
//		{
//			MethodName: "SessionClosed",
//			Handler:    _Member_SessionClosed_Handler,
//		},
//		{
//			MethodName: "CloseSession",
//			Handler:    _Member_CloseSession_Handler,
//		},
//	},
//	Streams:  []grpc.StreamDesc{},
//	Metadata: "cluster.proto",
//}
