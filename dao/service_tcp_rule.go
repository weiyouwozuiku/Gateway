package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port      int   `json:"port" gorm:"column:port" description:"端口"`
}

func (rule *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}
func (rule *TcpRule) Find(ctx *gin.Context, tx *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := tx.WithContext(ctx).Where(search).Find(model).Error
	return model, err
}
func (rule *TcpRule) Save(ctx *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Save(rule).Error; err != nil {
		return err
	}
	return nil
}
