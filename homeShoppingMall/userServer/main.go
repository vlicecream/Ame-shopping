package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"homeShoppingMall/userServer/consulRegister"
	"homeShoppingMall/userServer/dao/mysql"
	"homeShoppingMall/userServer/dao/redis"
	"homeShoppingMall/userServer/handler"
	"homeShoppingMall/userServer/logger"
	"homeShoppingMall/userServer/pkg/snowflake"
	"homeShoppingMall/userServer/proto"
	"homeShoppingMall/userServer/settings"
	"homeShoppingMall/userServer/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 加载配置文件
	if err := settings.Init(); err != nil {
		zap.L().Error("Main settings.Init failed")
		return
	}

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		zap.L().Error("Main logger.Init failed")
		return
	}
	defer zap.L().Sync()

	// 3. 初始化mysql
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		zap.L().Error("Main mysql.Init failed")
		return
	}
	defer mysql.Close()

	// 初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("Main mysql.Init failed")
		return
	}

	// 4. 初始化雪花算法
	if err := snowflake.Init(settings.Conf.SnowflakeConfig.StartTime, settings.Conf.SnowflakeConfig.MachineID); err != nil {
		zap.L().Error("Main snowflake.Init failed")
	}

	// 小插曲 初始化动态端口
	port, err := utils.DynamicPort()
	if err != nil {
		zap.L().Error("Main utils.DynamicPort failed")
	}

	if settings.Conf.Mode == "dev" {
		port = settings.Conf.Port
	}

	// 5. 创建grpc服务
	server := grpc.NewServer()
	proto.RegisterUserSeverServer(server, &handler.UserSeverServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Conf.Host, port))
	if err != nil {
		fmt.Println("网络错误")
	}
	fmt.Println(fmt.Sprintf("正在监听%d端口", port))

	// 6. 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 7. consul服务注册&健康检查
	client, serviceID := consulRegister.ConsulInit()

	// 8. grpc服务启动
	go func() {
		if err := server.Serve(listen); err != nil {
			zap.Error(err)
		}
	}()

	// 9. 退出注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.Error(err)
	}
	zap.L().Info("注销成功")
}
