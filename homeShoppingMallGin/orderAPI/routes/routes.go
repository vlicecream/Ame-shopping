package routes

import (
	"github.com/gin-gonic/gin"
	"homeShoppingMallGin/orderAPI/api"
	"homeShoppingMallGin/orderAPI/logger"
	"homeShoppingMallGin/orderAPI/middlewares"
)

func Init() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("api/v1").Use(middlewares.JaegerTrace()).Use(middlewares.Cors())
	// 购物车
	v1.GET("/getShopCar", middlewares.JWTAuthMiddleware(), api.GetShoppingCar)          // 查看用户所有购物车
	v1.POST("/createShopCar", middlewares.JWTAuthMiddleware(), api.CreateShoppingCar)   // 创建购物车
	v1.PUT("/updateShopCar", middlewares.JWTAuthMiddleware(), api.UpdateShoppingCar)    // 更新购物车
	v1.DELETE("/deleteShopCar", middlewares.JWTAuthMiddleware(), api.DeleteShoppingCar) // 删除购物车
	// 订单
	v1.GET("/getOrder", middlewares.JWTAuthMiddleware(), api.GetOrder)             // 获取用户所有订单
	v1.GET("/getOrderDetail", middlewares.JWTAuthMiddleware(), api.GetOrderDetail) // 获取订单详细信息
	v1.POST("/createOrder", middlewares.JWTAuthMiddleware(), api.CreateOrder)      // 创建订单
	v1.PUT("/updateOrder", middlewares.JWTAuthMiddleware(), api.UpdateOrder)       // 更新订单
	// 库存
	v1.POST("/setInventory", middlewares.JWTAuthMiddleware(), middlewares.Admin(), api.SetInventory) // 设置库存
	v1.GET("/getInventory", middlewares.JWTAuthMiddleware(), middlewares.Admin(), api.GetInventory)
	// 支付宝回调
	v1.POST("/alipay/notify", api.Notify) // 支付宝回调地址
	// 健康检查
	r.GET("/health")
	return r
}
