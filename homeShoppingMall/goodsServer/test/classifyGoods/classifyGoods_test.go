package classifyGoods

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"homeShoppingMall/goodsServer/proto"
	"log"
	"testing"
)

var conn *grpc.ClientConn
var err error
var goodsClient proto.GoodsServerClient

func grpcInit() {
	conn, err = grpc.Dial(
		"consul://192.168.198.200:8500/goods?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	goodsClient = proto.NewGoodsServerClient(conn)
}

//获取所有一级分类
func TestGetClassifyInfo(t *testing.T) {
	grpcInit()
	rsp, err := goodsClient.GetClassifyInfo(context.Background(), &proto.GoodsEmpty{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, values := range rsp.Info {
		fmt.Println(values.Id, values.Name, values.Pid)
	}
}

// 获取所有二级分类
func TestGetChildClassifyInfo(t *testing.T) {
	grpcInit()
	rsp, err := goodsClient.GetChildClassifyInfo(context.Background(), &proto.ClassifyChildInfoRequest{
		PName: "灯光",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Info)
	fmt.Println(rsp.ListInfo)
}

//新增商品分类
func TestCreateClassifyInfo(t *testing.T) {
	grpcInit()
	rsp, err := goodsClient.CreateClassifyInfo(context.Background(), &proto.ClassifyCreateInfoRequest{
		Name:  "砖头",
		PName: "地板",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Id, rsp.Name, rsp.Pid)
}

//更新商品分类
func TestUpdateClassifyInfo(t *testing.T) {
	grpcInit()
	if _, err := goodsClient.UpdateClassifyInfo(context.Background(), &proto.ClassifyUpdateInfoRequest{
		OldName: "砖头",
		NewName: "神砖",
		PName:   "地板",
	}); err != nil {
		fmt.Println(err)
	}
}

// 删除商品分类
func TestDeleteClassifyInfo(t *testing.T) {
	grpcInit()
	if _, err := goodsClient.DeleteClassifyInfo(context.Background(), &proto.ClassifyDeleteInfoRequest{Name: "砖头"}); err != nil {
		fmt.Println(err)
	}
}
