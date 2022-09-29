package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

func requestInLog(ctx *gin.Context) {
	traceContext := public.NewTrace()
	if traceId := ctx.Request.Header.Get("com-header-rid"); traceId != "" {
		traceContext.TraceId = traceId
	}
	if spanId := ctx.Request.Header.Get("com-header-spanid"); spanId != "" {
		traceContext.SpanId = spanId
	}
	ctx.Set("startExecTime", time.Now())
	ctx.Set("trace", traceContext)
	bodyBytes, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	public.Log.TagInfo(traceContext, public.LTagRequestIn, map[string]any{
		"uri":    ctx.Request.RequestURI,
		"method": ctx.Request.Method,
		"args":   ctx.Request.PostForm,
	})
}
func requestOutLog(ctx *gin.Context) {
	endExecTime := time.Now()
	response, _ := ctx.Get("response")
	st, _ := ctx.Get("startExecTime")
	startExecTime, _ := st.(time.Time)
	public.ComLogNotice(ctx, public.LTagRequestOut, map[string]any{
		"uri":       ctx.Request.RequestURI,
		"method":    ctx.Request.Method,
		"args":      ctx.Request.PostForm,
		"from":      ctx.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	})
}
func RequestLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if public.GetBoolConf("base.log.file_writer.on") {
			requestInLog(ctx)
			defer requestOutLog(ctx)
		}
		ctx.Next()
	}
}
