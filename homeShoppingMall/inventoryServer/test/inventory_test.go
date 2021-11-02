package test

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"homeShoppingMall/inventoryServer/proto"
	"log"
	"testing"
)

var conn *grpc.ClientConn
var err error
var goodsClient proto.InventoryClient

func grpcInit() {
	conn, err = grpc.Dial(
		"consul://192.168.198.200:8500/inventory?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	goodsClient = proto.NewInventoryClient(conn)
}

// 添加商品和库存
func TestSetGoodsInventory(t *testing.T) {
	grpcInit()
	// 调用接口
	if _, err := goodsClient.SetGoodsInventory(context.Background(), &proto.GoodsInfo{
		Goods:        "吊灯001",
		InventoryNum: 200,
	}); err != nil {
		fmt.Println(err)
	}
}

// 拿取商品库存
func TestGetGoodsInventory(t *testing.T) {
	grpcInit()
	rsp, err := goodsClient.GetGoodsInventory(context.Background(), &proto.GoodsInfo{
		Goods:        "吊灯001",
		InventoryNum: 0,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}

// 预减库存
func TestSell(t *testing.T) {
	grpcInit()
	_, err := goodsClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: []*proto.GoodsInfo{
		{Goods:        "吊灯001", InventoryNum: 1},
		{Goods:        "吊灯002", InventoryNum: 2},
	}})
	if err != nil {
		fmt.Println(err)
	}
}

// 归还库存
func TestReBack(t *testing.T) {
	grpcInit()
	_, err := goodsClient.ReBack(context.Background(), &proto.SellInfo{GoodsInfo: []*proto.GoodsInfo{
		{Goods:        "吊灯001", InventoryNum: 3},
		{Goods:        "吊灯002", InventoryNum: 3},
	}})
	if err != nil {
		fmt.Println(err)
	}
}
