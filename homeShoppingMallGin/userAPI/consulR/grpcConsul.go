package consulR

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"homeShoppingMallGin/userAPI/proto"
	"homeShoppingMallGin/userAPI/settings"
)

var UserSrvClient proto.UserSeverClient

func Init() {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port,
			settings.Conf.ConsulConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Error("consulConn grpc.Dial failed", zap.Error(err))
	}
	// 生成grpc的client并调用接口
	UserClient := proto.NewUserSeverClient(conn)
	UserSrvClient = UserClient
}
