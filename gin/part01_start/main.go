package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Start:https://github.com/gin-gonic/gin
// 字节跳动的hertz（更高性能）
func main() {

	server := gin.Default() // 创建Engine

	// GET请求
	server.GET("/hello", func(c *gin.Context) { // 路由注册
		// 第二个参数为handlers ...HandlerFunc
		//type HandlerFunc func(*Context)
		c.String(http.StatusOK, "hello world")
	})

	// POST请求
	server.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello POST")
	})

	go func() { // 逻辑服务器（Web服务器实例）（此处为一个进程监听两个端口）
		server1 := gin.Default()
		server1.GET("/hi", func(c *gin.Context) {
			c.String(http.StatusOK, "hello")
		})
		server1.Run(":8081")
	}()
	server.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
