package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/public"
)

func IPAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isMatched := false
		for _, host := range public.GetStringSliceConf("base.http.allow_ip") {
			if ctx.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(ctx, InnerErr, fmt.Errorf("%v, not in iplist", ctx.ClientIP()))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
