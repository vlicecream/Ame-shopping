package classifyGoods

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"homeShoppingMallGin/goodsAPI/consulR"
	"homeShoppingMallGin/goodsAPI/models"
	"homeShoppingMallGin/goodsAPI/myResponseCode"
	"homeShoppingMallGin/goodsAPI/proto"
	"homeShoppingMallGin/goodsAPI/validators"
)

// GetOneClassify 获取一级分类并限流
func GetOneClassify(c *gin.Context) {
	// 进行限流操作
	e, b := sentinel.Entry("goodsApi", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeTooManyRequests)
		return
	}
	// 直接调用接口
	rsp, err := consulR.GoodsSrvClient.GetClassifyInfo(context.Background(), &proto.GoodsEmpty{})
	if err != nil {
		zap.L().Error("api.classifyGoods.GetOneClassify consulR.GoodsSrvClient.GetClassifyInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 取消限流区域
	e.Exit()
	// 返回数据
	myResponseCode.ResponseSuccess(c, rsp.Info)
}

// GetChildClassify 获取二级分类并限流
func GetChildClassify(c *gin.Context) {
	// 进行限流操作
	e, b := sentinel.Entry("goodsApi", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeTooManyRequests)
		return
	}
	// 获取数据
	pName := c.Query("name")
	rsp, err := consulR.GoodsSrvClient.GetChildClassifyInfo(context.Background(), &proto.ClassifyChildInfoRequest{
		PName: pName,
	})
	if err != nil {
		zap.L().Error("api.classifyGoods.GetChildClassify consulR.GoodsSrvClient.GetChildClassifyInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 取消限流区域
	e.Exit()
	myResponseCode.ResponseSuccess(c, rsp.ListInfo)
}

// RegisterClassify 新增商品分类
func RegisterClassify(c *gin.Context) {
	// 初始化结构体
	var classifyInfo models.RegisterClassifyGoods
	// 获取数据
	if err := c.ShouldBindJSON(&classifyInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用新增接口
	rsp, err := consulR.GoodsSrvClient.CreateClassifyInfo(context.Background(), &proto.ClassifyCreateInfoRequest{
		Name:  classifyInfo.Name,
		PName: classifyInfo.PName,
	})
	if err != nil {
		zap.L().Error("api.classifyGoods.RegisterClassify consulR.GoodsSrvClient.CreateClassifyInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// UpdateClassify 更新商品分类
func UpdateClassify(c *gin.Context) {
	// 初始化结构体
	var classifyInfo models.UpdateClassifyGoods
	// 获取数据
	if err := c.ShouldBindJSON(&classifyInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用新增接口
	rsp, err := consulR.GoodsSrvClient.UpdateClassifyInfo(context.Background(), &proto.ClassifyUpdateInfoRequest{
		OldName: classifyInfo.OldName,
		NewName: classifyInfo.NewName,
		PName:   classifyInfo.PName,
	})
	if err != nil {
		zap.L().Error("api.classifyGoods.UpdateClassify consulR.GoodsSrvClient.UpdateClassifyInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// DeleteClassify 删除商品分类
func DeleteClassify(c *gin.Context) {
	// 获取数据
	deleteName := c.Query("name")
	// 调用接口
	if _, err := consulR.GoodsSrvClient.DeleteClassifyInfo(context.Background(), &proto.ClassifyDeleteInfoRequest{Name: deleteName}); err != nil {
		zap.L().Error("api.classifyGoods.DeleteClassify consulR.GoodsSrvClient.DeleteClassifyInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		fmt.Println(2)
		return
	}
	myResponseCode.ResponseSuccess(c, "Delete ok")
}
