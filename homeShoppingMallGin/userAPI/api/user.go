package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
	"homeShoppingMallGin/userAPI/consulR"
	"homeShoppingMallGin/userAPI/models"
	"homeShoppingMallGin/userAPI/myResponseCode"
	"homeShoppingMallGin/userAPI/pkg/jwt"
	"homeShoppingMallGin/userAPI/proto"
	"homeShoppingMallGin/userAPI/validators"
	"net/http"
	"strconv"
)

var ok *proto.IsRight

// GetAllUserList 获取所有用户信息并分页
func GetAllUserList(c *gin.Context) {
	// 获取分页信息
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	// 调用server端方法
	response, err := consulR.UserSrvClient.GetAllUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.L().Error("api.user consul.UserSrvClient.GetAllUserList failed", zap.Error(err))
		return
	}
	result := make([]interface{}, 0)
	for _, value := range response.UserList {
		user := models.UserResponse{
			UserId:   value.UserID,
			Nickname: value.NickName,
			Gender:   value.Gender,
			Mobile:   value.Mobile,
			Role:     value.Role,
		}
		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

// GetUserInfo 获取单个用户信息
func GetUserInfo(c *gin.Context) {
	// 通过url拿取数据
	mobile := c.Query("mobile")
	// 调用后端方法
	rsp, err := consulR.UserSrvClient.GetUserInfoByMobile(context.Background(), &proto.MobileInfo{Mobile: mobile})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
		return
	}
	myResponseCode.ResponseSuccess(c, rsp)
}

// UserLogin 用户登录接口
func UserLogin(c *gin.Context) {
	user := new(models.UserLoginInfo)
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(user); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 判断用户是否存在
	rsp, err := consulR.UserSrvClient.GetUserInfoByMobile(context.Background(), &proto.MobileInfo{Mobile: user.Mobile})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeMobileExist)
		return
	}
	// 判断验证码是否一致
	if ok, err = consulR.UserSrvClient.CheckAuthCode(context.Background(), &proto.AuthCodeInfo{
		Mobile:       rsp.Mobile,
		UserAuthCode: user.AuthCode,
	}); err != nil {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidAuthCode, "登陆失败")
		return
	}
	if !ok.Ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidAuthCode)
		return
	}
	// 判断密码是否正确
	if ok, err = consulR.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordInfo{
		Password:          user.Password,
		EncryptedPassword: rsp.Password,
	}); err != nil {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidPassword, "登陆失败")
		return
	}
	if !ok.Ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidPassword)
		return
	}
	// 密码校验正确，返回token
	token, err := jwt.GenToken(rsp.UserID, user.Mobile, rsp.Role)
	if err != nil {
		zap.L().Error("user jwt.GenToken failed", zap.Error(err))
	}
	myResponseCode.ResponseSuccess(c, token)
}

// SendAuthCode 短信验证码发送接口
func SendAuthCode(c *gin.Context) {
	user := new(models.UserRLCodeInfo)
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(user); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 发送短信接口
	CodeStr, err := consulR.UserSrvClient.SendAuthCode(context.Background(), &proto.MobileInfo{Mobile: user.Mobile})
	if err != nil {
		zap.L().Error("user consulR.UserSrvClient.SendAuthCode failed", zap.Error(err))
	}
	myResponseCode.ResponseSuccess(c, CodeStr)
}

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	user := new(models.UserRegisterInfo)
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(user); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 判断用户是否存在
	_, err := consulR.UserSrvClient.GetUserInfoByMobile(context.Background(), &proto.MobileInfo{Mobile: user.Mobile})
	if err == nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeMobileHave)
		return
	}

	// 判断验证码是否一致
	if ok, err = consulR.UserSrvClient.CheckAuthCode(context.Background(), &proto.AuthCodeInfo{
		Mobile:       user.Mobile,
		UserAuthCode: user.AuthCode,
	}); err != nil {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidAuthCode, "验证码错误")
		return
	}

	if !ok.Ok {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidAuthCode, "验证码错误")
		return
	}

	// 判断密码是否正确
	if user.Password != user.RePassword {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidPassword, "密码不一致")
		return
	}
	// 保存用户
	_, err = consulR.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		UserID:   0,
		Mobile:   user.Mobile,
		Password: user.Password,
	})
	if err != nil {
		zap.L().Error("api.user consulR.UserSrvClient.CreateUser failed", zap.Error(err))
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, "注册成功")
}

