// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// PriceClient is the client API for Price service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PriceClient interface {
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchReply, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error)
}

type priceClient struct {
	cc grpc.ClientConnInterface
}

func NewPriceClient(cc grpc.ClientConnInterface) PriceClient {
	return &priceClient{cc}
}

func (c *priceClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchReply, error) {
	out := new(FetchReply)
	err := c.cc.Invoke(ctx, "/proto.Price/Fetch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error) {
	out := new(ListReply)
	err := c.cc.Invoke(ctx, "/proto.Price/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PriceServer is the server API for Price service.
// All implementations must embed UnimplementedPriceServer
// for forward compatibility
type PriceServer interface {
	Fetch(context.Context, *FetchRequest) (*FetchReply, error)
	List(context.Context, *ListRequest) (*ListReply, error)
	mustEmbedUnimplementedPriceServer()
}

// UnimplementedPriceServer must be embedded to have forward compatible implementations.
type UnimplementedPriceServer struct {
}

func (UnimplementedPriceServer) Fetch(context.Context, *FetchRequest) (*FetchReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fetch not implemented")
}
func (UnimplementedPriceServer) List(context.Context, *ListRequest) (*ListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedPriceServer) mustEmbedUnimplementedPriceServer() {}

// UnsafePriceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PriceServer will
// result in compilation errors.
type UnsafePriceServer interface {
	mustEmbedUnimplementedPriceServer()
}

func RegisterPriceServer(s grpc.ServiceRegistrar, srv PriceServer) {
	s.RegisterService(&Price_ServiceDesc, srv)
}

func _Price_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PriceServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Price/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PriceServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Price_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PriceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Price/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PriceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Price_ServiceDesc is the grpc.ServiceDesc for Price service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Price_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Price",
	HandlerType: (*PriceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _Price_Fetch_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Price_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/price.proto",
}
