package routes

import (
	"github.com/gin-gonic/gin"
	"homeShoppingMallGin/userAPI/api"
	"homeShoppingMallGin/userAPI/logger"
	"homeShoppingMallGin/userAPI/middlewares"
)

func Init() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1").Use(middlewares.Cors())
	// 获取所有用户列表
	v1.GET("/user/getUserList", middlewares.JWTAuthMiddleware(), middlewares.Admin() ,api.GetAllUserList)
	// 利用手机号返回用户详细信息
	v1.GET("/user/getUserInfo", middlewares.JWTAuthMiddleware(), api.GetUserInfo)
	// 用户登录
	v1.POST("/user/login", api.UserLogin)
	// 发送短信验证码接口
	v1.POST("/user/sendCode", api.SendAuthCode)
	// 用户注册
	v1.POST("/user/register", api.UserRegister)
	// 用户更新手机号及密码
	v1.PUT("/user/updateMP", middlewares.JWTAuthMiddleware(), api.UserUpdateMobilePassword)
	// 用户更新昵称
	v1.PUT("/user/updateNickName", middlewares.JWTAuthMiddleware(), api.UserUpdateNickName)
	// 管理员更新权限
	v1.PUT("/user/adminRole", middlewares.JWTAuthMiddleware(), middlewares.Admin(), api.AdminRole)
	// HTTP健康检查
	r.GET("/health")
	return r
}
