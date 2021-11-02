package goods

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
	"strconv"
)

// GetGoodsList 过滤商品信息并分页并限流
func GetGoodsList(c *gin.Context) {
	// 初始化
	request := &proto.ClassifyGoodsInfoRequest{}
	// 拿取过滤的条件 分页行数 转换成int64
	pnStr := c.DefaultQuery("pn", "0")
	pnInt, err := strconv.Atoi(pnStr)
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	request.Pn = int64(pnInt)

	// 分页码数 转换成int64
	pSize := c.DefaultQuery("pSize", "0")
	pSizeInt, err := strconv.Atoi(pSize)
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	request.PSize = int64(pSizeInt)

	// 最低价格 转换成int64
	priceMin := c.DefaultQuery("pMin", "0")
	priceMinInt, err := strconv.Atoi(priceMin)
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	request.PriceMin = int64(priceMinInt)

	// 最高价格 转换成int64
	priceMax := c.DefaultQuery("pMax", "0")
	priceMaxInt, err := strconv.Atoi(priceMax)
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	request.PriceMax = int64(priceMaxInt)

	// 是否热销
	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	// 是否新品
	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	// 是否免运费
	isFree := c.DefaultQuery("iff", "0")
	if isFree == "1" {
		request.IsFreightFree = true
	}
	// 关键字
	keyword := c.DefaultQuery("kw", "")
	request.Name = keyword
	// 分类查询
	classifyId := c.DefaultQuery("cID", "0")
	request.TopClassify = classifyId

	// 对这个接口进行限流
	e, b := sentinel.Entry("goodsApi", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeTooManyRequests)
		return
	}
	// 对这个接口进行预热与冷启动
	a, p := sentinel.Entry("goodsApi-limiting", sentinel.WithTrafficType(base.Inbound))
	if p != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeTooManyRequests)
		return
	}
	// 调用接口
	rsp, err := consulR.GoodsSrvClient.GetClassifyGoods(context.WithValue(context.Background(), "ginContext", c), request)
	if err != nil {
		zap.L().Error("api.goods consulR.GoodsSrvClient.GetClassifyGoods failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 结束限流
	e.Exit()
	a.Exit()
	myResponseCode.ResponseSuccess(c, rsp)
}

// RegisterGoods 注册商品
func RegisterGoods(c *gin.Context) {
	// 初始化结构体
	var goodsInfo models.RegisterGoodsInfo
	if err := c.ShouldBindJSON(&goodsInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 直接调用SERVER接口
	rsp, err := consulR.GoodsSrvClient.CreateGoodsInfo(context.WithValue(context.Background(), "ginContext", c), &proto.GoodsCreateInfoRequest{
		Name:              goodsInfo.Name,
		GoodsPrice:        goodsInfo.GoodsPrice,
		PromotionPrice:    goodsInfo.PromotionPrice,
		GoodsIntroduction: goodsInfo.GoodsIntroduction,
		CreateTime:        goodsInfo.CreateTime,
		ClassifyGoods:     goodsInfo.ClassifyGoods,
		SalesVolume:       goodsInfo.SalesVolume,
		CollectNum:        goodsInfo.CollectNum,
		IsNew:             goodsInfo.IsNew,
		IsHot:             goodsInfo.IsHot,
		IsShow:            goodsInfo.IsShow,
		IsFreightFree:     goodsInfo.IsFreightFree,
		Image:             goodsInfo.Image,
	})
	if err != nil {
		zap.L().Error("api.goods.RegisterGoods consulR.GoodsSrvClient.CreateGoodsInfo failed", zap.Error(err))
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, err)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// UpdateGoods 更新商品
func UpdateGoods(c *gin.Context) {
	// 初始化结构体
	var goodsInfo models.RegisterGoodsInfo
	if err := c.ShouldBindJSON(&goodsInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 直接调用SERVER接口
	rsp, err := consulR.GoodsSrvClient.UpdateGoodsInfo(context.WithValue(context.Background(), "ginContext", c), &proto.GoodsCreateInfoRequest{
		Name:              goodsInfo.Name,
		GoodsPrice:        goodsInfo.GoodsPrice,
		PromotionPrice:    goodsInfo.PromotionPrice,
		GoodsIntroduction: goodsInfo.GoodsIntroduction,
		CreateTime:        goodsInfo.CreateTime,
		ClassifyGoods:     goodsInfo.ClassifyGoods,
		SalesVolume:       goodsInfo.SalesVolume,
		CollectNum:        goodsInfo.CollectNum,
		IsNew:             goodsInfo.IsNew,
		IsHot:             goodsInfo.IsHot,
		IsShow:            goodsInfo.IsShow,
		IsFreightFree:     goodsInfo.IsFreightFree,
		Image:             goodsInfo.Image,
	})
	if err != nil {
		zap.L().Error("api.goods.UpdateGoods consulR.GoodsSrvClient.UpdateGoodsInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// DeleteGoods 删除商品
func DeleteGoods(c *gin.Context) {
	// 获取数据
	name := c.Query("name")
	fmt.Println(name)
	// 调用proto
	if _, err := consulR.GoodsSrvClient.DeleteGoodsInfo(context.WithValue(context.Background(), "ginContext", c), &proto.GoodsDeleteInfoRequest{Name: name}); err != nil {
		zap.L().Error("api.goods.UpdateGoods consulR.GoodsSrvClient.UpdateGoodsInfo failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, "Delete ok")
}
