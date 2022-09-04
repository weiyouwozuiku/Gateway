package middleware

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

//func InitDBPool(path string) error {
//	DBConfMap := &MysqlMapConf{}
//	err := ParseConfig(path, DBConfMap)
//	if err != nil {
//		return err
//	}
//	if len(DBConfMap.List) == 0 {
//		fmt.Printf("[info] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
//	}
//	DBMapPool = map[string]*sql.DB{}
//	GORMMapPool = map[string]*gorm.DB{}
//	for confName, DBConf := range DBConfMap.List {
//		dbpool, err := sql.Open("mysql", DBConf.DataSourceName)
//		if err != nil {
//			return err
//		}
//		dbpool.SetMaxOpenConns(DBConf.MaxOpenConn)
//		dbpool.SetMaxIdleConns(DBConf.MaxIdleConn)
//		dbpool.SetConnMaxLifetime(time.Duration(DBConf.MaxConnLifeTime) * time.Second)
//		err = dbpool.Ping()
//		if err != nil {
//			return err
//		}
//		dbgorm, err := gorm.Open(mysql.Open(DBConf.DataSourceName), &gorm.Config{
//			NamingStrategy: schema.NamingStrategy{
//				SingularTable: true},
//				Logger: ,
//		})
//		db, _ := dbgorm.DB()
//		err = db.Ping()
//		if err != nil {
//			return err
//		}
//
//	}
//}

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

type MySQLGormLogger struct {
}
