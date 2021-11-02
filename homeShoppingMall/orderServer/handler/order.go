package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"homeShoppingMall/orderServer/consulRegister"
	"homeShoppingMall/orderServer/dao/mysql"
	"homeShoppingMall/orderServer/models"
	"homeShoppingMall/orderServer/pkg/snowflake"
	"homeShoppingMall/orderServer/proto"
	"homeShoppingMall/orderServer/settings"
	"strconv"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// OrderListener rocketmq发送本地事务消息
// OrderListener 结构体与结构体方法的不懂 可以通过rocketmq.NewTransactionProducer源码查询
type OrderListener struct {
	Code   codes.Code
	Detail error
	AllPrice int64
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	/*
		1. 查询购物车里面的信息生成订单基本信息
		2. 不要太信前端传的总金额，后端要自己查，验证一下，所以要用到商品服务 查看金额
		3. 用到库存服务扣减库存
		4. 存储商品基本信息
		5. 订单不管超时还是支付成功都删除订单
	*/
	// 编写sqlStr语句
	sqlStr := `select * from order_shopping_car where user_id = ? and selected = true;`
	sqlOrderGoodsInfoStr := `insert into order_goods_info(goods_sell_num, goods_order_num, goods, goods_price) values(?, ?, ?, ?) `
	sqlCreateStr := `insert into order_info(id, user_id, order_all_price, address, name, phone, goods_order_num) values(?, ?, ?, ?, ?, ?, ?)`
	// 初始化结构体
	var shoppingCarModels []models.ShoppingCarInfo
	var orderGoodsInfo []*models.OrderGoodsInfo
	var orderInfo *proto.OrderInfoResponse
	var inventoryInfo []*proto.GoodsInfo
	var goodsName []string
	var money int64
	GoodsNumsMap := make(map[string]int32)
	// 把序列化的订单信息反序列化出来
	if err := json.Unmarshal(msg.Body, &orderInfo); err != nil {
		o.Code = codes.Internal
		o.Detail = err
		zap.L().Error("orderServer order ExecuteLocalTransaction json.Unmarshal failed", zap.Error(err))
		return primitive.RollbackMessageState
	}
	// sqlx
	if err := mysql.DB.Select(&shoppingCarModels, sqlStr, orderInfo.UserID); err != nil {
		o.Code = codes.Internal
		o.Detail = err
		err = errors.New("购物车没有选中内容")
		return primitive.RollbackMessageState
	}

	// 循环取出商品名字数据
	for _, values := range shoppingCarModels {
		goodsName = append(goodsName, values.Goods)
		GoodsNumsMap[values.Goods] = values.GoodsNums
	}

	// 调用商品微服务的批量商品信息的接口
	rsp, err := consulRegister.GoodsClient.BatchGetGoods(context.Background(), &proto.BathGoodsNameInfoRequest{Name: goodsName})
	if err != nil {
		o.Code = codes.Internal
		zap.L().Error("计算失败 order.handler.order  CheckShoppingCar.mysql.DB.Select failed", zap.Error(err))
		return primitive.RollbackMessageState
	}

	// 循环取出数据并算出总金额,顺手保存需要的商品信息
	for _, values := range rsp.GoodsInfo {
		money = money + values.GoodsPrice*int64(GoodsNumsMap[values.Name])
		orderGoodsInfo = append(orderGoodsInfo, &models.OrderGoodsInfo{
			GoodsSellNum: int64(GoodsNumsMap[values.Name]),
			GoodsPrice:   values.GoodsPrice,
			Goods:        values.Name,
		})
		inventoryInfo = append(inventoryInfo, &proto.GoodsInfo{
			Goods:        values.Name,
			InventoryNum: int64(GoodsNumsMap[values.Name]),
		})
	}
	// 调用库存服务来扣减库存
	orderSn := strconv.Itoa(int(orderInfo.OrderSn))
	if _, err := consulRegister.InventoryClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: inventoryInfo, OrderGoodsNum: orderSn}); err != nil {
		o.Code = codes.Internal
		o.Detail = err
		zap.L().Error("库存扣减失败 order.handler.order  CheckShoppingCar.mysql.DB.Select failed", zap.Error(err))
		return primitive.RollbackMessageState
	}

	/*检测rocketmq事务回查后的库存归还的测试代码*/
	//o.Code = codes.Internal
	//o.Detail = errors.New("错误")
	//return primitive.CommitMessageState
	/*检测结束*/

	// 创建订单信息
	orderInfo.AllPrice = money
	o.AllPrice = money
	// 启动本地mysql事务来回滚
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		o.Code = codes.Internal
		o.Detail = err
		zap.L().Error("begin trans failed", zap.Error(err))
		return primitive.RollbackMessageState
	}
	// sqlx保存 订单信息
	if _, err = tx.Exec(sqlCreateStr, orderInfo.Id, orderInfo.UserID, orderInfo.AllPrice, orderInfo.Address, orderInfo.Name,
		orderInfo.Mobile, orderInfo.OrderSn); err != nil {
		if err := tx.Rollback(); err != nil {
			zap.L().Error("order tx.Rollback failed", zap.Error(err))
		}
		err = errors.New("订单创建失败")
		o.Code = codes.Internal
		o.Detail = err
		zap.L().Error("订单创建失败 order.handler.order  CheckShoppingCar.mysql.DB.Select failed", zap.Error(err))
		return primitive.CommitMessageState
	}
	// sqlx保存 订单商品信息
	for _, values := range orderGoodsInfo {
		if _, err = tx.Exec(sqlOrderGoodsInfoStr, values.GoodsSellNum, orderInfo.OrderSn, values.Goods,
			values.GoodsPrice); err != nil {
			if err := tx.Rollback(); err != nil {
				zap.L().Error("order tx.Rollback failed", zap.Error(err))
			}
			o.Code = codes.Internal
			o.Detail = err
			err = errors.New("订单商品信息创建失败")
			zap.L().Error("order.handler.order  CheckShoppingCar.mysql.DB.Select failed", zap.Error(err))
			return primitive.CommitMessageState
		}
	}

	// 发送延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{fmt.Sprintf("%s:%d", settings.Conf.RocketMQ.Host, settings.Conf.RocketMQ.Port)}))
	if err != nil {
		panic("生成producer失败")
	}

	//不要在一个进程中使用多个producer， 但是不要随便调用shutdown因为会影响其他的producer
	if err = p.Start(); err != nil {
		zap.L().Error("启动producer失败", zap.Error(err))
	}

	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(5)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.L().Error("发送延时消息失败", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			zap.L().Error("order tx.Rollback failed", zap.Error(err))
		}
		o.Code = codes.Internal
		o.Detail = errors.New("发送延时消息失败")
		return primitive.CommitMessageState
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("order tx.Commit failed", zap.Error(err))
	}
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

