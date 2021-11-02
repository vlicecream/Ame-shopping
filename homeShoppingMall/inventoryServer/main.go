package main

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"homeShoppingMall/inventoryServer/consulRegister"
	"homeShoppingMall/inventoryServer/dao/mysql"
	"homeShoppingMall/inventoryServer/dao/redis"
	"homeShoppingMall/inventoryServer/handler"
	"homeShoppingMall/inventoryServer/logger"
	"homeShoppingMall/inventoryServer/proto"
	"homeShoppingMall/inventoryServer/settings"
	"homeShoppingMall/inventoryServer/utils"
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
	if err := logger.Init(settings.Conf.Mode); err != nil {
		zap.L().Error("Main logger.Init failed")
		return
	}
	defer zap.L().Sync()

	// 3. 初始化mysql
	if err := mysql.Init(); err != nil {
		zap.L().Error("Main mysql.Init failed")
		return
	}
	defer mysql.Close()
	// 初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("Main mysql.Init failed")
		return
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
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
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
			zap.L().Error("Main rocketmq.NewPushConsumer failed",zap.Error(err) )
		}
	}()

	/*9. rocketmq  summer消费*/
	// 创建一个PushConsumer实例
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", settings.Conf.RocketMQ.Host, settings.Conf.RocketMQ.Port)}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("Ame"), // 组名字 同时也可以利用组名字来负载均衡
	)
	if err != nil {
		zap.L().Error("Main rocketmq.NewPushConsumer failed")
	}

	// 订阅一个主题（目前只支持一个主题），并定义您的消费功能
	if err := c.Subscribe("order", consumer.MessageSelector{}, handler.AutoBack); err != nil {
		zap.L().Error("Main  c.Subscribe failed")
	}

	// 启动
	if err := c.Start(); err != nil {
		zap.L().Error("Main  c.Start failed")
	}
	// 9. 退出注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.Error(err)
	}
	zap.L().Info("注销成功")
}
