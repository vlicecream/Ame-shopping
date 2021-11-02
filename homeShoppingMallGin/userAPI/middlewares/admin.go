package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"homeShoppingMallGin/goodsAPI/myResponseCode"
)

func Admin() func(c *gin.Context) {
	return func(c *gin.Context) {
		role, ok := c.Get("userRole")
		if !ok {
			zap.L().Error("middleware Admin c.Get failed", zap.Error(errors.New("jwt拿取用户角色失败")))
			myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
			c.Abort()
			return
		}
		if role == "user" {
			myResponseCode.ResponseError(c, myResponseCode.CodePermissionDenied)
			c.Abort()
			return
		}
		if role == "admin" {
			c.Next()
			return
		}
	}
}
