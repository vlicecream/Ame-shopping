package myResponseCode

type ResCode int64

// 定义各种报错的响应码
const (
	CodeSuccess ResCode = 100 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeMobileExist
	CodeMobileReadyExist
	CodeInvalidAuthCode
	CodeTooManyRequests
	CodePermissionDenied
)

// 用一个大map收集起来
var codeMsgMap = map[ResCode]string{
	CodeSuccess:          "success",
	CodeInvalidParam:     "无效的参数",
	CodeUserExist:        "用户名不存在",
	CodeUserNotExist:     "密码不存在",
	CodeInvalidPassword:  "无效的密码",
	CodeServerBusy:       "服务繁忙",
	CodeNeedLogin:        "需要登录",
	CodeInvalidToken:     "无效的token",
	CodeMobileExist:      "手机号不存在",
	CodeInvalidAuthCode:  "无效的验证码",
	CodeMobileReadyExist: "手机已存在",
	CodeTooManyRequests:  "请求次数过多，已被限流",
	CodePermissionDenied: "没有权限，请联系管理员",
}

func (r ResCode) Msg() string {
	msg, ok := codeMsgMap[r]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
