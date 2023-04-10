package dao

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dto"
	"github.com/weiyouwozuiku/Gateway/public"
	"gorm.io/gorm"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"salt"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (ad *Admin) TableName() string {
	return "gateway_admin"
}
func (ad *Admin) Find(ctx *gin.Context, db *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	if err := db.WithContext(ctx).Where(search).Find(out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
func (ad *Admin) Login(ctx *gin.Context, db *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := ad.Find(ctx, db, &Admin{UserName: param.Username, IsDelete: 0})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户信息不存在")
		} else {
			return nil, err
		}
	}
	saltPassword := public.GenSaltPasswd(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误，请重新输入")
	}
	return adminInfo, nil
}
func (ad *Admin) Save(ctx *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Save(ad).Error
}
