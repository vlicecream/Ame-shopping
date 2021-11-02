package goods

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
		"consul://120.26.67.141:8500/goodsServer?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	goodsClient = proto.NewGoodsServerClient(conn)
}

// 根据筛选条件来查询商品
//func TestGetClassifyGoods(t *testing.T) {
//	grpcInit()
//	rsp, err := goodsClient.GetClassifyGoods(context.Background(), &proto.ClassifyGoodsInfoRequest{
//		PriceMin:      0,
//		PriceMax:      0,
//		Name:          "",
//		TopClassify:   "",
//		IsNew:         false,
//		IsHot:         false,
//		IsShow:        false,
//		IsFreightFree: false,
//		Pn:            1,
//		PSize:         1,
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(rsp.GoodsInfo)
//}

// 批量查询商品
func TestBatchGetGoods(t *testing.T) {
	grpcInit()
	rsp, err := goodsClient.BatchGetGoods(context.Background(), &proto.BathGoodsNameInfoRequest{Name: []string{"吊灯001", "吊灯002"}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.GoodsInfo)
}

//新增商品
//func TestCreateGoodsInfo(t *testing.T) {
//	grpcInit()
//	rsp, err := goodsClient.CreateGoodsInfo(context.Background(), &proto.GoodsCreateInfoRequest{
//		Name:              "吊灯003",
//		GoodsPrice:        100,
//		PromotionPrice:    0,
//		GoodsIntroduction: "这是吊灯003",
//		CreateTime:        "",
//		ClassifyGoods:     "吊灯",
//		QuantityStock:     200,
//		SalesVolume:       0,
//		CollectNum:        0,
//		IsNew:             true,
//		IsHot:             false,
//		IsShow:            false,
//		IsFreightFree:     false,
//		Image:             nil,
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(rsp.Name)
//}

// 更新商品
//func TestUpdateGoodsInfo(t *testing.T) {
//	grpcInit()
//	if _, err := goodsClient.UpdateGoodsInfo(context.Background(), &proto.GoodsCreateInfoRequest{
//		Name:           "吊灯003",
//		GoodsPrice:     100,
//		PromotionPrice: 105,
//		GoodsIntroduction: "这是吊灯003",
//		CreateTime:    "",
//		ClassifyGoods: "",
//		QuantityStock: 100,
//		SalesVolume:   0,
//		CollectNum:    0,
//		IsNew:         false,
//		IsHot:         false,
//		IsFreightFree: false,
//		Image:         nil,
//	}); err != nil {
//		fmt.Println(err)
//	}
//}

// 删除商品
//func TestDeleteGoodsInfo(t *testing.T) {
//	grpcInit()
//	if _, err := goodsClient.DeleteGoodsInfo(context.Background(), &proto.GoodsDeleteInfoRequest{Name: "吊灯003"}); err != nil {
//		fmt.Println(err)
//	}
//}
