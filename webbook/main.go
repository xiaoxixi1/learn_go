package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	//"project_go/webbook/internal/repository"
	//"project_go/webbook/internal/repository/dao"
	//"project_go/webbook/internal/service"
	//"project_go/webbook/internal/web"
	"project_go/webbook/internal/web/middleware"
)

func main() {
	server := InitWebServer()
	server.Run(":8080")
}

//func InitWebServer(redisClient redis.Cmdable) *gin.Engine {
//	server := gin.Default()
//
//	/**
//	  解决跨域问题
//	  跨域问题：只要前端到后台的协议，域名+端口有一个不相同就会存在跨域问题
//	  gin提供middleware来解决跨域问题
//	  middleware也是面向切面编程AOP的解决方案
//	  执行顺序： 请求-->middlwware1-->middleware2--> ... -->业务逻辑
//	  middleware接入当时：Engine.Use ，这里Use的参数可以传入任意个HandleFunc
//	  HandlerFunc是func的衍生类型：type HandlerFunc func(*Context)
//	*/
//	//做压测去掉限流
//	//每秒100个请求，之所以使用redis进行限流，是针对集群多实例的场景
//	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//
//	useJWT(server)
//	return server
//}

//func InitUserHandler(db *gorm.DB, redisClient redis.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	uc := cache.NewUserCache(redisClient)
//	cc := cache.NewCodeCache(redisClient)
//	ur := repository.NewUseRepository(ud, uc)
//	cr := repository.NewCodeRepo(cc)
//	us := service.NewUserService(ur)
//	localSms := localmodel.NewService()
//	cs := service.NewCodeService(cr, localSms)
//	return web.NewUserHandler(us, cs)
//}

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

//func useJWT(server *gin.Engine) {
//	login := &middleware.LoginJWTMiddleware{}
//	server.Use(login.CheckLoginJWTBuild())
//}
