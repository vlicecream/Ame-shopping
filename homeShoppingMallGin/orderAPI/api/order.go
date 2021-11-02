package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"homeShoppingMallGin/orderAPI/consulR"
	"homeShoppingMallGin/orderAPI/models"
	"homeShoppingMallGin/orderAPI/myResponseCode"
	"homeShoppingMallGin/orderAPI/proto"
	"homeShoppingMallGin/orderAPI/settings"
	"homeShoppingMallGin/orderAPI/validators"
	"net/http"
	"strconv"
)

// 调用支付宝支付接口
func alipayInterface(notifyURL, returnURL int64) (str string, err error) {
	// 调用支付宝沙箱环境创建支付订单
	client, err := alipay.New(settings.Conf.AlipayConfig.AppID, settings.Conf.AlipayConfig.PrivateKey, false)
	if err != nil {
		zap.L().Error("api.order.CreateOrder alipay.New failed", zap.Error(err))
		return "", err
	}
	if err := client.LoadAliPayPublicKey(settings.Conf.AlipayConfig.AliPublicKey); err != nil {
		zap.L().Error("api.order.CreateOrder client.LoadAliPayPublicKey failed", zap.Error(err))
		return "", err
	}
	// 初始化结构体
	var p = alipay.TradePagePay{}
	// 把INT转成字符串类型
	orderSnInt := strconv.FormatInt(notifyURL, 10)
	allPriceInt := strconv.FormatInt(returnURL, 10)
	// 进行配置
	p.NotifyURL = settings.Conf.AlipayConfig.NotifyURL
	p.ReturnURL = settings.Conf.AlipayConfig.ReturnURL
	p.Subject = orderSnInt
	p.OutTradeNo = orderSnInt
	p.TotalAmount = allPriceInt
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	// 生成url
	url, err := client.TradePagePay(p)
	if err != nil {
		zap.L().Error("api.order.CreateOrder client.TradeWapPay failed", zap.Error(err))
		return "", err
	}
	return url.String(), nil
}

// 通过jwt拿到UserID
func getUserID(c *gin.Context) int64 {
	// 通过jwt拿到UserID
	userID, ok := c.Get("userID")
	if !ok {
		zap.L().Error("api.order.GetShoppingCar c.Get failed 拿不到userID")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
	}
	// 类型断言转换
	userIDInt, ok := userID.(int64)
	if !ok {
		zap.L().Error("api.order.DeleteShoppingCar userID.(int64) 类型转换失败")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
	}
	return userIDInt
}

