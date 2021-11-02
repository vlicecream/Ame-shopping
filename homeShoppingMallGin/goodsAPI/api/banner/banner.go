package banner

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"homeShoppingMallGin/goodsAPI/consulR"
	"homeShoppingMallGin/goodsAPI/models"
	"homeShoppingMallGin/goodsAPI/myResponseCode"
	"homeShoppingMallGin/goodsAPI/proto"
	"homeShoppingMallGin/goodsAPI/validators"
	"strconv"
)

// GetAllBanner 查看所有轮播图
func GetAllBanner(c *gin.Context) {
	// 调用接口
	rsp, err := consulR.GoodsSrvClient.GetBannerInfo(context.Background(), &proto.GoodsEmpty{})
	if err != nil {
		zap.L().Error("api.banner.GetAllBanner consulR.GoodsSrvClient.GetBannerInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// CreateBanner 新增轮播图
func CreateBanner(c *gin.Context) {
	// 初始化结构体
	var bannerInfo models.BannerInfo
	// 获取数据
	if err := c.ShouldBindJSON(&bannerInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用接口
	rsp, err := consulR.GoodsSrvClient.CreateBannerInfo(context.Background(), &proto.BannerCreateInfoRequest{
		ImageUrl:      bannerInfo.ImageUrl,
		ImageGoodsUrl: bannerInfo.ImageGoodsUrl,
		Level:         bannerInfo.Level,
	})
	if err != nil {
		zap.L().Error("api.banner.CreateBanner consulR.GoodsSrvClient.CreateBannerInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// UpdateBanner 更新轮播图
func UpdateBanner(c *gin.Context) {
	//  初始化结构体
	var bannerInfo models.BannerInfo
	// 获取数据
	if err := c.ShouldBindJSON(&bannerInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用接口
	if _, err := consulR.GoodsSrvClient.UpdateBannerInfo(context.Background(), &proto.BannerCreateInfoRequest{
		ImageUrl:      bannerInfo.ImageUrl,
		ImageGoodsUrl: bannerInfo.ImageGoodsUrl,
		Level:         bannerInfo.Level,
	}); err != nil {
		zap.L().Error("api.banner.UpdateBanner consulR.GoodsSrvClient.UpdateBannerInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, "Update ok")
}

// DeleteBanner 删除轮播图
func DeleteBanner(c *gin.Context) {
	// 获取数据
	deleteLevel := c.Query("level")
	// 把ID转成int64
	level, err := strconv.Atoi(deleteLevel)
	if err != nil {
		zap.L().Error("api.banner.DeleteBanner strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 调用接口
	if _, err := consulR.GoodsSrvClient.DeleteBannerInfo(context.Background(), &proto.BannerDeleteInfoRequest{Level: int64(level)}); err != nil{
		fmt.Println("goodsClient.DeleteBannerInfo failed")
		return
	}
	myResponseCode.ResponseSuccess(c, "Delete ok")
}