package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/dto"
	"gorm.io/gorm"
)

type App struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *App) TableName() string {
	return "gateway_app"
}

func (t *App) Find(c *gin.Context, tx *gorm.DB, search *App) (*App, error) {
	model := &App{}
	err := tx.WithContext(c).Where(search).Find(model).Error
	return model, err
}

func (t *App) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *App) AppList(c *gin.Context, tx *gorm.DB, param *dto.AppListInput) ([]App, int64, error) {
	var (
		list  []App
		count int64
	)
	pageNo := param.PageNo
	pageSize := param.PageSize
	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.WithContext(c).Table(t.TableName()).Select("*").Where("id_delete=?", 0)
	if param.Info != "" {
		query = query.Where(" (name like ? or app_id like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, errCount
	}
	return list, count, nil
}
