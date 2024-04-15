package main

import (
	"github.com/gin-gonic/gin" // 每次引入新的依赖之后，接的执行go mod tidy，来保证实际依赖和go.mod一致，保持依赖的整洁性
	"net/http"
)

func main() {
	// 创建一个Engine,相当于一个web服务器，监听一个端口,一个go进程可以创建多个Engine
	server := gin.Default()
	// 注册静态路由
	server.GET("/hello", func(cxt *gin.Context) {
		//context表示上下文，作用是处理请求，返回响应
		cxt.String(http.StatusOK, "hello world") // 返回响应
	})
	// 注册参数路由,路径参数
	server.GET("/hello/:name", func(cxt *gin.Context) {
		name := cxt.Param("name")
		cxt.String(http.StatusOK, "hello "+name)
	})
	// 查询参数
	// GET /order?id=1
	server.GET("/order", func(context *gin.Context) {
		id := context.Query("id")
		context.String(http.StatusOK, "订单id是 "+id)
	})
	// 通配符路由,注意的是*不能单独出现，不能注册路由/star/*
	server.GET("/view/*.html", func(context *gin.Context) {
		view := context.Param(".html")
		context.String(http.StatusOK, "view是"+view)
	})
	server.Run(":8080")
}
