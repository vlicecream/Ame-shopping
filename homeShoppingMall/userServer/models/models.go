package models

import "time"

// AllUserInfo 用户要显示的信息
type AllUserInfo struct {
	ID         int32     `db:"id"`
	UserID     int64     `db:"user_id"`
	NickName   string    `db:"nick_name"`
	Mobile     string    `db:"mobile"`
	Password   string    `db:"password"`
	Gender     string    `db:"gender"`
	Role       string    `db:"role"`
	CreateTime time.Time `db:"create_time"`
	DeleteTime time.Time `db:"delete_time"`
	IsDelete   bool `db:"is_delete"`
}
