package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"homeShoppingMall/goodsServer/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// Init 初始化mysql数据库
func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		settings.Conf.MysqlConfig.UserName,
		settings.Conf.MysqlConfig.Password,
		settings.Conf.MysqlConfig.Host,
		settings.Conf.MysqlConfig.Port,
		settings.Conf.MysqlConfig.DataBase)
	// 也可以使用MustConnect连接不成功就panic
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("connect DB failed, err", zap.Error(err))
		return
	}
	DB.SetMaxOpenConns(settings.Conf.MysqlConfig.SetMaxOpenConns)
	DB.SetMaxIdleConns(settings.Conf.MysqlConfig.SetMaxIdleConns)

	return
}

// Close 定义一个关闭mysql的向外暴露的接口
func Close() (err error) {
	err = DB.Close()
	return
}

