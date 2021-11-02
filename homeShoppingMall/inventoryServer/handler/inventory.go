package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"homeShoppingMall/inventoryServer/dao/mysql"
	"homeShoppingMall/inventoryServer/models"
	"homeShoppingMall/inventoryServer/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// SetGoodsInventory 设置商品库存
func (i *InventoryServer) SetGoodsInventory(ctx context.Context, in *proto.GoodsInfo) (*proto.InventoryEmpty, error) {
	// 编写sql语句
	sqlStr := `insert into inventory(goods, inventory_num) values(?, ?)`
	// sqlx
	if _, err := mysql.DB.Exec(sqlStr, in.Goods, in.InventoryNum); err != nil {
		zap.L().Error("inventory.SetGoodsInventory mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.InventoryEmpty{}, nil
}

// GetGoodsInventory 拿取商品库存
func (i *InventoryServer) GetGoodsInventory(ctx context.Context, in *proto.GoodsInfo) (*proto.GoodsInfo, error) {
	// 初始化结构体
	var inventoryInfo models.InventoryInfo
	// 编写sql语句
	sqlStr := `select * from inventory where goods = ?`
	// sqlx
	if err := mysql.DB.Get(&inventoryInfo, sqlStr, in.InventoryNum); err != nil {
		zap.L().Error("inventory.GetGoodsInventory mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	goodsInfo := &proto.GoodsInfo{
		Goods:        inventoryInfo.Goods,
		InventoryNum: inventoryInfo.InventoryNum,
	}
	return goodsInfo, nil
}

// Sell 商品库存预扣减
func (i *InventoryServer) Sell(ctx context.Context, in *proto.SellInfo) (*proto.InventoryEmpty, error) {
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		zap.L().Error("begin trans failed", zap.Error(err))
		return nil, err
	}
	// 编写sql语句
	sqlSellInfoStr := `insert into inventory_shell(orders_goods_num, status) values(?, ?)`
	//sqlx
	if _, err = mysql.DB.Exec(sqlSellInfoStr, in.OrderGoodsNum, "NO"); err != nil {
		if err = tx.Rollback(); err != nil {
			zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
		}
		zap.L().Error("inventory handler Shell mysql.DB.Exec", zap.Error(err))
		err = errors.New("创建数据失败")
		return nil, err
	}
	// 编写sql语句
	sqlSelectStr := `select * from inventory where goods = ?`
	sqlUpdateStr := `update inventory set inventory_num = ? where goods = ?`
	sqlCreateStr := `insert into inventory_shell_detail(goods, num, orders_goods_num) values(?, ?, ?)`
	// 循环取出数据
	for _, values := range in.GoodsInfo {
		// 初始化
		var inventoryInfo models.InventoryInfo
		// sqlx
		if err = tx.Get(&inventoryInfo, sqlSelectStr, values.Goods); err != nil {
			err = errors.New("没有这个商品")
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			return nil, err
		}
		// 判断库存够不够
		if inventoryInfo.InventoryNum < values.InventoryNum {
			err = errors.New("库存不够,请联系商家补充")
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			return nil, err
		}
		// 减去销售量
		inventoryInfo.InventoryNum -= values.InventoryNum
		// sqlx保存
		_, err = tx.Exec(sqlUpdateStr, inventoryInfo.InventoryNum, values.Goods)
		if err != nil {
			zap.L().Error("inventory.Sell mysql.DB.Exec failed", zap.Error(err))
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			return nil, err
		}
		_, err = tx.Exec(sqlCreateStr, values.Goods, values.InventoryNum, in.OrderGoodsNum)
		if err != nil {
			zap.L().Error("inventory.Sell mysql.DB.Exec failed", zap.Error(err))
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("inventory tx.Commit failed", zap.Error(err))
	}
	return &proto.InventoryEmpty{}, nil
}

// ReBack 商品库存归还
func (i *InventoryServer) ReBack(ctx context.Context, in *proto.SellInfo) (*proto.InventoryEmpty, error) {
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			zap.L().Error("inventoryServer.handler ReBack rollback", zap.Error(err))
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			zap.L().Error("inventoryServer.handler ReBack rollback", zap.Error(err))
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			zap.L().Warn("inventoryServer.handler ReBack commit")
		}
	}()
	// 循环取出数据
	for _, values := range in.GoodsInfo {
		// 编写sql语句查询是否存在
		sqlSelectStr := `select * from inventory where goods = ?`
		sqlUpdateStr := `update inventory set inventory_num = ? where goods = ?`
		// 初始化
		var inventoryInfo models.InventoryInfo
		// sqlx
		if err = tx.Get(&inventoryInfo, sqlSelectStr, values.Goods); err != nil {
			err = errors.New("没有这个商品")
			return nil, err
		}
		// 加去销售量
		inventoryInfo.InventoryNum += values.InventoryNum
		// sqlx保存
		rs, err := tx.Exec(sqlUpdateStr, inventoryInfo.InventoryNum, values.Goods)
		if err != nil {
			zap.L().Error("inventory.Sell mysql.DB.Exec failed", zap.Error(err))
			return nil, err
		}
		n, err := rs.RowsAffected()
		if err != nil {
			return nil, err
		}
		if n != 1 {
			err = errors.New("exec sqlStr1 failed")
			return nil, err
		}
	}
	return &proto.InventoryEmpty{}, nil
}

// AutoBack rocketmq
func AutoBack(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// 拿取订单编号
	type orderInfo struct {
		OrderSn int64
	}
	// 开启事务
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		if err = tx.Rollback(); err != nil {
			zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
		}
		return consumer.ConsumeSuccess, err
	}

	for i := range msg {
		/*
			既然是归还库存,应该知道每件商品归还多少的库存
			这里一定要注意幂等性，不能因为网络波动等原因，导致库存多还等问题
			所以要新建表，详细记录这些数据
		*/
		var orderDetail orderInfo
		if err = json.Unmarshal(msg[i].Body, &orderDetail); err != nil {
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			zap.L().Error("inventory.AutoBack json.Unmarshal failed", zap.Error(err))
			err = errors.New("反序列化失败")
			return consumer.ConsumeSuccess, err
		}
		// 把库存加回去，将status改为OK， 这个要在事务下进行
		var shellInfo models.SellInfo
		var ShellDetail []models.SellDetail
		// 编写sql语句
		sqlSelectStr := `select * from inventory_shell where orders_goods_num = ? and status = 'NO'` // 查询是否有没有这个数据
		sqlDetailStr := `select goods, num from inventory_shell_detail where orders_goods_num = ? `
		sqlUpdateShellStr := `update inventory_shell set status = 'ok' where orders_goods_num = ?`
		sqlUpdateInventoryStr := `update inventory set inventory_num = inventory_num+? where goods = ?`
		fmt.Println(orderDetail.OrderSn)
		// sqlx
		if err = tx.Get(&shellInfo, sqlSelectStr, orderDetail.OrderSn); err != nil {
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			zap.L().Error("inventory.AutoBack tx.Get failed", zap.Error(err))
			return consumer.ConsumeSuccess, err
		}
		if err = tx.Select(&ShellDetail, sqlDetailStr, orderDetail.OrderSn); err != nil {
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			zap.L().Error("inventory.AutoBack tx.Get failed", zap.Error(err))
			return consumer.ConsumeSuccess, err
		}
		// 循环取出商品名和商品购买数量
		for _, values := range ShellDetail {
			fmt.Println(values.Nums, values.Goods)
			if _, err = tx.Exec(sqlUpdateInventoryStr, values.Nums, values.Goods); err != nil {
				if err = tx.Rollback(); err != nil {
					zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
				}
				zap.L().Error("inventory.AutoBack mysql.DB.Exec failed", zap.Error(err))
				err = errors.New("更新inventory表失败")
				return consumer.ConsumeRetryLater, err
			}
		}

		if _, err = tx.Exec(sqlUpdateShellStr, orderDetail.OrderSn); err != nil {
			if err = tx.Rollback(); err != nil {
				zap.L().Error("inventory tx.Rollback failed", zap.Error(err))
			}
			zap.L().Error("inventory.AutoBack mysql.DB.Exec failed", zap.Error(err))
			err = errors.New("更新inventory_shell表失败")
			return consumer.ConsumeRetryLater, err
		}
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("inventory tx.Commit failed", zap.Error(err))
	}
	return consumer.ConsumeSuccess, nil
}
