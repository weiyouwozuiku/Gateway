package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

func (rule *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}
func (rule *GrpcRule) Find(ctx *gin.Context, tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	model := &GrpcRule{}
	err := tx.WithContext(ctx).Where(search).Find(model).Error
	return model, err
}
func (rule *GrpcRule) Save(ctx *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Save(rule).Error; err != nil {
		return err
	}
	return nil
}