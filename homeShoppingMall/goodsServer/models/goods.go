package models

// GoodsInfo 商品的信息
type GoodsInfo struct {
	Id                int64    `json:"id" db:"id"`
	QuantityStock     int64    `json:"quantity_stock" db:"quantity_stock"`
	SalesVolume       int64    `json:"sales_volume" db:"sales_volume"`
	CollectNum        int64    `json:"collect_num" db:"collect_num"`
	GoodsPrice        int64   `json:"goods_price" db:"goods_price"`
	PromotionPrice    int64   `json:"promotion_price" db:"promotion_price"`
	Name              string   `json:"name" db:"name"`
	GoodsIntroduction string   `json:"goods_introduction" db:"goods_introduction"`
	ClassifyName      string   `json:"classify_name" db:"classify_name"`
	CreateTime        string   `json:"create_time" db:"create_time"`
	ClassifyGoods     string   `json:"classify_goods" db:"classify_goods"`
	IsShow            bool     `json:"is_show" db:"is_show"`
	IsNew             bool     `json:"is_new" db:"is_new"`
	IsFreightFree     bool     `json:"is_freight_free" db:"is_freight_free"`
	IsHot             bool     `json:"is_hot" db:"is_hot"`
	Image             []string `json:"image" db:"image"`
}

// ClassifyGoodsInfo 商品分类的信息
type ClassifyGoodsInfo struct {
	ID   int64  `json:"id" db:"id"`
	PID  int64  `json:"pid" db:"pid"`
	Name string `json:"name" db:"name"`
}

// BannerInfo 轮播图信息
type BannerInfo struct {
	ID            int64  `db:"id"`
	Level         int64  `db:"level"`
	ImageUrl      string `db:"image_url"`
	ImageGoodsUrl string `db:"image_goods_url"`
}

// BannerCreateInfo 轮播图新增信息
type BannerCreateInfo struct {
	ImageUrl      string `json:"image_url" db:"image_url"`
	ImageGoodsUrl string `json:"image_goods_url" db:"image_goods_url"`
	Level         int64  `json:"level" db:"level"`
}

// GoodsImageInfo 商品图片信息
type GoodsImageInfo struct {
	ImageUrl string `json:"image_url" db:"image_url"`
	//Goods string `json:"goods" db:"goods_name"`
}