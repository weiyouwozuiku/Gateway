package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := ServiceController{}
	group.GET("/service_list", service.ServiceList)
}

// ServiceList godoc
// @Summary     服务列表
// @Description 服务列表
// @Tags        服务管理
// @ID          /service/service_list
// @Accept      json
// @Produce     json
// @Param       info      query    string                                          false "关键词"
// @Param       page_size query    int                                             true  "每页个数"
// @Param       page_no   query    int                                             true  "当前页数"
// @Success     200       {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router      /service/service_list [get]
func (service *ServiceController) ServiceList(ctx *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, middleware.InvalidParamsCode, err)
		return
	}
	tx, err := server.GetGORMPool(server.DBDefault)
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
	var outList []dto.ServiceListItemOutput
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(ctx, tx, &listItem)
		if err != nil {
			middleware.ResponseError(ctx, middleware.GormQueryFailed, err)
			return
		}

		// 1. http后缀接入 clusterIP+clusterPort+path
		// 2. http域名接入 domain
		// 3. tcp\grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := public.GetStringConf("base.cluster.cluster_ip")
		clusterPort := public.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := public.GetStringConf("base.cluster.cluster_ssl_port")
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIpListByModel()
		//counter,err:=

		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			// TODO 之后添加流量检测
			Qps:       0,
			Qpd:       0,
			TotalNode: len(ipList),
		}
		outList = append(outList, outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(ctx, out)
}

// ServiceDelete godoc
// @Summary     服务删除
// @Description 服务删除
// @Tags        服务管理
// @ID          /service/service_delete
// @Accept      json
// @Produce     json
// @Param       id  query    string                           true "服务ID"
// @Success     200 {object} middleware.Response{data=string} "success"
// @Router      /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := server.GetGORMPool(server.DBDefault)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
