package models

import "time"

// ShoppingCarInfo 购物车信息
type ShoppingCarInfo struct {
	Id         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	GoodsNums  int32     `json:"goods_nums" db:"goods_nums"`
	Goods      string    `json:"goods" db:"goods"`
	Selected   bool      `json:"selected" db:"selected"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	Id             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	OrderAllPrice  int64     `json:"order_all_price" db:"order_all_price"`
	GoodsOrderNum  int64     `json:"orderSn" db:"goods_order_num"`
	PayType        string    `json:"pay_type" db:"pay_type"`
	Status         string    `json:"status" db:"status"`
	AlipayOrderNum string    `json:"alipay_order_num" db:"alipay_order_num"`
	Mobile         string    `json:"mobile" db:"mobile"`
	Address        string    `json:"address" db:"address"`
	Name           string    `json:"name" db:"name"`
	Phone          string    `json:"phone" db:"phone"`
	Message        string    `json:"message" db:"message"`
	PayTime        string    `json:"pay_time" db:"pay_time"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
}

// OrderGoodsInfo 订单商品信息
type OrderGoodsInfo struct {
	Id            int64     `json:"id" db:"id"`
	GoodsSellNum  int64     `json:"goods_sell_num" db:"goods_sell_num"`
	GoodsPrice    int64     `json:"goods_price" db:"goods_price"`
	GoodsOrderNum string    `json:"goods_order_num" db:"goods_order_num"`
	Goods         string    `json:"goods" db:"goods"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}

