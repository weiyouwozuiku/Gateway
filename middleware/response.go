package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/middleware/lib"
)

// 1000以下状态通用码,1000以上为自定义错误码
type ResponseCode int32

const (
	SuccessCode ResponseCode = iota
	UndefErrCode
	ValidErrorCode
	InternalErrorCode
)

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      any          `json:"data"`
	TraceID   any          `json:"trace_id"`
	Stack     any          `json:"stack"`
}

func ResponseSuccess(ctx *gin.Context, data any) {
	trace, _ := ctx.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	resp := &Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  "",
		Data:      data,
		TraceID:   traceId,
		Stack:     nil,
	}
	ctx.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	ctx.Set("response", string(response))
}

func ResponseError(ctx *gin.Context, code ResponseCode, err error) {
	trace, _ := ctx.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	stack := ""
	if ctx.Query("is_debug") == "1" || GetConfEnv() == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}
	resp := &Response{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		Data:      "",
		TraceID:   traceId,
		Stack:     stack,
	}
	ctx.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	ctx.Set("response", string(response))
	ctx.AbortWithError(int(code), err)
}
