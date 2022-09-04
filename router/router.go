package router

import (
	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/docs"
	"github.com/weiyouwozuiku/gateway/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = middleware.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = middleware.GetStringConf("base.swagger.description")
	docs.SwaggerInfo.Version = middleware.GetStringConf("base.swagger.version")
	docs.SwaggerInfo.Host = middleware.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = middleware.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = middleware.GetStringSliceConf("base.swagger.schemes")
}
