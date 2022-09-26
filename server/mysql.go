package server

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/weiyouwozuiku/Gateway/public"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBMapPool              map[string]*sql.DB
	GORMMapPool            map[string]*gorm.DB
	DBDefaultPool          *sql.DB
	GORMDefaultPool        *gorm.DB
	DefaultMySQLGORMLogger = MySQLGORMLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}
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
	if err := public.ParseConfig(confName, DBConfMap); err != nil {
		return err
	}
	if len(DBConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(public.TimeFormat), " empty mysql config.")
	}
	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for k, v := range DBConfMap.List {
		// 原生db连接池
		sqlDB, err := sql.Open(v.DriverName, v.DataSourceName)
		if err != nil {
			return err
		}
		sqlDB.SetMaxOpenConns(v.MaxOpenConn)
		sqlDB.SetMaxIdleConns(v.MaxIdleConn)
		sqlDB.SetConnMaxLifetime(time.Duration(v.MaxConnLifeTime) * time.Second)
		if err := sqlDB.Ping(); err != nil {
			return err
		}
		// gorm db连接池
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger: &DefaultMySQLGORMLogger,
		})
		if err != nil {
			return err
		}
		if gormSQLDB, err := gormDB.DB(); err != nil {
			return err
		} else {
			if err := gormSQLDB.Ping(); err != nil {
				return err
			}
		}
		// TODO: gorm设置优化
		DBMapPool[k] = sqlDB
		GORMMapPool[k] = gormDB
	}
	// 配置默认数据库连接池
	if db, err := GetDBPool("default"); err != nil {
		DBDefaultPool = db
	}
	if db, err := GetGORMPool("default"); err != nil {
		GORMDefaultPool = db
	}
	return nil
}
func GetDBPool(name string) (*sql.DB, error) {
	if pool, ok := DBMapPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("get DBPool error")
}
func GetGORMPool(name string) (*gorm.DB, error) {
	if pool, ok := GORMMapPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("get GORMPool error")
}
func CloseDB() error {
	for _, db := range DBMapPool {
		if err := db.Close(); err != nil {
			return err
		}
	}
	for _, db := range GORMMapPool {
		if sqlDB, err := db.DB(); err != nil {
			return err
		} else {
			if err := sqlDB.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
func DBPoolQuery(trace *public.TraceContext, sqlDB *sql.DB, query string, args ...any) (*sql.Rows, error) {
	startExecTime := time.Now()
	rows, err := sqlDB.Query(query, args...)
	endExecTime := time.Now()
	if err != nil {
		public.Log.TagError(trace, public.LTagMySQLError, map[string]any{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		public.Log.TagInfo(trace, public.LTagMySQLInfo, map[string]any{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return rows, err
}
