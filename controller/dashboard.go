package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/handler"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
	"time"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	service := &DashboardController{}
	group.GET("/panel_group_data", service.PanelGroupData)
	group.GET("/flow_stat", service.FlowStat)
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
//
//	@Summary		指标统计
//	@Description	指标统计
//	@Tags			首页大盘
//	@ID				/dashboard/panel_group_data
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=dto.PanelGroupDataOutput}	"success"
//	@Router			/dashboard/panel_group_data [get]
func (s *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := handler.GetGORMPool(handler.DBDefault)
	if err != nil {
		middleware.ResponseError(c, middleware.GetGormPoolErr, err)
	}
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		return
	}
	app := &dao.App{}
	_, appNum, err := app.AppList(c, tx, &dto.AppListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		return
	}
	counter, err := middleware.FlowCounterHandler.GetFlowCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		return
	}
	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: counter.TotalCount,
		CurrentQPS:      counter.QPS,
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceStat godoc
//
//	@Summary		服务统计
//	@Description	服务统计
//	@Tags			首页大盘
//	@ID				/dashboard/service_stat
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=dto.DashServiceStatOutput}	"success"
//	@Router			/dashboard/service_stat [get]
func (s *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := handler.GetGORMPool(handler.DBDefault)
	if err != nil {
		middleware.ResponseError(c, middleware.GetGormPoolErr, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		return
	}
	var legend []string
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, middleware.InnerErr, errors.New(fmt.Sprintf("load_type not found,current load_type is %d", item.LoadType)))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
//
//	@Summary		服务统计
//	@Description	服务统计
//	@Tags			首页大盘
//	@ID				/dashboard/flow_stat
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middleware.Response{data=dto.ServiceStatOutput}	"success"
//	@Router			/dashboard/flow_stat [get]
func (s *DashboardController) FlowStat(c *gin.Context) {
	counter, err := middleware.FlowCounterHandler.GetFlowCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, middleware.InnerErr, err)
		return
	}
	var todayList []int64
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, public.TimeLocation)
		hourDate, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourDate)
	}
	var yesterdayList []int64
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, public.TimeLocation)
		hourDate, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourDate)
	}
	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}
