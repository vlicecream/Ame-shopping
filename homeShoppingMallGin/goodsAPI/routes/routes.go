package routes

import (
	"github.com/gin-gonic/gin"
	"homeShoppingMallGin/goodsAPI/api/banner"
	"homeShoppingMallGin/goodsAPI/api/classifyGoods"
	"homeShoppingMallGin/goodsAPI/api/goods"
	"homeShoppingMallGin/goodsAPI/logger"
	"homeShoppingMallGin/goodsAPI/middlewares"
)

func Init() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1").Use(middlewares.JaegerTrace()).Use(middlewares.Cors()).Use(middlewares.JWTAuthMiddleware())
	// 商品
	v1.GET("/goods/getGoods", goods.GetGoodsList)                        // 过滤获得商品信息
	v1.POST("/goods/register", middlewares.Admin(), goods.RegisterGoods) // 注册商品信息
	v1.PUT("/goods/update", middlewares.Admin(), goods.UpdateGoods)      // 更新商品信息
	v1.DELETE("/goods/delete", middlewares.Admin(), goods.DeleteGoods)   // 删除商品信息
	// 商品分类
	v1.GET("/classify/getOneClassify", classifyGoods.GetOneClassify)                   // 获取所有一级分类信息
	v1.GET("/classify/getChildClassify", classifyGoods.GetChildClassify)               // 获取所有的二级分类
	v1.POST("/classify/register", middlewares.Admin(), classifyGoods.RegisterClassify) // 注册分类信息
	v1.PUT("/classify/update", middlewares.Admin(), classifyGoods.UpdateClassify)      // 更新分类信息
	v1.DELETE("/classify/delete", middlewares.Admin(), classifyGoods.DeleteClassify)   // 删除分类信息
	// 轮播图分类
	v1.GET("/banner/getBanner", middlewares.Admin(), banner.GetAllBanner) // 获取所有轮播图信息
	v1.POST("/banner/register", middlewares.Admin(), banner.CreateBanner) // 新增轮播图
	v1.PUT("/banner/update", middlewares.Admin(), banner.UpdateBanner)    // 更新轮播图
	v1.DELETE("/banner/delete", middlewares.Admin(), banner.DeleteBanner) // 删除轮播图
	// 健康检查
	r.GET("/health")
	return r
}
