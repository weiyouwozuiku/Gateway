package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

type ResponseCode int

const (
	SuccessCode ResponseCode = iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000

	InvalidParamsCode       ResponseCode = 2000
	GetGormPoolFailed       ResponseCode = 2001
	GROUPALL_SAVE_FLOWERROR ResponseCode = 2002
	AdminLoginFailed        ResponseCode = 2003
	SessionParseFailed      ResponseCode = 2004
	GormQueryFailed         ResponseCode = 2005
	GormSaveFailed          ResponseCode = 2006
	SessionOptFailed        ResponseCode = 2207
)

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      any          `json:"data"`
	TraceId   any          `json:"trace_id"`
	Stack     any          `json:"stack"`
}

func ResponseError(ctx *gin.Context, code ResponseCode, err error) {
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
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		Data:      "",
		TraceId:   traceId,
		Stack:     stack,
	}
	ctx.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	ctx.Set("response", string(response))
	ctx.AbortWithError(http.StatusOK, err)
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
