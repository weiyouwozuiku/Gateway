package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/middleware"
)

type AdminLoginInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (input *AdminLoginInput) BindValidParam(ctx *gin.Context) error {
	return middleware.DefaultValidParams(ctx, input)
}
