package consulRegister

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"homeShoppingMall/orderServer/proto"
	"homeShoppingMall/orderServer/settings"
	"log"
)

var GoodsClient proto.GoodsServerClient
var InventoryClient proto.InventoryClient

func ConsulInit() (*api.Client, string) {
	// 配置客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port) // consul WEB IP

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}


	// 生成检查对象
	check := &api.AgentServiceCheck{
		Timeout:                        "5s", // 超时5秒错误
		Interval:                       "5s", // 5s检查一次
		GRPC:                           fmt.Sprintf("%s:%d", settings.Conf.ConsulConfig.Host, settings.Conf.Port),
		DeregisterCriticalServiceAfter: "10s",
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Address = settings.Conf.ConsulConfig.Host
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = settings.Conf.Port
	registration.Name = settings.Conf.ConsulConfig.Name
	registration.Tags = settings.Conf.ConsulConfig.Tag
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	return client, serviceID
}

func GrpcInit() error {
	// 商品
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port,
			settings.Conf.GoodsClient.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	GoodsClient = proto.NewGoodsServerClient(goodsConn)

	// 库存
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port,
			settings.Conf.InventoryClient.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	InventoryClient = proto.NewInventoryClient(inventoryConn)
	return err
}
