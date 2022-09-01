package middleware

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

func InitDBPool(path string) error {
	//普通DB方式
	DbConfMap := &MysqlMapConf{}
	err := ParseConfig(path, DbConfMap)
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get db pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get db pool error")
}