// GetShoppingCar 查看购物车
func GetShoppingCar(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 调用查看购物车接口
	rsp, err := consulR.OrderClient.CheckShoppingCar(context.WithValue(context.Background(), "ginContext", c), &proto.UserInfo{UserID: userID})
	if err != nil {
		zap.L().Error("api.order.GetShoppingCar consulR.OrderClient.CheckShoppingCar failed 查看购物车失败")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	if rsp == nil {
		myResponseCode.ResponseSuccess(c, "购物车为空")
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// CreateShoppingCar 创建购物车
func CreateShoppingCar(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 初始化结构体
	var createShoppingCar models.CreateShoppingCar
	// 获取数据
	if err := c.ShouldBindJSON(&createShoppingCar); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用创建购物车接口
	rsp, err := consulR.OrderClient.CreateShoppingCar(context.WithValue(context.Background(), "ginContext", c), &proto.CreateCarRequest{
		UserID:   userID,
		Goods:    createShoppingCar.Goods,
		Nums:     createShoppingCar.GoodsNum,
		Selected: true,
	})
	if err != nil {
		zap.L().Error("api.order.GetShoppingCar c.Get failed 类型转换失败")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// UpdateShoppingCar 更新购物车
func UpdateShoppingCar(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 初始化结构体
	var createShoppingCar models.CreateShoppingCar
	// 获取数据
	if err := c.ShouldBindJSON(&createShoppingCar); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用更新购物车接口
	rsp, err := consulR.OrderClient.UpdateShoppingCar(context.WithValue(context.Background(), "ginContext", c), &proto.CreateCarRequest{
		UserID:   userID,
		Goods:    createShoppingCar.Goods,
		Nums:     createShoppingCar.GoodsNum,
		Selected: createShoppingCar.Selected,
	})
	if err != nil {
		zap.L().Error("api.order.GetShoppingCar c.Get failed 类型转换失败")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// DeleteShoppingCar 删除购物车
func DeleteShoppingCar(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 初始化结构体
	var createShoppingCar models.DeleteShoppingCar
	// 获取数据
	if err := c.ShouldBindJSON(&createShoppingCar); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用删除购物车接口
	if _, err := consulR.OrderClient.DeleteShoppingCar(context.WithValue(context.Background(), "ginContext", c), &proto.DeleteCarRequest{
		UserID: userID,
		Goods:  createShoppingCar.Goods,
	}); err != nil {
		zap.L().Error("api.order.DeleteShoppingCar consulR.OrderClient.DeleteShoppingCar failed")
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeSuccess)
}

// GetOrder 获取用户所有订单
func GetOrder(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 拿取过滤的条件 分页行数 转换成int64
	pnStr := c.DefaultQuery("pn", "0")
	pnInt, err := strconv.Atoi(pnStr)
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}

	// 分页码数 转换成int64
	pSize := c.DefaultQuery("pSize", "10")
	pSizeInt, err := strconv.Atoi(pSize)
	if err != nil {
		zap.L().Error("api.goods GetOrder.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}

	// 调用查询订单接口
	rsp, err := consulR.OrderClient.CheckOrder(context.WithValue(context.Background(), "ginContext", c), &proto.OrderFilterInfo{
		UserID: userID,
		Pn:     int64(pnInt),
		PSize:  int64(pSizeInt),
	})
	if err != nil {
		zap.L().Error("api.goods GetGoodsList.strconv.Atoi failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// GetOrderDetail 获取订单详细信息
func GetOrderDetail(c *gin.Context) {
	// 拿到商品订单号
	var orderGoodsNum models.OrderNum
	if err := c.ShouldBindJSON(&orderGoodsNum); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 调用订单详情接口
	rsp, err := consulR.OrderClient.CheckOrderDetail(context.WithValue(context.Background(), "ginContext", c), &proto.OrderDetailInfoRequest{
		OrderGoodsNum: orderGoodsNum.OrderGoodsNum,
	})
	if err != nil {
		zap.L().Error("api.order.GetGoodsList consulR.OrderClient.CheckOrderDetail failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 因为用户万一在生成订单的时候不支付，所以就得在订单详细信息这里放支付的url
	url, err := alipayInterface(rsp.OrderInfoResponse.OrderSn, rsp.OrderInfoResponse.AllPrice)
	if err != nil {
		zap.L().Error("api.order.GetGoodsList consulR.OrderClient.CheckOrderDetail failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 初始化map，把信息都装进去
	var m = make(map[string]interface{})
	m["data"] = rsp
	m["url"] = url
	myResponseCode.ResponseSuccess(c, rsp)
}

// CreateOrder 新建订单
func CreateOrder(c *gin.Context) {
	// 拿到UserID
	userID := getUserID(c)
	// 初始化结构体
	var orderInfo models.OrderInfo
	// 获取数据
	if err := c.ShouldBindJSON(&orderInfo); err != nil {
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
	rsp, err := consulR.OrderClient.CreateOrder(context.WithValue(context.Background(), "ginContext", c), &proto.CreateOrderInfo{
		UserID:  userID,
		Address: orderInfo.Address,
		Name:    orderInfo.Name,
		Mobile:  orderInfo.Mobile,
		Message: orderInfo.Message,
	})
	if err != nil {
		zap.L().Error("api.order.CreateOrder consulR.OrderClient.CreateOrder failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 调用支付宝沙箱环境创建支付订单
	url, err := alipayInterface(rsp.OrderSn, rsp.AllPrice)
	if err != nil {
		zap.L().Error("api.order.CreateOrder consulR.OrderClient.CreateOrder failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 初始化map，把信息都装进去
	var m = make(map[string]interface{})
	m["data"] = rsp
	m["url"] = url
	myResponseCode.ResponseSuccess(c, m)
}

// UpdateOrder 更新订单
func UpdateOrder(c *gin.Context) {
	// 初始化结构体
	var orderNum models.OrderNum
	// 获取数据
	if err := c.ShouldBindJSON(&orderNum); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}

	if _, err := consulR.OrderClient.UpdateOrderStatus(context.WithValue(context.Background(), "ginContext", c), &proto.OrderInfo{
		OrderSn: orderNum.OrderGoodsNum,
		Status:  orderNum.Status,
	}); err != nil {
		zap.L().Error("api.order.UpdateOrder consulR.OrderClient.UpdateOrderStatus failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeSuccess)
}

// Notify 支付宝订单回调
func Notify(c *gin.Context) {
	// 调用支付宝沙箱环境创建支付订单
	client, err := alipay.New(settings.Conf.AlipayConfig.AppID, settings.Conf.AlipayConfig.PrivateKey, false)
	if err != nil {
		zap.L().Error("api.order.CreateOrder alipay.New failed", zap.Error(err))
		return
	}
	if err := client.LoadAliPayPublicKey(settings.Conf.AlipayConfig.AliPublicKey); err != nil {
		zap.L().Error("api.order.CreateOrder client.LoadAliPayPublicKey failed", zap.Error(err))
		return
	}

	noti, err := client.GetTradeNotification(c.Request)
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	fmt.Println(noti.TradeStatus)
	fmt.Println(noti.OutTradeNo)
	_, err = consulR.OrderClient.UpdateOrderStatus(context.WithValue(context.Background(), "ginContext", c), &proto.OrderInfo{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})

	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	var w http.ResponseWriter
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	alipay.AckNotification(w)
}

// SetInventory 设置库存
func SetInventory(c *gin.Context) {
	// 初始化结构体
	var inventoryInfo models.InventoryInfo
	// 获取数据
	if err := c.ShouldBindJSON(&inventoryInfo); err != nil {
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
	if _, err := consulR.InventoryClient.SetGoodsInventory(context.Background(), &proto.GoodsInfo{
		Goods:        inventoryInfo.GoodsName,
		InventoryNum: inventoryInfo.Nums,
	}); err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeServerBusy)
}

// GetInventory 查看库存
func GetInventory(c *gin.Context) {
	// 初始化结构体
	var inventoryInfo models.GetInventoryInfo
	// 获取数据
	if err := c.ShouldBindJSON(&inventoryInfo); err != nil {
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
	rsp, err := consulR.InventoryClient.GetGoodsInventory(context.Background(), &proto.GoodsInfo{
		Goods:        inventoryInfo.GoodsName,
	})
	if err != nil {
		zap.Error(err)
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}