// CheckLocalTransaction 消息回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo *proto.OrderInfoResponse
	// 把序列化的订单信息反序列化出来
	if err := json.Unmarshal(msg.Body, &orderInfo); err != nil {
		o.Detail = err
		zap.L().Error("orderServer order CheckLocalTransaction json.Unmarshal failed", zap.Error(err))
		return primitive.CommitMessageState
	}

	/*
		通过查询订单号来确定这个确实是成功执行了
		但尼不能就确定说明执行完就是库存已经扣减了
	*/
	// 编写查询sql语句
	sqlStr := `select count(goods_order_num) from order_info where goods_order_num = ?`
	var count int
	// sqlx
	if err := mysql.DB.Get(&count, sqlStr, orderInfo.OrderSn); err != nil {
		o.Detail = err
		zap.L().Error("orderServer order CheckLocalTransaction mysql.DB.Get failed", zap.Error(err))
		return primitive.CommitMessageState
	}

	if count == 0 {
		err := errors.New("订单不存在，回查消息不成功")
		o.Detail = err
		zap.L().Error("orderServer order CheckLocalTransaction mysql.DB.Get failed", zap.Error(err))
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}

// CheckShoppingCar 查看购物车
func (o *OrderServer) CheckShoppingCar(ctx context.Context, in *proto.UserInfo) (*proto.ShoppingCarListResponse, error) {
	// 编写sql语句
	sqlStr := `select * from order_shopping_car where user_id = ? and selected = true`
	// 初始化结构体
	var shoppingCarModels []models.ShoppingCarInfo
	var shoppingCarInfo []*proto.ShoppingCarInfo
	shoppingCarList := proto.ShoppingCarListResponse{}
	// sqlx
	if err := mysql.DB.Select(&shoppingCarModels, sqlStr, in.UserID); err != nil {
		zap.L().Error("order.handler.order  CheckShoppingCar.mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 循环取出数据
	for key, values := range shoppingCarModels {
		shoppingCarInfo = append(shoppingCarInfo, &proto.ShoppingCarInfo{
			UserID:   values.UserID,
			GoodsNum: values.GoodsNums,
			Goods:    values.Goods,
			Selected: values.Selected,
		})
		shoppingCarList.Total = int32(key) + 1
	}
	shoppingCarList.ShoppingCarInfo = shoppingCarInfo
	return &shoppingCarList, nil
}

// CreateShoppingCar 创建购物车
func (o *OrderServer) CreateShoppingCar(ctx context.Context, in *proto.CreateCarRequest) (*proto.ShoppingCarInfo, error) {
	// 编写sql语句
	sqlStr := `insert into order_shopping_car(user_id, goods, goods_nums, selected) values(?, ?, ?, ?)` // 购物车不存在这个商品的时候就新建数据
	sqlSelectStr := `select * from order_shopping_car where user_id = ? and goods = ?`                  // 查询是否有购物车
	sqlUpdateStr := `update order_shopping_car set goods_nums = ? where user_id = ? and goods = ?`      // 如果购物车有这个商品就去数目+1
	// 判断这个商品是否存在
	var shoppingCarModels models.ShoppingCarInfo
	// sqlx查询
	if err := mysql.DB.Get(&shoppingCarModels, sqlSelectStr, in.UserID, in.Goods); err != nil {
		shoppingCarModels.UserID = in.UserID
		shoppingCarModels.GoodsNums = in.Nums
		shoppingCarModels.Goods = in.Goods
		shoppingCarModels.Selected = false
		// 查不到就去新建购物车数据
		if _, err := mysql.DB.Exec(sqlStr, shoppingCarModels.UserID, shoppingCarModels.Goods,
			shoppingCarModels.GoodsNums, shoppingCarModels.Selected); err != nil {
			zap.L().Error("order.handler.order  CreateShoppingCar.mysql.DB.Exec failed", zap.Error(err))
			return nil, err
		}
	} else {
		// 查到了已有数据就去更新数量
		shoppingCarModels.GoodsNums += in.Nums
		if _, err := mysql.DB.Exec(sqlUpdateStr, shoppingCarModels.GoodsNums, shoppingCarModels.UserID, shoppingCarModels.Goods); err != nil {
			zap.L().Error("order.handler.order  CreateShoppingCar.mysql.DB.Exec failed", zap.Error(err))
			return nil, err
		}
	}

	shoppingCarInfo := &proto.ShoppingCarInfo{
		UserID:   shoppingCarModels.UserID,
		GoodsNum: shoppingCarModels.GoodsNums,
		Goods:    shoppingCarModels.Goods,
		Selected: false,
	}
	return shoppingCarInfo, nil
}

// UpdateShoppingCar 更新购物车
func (o *OrderServer) UpdateShoppingCar(ctx context.Context, in *proto.CreateCarRequest) (*proto.OrderEmpty, error) {
	// 编写sql语句
	sqlStr := `update order_shopping_car set goods_nums = ?, selected = ? where user_id = ? and goods = ?`
	sqlSelectStr := `select * from order_shopping_car where user_id = ? and goods = ?`
	// 判断这个商品是否存在
	var shoppingCarModels models.ShoppingCarInfo
	// sqlx查询
	if err := mysql.DB.Get(&shoppingCarModels, sqlSelectStr, in.UserID, in.Goods); err != nil {
		err = errors.New("您要更新的商品没有加入购物车")
		return nil, err
	} else {
		shoppingCarModels.GoodsNums = in.Nums
		shoppingCarModels.Selected = in.Selected
	}

	if _, err := mysql.DB.Exec(sqlStr, shoppingCarModels.GoodsNums, shoppingCarModels.Selected, shoppingCarModels.UserID, shoppingCarModels.Goods); err != nil {
		zap.L().Error("order.handler.order  CreateShoppingCar.mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}

	return &proto.OrderEmpty{}, nil
}

// DeleteShoppingCar 删除购物车
func (o *OrderServer) DeleteShoppingCar(ctx context.Context, in *proto.DeleteCarRequest) (*proto.OrderEmpty, error) {
	for _, values := range in.Goods {
		// 编写sql语句
		sqlStr := `delete from order_shopping_car where goods = ? and user_id = ?`
		// sqlx
		if _, err := mysql.DB.Exec(sqlStr, values, in.UserID); err != nil {
			zap.L().Error("order.handler.order  DeleteShoppingCar.mysql.DB.Exec failed", zap.Error(err))
			return nil, err
		}
	}
	return &proto.OrderEmpty{}, nil
}

// CheckOrder 获取用户所有订单
func (o *OrderServer) CheckOrder(ctx context.Context, in *proto.OrderFilterInfo) (*proto.OrderListResponse, error) {
	// 编写sql语句
	sqlStr := `select * from order_info where user_id = ? limit ?, ?`
	// 初始化结构体
	var orderModel []models.OrderInfo
	var orderInfo []*proto.OrderInfoResponse
	shoppingCarList := proto.OrderListResponse{}
	fmt.Println(in.UserID)
	// sqlx
	if err := mysql.DB.Select(&orderModel, sqlStr, in.UserID, in.Pn, in.PSize); err != nil {
		zap.L().Error("order.handler.order  DeleteShoppingCar.mysql.DB.Exec failed", zap.Error(err))
		err = errors.New("没有任何订单")
		return nil, err
	}
	fmt.Println(orderModel)
	// 循环取出数据
	for key, values := range orderModel {
		orderInfo = append(orderInfo, &proto.OrderInfoResponse{
			Id:       values.Id,
			UserID:   values.UserID,
			Address:  values.Address,
			Name:     values.Name,
			Mobile:   values.Phone,
			Message:  values.Message,
			OrderSn:  values.GoodsOrderNum,
			PayType:  values.PayType,
			AllPrice: values.OrderAllPrice,
			Status:   values.Status,
		})
		shoppingCarList.Total = int64(key) + 1
	}
	shoppingCarList.OrderInfoResponse = orderInfo
	return &shoppingCarList, nil
}

// CreateOrder 新建订单
func (o *OrderServer) CreateOrder(ctx context.Context, in *proto.CreateOrderInfo) (*proto.OrderInfoResponse, error) {
	// rocketmq 发送事务消息
	orderListener := OrderListener{}
	// 创建一个新producer
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", settings.Conf.RocketMQ.Host, settings.Conf.RocketMQ.Port)}),
	)
	if err != nil {
		zap.L().Error("orderServer order CreateOrder rocketmq.NewTransactionProducer failed", zap.Error(err))
		return nil, err
	}
	if err := p.Start(); err != nil {
		zap.L().Error("orderServer order CreateOrder p.Start failed", zap.Error(err))
		return nil, err
	}

	// 创建订单信息
	orderInfo := &proto.OrderInfoResponse{
		UserID:  in.UserID,
		Address: in.Address,
		Name:    in.Name,
		Mobile:  in.Mobile,
		Message: in.Message,
		OrderSn: snowflake.GenID(),
	}


	// 把订单信息json序列化然后发送出去
	orderInfoByte, err := json.Marshal(orderInfo)
	if err != nil {
		zap.L().Error("orderServer order CreateOrder json.Marshal failed", zap.Error(err))
		return nil, err
	}

	// 发送事务消息
	_, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order", orderInfoByte))
	if err != nil {
		zap.L().Error("orderServer order CreateOrder p.SendMessageInTransaction failed", zap.Error(err))
		return nil, err
	}
	// 把总金额补上
	orderInfo.AllPrice = orderListener.AllPrice
	// 判断事务消息是否成功进行完成
	if orderListener.Code != codes.OK {
		return nil, orderListener.Detail
	}
	return orderInfo, nil
}

// CheckOrderDetail 获取用户订单详细信息
func (o *OrderServer) CheckOrderDetail(ctx context.Context, in *proto.OrderDetailInfoRequest) (*proto.OrderDetailResponse, error) {
	// 编写sql语句
	sqlStr := `select * from order_info where goods_order_num = ?`           // 查询订单信息
	sqlGoodsStr := `select * from order_goods_info where goods_order_num =?` // 查询订单商品信息
	// 初始化结构体
	var orderModel models.OrderInfo
	var orderList proto.OrderDetailResponse
	var orderGoodsModels []models.OrderGoodsInfo
	// sqlx
	if err := mysql.DB.Get(&orderModel, sqlStr, in.OrderGoodsNum); err != nil {
		zap.L().Error("order.handler.order  DeleteShoppingCar.mysql.DB.Exec failed", zap.Error(err))
		err = errors.New("订单不存在")
		return nil, err
	}
	// 把 models.OrderInfo 转变成 *proto.OrderDetailResponse
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = orderModel.Id
	orderInfo.OrderSn = orderModel.GoodsOrderNum
	orderInfo.Name = orderModel.Name
	orderInfo.Message = orderModel.Message
	orderInfo.Status = orderModel.Status
	orderInfo.PayType = orderModel.PayType
	orderInfo.Mobile = orderModel.Phone
	orderInfo.Address = orderModel.Address
	orderInfo.AllPrice = orderModel.OrderAllPrice
	orderInfo.UserID = orderModel.UserID

	orderList.OrderInfoResponse = &orderInfo
	// 查询订单商品信息
	if err := mysql.DB.Select(&orderGoodsModels, sqlGoodsStr, in.OrderGoodsNum); err != nil {
		zap.L().Error("order.handler.order  DeleteShoppingCar.mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 循环取出信息
	for _, values := range orderGoodsModels {
		orderList.Goods = append(orderList.Goods, &proto.OrderItemResponse{
			OrdersID:  values.GoodsOrderNum,
			GoodsName: values.Goods,
			Price:     values.GoodsPrice,
			Nums:      values.GoodsSellNum,
		})
	}
	return &orderList, nil
}

// UpdateOrderStatus 修改订单状态
func (o *OrderServer) UpdateOrderStatus(ctx context.Context, in *proto.OrderInfo) (*proto.OrderEmpty, error) {
	// 编写sql语句
	sqlUpdateStr := `update order_info set status = ? where goods_order_num = ?`
	sqlSelectStr := `select count(goods_order_num) from order_info where goods_order_num = ?`

	// 初始化
	var count int32
	if err := mysql.DB.Get(&count, sqlSelectStr, in.OrderSn); err != nil {
		zap.L().Error("order.handler.order  UpdateOrderSelected.mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count == 0 {
		err := errors.New("订单不存在")
		zap.L().Error("订单不存在 order.handler.order", zap.Error(err))
		return nil, err
	}

	// sqlx 修改订单状态
	if _, err := mysql.DB.Exec(sqlUpdateStr, in.Status, in.OrderSn); err != nil {
		zap.L().Error("修改订单失败 order.handler.order  UpdateOrderSelected.mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.OrderEmpty{}, nil
}

// OrderTimeout rocketmq order超时
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo models.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		//查询订单的支付状态，如果已支付什么都不做，如果未支付，归还库存
		var order models.OrderInfo
		// 编写查询sql语句
		sqlStr := `select * from order_info where goods_order_num = ?`
		// sqlx
		if err := mysql.DB.Get(&order, sqlStr, orderInfo.GoodsOrderNum); err != nil {
			zap.L().Error("数据库获取数据失败", zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		// 如果是已经支付 什么事都不做
		if order.Status == "TRADE_SUCCESS" {
			return consumer.ConsumeSuccess, nil
		}
		// 如果没支付，归还库存
		if order.Status != "TRADE_SUCCESS" {
			// 启动本地mysql事务来回滚
			tx, err := mysql.DB.Beginx() // 开启事务
			if err != nil {
				zap.L().Error("数据库启动事务失败", zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
			//归还库存，我们可以模仿order中发送一个消息到 order_back中去
			//修改订单的状态为已支付,编写sql语句
			sqlStatusStr := `update order_info set status = ? where goods_order_num = ?`
			order.Status = "TRADE_CLOSED"
			fmt.Println(order.Status, orderInfo.GoodsOrderNum)
			// sqlx
			if _, err = tx.Exec(sqlStatusStr, order.Status, orderInfo.GoodsOrderNum); err != nil {
				if err = tx.Rollback(); err != nil {
					zap.L().Error("p.Start() failed", zap.Error(err))
					return consumer.ConsumeRetryLater, err
				}
				zap.L().Error("数据库查询失败", zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{fmt.Sprintf("%s:%d", settings.Conf.RocketMQ.Host, settings.Conf.RocketMQ.Port)}))
			if err != nil {
				if err = tx.Rollback(); err != nil {
					zap.L().Error("p.Start() failed", zap.Error(err))
					return consumer.ConsumeRetryLater, err
				}
				zap.L().Error("rocketmq.NewProducer failed", zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}

			if err = p.Start(); err != nil {
				if err = tx.Rollback(); err != nil {
					zap.L().Error("p.Start() failed", zap.Error(err))
					return consumer.ConsumeRetryLater, err
				}
				zap.L().Error("p.Start() failed", zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}

			if _, err = p.SendSync(context.Background(), primitive.NewMessage("order", msgs[i].Body)); err != nil {
				if err = tx.Rollback(); err != nil {
					zap.L().Error("p.Start() failed", zap.Error(err))
					return consumer.ConsumeRetryLater, err
				}
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			if err = tx.Commit(); err != nil {
				zap.L().Error("tx.Commit() failed", zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
		}
	}
	return consumer.ConsumeSuccess, nil
}
