package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project_go/webbook/internal/repository"
	"project_go/webbook/internal/repository/dao"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/web"
	"time"
)

func main() {
	db := InitDb()
	server := InitWebServer()
	ud := dao.NewUserDao(db)
	ur := repository.NewUseRepository(ud)
	us := service.NewUserService(ur)
	userHandler := web.NewUserHandler(us)
	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:Jike1504240602*@tcp(localhost:13306)/webook"))
	if err != nil {
		panic("failed to connect database")
	}
	dao.InitTables(db)
	return db
}

func InitWebServer() *gin.Engine {
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
	return server
}
