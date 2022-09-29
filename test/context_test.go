package test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

func Test_GetTraceContext(t *testing.T) {
	ctx := context.Background()
	trace := public.GetTraceContext(ctx)
	ctx = public.SetTraceContext(ctx, trace)
	if v, exists := ctx.Value(public.TraceKey).(*public.TraceContext); exists {
		t.Logf("%#v", v)
	}
}
func Test_GetGinTraceContext(t *testing.T) {
	ctx := &gin.Context{}
	trace := public.GetTraceContext(ctx)
	public.SetGinTraceContext(ctx, trace)
	if v, exists := ctx.Get(public.TraceKey); exists {
		t.Logf("%#v", v)
	}
}
