package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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
	trace := GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	Log.TagInfo(trace, LTagMySQLInfo, params)
}
func (mgl *MySQLGORMLogger) Warn(ctx context.Context, message string, values ...any) {
	trace := GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	Log.TagWarn(trace, LTagMySQLWarn, params)
}
func (mgl *MySQLGORMLogger) Error(ctx context.Context, message string, values ...any) {
	trace := GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	Log.TagError(trace, LTagMySQLError, params)
}
func (mgl *MySQLGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	trace := GetTraceContext(ctx)
	if mgl.LogLevel <= logger.Silent {
		return
	}
	sqlStr, rows := fc()
	currentTime := begin.Format(TimeFormat)
	elapsed := time.Since(begin)
	msg := map[string]any{
		"FileWithLineNum": utils.FileWithLineNum(),
		"sql":             sqlStr,
		"rows":            "-",
		"proc_time":       float64(elapsed.Milliseconds()),
		"current_time":    currentTime,
	}
	switch {
	case err != nil && mgl.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		msg["err"] = err
		if rows == -1 {
			Log.TagInfo(trace, LTagMySQLFailed, msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, LTagMySQLFailed, msg)
		}
	case elapsed > mgl.SlowThreshold && mgl.SlowThreshold != 0 && mgl.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
		msg["slowLog"] = slowLog
		if rows == -1 {
			Log.TagInfo(trace, LTagMySQLSlow, msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, LTagMySQLSlow, msg)
		}
	case mgl.LogLevel == logger.Info:
		if rows == -1 {
			Log.TagInfo(trace, LTagMySQLSuccess, msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, LTagMySQLSuccess, msg)
		}
	}
}
