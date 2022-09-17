package middleware

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	DBMapPool   map[string]*sql.DB
	GORMMapPool map[string]*gorm.DB
)

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

type MySQLMapConf struct {
	List map[string]*MySQLConf `mapstructure:"list"`
}

func InitDBConf(confName string) error {
	DBConfMap := &MySQLMapConf{}
	if err := ParseConfig(confName, DBConfMap); err != nil {
		return err
	}
	if len(DBConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}
	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for k, v := range DBConfMap.List {
		// 1. 创建mysql连接池
		mysqlPool, err := sql.Open("mysql", v.DataSourceName)
		if err != nil {
			return err
		}
		mysqlPool.SetMaxOpenConns(v.MaxOpenConn)
		mysqlPool.SetMaxIdleConns(v.MaxIdleConn)
		mysqlPool.SetConnMaxLifetime(time.Duration(v.MaxConnLifeTime) * time.Second)
		if err := mysqlPool.Ping(); err != nil {
			return err
		}
		// 2. gorm使用已有连接池创建
		gormPool, err := gorm.Open()
	}
	return nil
}
