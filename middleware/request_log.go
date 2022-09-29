package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

func RequestInLog(ctx *gin.Context) {
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
