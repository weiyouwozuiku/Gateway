package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dto"
	"gorm.io/gorm"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
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
	var list []ServiceInfo
	offset := (param.PageNo - 1) * param.PageSize
	query := tx.WithContext(ctx).Table(info.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Count(&total)
	return list, total, nil
}
func (info *ServiceInfo) ServiceDetail(ctx *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := info.Find(ctx, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	resultChan := make(chan any, 5)
	findHttpRule(ctx, tx, search.ID, resultChan)
	findTcpRule(ctx, tx, search.ID, resultChan)
	findGrpcRule(ctx, tx, search.ID, resultChan)
	findAccess(ctx, tx, search.ID, resultChan)
	findLoad(ctx, tx, search.ID, resultChan)

	detail := &ServiceDetail{Info: search}
	for i := 0; i < 5; i++ {
		result := <-resultChan
		switch result.(type) {
		case *HttpRule:
			detail.HTTPRule = result.(*HttpRule)
		case *TcpRule:
			detail.TCPRule = result.(*TcpRule)
		case *GrpcRule:
			detail.GRPCRule = result.(*GrpcRule)
		case *LoadBalance:
			detail.LoadBalance = result.(*LoadBalance)
		case *AccessControl:
			detail.AccessControl = result.(*AccessControl)
		case error:
			return nil, result.(error)
		}
	}
	return detail, nil
}
func (info *ServiceInfo) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	list := []dto.DashServiceStatItemOutput{}
	if err := tx.WithContext(c).Table(info.TableName()).Where("is_delete=0").Select("load_type,count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func findHttpRule(ctx *gin.Context, tx *gorm.DB, searchId int64, resultChan chan<- any) {
	var err error
	item := &HttpRule{ServiceID: searchId}
	item, err = item.Find(ctx, tx, item)
	if err != nil && err != gorm.ErrRecordNotFound {
		resultChan <- err
	} else {
		resultChan <- item
	}
}
func findTcpRule(ctx *gin.Context, tx *gorm.DB, searchId int64, resultChan chan<- any) {
	var err error
	item := &TcpRule{ServiceID: searchId}
	item, err = item.Find(ctx, tx, item)
	if err != nil && err != gorm.ErrRecordNotFound {
		resultChan <- err
	} else {
		resultChan <- item
	}
}
func findGrpcRule(ctx *gin.Context, tx *gorm.DB, searchId int64, resultChan chan<- any) {
	var err error
	item := &GrpcRule{ServiceID: searchId}
	item, err = item.Find(ctx, tx, item)
	if err != nil && err != gorm.ErrRecordNotFound {
		resultChan <- err
	} else {
		resultChan <- item
	}
}
func findAccess(ctx *gin.Context, tx *gorm.DB, searchId int64, resultChan chan<- any) {
	var err error
	item := &AccessControl{ServiceID: searchId}
	item, err = item.Find(ctx, tx, item)
	if err != nil && err != gorm.ErrRecordNotFound {
		resultChan <- err
	} else {
		resultChan <- item
	}
}
func findLoad(ctx *gin.Context, tx *gorm.DB, searchId int64, resultChan chan<- any) {
	var err error
	item := &LoadBalance{ServiceID: searchId}
	item, err = item.Find(ctx, tx, item)
	if err != nil && err != gorm.ErrRecordNotFound {
		resultChan <- err
	} else {
		resultChan <- item
	}
}
