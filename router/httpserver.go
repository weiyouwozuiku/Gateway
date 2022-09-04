package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/gateway/middleware"
)

var HttpSrvHandler *http.Server

func HttpServerRun() {
	gin.SetMode(middleware.GetStringConf("base.base.debug_mode"))
	HttpSrvHandler = &http.Server{
		Addr:        middleware.GetStringConf("base.base.http_addr"),
		Handler:     InitRouter(),
		ReadTimeout: time.Duration(middleware.GetIntConf("base.base.read_timeout")) * time.Second,
	}
}
