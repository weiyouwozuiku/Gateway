package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (ad *AdminController) AdminInfo(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	// 1. 读取sessionKey对应json 转换为结构体
	// 2. 取出数据然后封装输出到结构体
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(ctx, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (ad *AdminController) ChangePwd(ctx *gin.Context) {
	param := &dto.ChangePwdInput{}
	if err := param.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	// 1. session读取用户信息到结构体 adminSessionInfo
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	// 2. sessInfo.ID读取数据库信息 adminInfo
	tx, err := server.GetGORMPool(server.DBDefault)
	if err != nil {
		middleware.ResponseError(ctx, middleware.GetGormPoolFailed, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(ctx, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(ctx, middleware.GormQueryFailed, err)
		return
	}
	// 3. param.password+adminInfo.salt sha256 saltPassword
	adminInfo.Password = public.GenSaltPasswd(adminInfo.Salt, param.Password)
	// 4. saltPassword==>adminInfo.password 执行数据保存
	if err := adminInfo.Save(ctx, tx); err != nil {
		middleware.ResponseError(ctx, middleware.GormSaveFailed, err)
		return
	}
	middleware.ResponseSuccess(ctx, "密码修改成功")
}
