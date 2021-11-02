package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"homeShoppingMall/orderServer/settings"
)

// Rdb 声明一个全局的rdb变量
var Rdb *redis.Client

// InitClient 初始化连接
func InitClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Conf.RedisConfig.Host, settings.Conf.RedisConfig.Port),
		Password: settings.Conf.RedisConfig.Password, // no password set
		DB:       settings.Conf.RedisConfig.DataBase,  // use default DB
	})

	_, err = Rdb.Ping().Result()
	return err
}