package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/handler"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
)

type AppController struct{}

// AppControllerRegister admin路由注册
func AppRegister(router *gin.RouterGroup) {
	admin := AppController{}
	router.GET("/app_list", admin.AppList)
	router.GET("/app_detail", admin.AppDetail)
	router.GET("/app_stat", admin.AppStatistics)
	router.GET("/app_delete", admin.AppDelete)
	router.POST("/app_add", admin.AppAdd)
	router.POST("/app_update", admin.AppUpdate)
}

// AppList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.AppListOutput} "success"
// @Router /app/app_list [get]
func (admin *AppController) AppList(c *gin.Context) {
	params := &dto.AppListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}
	info := &dao.App{}
	list, total, err := info.AppList(c, handler.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, middleware.QueryGormErr, err)
		return
	}

	outputList := []dto.AppListItemOutput{}
	for _, item := range list {
		appCounter, err := middleware.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + item.AppID)
		if err != nil {
			middleware.ResponseError(c, middleware.InnerErr, err)
			c.Abort()
			return
		}
		outputList = append(outputList, dto.AppListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
		})
	}
	output := dto.AppListOutput{
		List:  outputList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// AppDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (admin *AppController) AppDetail(c *gin.Context) {
	params := &dto.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, handler.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, middleware.QueryGormErr, err)
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}

// AppDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (admin *AppController) AppDelete(c *gin.Context) {
	params := &dto.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, handler.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, middleware.QueryGormErr, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, handler.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, middleware.SaveGormErr, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.AppAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (admin *AppController) AppAdd(c *gin.Context) {
	params := &dto.AppAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}

	//验证app_id是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := search.Find(c, handler.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, middleware.InnerErr, errors.New("租户ID被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	tx := handler.GORMDefaultPool
	info := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, middleware.SaveGormErr, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.AppUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (admin *AppController) AppUpdate(c *gin.Context) {
	params := &dto.AppUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, handler.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, middleware.QueryGormErr, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(c, handler.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, middleware.SaveGormErr, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppStatistics godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.StatisticsOutput} "success"
// @Router /app/app_stat [get]
func (admin *AppController) AppStatistics(c *gin.Context) {
	params := &dto.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamErr, err)
		return
	}

	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, handler.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, middleware.QueryGormErr, err)
		return
	}

	//今日流量全天小时级访问统计
	todayStat := []int64{}
	counter, err := middleware.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + detail.AppID)
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		c.Abort()
		return
	}
	currentTime := time.Now()
	for i := 0; i <= time.Now().In(public.TimeLocation).Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, public.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayStat = append(todayStat, hourData)
	}

	//昨日流量全天小时级访问统计
	yesterdayStat := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, public.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayStat = append(yesterdayStat, hourData)
	}
	stat := dto.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	middleware.ResponseSuccess(c, stat)
	return
}
