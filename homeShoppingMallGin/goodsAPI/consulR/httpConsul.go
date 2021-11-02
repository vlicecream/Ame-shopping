package consulR

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

/*以下代码利用了Go的面向接口概念*/

type Registry struct {
	Port int
	Host string
}

type RegistryClient interface {
	Register(address, name, id string, port int, tags []string) (err error)
	DeRegister(serviceID string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Port: port,
		Host: host,
	}
}

func (r *Registry) Register(address, name, id string, port int, tags []string) (err error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry) DeRegister(serviceID string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	_ = client.Agent().ServiceDeregister(serviceID)
	return nil
}