package consulR

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"homeShoppingMallGin/orderAPI/proto"
	"homeShoppingMallGin/orderAPI/settings"
	"homeShoppingMallGin/orderAPI/utils/otgrpc"
)

var OrderClient proto.OrderClient
var InventoryClient proto.InventoryClient

func Init() {
	// 连接orderClient
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port,
			settings.Conf.ConsulConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.L().Error("consulConn grpc.Dial failed", zap.Error(err))
	}
	// 生成grpc的client并调用接口
	OrderClient = proto.NewOrderClient(orderConn)

	// 连接orderClient
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port,
			settings.Conf.InventoryConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.L().Error("consulConn grpc.Dial failed", zap.Error(err))
	}
	// 生成grpc的client并调用接口
	InventoryClient = proto.NewInventoryClient(inventoryConn)
}
