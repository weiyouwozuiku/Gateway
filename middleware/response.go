package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	traceContext, _ := trace.(*TraceContext)
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

func ResponseError(ctx *gin.Context, data any) {

}