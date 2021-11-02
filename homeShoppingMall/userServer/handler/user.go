package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
	"homeShoppingMall/userServer/dao/mysql"
	"homeShoppingMall/userServer/dao/redis"
	"homeShoppingMall/userServer/models"
	"homeShoppingMall/userServer/pkg/snowflake"
	"homeShoppingMall/userServer/proto"
	"homeShoppingMall/userServer/settings"
	"math/rand"
	"os"
	"time"
)

var smallSecret = "夏天夏天悄悄过去留下"

const char = "0123456789"

type UserSeverServer struct {
	*proto.UnimplementedUserSeverServer
}

// RandChar 随机验证码
func RandChar() string {
	/*
		伪随机代码
			rand.NewSource(time.Now().UnixNano()) // 产生随机种子
			var s bytes.Buffer
			for i := 0; i < 5; i++ {
				s.WriteByte(char[rand.Int63()%int64(len(char))])
			}
			return s.String()
	*/

	/*真随机代码*/
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%04v", rnd.Int31n(10000))
	return code
}

// MD5Password MD5密码加密
func MD5Password(password string) string {
	p := md5.New()
	p.Write([]byte(smallSecret))
	return hex.EncodeToString(p.Sum([]byte(password)))
}

// Model2UserInfo 转换proto.UserInfoResponse
func Model2UserInfo(user models.AllUserInfo) *proto.UserInfoResponse {
	userInfoResponse := proto.UserInfoResponse{
		Id:       user.ID,
		UserID:   user.UserID,
		Mobile:   user.Mobile,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	return &userInfoResponse
}

// GetAllUserList 查询所有用户信息并分页
func (u *UserSeverServer) GetAllUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListResponse, error) {
	// 写查询sql语句
	sqlStr := `select * from user limit ?,?`
	// 初始化结构体
	var user []models.AllUserInfo
	var userL proto.UserListResponse
	// sqlx
	if err := mysql.DB.Select(&user, sqlStr, in.Pn, in.PSize); err != nil {
		zap.L().Warn("handler.user.GetAllUserList mysql.DB.Select failed", zap.Error(errors.New("搜素不到用户数据")))
		return nil, err
	}
	// 循环遍历取出user 转成UserInfo类型
	for _, user1 := range user {
		userInfoResponse := Model2UserInfo(user1)
		userL.UserList = append(userL.UserList, userInfoResponse)
	}
	return &userL, nil
}

// GetUserInfoByMobile 通过手机号来查询用户信息
func (u *UserSeverServer) GetUserInfoByMobile(ctx context.Context, in *proto.MobileInfo) (*proto.UserInfoResponse, error) {
	// 查询sql语句
	sqlStr := `select * from user where mobile = ?`
	// 初始化结构体
	var user models.AllUserInfo
	// sqlx
	if err := mysql.DB.Get(&user, sqlStr, in.Mobile); err != nil {
		zap.L().Warn("handler.user.GetUserInfoByMobile mysql.DB.Get failed", zap.Error(errors.New("搜素不到数据")))
		return nil, err
	}
	// 转换类型 返回数据
	userResponse := Model2UserInfo(user)
	return userResponse, nil
}

// GetUserInfoByUserID 通过用户ID来查询用户信息
func (u *UserSeverServer) GetUserInfoByUserID(ctx context.Context, in *proto.UserID) (*proto.UserInfoResponse, error) {
	// 查询sql语句
	sqlStr := `select * from user where user_id = ?`
	// 初始化结构体
	var user models.AllUserInfo
	// sqlx
	if err := mysql.DB.Get(&user, sqlStr, in.UserID); err != nil {
		zap.L().Warn("handler.user.GetUserInfoByUserID mysql.DB.Get failed", zap.Error(errors.New("搜素不到数据")))
		return nil, err
	}
	// 转换类型 返回数据
	userResponse := Model2UserInfo(user)
	return userResponse, nil
}

