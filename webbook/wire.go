//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"project_go/webbook/internal/repository"
	"project_go/webbook/internal/repository/cache/freecache"
	"project_go/webbook/internal/repository/cache/redis"
	"project_go/webbook/internal/repository/dao"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/web"
	"project_go/webbook/ioc"
)

/*
*

	这里有一个冲突的点是：在使用wire的时候，初始化方法NewXXXX最好返回接口
	但是go的推荐做法是返回具体类型，这和wire是冲突的
*/
func InitWebServer() *gin.Engine {
	wire.Build(
		// 三方件
		ioc.InitDb,
		ioc.InitRedis,
		ioc.InitFreeCache,
		// dao和cache
		dao.NewUserDao,
		freecache.NewCodeFreeCache,
		redis.NewUserCache,
		//repository
		repository.NewUseRepository,
		repository.NewCodeRepo,
		//service
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSmsSendService,
		// Handler
		web.NewUserHandler,
		// middleware
		ioc.InitGinMiddlewares,
		//web
		ioc.InitWebServer,
	)
	return &gin.Engine{}
}
