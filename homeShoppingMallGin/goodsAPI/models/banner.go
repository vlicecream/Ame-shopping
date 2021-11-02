package models

// BannerInfo 新增轮播图信息
type BannerInfo struct {
	ImageUrl string `json:"image_url" binding:"required"`
	ImageGoodsUrl string `json:"image_goods_url"  binding:"required"`
	Level int64 `json:"level"  binding:"required"`
}