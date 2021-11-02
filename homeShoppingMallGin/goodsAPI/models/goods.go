package models

// RegisterGoodsInfo 注册商品的信息
type RegisterGoodsInfo struct {
	Name              string   `json:"name" binding:"required"`
	GoodsIntroduction string   `json:"goods_introduction" binding:"required"`
	ClassifyGoods     string   `json:"classify_goods" binding:"required"`
	CreateTime        string   `json:"create_time"`
	GoodsPrice        int64   `json:"goods_price" binding:"required"`
	PromotionPrice    int64   `json:"promotion_price"`
	SalesVolume       int64    `json:"sales_volume"`
	CollectNum        int64    `json:"collect_num"`
	IsNew             bool     `json:"is_new"`
	IsHot             bool     `json:"is_hot"`
	IsShow            bool     `json:"is_show"`
	IsFreightFree     bool     `json:"is_freight_free"`
	Image             []string `json:"image" binding:"required"`
}
