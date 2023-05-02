package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/weiyouwozuiku/Gateway/public"
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
	trace := public.GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	public.Log.TagInfo(trace, public.LTagMySQLInfo, params)
}
func (mgl *MySQLGORMLogger) Warn(ctx context.Context, message string, values ...any) {
	trace := public.GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	public.Log.TagWarn(trace, public.LTagMySQLWarn, params)
}
func (mgl *MySQLGORMLogger) Error(ctx context.Context, message string, values ...any) {
	trace := public.GetTraceContext(ctx)
	params := make(map[string]any)
	params["message"] = message
	params["values"] = fmt.Sprint(values...)
	public.Log.TagError(trace, public.LTagMySQLError, params)
}
func (mgl *MySQLGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	trace := public.GetTraceContext(ctx)
	if mgl.LogLevel <= logger.Silent {
		return
	}
	sqlStr, rows := fc()
	currentTime := begin.Format(public.TimeFormat)
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
			public.Log.TagInfo(trace, public.LTagMySQLFailed, msg)
		} else {
			msg["rows"] = rows
			public.Log.TagInfo(trace, public.LTagMySQLFailed, msg)
		}
	case elapsed > mgl.SlowThreshold && mgl.SlowThreshold != 0 && mgl.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
		msg["slowLog"] = slowLog
		if rows == -1 {
			public.Log.TagInfo(trace, public.LTagMySQLSlow, msg)
		} else {
			msg["rows"] = rows
			public.Log.TagInfo(trace, public.LTagMySQLSlow, msg)
		}
	case mgl.LogLevel == logger.Info:
		if rows == -1 {
			public.Log.TagInfo(trace, public.LTagMySQLSuccess, msg)
		} else {
			msg["rows"] = rows
			public.Log.TagInfo(trace, public.LTagMySQLSuccess, msg)
		}
	}
}
