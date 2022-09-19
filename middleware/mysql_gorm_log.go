package middleware

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type MySQLGORMLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (mgl *MySQLGORMLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	mgl.LogLevel = logLevel
	return mgl
}
func (mgl *MySQLGORMLogger) Info(ctx context.Context, message string, values ...any) {

}
func (mgl *MySQLGORMLogger) Warn(context.Context, string, ...any)  {}
func (mgl *MySQLGORMLogger) Error(context.Context, string, ...any) {}
func (mgl *MySQLGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}
