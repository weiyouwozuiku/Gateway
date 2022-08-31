package controller

import "github.com/gin-gonic/gin"

type AdminController struct{}

func AdminRegister(router *gin.RouterGroup) {
	admin := &AdminController{}
	router.GET("/admin_info", admin.AdminInfo)
	router.POST("/change_pwd", admin.ChangePwd)
}
func (admin *AdminController) AdminInfo(ctx *gin.Context) {

}
func (admin *AdminController) ChangePwd(ctx *gin.Context) {

}
