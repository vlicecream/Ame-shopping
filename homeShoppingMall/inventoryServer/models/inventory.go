package models

// InventoryInfo 库存信息
type InventoryInfo struct {
	Goods        string `json:"goods" db:"goods"`
	InventoryNum int64  `json:"inventory_num" db:"inventory_num"`
	ID           int64  `json:"id" db:"id"`
}

// SellInfo 归还库存信息
type SellInfo struct {
	ID            int    `json:"id" db:"id"`
	OrderGoodsNum string `json:"order_goods_num" db:"orders_goods_num"`
	Status        string `json:"status" db:"status"`
}

type SellDetail struct {
	Goods         string `json:"goods" db:"goods"`
	Nums          int64  `json:"nums" db:"num"`
}
