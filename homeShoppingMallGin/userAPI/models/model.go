package models

// UserResponse 用户返回响应信息
type UserResponse struct {
	UserId   int64  `db:"user_id"`
	Mobile   string `db:"mobile"`
	Nickname string `db:"nickname"`
	Gender   string `db:"gender"`
	Role     string `db:"role"`
}

// UserLoginInfo 用户登录提交信息
type UserLoginInfo struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	AuthCode string `json:"auth_code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRLCodeInfo 用户注册&登录获取验证码
type UserRLCodeInfo struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
}

// UserRegisterInfo UserLoginInfo 用户登录提交信息
type UserRegisterInfo struct {
	Mobile     string `json:"mobile" binding:"required,mobile"`
	AuthCode   string `json:"auth_code" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// UserUpdateInfo 用户更新提交信息
type UserUpdateInfo struct {
	Mobile      string `json:"mobile" binding:"required"`
	AuthCode    string `json:"auth_code" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	RePassword  string `json:"re_password" binding:"required"`
}

// UserUpdateNickName 用户修改
type UserUpdateNickName struct {
	NickName string `json:"nick_name" binding:"required"`
}

// AdminUserUpdateRole 修改管理权限，只有管理员才能做
type AdminUserUpdateRole struct {
	Mobile string `json:"mobile" binding:"required"`
	Role   string `json:"role" binding:"required"`
}
