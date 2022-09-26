package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/weiyouwozuiku/Gateway/docs"
	"github.com/weiyouwozuiku/Gateway/public"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = public.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = public.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Host = public.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = public.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Version = public.GetStringConf("base.swagger.version")
	docs.SwaggerInfo.Schemes = public.GetStringSliceConf("base.wagger.schemes")

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// TODO 后续增加router
	return router
}
