package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/dto"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login")
}
func (admin *AdminLoginController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(ctx); err != nil {

	}
}
