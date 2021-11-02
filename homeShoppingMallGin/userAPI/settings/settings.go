package settings

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conf = new(AppConfig) // 把这些配置保存到一个变量

type AppConfig struct {
	Mode         string        `mapstructure:"mode" json:"mode"`
	Host         string        `mapstructure:"host" json:"host"`
	Port         int           `mapstructure:"port" json:"port"`
	LogConfig    *LogConfig    `mapstructure:"logger" json:"logger"`
	ConsulConfig *ConsulConfig `mapstructure:"consul" json:"consul"`
	NacosConfig  *NacosConfig  `mapstructure:"nacos" json:"nacosConfig"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	FileName   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups"`
}

type ConsulConfig struct {
	Port       int      `mapstructure:"port" json:"port"`
	Host       string   `mapstructure:"host" json:"host"`
	Name       string   `mapstructure:"name" json:"name"`
	GinName    string   `mapstructure:"GinName" json:"GinName"`
	Tag        []string `mapstructure:"tag" json:"tag"`
}

type NacosConfig struct {
	NamespaceID string `mapstructure:"NamespaceID"`
	IP          string `mapstructure:"IP"`
	DataID      string `mapstructure:"DataID"`
	User        string `maostructure:"user"`
	Password    string `maostructure:"password"`
	Group       string `mapstructure:"Group"`
	Port        uint64 `mapstructure:"Port"`
}

// Init 初始化配置文件
func Init() (err error) {
	// 读取配置文件
	viper.SetConfigName("config")     // 配置文件名称(无扩展名)
	viper.AddConfigPath("./settings") // 查找配置文件所在的路径
	err = viper.ReadInConfig()        // 查找并读取配置文件
	if err != nil {                   // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置文件反序列化到配置变量Conf中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%s", err)
	}
	// 监控配置文件的改变
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", in.Name)
	})

	// 为nacos初始化
	NacosConfigInfo()
	return
}

func NacosConfigInfo() {
	// 一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: Conf.NacosConfig.IP,
			Port:   Conf.NacosConfig.Port,
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         Conf.NacosConfig.NamespaceID, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "userServer/tmp/nacos/log",
		CacheDir:            "userServer/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	iClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.L().Error("settings clients.CreateConfigClient failed", zap.Error(err))
	}
	// 获取配置信息
	content, err := iClient.GetConfig(vo.ConfigParam{
		DataId: Conf.NacosConfig.DataID,
		Group:  Conf.NacosConfig.Group,
	})
	// 反序列化信息
	if err = json.Unmarshal([]byte(content), &Conf); err != nil {
		zap.L().Error("settings json.Unmarshal failed", zap.Error(err))
	}
}
