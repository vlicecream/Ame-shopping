package banner

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
var GoodsClient proto.GoodsServerClient

func GrpcInit() {
	conn, err = grpc.Dial(
		"consul://192.168.198.200:8500/goods?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	GoodsClient = proto.NewGoodsServerClient(conn)
}

// TestGetBannerInfo 查询所有轮播图
func TestGetBannerInfo(t *testing.T) {
	GrpcInit()
	// 调用查询所有轮播图
	rsp, err := GoodsClient.GetBannerInfo(context.Background(), &proto.GoodsEmpty{})
	if err != nil {
		fmt.Println("goodsClient.GetBannerInfo failed")
	}
	for _, values := range rsp.BannerInfo {
		fmt.Println(values.ImageUrl, values.ImageGoodsUrl, values.Level)
	}
}

// 创建轮播图
func TestCreateBannerInfo(t *testing.T) {
	GrpcInit()
	rsp, err := GoodsClient.CreateBannerInfo(context.Background(), &proto.BannerCreateInfoRequest{
		ImageUrl:      "https://www.Ame4.com",
		ImageGoodsUrl: "https://www.Ame4goods.com",
		Level:         4,
	})
	if err != nil {
		fmt.Println("goodsClient.CreateBannerInfo failed")
	}
	fmt.Println(rsp)
}

// 更新轮播图
func TestUpdateBannerInfo(t *testing.T) {
	GrpcInit()
	if _, err := GoodsClient.UpdateBannerInfo(context.Background(), &proto.BannerCreateInfoRequest{
		ImageUrl:      "https://www.Ame3.com",
		ImageGoodsUrl: "https://www.Ame3goods.com",
		Level:         3,
	}); err != nil {
		fmt.Println("goodsClient.UpdateBannerInfo failed")
	}
}

// 删除轮播图
//func TestDeleteBannerInfo(t *testing.T) {
//	GrpcInit()
//	if _, err := GoodsClient.DeleteBannerInfo(context.Background(), &proto.BannerDeleteInfoRequest{Id: 4}); err != nil{
//		fmt.Println("goodsClient.DeleteBannerInfo failed")
//	}
//}