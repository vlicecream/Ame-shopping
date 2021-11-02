package validators

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, err := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if err != nil {
		zap.L().Error("myValidator  regexp.MatchString failed", zap.Error(err))
		return false
	}
	if !ok {
		return false
	}
	return true
}

//// RegisterTranslator 为自定义字段添加翻译功能
//func RegisterTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
//	return func(trans ut.Translator) error {
//		if err := trans.Add(tag, msg, false); err != nil {
//			return err
//		}
//		return nil
//	}
//}

//// Translate 自定义字段的翻译方法
//func Translate(trans ut.Translator, fe validator.FieldError) string {
//	msg, err := trans.T(fe.Tag(), fe.Field())
//	if err != nil {
//		panic(fe.(error).Error())
//	}
//	return msg
//}