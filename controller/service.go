package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/server"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := ServiceController{}

}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(ctx *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	tx, err := server.GetGORMPool("default")
	if err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	// 1. 从db获取信息
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(ctx, tx, param)
	if err != nil {
		middleware.ResponseError(ctx, middleware.GormQueryFailed, err)
		return
	}
	// 2. 格式化输出
	outList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := listItem.Ser
	}
}
