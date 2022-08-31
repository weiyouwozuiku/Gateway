package controller

import "github.com/gin-gonic/gin"

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login")
}
func (admin *AdminLoginController) AdminLogin(ctx *gin.Context) {

}
