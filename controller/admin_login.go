package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/dto"
	"github.com/weiyouwozuiku/gateway/middleware"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login")
}
func (admin *AdminLoginController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 10001, err)
	}
}
