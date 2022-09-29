package middleware

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/weiyouwozuiku/Gateway/public"
)

func customValidation(val *validator.Validate) {
	// 自定义验证方法
	val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "admin"
	})
	val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
		return matched
	})
	val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
		return matched
	})
	val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 2 {
				return false
			}
		}
		return true
	})
	val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 3 {
				return false
			}
		}
		return true
	})
	val.RegisterValidation("valid_ipportlist", func(fl validator.FieldLevel) bool {
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
				return false
			}
		}
		return true
	})
	val.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, item := range strings.Split(fl.Field().String(), ",") {
			matched, _ := regexp.Match(`\S+`, []byte(item)) //ip_addr
			if !matched {
				return false
			}
		}
		return true
	})
	val.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
		fmt.Println(fl.Field().String())
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched, _ := regexp.Match(`^\d+$`, []byte(ms)); !matched {
				return false
			}
		}
		return true
	})
}
func customTranslation(val *validator.Validate, trans ut.Translator) {
	// 自定义翻译器
	val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
		return ut.Add("valid_username", "{0} 填写不正确哦", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_username", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
		return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_service_name", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
		return ut.Add("valid_rule", "{0} 必须是非空字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_rule", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
		return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_url_rewrite", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_header_transfor", trans, func(ut ut.Translator) error {
		return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_header_transfor", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_ipportlist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_ipportlist", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_iplist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_iplist", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_weightlist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_weightlist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_weightlist", fe.Field())
		return t
	})
}
func ValidtorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置支持语言
		en := en.New()
		zh := zh.New()
		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		// 设置验证器
		val := validator.New()
		// 根据参数获取翻译器实例
		locale := ctx.DefaultQuery("locale", "zh")
		trans, found := uni.GetTranslator(locale)
		if !found {
			trans, _ = uni.GetTranslator("zh")
		}
		// 翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(field reflect.StructField) string {
				return field.Tag.Get("en_comment")
			})
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(field reflect.StructField) string {
				return field.Tag.Get("comment")
			})
			customValidation(val)
			customTranslation(val, trans)
		}
		ctx.Set(public.TranslatorKey, trans)
		ctx.Set(public.ValidtorKey, val)
		ctx.Next()
	}
}
