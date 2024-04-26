//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"project_go/webbook/internal/repository"
	"project_go/webbook/internal/repository/cache"
	"project_go/webbook/internal/repository/dao"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/web"
	"project_go/webbook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 三方件
		ioc.InitDb,
		ioc.InitRedis,
		// dao和cache
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
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
