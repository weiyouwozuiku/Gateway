package dao

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/public"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	ID                  int64     `json:"id" gorm:"primary_key"`
	ServiceType         int       `json:"service_type" gorm:"column:service_type" description:"服务类型 0=http 1=tcp 2=grpc"`
	ServiceName         string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc         string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	ServicePort         int       `json:"service_port" gorm:"column:service_port" description:"服务端口(只针对 tcp/grpc)"`
	HTTPHosts           string    `json:"http_hosts" gorm:"column:http_hosts" description:"域名信息"`
	HTTPPaths           string    `json:"http_paths" gorm:"column:http_paths" description:"路径信息"`
	HttpStripPrefix     int       `json:"http_strip_prefix" gorm:"column:http_strip_prefix" description:"http转发前剥离前缀"`
	LoadBalanceStrategy string    `json:"load_balance_strategy" gorm:"column:load_balance_strategy" description:"负载策略"`
	LoadBalanceType     string    `json:"load_balance_type" gorm:"column:load_balance_type" description:"负载类型"`
	AuthType            string    `json:"auth_type" gorm:"column:auth_type" description:"鉴权类型"`
	UpstreamList        string    `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string    `json:"plugin_conf" gorm:"column:plugin_conf" description:"插件配置"`
	CreatedAt           time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	UpdatedAt           time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete            int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (info *ServiceInfo) TableName() string {
	return "gateway_service_info"
}
func (info *ServiceInfo) Find(ctx *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.WithContext(ctx).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (info *ServiceInfo) Save(ctx *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Save(info).Error
}
func (info *ServiceInfo) Delete(ctx *gin.Context, tx *gorm.DB, search *ServiceInfo) error {
	if err := tx.WithContext(ctx).Where("id=?", search.ID).Delete(search).Error; err != nil {
		return err
	}
	return nil
}
func (info *ServiceInfo) PageList(ctx *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize
	query := tx.WithContext(ctx)
	query = query.Table(info.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Count(&total)
	return list, total, nil
}
func (info *ServiceInfo) ServiceDetail(ctx *gin.Context, tx *gorm.DB, serInfo *ServiceInfo) (*ServiceDetail, error) {
	if serInfo.ServiceName == "" {
		fmt.Println("FInd", public.Obj2Json(serInfo))
		if tmp, err := info.Find(ctx, tx, serInfo); err != nil {
			return nil, err
		} else {
			serInfo = tmp
		}

	}
	pluginConf:=simplejson.New()
	if tmp,err:=simplejson.NewJson([]byte(serInfo.))
}
