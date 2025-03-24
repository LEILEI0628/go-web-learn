package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.New()
	// 静态路由
	server.GET("/static", func(context *gin.Context) {
		// 重定向问题："/static/"重定向为BASE_URL:PORT/static/
		context.String(http.StatusOK, "这是静态路由")
	})
	// 测试用例：
	// http://localhost:8080/static *PASS
	// http://localhost:8080/static/ *PASS
	// http://localhost:8080/static/a *Fail

	// 参数路由
	server.GET("/user/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "hello,%s", name)
	})
	// 测试用例：
	// http://localhost:8080/user/LEILEI *PASS

	// 通配符匹配 TODO:测试用例2的匹配问题和Param参数问题
	server.GET("/views/*.html", func(context *gin.Context) {
		page := context.Param(".html") // “*”不能单独出现，如：/views/*，/views/*/page等
		context.String(http.StatusOK, "匹配值为：%s", page)
	})
	// 测试用例：
	// http://localhost:8080/views/index.html *PASS
	// http://localhost:8080/views/aaa *PASS(?????)

	// 查询参数
	server.GET("/order", func(context *gin.Context) {
		id := context.Query("id") // 查询参数使用context.Query
		context.String(http.StatusOK, "id:%s", id)
	})
	// 测试用例：
	// http://localhost:8080/order?id=12345 *PASS

	server.Run(":8080")
}
