package middleware_test

import (
	"context"
	"testing"

	"github.com/weiyouwozuiku/Gateway/middleware"

	"github.com/gin-gonic/gin"
)

func Test_GetTraceContext(t *testing.T) {
	ctx := context.Background()
	trace := middleware.GetTraceContext(ctx)
	ctx = middleware.SetTraceContext(ctx, trace)
	if v, exists := ctx.Value("trace").(*middleware.TraceContext); exists {
		t.Logf("%#v", v)
	}
}
func Test_GetGinTraceContext(t *testing.T) {
	ctx := &gin.Context{}
	trace := middleware.GetTraceContext(ctx)
	middleware.SetGinTraceContext(ctx, trace)
	if v, exists := ctx.Get("trace"); exists {
		t.Logf("%#v", v)
	}
}
