package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var HttpSrvHandler *http.Server

func HttpServerRun() {
	gin.SetMode(Get)
}
