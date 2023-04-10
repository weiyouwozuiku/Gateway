package initialize

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/weiyouwozuiku/Gateway/controller"
	"github.com/weiyouwozuiku/Gateway/docs"
	"github.com/weiyouwozuiku/Gateway/middleware"
	"github.com/weiyouwozuiku/Gateway/public"
)

var SessionKey = []byte("secret")

const AdminSession = "adminSession"

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = public.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = public.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Host = public.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = public.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Version = public.GetStringConf("base.swagger.version")
	docs.SwaggerInfo.Schemes = public.GetStringSliceConf("base.swagger.schemes")

	var router *gin.Engine
	if public.ConfBase.DebugMode == "debug" {
		router = gin.Default()
	} else {
		router = gin.New()
	}
	router.Use(middlewares...)

	// 登录session存放redis
	store, err := sessions.NewRedisStore(10, "tcp", public.GetStringConf("base.session.redis_server"), public.GetStringConf("base.session.redis_password"), SessionKey)
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}

	// 探活接口
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// admin_login
	adminLoginRouter := router.Group("/admin_login")

	adminLoginRouter.Use(
		sessions.Sessions(AdminSession, store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.ValidtorMiddleware(),
	)
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}

	// admin
	adminRouter := router.Group("/admin")
	adminRouter.Use(
		sessions.Sessions(AdminSession, store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.ValidtorMiddleware(),
	)
	{
		controller.AdminRegister(adminRouter)
	}

	serviceRouter := router.Group("/service")
	serviceRouter.Use(
		sessions.Sessions(AdminSession, store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.ValidtorMiddleware())
	{
		controller.ServiceRegister(serviceRouter)
	}
	// TODO 后续增加router
	return router
}
