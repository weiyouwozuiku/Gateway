package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	// 定义一个路径为 /ping 的 GET 格式路由，并返回 JSON 数据
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World !",
		})
	})
	r.Run() // 启动服务，并监听 8080 端口
}