// UserUpdateMobilePassword 用户更新信息接口 只更新手机和老密码
func UserUpdateMobilePassword(c *gin.Context) {
	// 通过JWT拿到用户ID
	userID, ok := c.Get("userID")
	if !ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	id := userID.(int64)
	// 实例化结构体
	var userInfo models.UserUpdateInfo
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 这里有两种方法，一是可以通过JWT拿，但是我这里没有加进去，所以我选择了直接调用查询用户ID接口，来获取老密码
	rsp, err := consulR.UserSrvClient.GetUserInfoByUserID(context.Background(), &proto.UserID{UserID: id})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	// 判断手机是否输入一样
	if rsp.Mobile != userInfo.Mobile {
		myResponseCode.ResponseError(c, myResponseCode.CodeMobileExist)
		return
	}
	// 判断老密码是否一致
	isRight, err := consulR.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordInfo{
		Password:          userInfo.OldPassword,
		EncryptedPassword: rsp.Password,
	})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	if !isRight.Ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidPassword)
		return
	}

	// 判断验证码是否一致
	if isRight, err = consulR.UserSrvClient.CheckAuthCode(context.Background(), &proto.AuthCodeInfo{
		Mobile:       userInfo.Mobile,
		UserAuthCode: userInfo.AuthCode,
	}); err != nil {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidAuthCode, "登陆失败")
		return
	}

	if !isRight.Ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeInvalidAuthCode)
		return
	}
	// 判断密码是否正确
	if userInfo.NewPassword != userInfo.RePassword {
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidPassword, "密码不一致")
		return
	}
	// 调用更新接口
	userID1 := strconv.FormatInt(rsp.UserID, 10)
	if _, err = consulR.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Mobile:   userInfo.Mobile,
		NickName: rsp.NickName,
		UserID:   userID1,
		Password: userInfo.NewPassword,
		Gender:   rsp.Gender,
		Role:     rsp.Role,
	}); err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeSuccess)
}

// UserUpdateNickName 用户更新昵称
func UserUpdateNickName(c *gin.Context) {
	// 通过JWT拿到用户ID
	userID, ok := c.Get("userID")
	if !ok {
		myResponseCode.ResponseError(c, myResponseCode.CodeNeedLogin)
		return
	}
	id := userID.(int64)
	// 实例化结构体
	var userInfo models.UserUpdateNickName
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 拿到所有用户数据
	rsp, err := consulR.UserSrvClient.GetUserInfoByUserID(context.Background(), &proto.UserID{UserID: id})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeMobileExist)
		return
	}
	// 调用更新接口
	userID1 := strconv.FormatInt(rsp.UserID, 10)
	if _, err = consulR.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Mobile:   rsp.Mobile,
		NickName: userInfo.NickName,
		UserID:   userID1,
		Password: rsp.Password,
		Gender:   rsp.Gender,
		Role:     rsp.Role,
	}); err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeSuccess)
}

// AdminRole 管理员修改用户权限接口
func AdminRole(c *gin.Context) {
	// 实例化结构体
	var userInfo models.AdminUserUpdateRole
	// 拿到用户post信息，并判断是不是ValidationErrors进行错误处理
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		// 判断发生的错误是不是validator.ValidationErrors类型
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			myResponseCode.ResponseError(c, myResponseCode.CodeInvalidParam)
			return
		}
		myResponseCode.ResponseErrorWithMsg(c, myResponseCode.CodeInvalidParam, validators.RemoveTopStruct(err.Translate(validators.Trans)))
		return
	}
	// 拿到用户数据
	rsp, err := consulR.UserSrvClient.GetUserInfoByMobile(context.Background(), &proto.MobileInfo{Mobile: userInfo.Mobile})
	if err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeMobileExist)
		return
	}
	// 调用更新接口
	userID1 := strconv.FormatInt(rsp.UserID, 10)
	if _, err = consulR.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Mobile:   rsp.Mobile,
		NickName: rsp.NickName,
		UserID:   userID1,
		Password: rsp.Password,
		Gender:   rsp.Gender,
		Role:     userInfo.Role,
	}); err != nil {
		myResponseCode.ResponseError(c, myResponseCode.CodeServerBusy)
		return
	}
	myResponseCode.ResponseSuccess(c, myResponseCode.CodeSuccess)
}
