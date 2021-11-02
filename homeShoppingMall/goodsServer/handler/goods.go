package handler

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"homeShoppingMall/goodsServer/dao/mysql"
	"homeShoppingMall/goodsServer/models"
	"homeShoppingMall/goodsServer/proto"
)

type GoodsServerServer struct {
	proto.UnimplementedGoodsServerServer
}

// GetPID 根据一级品牌名拿到他的ID
func GetPID(name string) int64 {
	// 编写查询语句
	sqlStr := `select id from classify_goods where name = ?`
	// 添加顶级分类的时候直接传PID为0
	if name == "" {
		return 0
	}
	// 初始化结构体
	var id int64
	// sqlx
	if err := mysql.DB.Get(&id, sqlStr, name); err != nil {
		zap.L().Error("handler.goods.GetPID failed", zap.Error(err))
		return 0
	}
	return id
}

// GetClassifyGoods 根据过滤条件搜索商品
func (g *GoodsServerServer) GetClassifyGoods(ctx context.Context, in *proto.ClassifyGoodsInfoRequest) (*proto.GoodsListResponse, error) {
	// 创建一个局部mysql db
	//GoodsDB := mysql.DB
	// 初始化结构体
	var sqlStr string
	var goodsInfo []models.GoodsInfo
	//var goodsInfoList []*proto.GoodsInfoResponse
	var goodsList proto.GoodsListResponse
	var goodsImageInfo []*models.GoodsImageInfo
	var goodsImageInfoResponse []*proto.GoodsImageResponse
	// 根据过滤条件的不同创建查询
	if in.IsHot == true {
		sqlStr = fmt.Sprintf(`select * from goods where is_hot = true and is_show = true limit %d,%d`, in.Pn, in.PSize)
	} else if in.IsNew == true {
		sqlStr = fmt.Sprintf(`select * from goods where is_new = true and is_show = true limit %d, %d`, in.Pn, in.PSize)
	} else if in.IsFreightFree == true {
		sqlStr = fmt.Sprintf(`select * from goods where is_freight_free = true and is_show = true limit %d, %d`, in.Pn, in.PSize)
	} else if in.Name != "" {
		sqlStr = fmt.Sprintf(`select * from goods where name like '%%%s%%' and is_show = true limit %d, %d`, in.Name, in.Pn, in.PSize)
	} else if in.PriceMin > 0 {
		sqlStr = fmt.Sprintf(`select * from goods where goods_price > %d and is_show = true limit %d, %d`, in.PriceMin, in.Pn, in.PSize)
	} else if in.PriceMax > 0 {
		sqlStr = fmt.Sprintf(`select * from goods where goods_price > %d and is_show = true limit %d, %d`, in.PriceMax, in.Pn, in.PSize)
	} else if in.TopClassify == "1" {
		PID := GetPID(in.Name)
		sqlStr = fmt.Sprintf(`select * from goods where is_show = true and classify_name in 
		(select name from classify_goods where pid = %d)`, PID)
	} else if in.TopClassify == "2" {
		sqlStr = fmt.Sprintf(`select * from goods where classify_name = '%s' and is_show = true`, in.Name)
	} else {
		sqlStr = `select * from goods where is_show = true limit 10`
	}
	// sqlx
	if err := mysql.DB.Select(&goodsInfo, sqlStr); err != nil {
		zap.L().Warn("handler.good.GetClassifyGoods mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 取出商品url图片的sql语句
	sqlImageStr := `select image_url from goods_image where goods_name = ?`
	// 循环取出数据
	for total, values := range goodsInfo {
		// sqlx把商品图片寻找出来
		if err := mysql.DB.Select(&goodsImageInfo, sqlImageStr, values.Name); err != nil {
			zap.L().Warn("handler.good.GetClassifyGoods mysql.DB.Select failed", zap.Error(err))
			return nil, err
		}
		// 循环取出商品图片数据并保存
		for _, imageValues := range goodsImageInfo {
			goodsImageInfoResponse = append(goodsImageInfoResponse, &proto.GoodsImageResponse{
				ImageUrl: imageValues.ImageUrl,
			})
		}
		// 保存商品信息
		goodsList.GoodsInfo = append(goodsList.GoodsInfo, &proto.GoodsInfoResponse{
			Id:                values.Id,
			SalesVolume:       values.SalesVolume,
			CollectNum:        values.CollectNum,
			Name:              values.Name,
			GoodsPrice:        values.GoodsPrice,
			PromotionPrice:    values.PromotionPrice,
			GoodsIntroduction: values.GoodsIntroduction,
			CreateTime:        values.CreateTime,
			IsShow:            values.IsShow,
			IsNew:             values.IsNew,
			IsFreightFree:     values.IsFreightFree,
			IsHot:             values.IsHot,
			ClassifyGoods:     values.ClassifyGoods,
			GoodsImageInfo:    goodsImageInfoResponse,
		})
		// 把切片清空
		goodsImageInfoResponse = goodsImageInfoResponse[0:0]
		goodsImageInfo = goodsImageInfo[0:0]
		// 存储搜索条数
		goodsList.Total = int64(total + 1)
	}
	return &goodsList, nil
}

// BatchGetGoods 批量查询商品
func (g *GoodsServerServer) BatchGetGoods(ctx context.Context, in *proto.BathGoodsNameInfoRequest) (*proto.GoodsListResponse, error) {
	// 初始化结构体
	var goodsInfo []models.GoodsInfo
	//var goodsInfoList []*proto.GoodsInfoResponse
	var goodsList proto.GoodsListResponse
	var goodsUrlInfo []*models.GoodsImageInfo
	var goodsImageInfoResponse []*proto.GoodsImageResponse
	// 商品url  sql语句
	sqlUrlStr := `select image_url from goods_image where goods_name = ?`
	for _, value := range in.Name {
		// 编写sql语句, 因为这是专门给购物车使用，可以只获得购物车展示信息
		sqlStr := `select * from goods where name = ? and is_show = true`
		// sqlx
		if err := mysql.DB.Select(&goodsInfo, sqlStr, value); err != nil {
			zap.L().Warn("handler.good.BatchGetGoods mysql.DB.Select failed", zap.Error(err))
			return nil, err
		}
	}
	// 循环取出数据
	for total, values := range goodsInfo {
		// 取出商品图片信息
		if err := mysql.DB.Select(&goodsUrlInfo, sqlUrlStr, values.Name); err != nil {
			zap.L().Warn("handler.good.BatchGetGoods mysql.DB.Select failed", zap.Error(err))
			return nil, err
		}
		for _, urlValue := range goodsUrlInfo {
			goodsImageInfoResponse = append(goodsImageInfoResponse, &proto.GoodsImageResponse{
				ImageUrl: urlValue.ImageUrl,
			})
		}
		goodsList.GoodsInfo = append(goodsList.GoodsInfo, &proto.GoodsInfoResponse{
			Id:                values.Id,
			SalesVolume:       values.SalesVolume,
			CollectNum:        values.CollectNum,
			Name:              values.Name,
			GoodsPrice:        values.GoodsPrice,
			PromotionPrice:    values.PromotionPrice,
			GoodsIntroduction: values.GoodsIntroduction,
			CreateTime:        values.CreateTime,
			IsShow:            values.IsShow,
			IsNew:             values.IsNew,
			IsFreightFree:     values.IsFreightFree,
			IsHot:             values.IsHot,
			ClassifyGoods:     values.ClassifyGoods,
			GoodsImageInfo:    goodsImageInfoResponse,
		})
		// 把切片清空
		goodsImageInfoResponse = goodsImageInfoResponse[0:0]
		goodsUrlInfo = goodsUrlInfo[0:0]
		// 存储搜索条数
		goodsList.Total = int64(total + 1)
	}
	return &goodsList, nil
}

// CreateGoodsInfo 新增商品
func (g *GoodsServerServer) CreateGoodsInfo(ctx context.Context, in *proto.GoodsCreateInfoRequest) (*proto.GoodsInfoResponse, error) {
	// 编写sql语句
	sqlStr := `insert into goods(name, goods_price, promotion_price, classify_name,
                  goods_introduction, is_show, is_new, is_freight_free,
                  is_hot) values(?, ?, ?, ?, ?, ?, ?, ?, ?)`   // 增加
	sqlSelectStr := `select count(name) from goods where name = ?` // 查询重复
	sqlImageStr := `insert into goods_image(image_url, goods_name) values(?, ?)`
	// 初始化结构体
	var count int

	// 首先查询是否已经存在这个分类
	if err := mysql.DB.Get(&count, sqlSelectStr, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateGoodsInfo mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("商品已存在")
	}
	// 不存在就创建
	if _, err := mysql.DB.Exec(sqlStr, in.Name, in.GoodsPrice, in.PromotionPrice, in.ClassifyGoods,
		in.GoodsIntroduction, in.IsShow, in.IsNew, in.IsFreightFree, in.IsHot); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	// 转换类型
	GoodsInfoResponse := &proto.GoodsInfoResponse{
		Id:                0,
		SalesVolume:       0,
		CollectNum:        0,
		Name:              in.Name,
		GoodsPrice:        in.GoodsPrice,
		PromotionPrice:    in.PromotionPrice,
		GoodsIntroduction: in.GoodsIntroduction,
		CreateTime:        "",
		IsShow:            in.IsShow,
		IsNew:             in.IsNew,
		IsFreightFree:     in.IsFreightFree,
		IsHot:             in.IsHot,
		ClassifyGoods:     in.ClassifyGoods,
	}
	// 把图片保存至图片表
	for _, values := range in.Image {
		if _, err := mysql.DB.Exec(sqlImageStr, values, in.Name); err != nil {
			zap.L().Warn("handler.good.CreateClassifyInfo.sqlImageStr mysql.DB.Exec failed", zap.Error(err))
			return nil, err
		}
	}

	return GoodsInfoResponse, nil
}

// UpdateGoodsInfo 更新商品信息
func (g *GoodsServerServer) UpdateGoodsInfo(ctx context.Context, in *proto.GoodsCreateInfoRequest) (*proto.GoodsEmpty, error) {
	// 编写sql语句
	sqlStr := `update goods set goods_price=?, promotion_price=?, classify_name=?,
                  goods_introduction=?, is_new=?, is_freight_free=?,
                  is_hot=? where name = ?`                         // 增加
	sqlSelectStr := `select count(name) from goods where name = ?` // 查询重复
	// 初始化结构体
	var count int

	// 首先查询是否已经存在这个分类
	if err := mysql.DB.Get(&count, sqlSelectStr, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateGoodsInfo mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("商品不存在")
	}
	// 更新存在就创建
	if _, err := mysql.DB.Exec(sqlStr, in.GoodsPrice, in.PromotionPrice, in.ClassifyGoods,
		in.GoodsIntroduction, in.IsNew, in.IsFreightFree, in.IsHot, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}

// DeleteGoodsInfo 删除商品信息,其实就是更改is_show字段
func (g *GoodsServerServer) DeleteGoodsInfo(ctx context.Context, in *proto.GoodsDeleteInfoRequest) (*proto.GoodsEmpty, error) {
	// 编写sql语句
	sqlStr := `update goods set is_show = false where name = ?`    // 增加
	sqlSelectStr := `select count(name) from goods where name = ?` // 查询重复
	// 初始化结构体
	var count int
	// 首先查询是否已经存在这个分类
	if err := mysql.DB.Get(&count, sqlSelectStr, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateGoodsInfo mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("商品不存在")
	}
	// 更新存在就创建
	if _, err := mysql.DB.Exec(sqlStr, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}

// GetClassifyInfo 查询一级品牌分类
func (g *GoodsServerServer) GetClassifyInfo(context.Context, *proto.GoodsEmpty) (*proto.ClassifyListResponse, error) {
	// 编写sql查询语句
	sqlStr := `select * from classify_goods where pid = 0`
	// 初始化结构体
	var classifyGoods []models.ClassifyGoodsInfo
	var classifyGoodsInfo []*proto.ClassifyInfoResponse
	var classifyGoodsList proto.ClassifyListResponse

	// sqlx
	if err := mysql.DB.Select(&classifyGoods, sqlStr); err != nil {
		zap.L().Warn("handler.good.GetClassifyInfo mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 循环取出数据然后转换类型
	for _, value := range classifyGoods {
		classifyGoodsInfo = append(classifyGoodsInfo, &proto.ClassifyInfoResponse{
			Id:   value.ID,
			Pid:  value.PID,
			Name: value.Name,
		})
	}
	classifyGoodsList.Info = classifyGoodsInfo
	return &classifyGoodsList, nil
}

// GetChildClassifyInfo 查询子品牌分类
func (g *GoodsServerServer) GetChildClassifyInfo(ctx context.Context, in *proto.ClassifyChildInfoRequest) (*proto.ChildClassifyListResponse, error) {
	// 编写sqlStr语句
	sqlStr := `select name from classify_goods where pid = ?`
	// 通过父类名来搜索子类的PID
	cPID := GetPID(in.PName)
	// 初始化结构体
	var classifyGoods []models.ClassifyGoodsInfo
	var classifyGoodsInfo []*proto.ClassifyChildInfoRequest
	var classifyGoodsList proto.ChildClassifyListResponse
	// sqlx
	if err := mysql.DB.Select(&classifyGoods, sqlStr, cPID); err != nil {
		zap.L().Warn("handler.good.GetChildClassifyInfo mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 循环取出数据并放入[]*proto.ClassifyChildInfoResponse
	for _, values := range classifyGoods {
		classifyGoodsInfo = append(classifyGoodsInfo, &proto.ClassifyChildInfoRequest{
			PName: "",
			Name:  values.Name,
		})
	}
	classifyGoodsList.Info = in
	classifyGoodsList.ListInfo = classifyGoodsInfo
	return &classifyGoodsList, nil
}

// CreateClassifyInfo 创建品牌分类
func (g *GoodsServerServer) CreateClassifyInfo(ctx context.Context, in *proto.ClassifyCreateInfoRequest) (*proto.ClassifyInfoResponse, error) {
	// 拿到上一级品牌分类的ID
	PID := GetPID(in.PName)
	// 编写sql语句
	sqlStr := `insert into classify_goods(name, pid) values(?, ?)`          // 增加
	sqlSelectStr := `select count(name) from classify_goods where name = ?` // 查询重复
	// 初始化结构体
	var count int
	// 首先查询是否已经存在这个分类
	if err := mysql.DB.Get(&count, sqlSelectStr, in.Name); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("商品已存在")
	}
	// 不存在就创建
	if _, err := mysql.DB.Exec(sqlStr, in.Name, PID); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	// 转成返回类型
	classifyInfoResponse := &proto.ClassifyInfoResponse{
		Id:   0,
		Pid:  PID,
		Name: in.Name,
	}
	return classifyInfoResponse, nil
}

// UpdateClassifyInfo 更新品牌分类
func (g *GoodsServerServer) UpdateClassifyInfo(ctx context.Context, in *proto.ClassifyUpdateInfoRequest) (*proto.GoodsEmpty, error) {
	// 拿到上一级父类的PID
	PID := GetPID(in.PName)
	// 编写sql语句
	sqlStr := `update  classify_goods set name = ?, pid = ? where name = ?`             // 增加
	sqlSelectStr := `select count(name) from classify_goods where name = ? and pid = ?` // 查询重复
	// 初始化结构体
	var count int
	// 首先查询是否已经存在这个分类
	if err := mysql.DB.Get(&count, sqlSelectStr, in.OldName, PID); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("目标用户不存在")
	}
	// 不存在就创建
	if _, err := mysql.DB.Exec(sqlStr, in.NewName, PID, in.OldName); err != nil {
		zap.L().Warn("handler.good.CreateClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}

// DeleteClassifyInfo 删除品牌分类
func (g *GoodsServerServer) DeleteClassifyInfo(ctx context.Context, in *proto.ClassifyDeleteInfoRequest) (*proto.GoodsEmpty, error) {
	// 编写sql语句
	sqlStr := `delete from classify_goods where name = ?` // 增加
	res, err := mysql.DB.Exec(sqlStr, in.Name)
	if err != nil {
		zap.L().Warn("handler.good.DeleteClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	row, err := res.RowsAffected()
	if row != 1 {
		zap.L().Warn("handler.good.DeleteClassifyInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}

// GetBannerInfo 查询所有轮播图
func (g *GoodsServerServer) GetBannerInfo(context.Context, *proto.GoodsEmpty) (*proto.BannerListResponse, error) {
	// 编写查询sql语句
	sqlStr := `select * from banner`
	// 初始化结构体
	var bannerInfo []models.BannerInfo
	var bannerInfoResponse []*proto.BannerInfoResponse
	bannerListResponse := proto.BannerListResponse{}
	// sqlx
	err := mysql.DB.Select(&bannerInfo, sqlStr)
	if err != nil {
		zap.L().Warn("handler.good.GetBannerInfo mysql.DB.Select failed", zap.Error(err))
		return nil, err
	}
	// 循环取出+转换类型
	for _, values := range bannerInfo {
		bannerInfoResponse = append(bannerInfoResponse, &proto.BannerInfoResponse{
			Id:            values.ID,
			Level:         values.Level,
			ImageUrl:      values.ImageUrl,
			ImageGoodsUrl: values.ImageGoodsUrl,
		})
	}
	bannerListResponse.BannerInfo = bannerInfoResponse
	return &bannerListResponse, nil
}

// CreateBannerInfo 创建轮播图
func (g *GoodsServerServer) CreateBannerInfo(ctx context.Context, in *proto.BannerCreateInfoRequest) (*proto.BannerInfoResponse, error) {
	// 编写sqlStr语句,在这里不做查询是否存在的选择，因为后台管理可以一目了然，轮播图不可能那么多
	sqlStr := `insert into banner(image_url, image_goods_url, level) values(?, ?, ?)`
	// sqlx
	_, err := mysql.DB.Exec(sqlStr, in.ImageUrl, in.ImageGoodsUrl, in.Level)
	if err != nil {
		zap.L().Warn("handler.good.CreateBannerInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	// 转换成*proto.BannerInfoResponse
	BannerInfoResponse := proto.BannerInfoResponse{
		Id:            0,
		ImageUrl:      in.ImageUrl,
		ImageGoodsUrl: in.ImageGoodsUrl,
	}
	return &BannerInfoResponse, nil
}

// UpdateBannerInfo 更新轮播图图信息
func (g *GoodsServerServer) UpdateBannerInfo(ctx context.Context, in *proto.BannerCreateInfoRequest) (*proto.GoodsEmpty, error) {
	// 编写sqlStr语句,在这里不做查询是否存在的选择，因为后台管理可以一目了然，轮播图不可能那么多
	sqlStr := `update banner set image_url = ?, image_goods_url = ? where level = ?`
	// sqlx
	_, err := mysql.DB.Exec(sqlStr, in.ImageUrl, in.ImageGoodsUrl, in.Level)
	if err != nil {
		zap.L().Warn("handler.good.UpdateBannerInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}

// DeleteBannerInfo 删除轮播图
func (g *GoodsServerServer) DeleteBannerInfo(ctx context.Context, in *proto.BannerDeleteInfoRequest) (*proto.GoodsEmpty, error) {
	// 编写sqlStr语句,在这里不做查询是否存在的选择，因为后台管理可以一目了然，轮播图不可能那么多
	sqlStr := `delete from banner where level = ?`
	// sqlx
	res, err := mysql.DB.Exec(sqlStr, in.Level)
	if err != nil {
		zap.L().Warn("handler.good.DeleteBannerInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	row, err := res.RowsAffected()
	if row != 1 {
		zap.L().Warn("handler.good.DeleteBannerInfo mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.GoodsEmpty{}, nil
}
