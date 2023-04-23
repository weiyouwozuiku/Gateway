package public

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 获取验证器
func GetValidator(ctx *gin.Context) (*validator.Validate, error) {
	valid, ok := ctx.Get(ValidtorKey)
	if !ok {
		return nil, errors.New("未设置验证器")
	}
	validtor, ok := valid.(*validator.Validate)
	if !ok {
		return nil, errors.New("获取验证器失败")
	}
	return validtor, nil
}

// 获取翻译器
func GetTranslation(ctx *gin.Context) (ut.Translator, error) {
	trans, ok := ctx.Get(TranslatorKey)
	if !ok {
		return nil, errors.New("未设置翻译器")
	}
	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("获取翻译器失败")
	}
	return translator, nil
}

func DefaultGetValidParams(ctx *gin.Context, params any) error {
	// 接收参数
	if err := ctx.ShouldBind(params); err != nil {
		return err
	}
	// 获取验证器
	valid, err := GetValidator(ctx)
	if err != nil {
		return err
	}
	// 获取翻译器
	trans, err := GetTranslation(ctx)
	if err != nil {
		return err
	}
	err = valid.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return nil
}
