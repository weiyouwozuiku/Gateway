package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

type ServiceAddInput struct {
	ServiceName         string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` //服务名称
	ServiceDesc         string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      //服务描述
	LoadType            int    `json:"load_type" form:"load_type" comment:"负载类型 0=http 1=tcp 2=grpc"`
	HTTPHosts           string `json:"http_hosts" form:"http_hosts" comment:"域名信息" example:"" validate:""`
	HTTPPaths           string `json:"http_paths" form:"http_paths" comment:"路径信息" example:"" validate:""`
	NeedStripUri        int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:""` //启用strip_uri
	Port                int    `json:"port" form:"port" comment:"tcp/grpc端口"`
	LoadBalanceStrategy string `json:"load_balance_strategy" form:"load_balance_strategy" comment:"负载策略"`
	LoadBalanceType     string `json:"load_balance_type" form:"load_balance_type" comment:"负载类型"`
	AuthType            string `json:"auth_type" form:"auth_type" comment:"鉴权类型"`
	UpstreamList        string `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string `json:"plugin_conf" form:"plugin_conf" comment:"插件配置"`
}
type ServiceUpdateInput struct {
	ID                  int64  `json:"id" form:"id" comment:"服务ID" example:"62" validate:"required,min=1"`                                //服务ID
	ServiceName         string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` //服务名称
	ServiceDesc         string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      //服务描述
	LoadType            int    `json:"load_type" form:"load_type" comment:"负载类型 0=http 1=tcp 2=grpc"`
	HTTPHosts           string `json:"http_hosts" form:"http_hosts" comment:"域名信息" example:""`
	HTTPPaths           string `json:"http_paths" form:"http_paths" comment:"路径信息" example:""`
	NeedStripUri        int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:""` //启用strip_uri  max=1,min=0
	Port                int    `json:"port" form:"port" comment:"tcp/grpc端口"`
	LoadBalanceStrategy string `json:"load_balance_strategy" form:"load_balance_strategy" comment:"负载策略"`
	LoadBalanceType     string `json:"load_balance_type" form:"load_balance_type" comment:"负载类型"`
	AuthType            string `json:"auth_type" form:"auth_type" comment:"鉴权类型"`
	UpstreamList        string `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string `json:"plugin_conf" form:"plugin_conf" comment:"插件配置"`
}
type ServiceDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"1" validate:"required"` //服务ID
}
type ServiceDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` //服务ID
}
type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` //每页条数
}
type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"0" validate:""` //总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表" validate:""`               //列表
}
type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	Qps         int64  `json:"qps" form:"qps"`                   //qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   //qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
	ActiveNode  int    `json:"active_node" form:"active_node"`   //活跃的节点数
}
type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" example:"0" validate:""`         //列表
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"0" validate:""` //列表
	//LastWeek  []int64 `json:"last_week" form:"last_week" comment:"七日流量" example:"0" validate:""` //列表
}

func (param *ServiceAddInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
func (param *ServiceUpdateInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
func (param *ServiceDeleteInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
func (param *ServiceListInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
func (param *ServiceDetailInput) BindValidParam(ctx *gin.Context) error {
	return public.DefaultGetValidParams(ctx, param)
}
