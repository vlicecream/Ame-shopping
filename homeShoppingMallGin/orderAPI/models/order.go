package models

// CreateShoppingCar 用户传递创建购物车信息
type CreateShoppingCar struct {
	Goods    string `json:"goods" binding:"required"`
	GoodsNum int32  `json:"goods_num" binding:"required"`
	Selected bool   `json:"selected"`
}

// DeleteShoppingCar 用户传递创建购物车信息
type DeleteShoppingCar struct {
	Goods []string `json:"goods" binding:"required"`
}

type OrderNum struct {
	OrderGoodsNum string `json:"order_goods_num" binding:"required"`
	Status        string `json:"status"`
}

// OrderInfo 创建订单信息
type OrderInfo struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Mobile  string `json:"mobile" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// InventoryInfo 设置库存信息
type InventoryInfo struct {
	GoodsName string `json:"goods_name" binding:"required"`
	Nums      int64  `json:"nums" binding:"required"`
}

// GetInventoryInfo 拿取库存信息
type GetInventoryInfo struct {
	GoodsName string `json:"goods_name" binding:"required"`
}