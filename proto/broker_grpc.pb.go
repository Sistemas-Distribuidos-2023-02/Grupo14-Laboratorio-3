// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: broker.proto

package proto

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
	Broker_RedirectInformant_FullMethodName = "/broker.Broker/RedirectInformant"
	Broker_Mediate_FullMethodName           = "/broker.Broker/Mediate"
)

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	RedirectInformant(ctx context.Context, in *InformantRequest, opts ...grpc.CallOption) (*FulcrumAddress, error)
	Mediate(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Acknowledgement, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) RedirectInformant(ctx context.Context, in *InformantRequest, opts ...grpc.CallOption) (*FulcrumAddress, error) {
	out := new(FulcrumAddress)
	err := c.cc.Invoke(ctx, Broker_RedirectInformant_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) Mediate(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Acknowledgement, error) {
	out := new(Acknowledgement)
	err := c.cc.Invoke(ctx, Broker_Mediate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	RedirectInformant(context.Context, *InformantRequest) (*FulcrumAddress, error)
	Mediate(context.Context, *Message) (*Acknowledgement, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) RedirectInformant(context.Context, *InformantRequest) (*FulcrumAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RedirectInformant not implemented")
}
func (UnimplementedBrokerServer) Mediate(context.Context, *Message) (*Acknowledgement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mediate not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_RedirectInformant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InformantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).RedirectInformant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_RedirectInformant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).RedirectInformant(ctx, req.(*InformantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_Mediate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Mediate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_Mediate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Mediate(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "broker.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RedirectInformant",
			Handler:    _Broker_RedirectInformant_Handler,
		},
		{
			MethodName: "Mediate",
			Handler:    _Broker_Mediate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "broker.proto",
}
