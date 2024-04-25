package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project_go/webbook/internal/config"
	"project_go/webbook/internal/repository"
	"project_go/webbook/internal/repository/cache"
	"project_go/webbook/internal/repository/dao"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/service/sms/localmodel"
	"project_go/webbook/internal/web"
	"project_go/webbook/pkg/ginx/middleware/ratelimit"

	//"project_go/webbook/internal/repository"
	//"project_go/webbook/internal/repository/dao"
	//"project_go/webbook/internal/service"
	//"project_go/webbook/internal/web"
	"project_go/webbook/internal/web/middleware"
	"time"
)

func main() {
	db := InitDb()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server := InitWebServer(redisClient)
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(redisClient)
	cc := cache.NewCodeCache(redisClient)
	ur := repository.NewUseRepository(ud, uc)
	cr := repository.NewCodeRepo(cc)
	us := service.NewUserService(ur)
	localSms := localmodel.NewService()
	cs := service.NewCodeService(cr, localSms)
	userHandler := web.NewUserHandler(us, cs)
	userHandler.RegisterRoutes(server)
	//server := gin.Default()
	//server.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello，启动成功了")
	//})
	server.Run(":8080")
}

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic("failed to connect database")
	}
	dao.InitTables(db)
	return db
}

func InitWebServer(redisClient redis.Cmdable) *gin.Engine {
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
	//做压测去掉限流
	//每秒100个请求，之所以使用redis进行限流，是针对集群多实例的场景
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
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
	})
	useJWT(server)
	return server
}

/**
  从这里session的存储方法可以发现，这里存储session数据被抽象成了一个接口，
  GIN提供了不同的实现，所以就可以自由切换了
  所以当设计核心系统的时候，或者打算提供什么功能给用户的时候，一定要问问自己，将来有没有可能需要不同的实现
*/

func useSession(server *gin.Engine) {
	// 创建两个middleware,一个初始化session,一个取session检验是否登录
	// 基于cookie存储session数据
	store := cookie.NewStore([]byte("secret"))
	// 基于memstore存储session数据
	// 第一个参数是authentication key，最好是32位或者64位,是指身份认证
	// 第二个参数是encryption key，是指数据加密
	//store := memstore.NewStore([]byte("FBP932qsW4e3STABJLbKAjpGS5EVg6sqNTLHXw5CCLYkHp84SUxFAnh3zfUfRmVd"), []byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8v3pnVdfLAXR2JSY5fuhJZVGNgK5e9hc9Gh"))
	//基于redis存储session数据
	// 第一个参数是：最大连接数
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379", "",
	//	[]byte("FBP932qsW4e3STABJLbKAjpGS5EVg6sqNTLHXw5CCLYkHp84SUxFAnh3zfUfRmVd"),
	//	[]byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8v3"))
	//if err != nil {
	//	panic(err)
	//}
	/**
	cookie,memstore,redis三个实现中，都需要传入两个key：
	    	authentication key：身份认证
	        encryption key：数据加密
	  这两个key再加上授权（权限控制）就是信息安全的三个核心概念
	*/
	server.Use(sessions.Sessions("ssid", store))
	login := &middleware.LoginMiddleware{}
	server.Use(login.CheckLoginBuild())
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddleware{}
	server.Use(login.CheckLoginJWTBuild())
}
