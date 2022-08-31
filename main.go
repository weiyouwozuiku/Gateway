package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	pprof.Register(r)
	// 定义一个路径为 /ping 的 GET 格式路由，并返回 JSON 数据
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World !",
		})
	})
	r.Run(":8080") // 启动服务，并监听 8000 端口
}
