package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

// RecoveryMiddleware捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 日志记录
				fmt.Println(string(debug.Stack()))
				public.ComLogWarning(ctx, public.LTagPanic, map[string]any{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})
				if public.ConfBase.DebugMode != "debug" {
					ResponseError(ctx, InnerErr, errors.New("内部错误"))
				} else {
					ResponseError(ctx, InnerErr, errors.New(fmt.Sprint(err)))
				}
			}
		}()
		ctx.Next()
	}
}
