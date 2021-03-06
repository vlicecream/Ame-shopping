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

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	// 购物车
	CheckShoppingCar(ctx context.Context, in *UserInfo, opts ...grpc.CallOption) (*ShoppingCarListResponse, error)
	CreateShoppingCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*ShoppingCarInfo, error)
	UpdateShoppingCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*OrderEmpty, error)
	DeleteShoppingCar(ctx context.Context, in *DeleteCarRequest, opts ...grpc.CallOption) (*OrderEmpty, error)
	// 订单
	CheckOrder(ctx context.Context, in *OrderFilterInfo, opts ...grpc.CallOption) (*OrderListResponse, error)
	CreateOrder(ctx context.Context, in *CreateOrderInfo, opts ...grpc.CallOption) (*OrderInfoResponse, error)
	CheckOrderDetail(ctx context.Context, in *OrderDetailInfoRequest, opts ...grpc.CallOption) (*OrderDetailResponse, error)
	UpdateOrderStatus(ctx context.Context, in *OrderInfo, opts ...grpc.CallOption) (*OrderEmpty, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) CheckShoppingCar(ctx context.Context, in *UserInfo, opts ...grpc.CallOption) (*ShoppingCarListResponse, error) {
	out := new(ShoppingCarListResponse)
	err := c.cc.Invoke(ctx, "/order/CheckShoppingCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) CreateShoppingCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*ShoppingCarInfo, error) {
	out := new(ShoppingCarInfo)
	err := c.cc.Invoke(ctx, "/order/CreateShoppingCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) UpdateShoppingCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*OrderEmpty, error) {
	out := new(OrderEmpty)
	err := c.cc.Invoke(ctx, "/order/UpdateShoppingCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) DeleteShoppingCar(ctx context.Context, in *DeleteCarRequest, opts ...grpc.CallOption) (*OrderEmpty, error) {
	out := new(OrderEmpty)
	err := c.cc.Invoke(ctx, "/order/DeleteShoppingCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) CheckOrder(ctx context.Context, in *OrderFilterInfo, opts ...grpc.CallOption) (*OrderListResponse, error) {
	out := new(OrderListResponse)
	err := c.cc.Invoke(ctx, "/order/CheckOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) CreateOrder(ctx context.Context, in *CreateOrderInfo, opts ...grpc.CallOption) (*OrderInfoResponse, error) {
	out := new(OrderInfoResponse)
	err := c.cc.Invoke(ctx, "/order/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) CheckOrderDetail(ctx context.Context, in *OrderDetailInfoRequest, opts ...grpc.CallOption) (*OrderDetailResponse, error) {
	out := new(OrderDetailResponse)
	err := c.cc.Invoke(ctx, "/order/CheckOrderDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) UpdateOrderStatus(ctx context.Context, in *OrderInfo, opts ...grpc.CallOption) (*OrderEmpty, error) {
	out := new(OrderEmpty)
	err := c.cc.Invoke(ctx, "/order/UpdateOrderStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility
type OrderServer interface {
	// 购物车
	CheckShoppingCar(context.Context, *UserInfo) (*ShoppingCarListResponse, error)
	CreateShoppingCar(context.Context, *CreateCarRequest) (*ShoppingCarInfo, error)
	UpdateShoppingCar(context.Context, *CreateCarRequest) (*OrderEmpty, error)
	DeleteShoppingCar(context.Context, *DeleteCarRequest) (*OrderEmpty, error)
	// 订单
	CheckOrder(context.Context, *OrderFilterInfo) (*OrderListResponse, error)
	CreateOrder(context.Context, *CreateOrderInfo) (*OrderInfoResponse, error)
	CheckOrderDetail(context.Context, *OrderDetailInfoRequest) (*OrderDetailResponse, error)
	UpdateOrderStatus(context.Context, *OrderInfo) (*OrderEmpty, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServer struct {
}

func (UnimplementedOrderServer) CheckShoppingCar(context.Context, *UserInfo) (*ShoppingCarListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckShoppingCar not implemented")
}
func (UnimplementedOrderServer) CreateShoppingCar(context.Context, *CreateCarRequest) (*ShoppingCarInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShoppingCar not implemented")
}
func (UnimplementedOrderServer) UpdateShoppingCar(context.Context, *CreateCarRequest) (*OrderEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateShoppingCar not implemented")
}
func (UnimplementedOrderServer) DeleteShoppingCar(context.Context, *DeleteCarRequest) (*OrderEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShoppingCar not implemented")
}
func (UnimplementedOrderServer) CheckOrder(context.Context, *OrderFilterInfo) (*OrderListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckOrder not implemented")
}
func (UnimplementedOrderServer) CreateOrder(context.Context, *CreateOrderInfo) (*OrderInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderServer) CheckOrderDetail(context.Context, *OrderDetailInfoRequest) (*OrderDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckOrderDetail not implemented")
}
func (UnimplementedOrderServer) UpdateOrderStatus(context.Context, *OrderInfo) (*OrderEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrderStatus not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_CheckShoppingCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CheckShoppingCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/CheckShoppingCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CheckShoppingCar(ctx, req.(*UserInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_CreateShoppingCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CreateShoppingCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/CreateShoppingCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CreateShoppingCar(ctx, req.(*CreateCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_UpdateShoppingCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).UpdateShoppingCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/UpdateShoppingCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).UpdateShoppingCar(ctx, req.(*CreateCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_DeleteShoppingCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).DeleteShoppingCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/DeleteShoppingCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).DeleteShoppingCar(ctx, req.(*DeleteCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_CheckOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderFilterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CheckOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/CheckOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CheckOrder(ctx, req.(*OrderFilterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CreateOrder(ctx, req.(*CreateOrderInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_CheckOrderDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderDetailInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CheckOrderDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/CheckOrderDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CheckOrderDetail(ctx, req.(*OrderDetailInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_UpdateOrderStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).UpdateOrderStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order/UpdateOrderStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).UpdateOrderStatus(ctx, req.(*OrderInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckShoppingCar",
			Handler:    _Order_CheckShoppingCar_Handler,
		},
		{
			MethodName: "CreateShoppingCar",
			Handler:    _Order_CreateShoppingCar_Handler,
		},
		{
			MethodName: "UpdateShoppingCar",
			Handler:    _Order_UpdateShoppingCar_Handler,
		},
		{
			MethodName: "DeleteShoppingCar",
			Handler:    _Order_DeleteShoppingCar_Handler,
		},
		{
			MethodName: "CheckOrder",
			Handler:    _Order_CheckOrder_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _Order_CreateOrder_Handler,
		},
		{
			MethodName: "CheckOrderDetail",
			Handler:    _Order_CheckOrderDetail_Handler,
		},
		{
			MethodName: "UpdateOrderStatus",
			Handler:    _Order_UpdateOrderStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
