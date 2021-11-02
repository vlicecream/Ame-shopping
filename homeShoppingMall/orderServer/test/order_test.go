package test

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"homeShoppingMall/orderServer/proto"
	"log"
	"testing"
)

var conn *grpc.ClientConn
var err error
var orderClient proto.OrderClient

func grpcInit() {
	conn, err = grpc.Dial(
		"consul://192.168.198.200:8500/order?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	orderClient = proto.NewOrderClient(conn)
}

// 查看购物车
func TestCheckShoppingCar(t *testing.T) {
	grpcInit()
	rsp, err := orderClient.CheckShoppingCar(context.Background(), &proto.UserInfo{UserID: 99781607077974016})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.ShoppingCarInfo)
}

// 添加购物车
func TestCreateShoppingCar(t *testing.T) {
	grpcInit()
	rsp, err := orderClient.CreateShoppingCar(context.Background(), &proto.CreateCarRequest{
		UserID:   99781607077974016,
		Goods:    "吊灯003",
		Nums:     1,
		Selected: false,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}

// 更新购物车
func TestUpdateShoppingCar(t *testing.T) {
	grpcInit()
	if _, err := orderClient.UpdateShoppingCar(context.Background(), &proto.CreateCarRequest{
		UserID:   99781607077974016,
		Goods:    "吊灯003",
		Nums:     1,
		Selected: true,
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Println("success update")
}

// 删除购物车
func TestDeleteShoppingCar(t *testing.T) {
	grpcInit()
	if _, err := orderClient.DeleteShoppingCar(context.Background(), &proto.DeleteCarRequest{
		UserID: 99781607077974016,
		Goods:  []string{"吊灯003", "吊灯001"},
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Println("success delete")
}

// 查看订单
func TestCheckOrder(t *testing.T) {
	grpcInit()
	rsp, err := orderClient.CheckOrder(context.Background(), &proto.OrderFilterInfo{
		UserID: 99781607077974016,
		Pn:     0,
		PSize:  1,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.OrderInfoResponse)
}

// 获取订单详细信息
func TestCheckOrderDetail(t *testing.T) {
	grpcInit()
	rsp, err := orderClient.CheckOrderDetail(context.Background(), &proto.OrderDetailInfoRequest{
		OrderGoodsNum:  "99781607077974016",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.OrderInfoResponse)
	fmt.Println(rsp.Goods)
}

// 创建订单
func TestCreateOrder(t *testing.T) {
	grpcInit()
	rsp, err := orderClient.CreateOrder(context.Background(), &proto.CreateOrderInfo{
		UserID:  99781607077974016,
		Address: "北京",
		Name:    "Ame1",
		Mobile:  "1383838388",
		Message: "message",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

// 更新订单
func TestUpdateOrder(t *testing.T) {
	grpcInit()
	if _, err := orderClient.UpdateOrderStatus(context.Background(), &proto.OrderInfo{
		Name:    "吊灯001",
		OrderSn: "102056760931520512",
		Status:  "overtime",
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Update success")
}