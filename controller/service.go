package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/weiyouwozuiku/Gateway/dao"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
	"gorm.io/gorm"
	"strings"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)

	group.POST("/service_add_tcp", service.ServiceAddTcp)
	group.POST("/service_update_tcp", service.ServiceUpdateTcp)
	group.POST("/service_add_grpc", service.ServiceAddGrpc)
	group.POST("/service_update_grpc", service.ServiceUpdateGrpc)
}

// ServiceList godoc
//
//	@Summary		服务列表
//	@Description	服务列表
//	@Tags			服务管理
//	@ID				/service/service_list
//	@Accept			json
//	@Produce		json
//	@Param			info		query		string											false	"关键词"
//	@Param			page_size	query		int												true	"每页个数"
//	@Param			page_no		query		int												true	"当前页数"
//	@Success		200			{object}	middleware.Response{data=dto.ServiceListOutput}	"success"
//	@Router			/service/service_list [get]
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
//
//	@Summary		服务删除
//	@Description	服务删除
//	@Tags			服务管理
//	@ID				/service/service_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string								true	"服务ID"
//	@Success		200	{object}	middleware.Response{data=string}	"success"
//	@Router			/service/service_delete [get]
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

// ServiceAddHTTP godoc
//
//	@Summary		添加HTTP服务
//	@Description	添加HTTP服务
//	@Tags			服务管理
//	@ID				/service/service_add_http
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ServiceAddHTTPInput				true	"body"
//	@Success		200		{object}	middleware.Response{data=string}	"success"
//	@Router			/service/service_add_http [post]
func (service *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2004, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx, err := server.GetGORMPool(server.DBDefault)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	tx = tx.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	} else {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			middleware.ResponseError(c, 2002, errors.New("服务名查询失败"))
			return
		}
	}

	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Find(c, tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("服务接入前缀或域名已存在"))
		return
	} else {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			middleware.ResponseError(c, 2003, errors.New("服务接入前缀或域名查询失败"))
			return
		}
	}

	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	//serviceModel.ID
	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	loadbalance := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadbalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
}
