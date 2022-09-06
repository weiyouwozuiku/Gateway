package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/weiyouwozuiku/gateway/middleware/lib"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

type MySQLGormLogger struct {
	LogLevel logger.LogLevel
	// 慢日志阈值
	SlowThreshold time.Duration
	Trace         *lib.TraceContext
}

var (
	DefaultMySQLGormLogger = MySQLGormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}
)

func (l *MySQLGormLogger) LogMode(level logger.LogLevel) logger.Interface {

}

func InitDBPool(path string) error {
	DBConfMap := &MySQLMapConf{}
	err := ParseConfig(path, DBConfMap)
	if err != nil {
		return err
	}
	if len(DBConfMap.Map) == 0 {
		fmt.Printf("[info] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}
	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for confName, DBConf := range DBConfMap.Map {
		dbpool, err := sql.Open("mysql", DBConf.DataSourceName)
		if err != nil {
			return err
		}
		dbpool.SetMaxOpenConns(DBConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DBConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DBConf.MaxConnLifeTime) * time.Second)
		err = dbpool.Ping()
		if err != nil {
			return err
		}
		dbgorm, err := gorm.Open(mysql.Open(DBConf.DataSourceName), &gorm.Config{
			Logger: &DefaultMySQLGormLogger,
		})
		if err != nil {
			return err
		}
		DBMapPool[confName] = dbpool
		GORMMapPool[confName] = dbgorm
	}
	if dbpool, err := GetDBPool("default"); err == nil {
		DBDefaultPool = dbpool
	}
	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}
	return nil
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

func CloseDB() error {
	for _, dbPool := range DBMapPool {
		dbPool.Close()
	}
	DBMapPool = make(map[string]*sql.DB)
	GORMMapPool = make(map[string]*gorm.DB)
	return nil
}
