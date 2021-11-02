package models

// RegisterClassifyGoods 新增商品的信息
type RegisterClassifyGoods struct {
	Name  string `json:"name" binding:"required"`
	PName string `json:"p_name"`
}

// UpdateClassifyGoods 更新商品提交信息
type UpdateClassifyGoods struct {
	OldName string `json:"old_name" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
	PName   string `json:"p_name" binding:"required"`
}
