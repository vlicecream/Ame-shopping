package utils

import (
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"homeShoppingMall/userServer/dao/mysql"
	"homeShoppingMall/userServer/handler"
	"homeShoppingMall/userServer/pkg/snowflake"
)

type User struct {
	UserID   int64  `db:"user_id"`
	Mobile   string `db:"mobile"`
	Password string `db:"password"`
}

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.UserID, u.Mobile, u.Password}, nil
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (user_id, mobile, password) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := mysql.DB.Exec(query, args...)
	return err
}

func Init() {
	u1 := User{UserID: snowflake.GenID(), Mobile: "1809777821", Password: handler.MD5Password("123456")}
	u2 := User{UserID: snowflake.GenID(), Mobile: "1809777822", Password: handler.MD5Password("123456")}
	u3 := User{UserID: snowflake.GenID(), Mobile: "1809777823", Password: handler.MD5Password("123456")}
	fmt.Println(u1.UserID, u2.UserID, u3.UserID)
	users := []interface{}{u1, u2, u3}
	err := BatchInsertUsers2(users)
	if err != nil {
		fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
	}
}
