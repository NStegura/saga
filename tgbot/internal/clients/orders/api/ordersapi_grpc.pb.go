// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: ordersapi.proto

package api

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OrdersApi_CreateOrder_FullMethodName    = "/ordersapi.OrdersApi/CreateOrder"
	OrdersApi_GetOrder_FullMethodName       = "/ordersapi.OrdersApi/GetOrder"
	OrdersApi_GetOrderStates_FullMethodName = "/ordersapi.OrdersApi/GetOrderStates"
	OrdersApi_GetOrders_FullMethodName      = "/ordersapi.OrdersApi/GetOrders"
	OrdersApi_GetPing_FullMethodName        = "/ordersapi.OrdersApi/GetPing"
)

// OrdersApiClient is the client API for OrdersApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrdersApiClient interface {
	CreateOrder(ctx context.Context, in *OrderIn, opts ...grpc.CallOption) (*OrderId, error)
	GetOrder(ctx context.Context, in *OrderId, opts ...grpc.CallOption) (*OrderOut, error)
	GetOrderStates(ctx context.Context, in *OrderId, opts ...grpc.CallOption) (*States, error)
	GetOrders(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*Orders, error)
	GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Pong, error)
}

type ordersApiClient struct {
	cc grpc.ClientConnInterface
}

func NewOrdersApiClient(cc grpc.ClientConnInterface) OrdersApiClient {
	return &ordersApiClient{cc}
}

func (c *ordersApiClient) CreateOrder(ctx context.Context, in *OrderIn, opts ...grpc.CallOption) (*OrderId, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderId)
	err := c.cc.Invoke(ctx, OrdersApi_CreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordersApiClient) GetOrder(ctx context.Context, in *OrderId, opts ...grpc.CallOption) (*OrderOut, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderOut)
	err := c.cc.Invoke(ctx, OrdersApi_GetOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordersApiClient) GetOrderStates(ctx context.Context, in *OrderId, opts ...grpc.CallOption) (*States, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(States)
	err := c.cc.Invoke(ctx, OrdersApi_GetOrderStates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordersApiClient) GetOrders(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*Orders, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Orders)
	err := c.cc.Invoke(ctx, OrdersApi_GetOrders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordersApiClient) GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Pong, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Pong)
	err := c.cc.Invoke(ctx, OrdersApi_GetPing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdersApiServer is the server API for OrdersApi service.
// All implementations must embed UnimplementedOrdersApiServer
// for forward compatibility.
type OrdersApiServer interface {
	CreateOrder(context.Context, *OrderIn) (*OrderId, error)
	GetOrder(context.Context, *OrderId) (*OrderOut, error)
	GetOrderStates(context.Context, *OrderId) (*States, error)
	GetOrders(context.Context, *UserId) (*Orders, error)
	GetPing(context.Context, *empty.Empty) (*Pong, error)
	mustEmbedUnimplementedOrdersApiServer()
}

// UnimplementedOrdersApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrdersApiServer struct{}

func (UnimplementedOrdersApiServer) CreateOrder(context.Context, *OrderIn) (*OrderId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrdersApiServer) GetOrder(context.Context, *OrderId) (*OrderOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
func (UnimplementedOrdersApiServer) GetOrderStates(context.Context, *OrderId) (*States, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderStates not implemented")
}
func (UnimplementedOrdersApiServer) GetOrders(context.Context, *UserId) (*Orders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrders not implemented")
}
func (UnimplementedOrdersApiServer) GetPing(context.Context, *empty.Empty) (*Pong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPing not implemented")
}
func (UnimplementedOrdersApiServer) mustEmbedUnimplementedOrdersApiServer() {}
func (UnimplementedOrdersApiServer) testEmbeddedByValue()                   {}

// UnsafeOrdersApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrdersApiServer will
// result in compilation errors.
type UnsafeOrdersApiServer interface {
	mustEmbedUnimplementedOrdersApiServer()
}

func RegisterOrdersApiServer(s grpc.ServiceRegistrar, srv OrdersApiServer) {
	// If the following call pancis, it indicates UnimplementedOrdersApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrdersApi_ServiceDesc, srv)
}

func _OrdersApi_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersApiServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrdersApi_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersApiServer).CreateOrder(ctx, req.(*OrderIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrdersApi_GetOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersApiServer).GetOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrdersApi_GetOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersApiServer).GetOrder(ctx, req.(*OrderId))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrdersApi_GetOrderStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersApiServer).GetOrderStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrdersApi_GetOrderStates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersApiServer).GetOrderStates(ctx, req.(*OrderId))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrdersApi_GetOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersApiServer).GetOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrdersApi_GetOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersApiServer).GetOrders(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrdersApi_GetPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersApiServer).GetPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrdersApi_GetPing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersApiServer).GetPing(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// OrdersApi_ServiceDesc is the grpc.ServiceDesc for OrdersApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrdersApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ordersapi.OrdersApi",
	HandlerType: (*OrdersApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _OrdersApi_CreateOrder_Handler,
		},
		{
			MethodName: "GetOrder",
			Handler:    _OrdersApi_GetOrder_Handler,
		},
		{
			MethodName: "GetOrderStates",
			Handler:    _OrdersApi_GetOrderStates_Handler,
		},
		{
			MethodName: "GetOrders",
			Handler:    _OrdersApi_GetOrders_Handler,
		},
		{
			MethodName: "GetPing",
			Handler:    _OrdersApi_GetPing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ordersapi.proto",
}
