package middleware

import (
	"errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if adminInfo, ok := session.Get(public.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			ResponseError(ctx, InternalErrorCode, errors.New("user not login"))
			ctx.Abort()
			return
		}
		// TODO 可以增加更细致的权限管理
		ctx.Next()
	}
}
