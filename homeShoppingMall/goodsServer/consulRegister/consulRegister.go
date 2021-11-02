package consulRegister

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"homeShoppingMall/goodsServer/settings"
)

func ConsulInit() (*api.Client, string){
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
