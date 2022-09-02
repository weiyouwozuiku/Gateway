package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/middleware"
)

var HttpSrvHandler *http.Server

func HttpServerRun() {
	gin.SetMode(middleware.GetStringConf("base.base.debug_mode"))

}
