package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"project_go/webbook/internal/web"
	"project_go/webbook/internal/web/middleware"
	"project_go/webbook/pkg/ginx/middleware/ratelimit"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHand *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHand.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowHeaders:     []string{"authorization", "Content-Type"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				//	if strings.Contains(origin, "localhost") {
				//		return true
				//	}
				return true
			},
			// 因为自定义了头部，所以跨域要加上配置
			// 约定后端放入x-jwt-token中，前端通过authorization 传回后端
			ExposeHeaders: []string{"x-jwt-token"},
			MaxAge:        12 * time.Hour,
		}), func(ctx *gin.Context) {
			println("这里执行一个middleware")
		},
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		(&middleware.LoginJWTMiddleware{}).CheckLoginJWTBuild(),
	}
}
