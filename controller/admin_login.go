package controller

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/weiyouwozuiku/Gateway/handler"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
//
//	@Summary		管理员登陆
//	@Description	管理员登陆
//	@Tags			管理员接口
//	@ID				/admin_login/login
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.AdminLoginInput								true	"body"
//	@Success		200		{object}	middleware.Response{data=dto.AdminLoginOutput}	"success"
//	@Router			/admin_login/login [post]
func (ad *AdminLoginController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, middleware.ParamErr, err)
		return
	}
	// 1. params.UserName获取管理员信息
	db, err := handler.GetGORMPool(handler.DBDefault)
	if err != nil {
		middleware.ResponseError(ctx, middleware.GetGormPoolErr, err)
	}
	admin := &dao.Admin{}
	admin, err = admin.Login(ctx, db, params)
	if err != nil {
		middleware.ResponseError(ctx, middleware.AdminLoginErr, err)
		return
	}
	// 设置session
	sessionInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(ctx, middleware.SessionParseErr, err)
		return
	}
	sess := sessions.Default(ctx)
	// session一天过期
	sess.Options(sessions.Options{MaxAge: 24 * 60 * 60})
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	err = sess.Save()
	if err != nil {
		middleware.ResponseError(ctx, middleware.SessionOptErr, errors.New("session保存失败"))
		return
	}
	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(ctx, out)
}

// AdminLoginOut godoc
//
//	@Summary		管理员退出
//	@Description	管理员退出
//	@Tags			管理员接口
//	@ID				/admin_login/logout
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=string}	"success"
//	@Router			/admin_login/logout [get]
func (ad *AdminLoginController) AdminLoginOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	if info := sess.Get(public.AdminSessionInfoKey); info == nil {
		middleware.ResponseError(ctx, middleware.AdminLoginErr, errors.New("当前用户未登录"))
		return
	}
	sess.Delete(public.AdminSessionInfoKey)
	err := sess.Save()
	if err != nil {
		middleware.ResponseError(ctx, middleware.SessionOptErr, errors.New("session保存失败"))
		return
	}
	middleware.ResponseSuccess(ctx, "退出成功")
}
