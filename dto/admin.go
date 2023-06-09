package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}
type ChangePwdInput struct {
	OriginPassword string `json:"origin_password" comment:"原密码" example:"123456" validate:"required"`
	Password       string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

func (param *ChangePwdInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
