package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"project_go/webbook/internal/web"
	"time"
)

func main() {
	server := gin.Default()
	/**
	  解决跨域问题
	  跨域问题：只要前端到后台的协议，域名+端口有一个不相同就会存在跨域问题
	  gin提供middleware来解决跨域问题
	  middleware也是面向切面编程AOP的解决方案
	  执行顺序： 请求-->middlwware1-->middleware2--> ... -->业务逻辑
	  middleware接入当时：Engine.Use ，这里Use的参数可以传入任意个HandleFunc
	  HandlerFunc是func的衍生类型：type HandlerFunc func(*Context)
	*/
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//	if strings.Contains(origin, "localhost") {
			//		return true
			//	}
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这里执行一个middleware")
	})
	userHandler := web.NewUserHandler()
	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}
