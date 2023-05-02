package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

const ()

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      any          `json:"data"`
	TraceId   any          `json:"trace_id"`
	Stack     any          `json:"stack"`
}

func ResponseError(ctx *gin.Context, info Error, err error) {
	trace, _ := ctx.Get(public.TraceKey)
	traceContext, _ := trace.(*public.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	stack := ""
	if ctx.Query("is_debug") == "1" || public.ConfEnv == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}
	resp := &Response{
		ErrorCode: info.errno,
		ErrorMsg:  fmt.Sprintf("%s||error=%v", info.errmsg, err),
		Data:      "",
		TraceId:   traceId,
		Stack:     stack,
	}
	ctx.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	ctx.Set("response", string(response))
	if err = ctx.AbortWithError(http.StatusOK, err); err != nil {
		public.ComLogError(ctx, public.LTagRespErr, map[string]any{"error": err})
	}
}

func ResponseSuccess(ctx *gin.Context, data any) {
	trace, _ := ctx.Get(public.TraceKey)
	traceContext, _ := trace.(*public.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	resp := &Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  "",
		Data:      data,
		TraceId:   traceId,
	}
	ctx.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	ctx.Set("response", string(response))
}
