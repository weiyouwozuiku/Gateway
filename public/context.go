package public

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	TimeLocation *time.Location
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
	LocalIP      = net.ParseIP("127.0.0.1")
)

type TraceContext struct {
	OpenTrace
	CSpanId string
}

func SetGinTraceContext(c *gin.Context, trace *TraceContext) error {
	if trace == nil || c == nil {
		return errors.New("context is nil")
	}
	c.Set("trace", trace)
	return nil
}
func SetTraceContext(ctx context.Context, trace *TraceContext) context.Context {
	if trace == nil {
		return ctx
	}
	return context.WithValue(ctx, "trace", trace)
}
func GetTraceContext(ctx context.Context) *TraceContext {
	if ctx == nil {
		return NewTrace()
	}
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtxTrace, exists := ginCtx.Get("trace")
		if !exists {
			return NewTrace()
		}
		if traceConext, ok := ginCtxTrace.(*TraceContext); ok {
			if traceConext != nil {
				return traceConext
			}
		}
		return NewTrace()
	}
	if traceContext, ok := ctx.Value("trace").(*TraceContext); ok {
		if traceContext != nil {
			return traceContext
		}
	}
	return NewTrace()
}
