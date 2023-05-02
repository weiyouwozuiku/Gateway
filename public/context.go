package public

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

type TraceContext struct {
	OpenTrace
	CSpanId string
}

func SetGinTraceContext(c *gin.Context, trace *TraceContext) error {
	if trace == nil || c == nil {
		return errors.New("context is nil")
	}
	c.Set(TraceKey, trace)
	return nil
}
func GetGinTraceContext(c *gin.Context) *TraceContext {
	if c == nil {
		return NewTrace()
	}
	traceContext, exiests := c.Get(TraceKey)
	if exiests {
		if tc, ok := traceContext.(*TraceContext); ok {
			return tc
		}
	}
	return NewTrace()
}
func SetTraceContext(ctx context.Context, trace *TraceContext) context.Context {
	if trace == nil {
		return ctx
	}
	return context.WithValue(ctx, TraceKey, trace)
}
func GetTraceContext(ctx context.Context) *TraceContext {
	if ctx == nil {
		return NewTrace()
	}
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtxTrace, exists := ginCtx.Get(TraceKey)
		if !exists {
			return NewTrace()
		}
		if traceContext, ok := ginCtxTrace.(*TraceContext); ok {
			if traceContext != nil {
				return traceContext
			}
		}
		return NewTrace()
	}
	if traceContext, ok := ctx.Value(TraceKey).(*TraceContext); ok {
		if traceContext != nil {
			return traceContext
		}
	}
	return NewTrace()
}
