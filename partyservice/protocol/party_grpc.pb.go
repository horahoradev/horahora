// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

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

// PartyserviceClient is the client API for Partyservice service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PartyserviceClient interface {
	NewWatchParty(ctx context.Context, in *NewPartyRequest, opts ...grpc.CallOption) (*NewPartyResponse, error)
	BecomeLeader(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*LeaderResponse, error)
	JoinParty(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error)
	HeartBeat(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error)
	GetPartyState(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*PartyState, error)
	NextVideo(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error)
}

type partyserviceClient struct {
	cc grpc.ClientConnInterface
}

func NewPartyserviceClient(cc grpc.ClientConnInterface) PartyserviceClient {
	return &partyserviceClient{cc}
}

func (c *partyserviceClient) NewWatchParty(ctx context.Context, in *NewPartyRequest, opts ...grpc.CallOption) (*NewPartyResponse, error) {
	out := new(NewPartyResponse)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/NewWatchParty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyserviceClient) BecomeLeader(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*LeaderResponse, error) {
	out := new(LeaderResponse)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/BecomeLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyserviceClient) JoinParty(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/JoinParty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyserviceClient) HeartBeat(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyserviceClient) GetPartyState(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*PartyState, error) {
	out := new(PartyState)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/GetPartyState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyserviceClient) NextVideo(ctx context.Context, in *PartyRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Partyservice/NextVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PartyserviceServer is the server API for Partyservice service.
// All implementations must embed UnimplementedPartyserviceServer
// for forward compatibility
type PartyserviceServer interface {
	NewWatchParty(context.Context, *NewPartyRequest) (*NewPartyResponse, error)
	BecomeLeader(context.Context, *PartyRequest) (*LeaderResponse, error)
	JoinParty(context.Context, *PartyRequest) (*Empty, error)
	HeartBeat(context.Context, *PartyRequest) (*Empty, error)
	GetPartyState(context.Context, *PartyRequest) (*PartyState, error)
	NextVideo(context.Context, *PartyRequest) (*Empty, error)
	mustEmbedUnimplementedPartyserviceServer()
}

// UnimplementedPartyserviceServer must be embedded to have forward compatible implementations.
type UnimplementedPartyserviceServer struct {
}

func (UnimplementedPartyserviceServer) NewWatchParty(context.Context, *NewPartyRequest) (*NewPartyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewWatchParty not implemented")
}
func (UnimplementedPartyserviceServer) BecomeLeader(context.Context, *PartyRequest) (*LeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BecomeLeader not implemented")
}
func (UnimplementedPartyserviceServer) JoinParty(context.Context, *PartyRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinParty not implemented")
}
func (UnimplementedPartyserviceServer) HeartBeat(context.Context, *PartyRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (UnimplementedPartyserviceServer) GetPartyState(context.Context, *PartyRequest) (*PartyState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPartyState not implemented")
}
func (UnimplementedPartyserviceServer) NextVideo(context.Context, *PartyRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextVideo not implemented")
}
func (UnimplementedPartyserviceServer) mustEmbedUnimplementedPartyserviceServer() {}

// UnsafePartyserviceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PartyserviceServer will
// result in compilation errors.
type UnsafePartyserviceServer interface {
	mustEmbedUnimplementedPartyserviceServer()
}

func RegisterPartyserviceServer(s grpc.ServiceRegistrar, srv PartyserviceServer) {
	s.RegisterService(&Partyservice_ServiceDesc, srv)
}

func _Partyservice_NewWatchParty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewPartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).NewWatchParty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/NewWatchParty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).NewWatchParty(ctx, req.(*NewPartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Partyservice_BecomeLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).BecomeLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/BecomeLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).BecomeLeader(ctx, req.(*PartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Partyservice_JoinParty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).JoinParty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/JoinParty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).JoinParty(ctx, req.(*PartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Partyservice_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).HeartBeat(ctx, req.(*PartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Partyservice_GetPartyState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).GetPartyState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/GetPartyState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).GetPartyState(ctx, req.(*PartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Partyservice_NextVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyserviceServer).NextVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Partyservice/NextVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyserviceServer).NextVideo(ctx, req.(*PartyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Partyservice_ServiceDesc is the grpc.ServiceDesc for Partyservice service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Partyservice_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Partyservice",
	HandlerType: (*PartyserviceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewWatchParty",
			Handler:    _Partyservice_NewWatchParty_Handler,
		},
		{
			MethodName: "BecomeLeader",
			Handler:    _Partyservice_BecomeLeader_Handler,
		},
		{
			MethodName: "JoinParty",
			Handler:    _Partyservice_JoinParty_Handler,
		},
		{
			MethodName: "HeartBeat",
			Handler:    _Partyservice_HeartBeat_Handler,
		},
		{
			MethodName: "GetPartyState",
			Handler:    _Partyservice_GetPartyState_Handler,
		},
		{
			MethodName: "NextVideo",
			Handler:    _Partyservice_NextVideo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "party.proto",
}