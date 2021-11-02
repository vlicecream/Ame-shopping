package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"homeShoppingMallGin/userAPI/consulR"
	"homeShoppingMallGin/userAPI/utils"
	"homeShoppingMallGin/userAPI/validators"

	"homeShoppingMallGin/userAPI/logger"
	"homeShoppingMallGin/userAPI/routes"
	"homeShoppingMallGin/userAPI/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置文件
	if err := settings.Init(); err != nil {
		zap.L().Debug("Main settings.Init failed")
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		zap.L().Debug("Main logger.Init failed")
		return
	}
	defer zap.L().Sync()

	// 3. 初始化grpc&consul服务发现
	consulR.Init()

	// 4. 初始化validators翻译器
	if err := validators.InitTrans("zh"); err != nil {
		zap.L().Error("init controller.Init failed, err:", zap.Error(err))
	}

	// 5. 注册自定义翻译器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("mobile", validators.ValidateMobile); err != nil {
			zap.L().Error("main v.RegisterValidation failed", zap.Error(err))
		}
		_ = v.RegisterTranslation("mobile", validators.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "手机格式不正确", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 3. 配置路由
	r := routes.Init()

	// 4. 启动服务(优雅关机)
	port, err := utils.DynamicPort()
	if err != nil {
		zap.L().Error("main utils.DynamicPort failed", zap.Error(err))
	}

	// 如果是测试环境就换配置端口
	if settings.Conf.Mode == "dev" {
		port = settings.Conf.Port
	}

	// 注册商品的http服务
	reClient := consulR.NewRegistryClient(settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port)
	// 拿到随机ID(uuid)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	if err := reClient.Register(settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.GinName,
		serviceID, port, settings.Conf.ConsulConfig.Tag); err != nil {
		zap.L().Debug("Main consulGin.NewRegistryClient failed")
	}


	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			zap.L().Error("main r.Run failed", zap.Error(err))
		}
	}()

	/*开始服务注销*/
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := reClient.DeRegister(serviceID); err != nil {
		zap.L().Error("main reClient.DeRegister failed", zap.Error(err))
	} else {
		zap.L().Info("服务注销成功")
	}
	/*服务注销结束*/

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit = make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