// CreateUser 创建用户
func (u *UserSeverServer) CreateUser(ctx context.Context, in *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// sql查询与创建语句创建
	sqlCheckStr := `select count(mobile) from user where mobile = ?`
	sqlCreateStr := `insert into user(user_id, mobile, password) values(?,?,?)`
	// 查询这个用户是否存在
	var count int
	if err := mysql.DB.Get(&count, sqlCheckStr, in.Mobile); err != nil {
		zap.L().Error("handler.user.CreateUser mysql.DB.Get failed", zap.Error(err))
		return nil, err
	}
	if count >= 1 {
		return nil, errors.New("用户已存在")
	}
	// 雪花算法生成用户ID
	userID := snowflake.GenID()
	// 对密码进行加密
	newPassword := MD5Password(in.Password)

	// 初始化结构体并保存数据
	var user models.AllUserInfo
	user.UserID = userID
	user.Mobile = in.Mobile
	user.Password = newPassword

	// sqlx保存
	_, err := mysql.DB.Exec(sqlCreateStr, user.UserID, user.Mobile, user.Password)
	if err != nil {
		zap.L().Error("handler.user.CreateUser mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	// 转换并返回
	userResponse := Model2UserInfo(user)
	return userResponse, nil
}

// UpdateUser 更新用户数据
func (u *UserSeverServer) UpdateUser(ctx context.Context, in *proto.UpdateUserInfo) (*proto.Empty, error) {
	// 编写sql更新语句
	sqlStr := `update user set nick_name=?, mobile=?, password=?, gender=?, role=? where user_id = ?`
	// sqlx更新
	_, err := mysql.DB.Exec(sqlStr, in.NickName, in.Mobile, in.Password, in.Gender, in.Role, in.UserID)
	if err != nil {
		zap.L().Error("handler.user.UpdateUser mysql.DB.Exec failed", zap.Error(err))
		return nil, err
	}
	return &proto.Empty{}, nil
}

// CheckPassword 查看并对比密码
func (u *UserSeverServer) CheckPassword(ctx context.Context, in *proto.PasswordInfo) (*proto.IsRight, error) {
	newPassword := MD5Password(in.Password)
	if newPassword == in.EncryptedPassword {
		return &proto.IsRight{Ok: true}, nil
	} else {
		return &proto.IsRight{Ok: false}, nil
	}
}

// CheckAuthCode 查看验证码并对比
func (u *UserSeverServer) CheckAuthCode(ctx context.Context, in *proto.AuthCodeInfo) (*proto.IsRight, error) {
	redisValue := redis.Rdb.Get(in.Mobile)
	if settings.Conf.Mode == "dev" {
		if in.UserAuthCode == "1111" {
			return &proto.IsRight{Ok: true}, nil
		} else {
			return &proto.IsRight{Ok: false}, nil
		}
	}else {
		if redisValue.Val() == in.UserAuthCode {
			return &proto.IsRight{Ok: true}, nil
		} else {
			return &proto.IsRight{Ok: false}, nil
		}
	}
}

// SendAuthCode 生成随机验证码存入数据库并发送短信
func (u *UserSeverServer) SendAuthCode(ctx context.Context, in *proto.MobileInfo) (*proto.RandomAuthCodeInfo, error) {
	// 生成一个随机字符串
	str := RandChar()
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
	 * 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	 * 以免泄露密钥对危及你的财产安全。
	 * CAM密匙查询: https://console.cloud.tencent.com/cam/capi*/
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		//"AKIDay6xXrvkLjriIcfSyDCuLZbNsFSQHvLW",
		//"wmPn0ziBNODqFPKyUu87pHSO9i4B2d07",
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5

	/* SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务
	 * 则必须手动指定域名，例如sms的上海金融区域名： sms.ap-shanghai-fsi.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，或者引用预设的常量 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 你可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/smsv2
	 * sms helper: https://cloud.tencent.com/document/product/382/3773 */

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr("1400578683")
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
	request.SignName = common.StringPtr("Ame林汀")
	/* 国际/港澳台短信 SenderId: 国内短信填空，默认未开通，如需开通请联系 [sms helper] */
	request.SenderId = common.StringPtr("")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	request.SessionContext = common.StringPtr("xxx")
	/* 短信码号扩展号: 默认未开通，如需开通请联系 [sms helper] */
	request.ExtendCode = common.StringPtr("")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs([]string{str})
	/* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
	request.TemplateId = common.StringPtr("1139487")
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs([]string{fmt.Sprintf("+86%s", in.Mobile)})

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*es.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, nil
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s\n", b)

	// 保存入redis数据库
	redis.Rdb.Set(in.Mobile, str, 180000000000)
	strResponse := &proto.RandomAuthCodeInfo{UserAuthCode: str}
	return strResponse, nil
}
