package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/service"

	//"regexp" // 官方正则表达式不支持?=的写法
	regexp "github.com/dlclark/regexp2"
)

/*
*

		将所有用户的路由请求定义在UserHandler上
	    同时定义一个RegisterRoutes来注册所有的路由
*/
type UserHandler struct {
	// 使用正则表达式预编译来提高性能
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           *service.UserService
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
	}
}

// 专门用于注册路由
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	/**
	server.POST("/users/signup", h.SignUp)
	server.POST("/users/login", h.Login)
	server.GET("/users/profile", h.Profile)
	server.POST("/users/edit", h.Edit)
	*/
	// 上面可以简化成分组路由
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
}

// 注册用户
func (h *UserHandler) SignUp(cxt *gin.Context) {
	// 在内部定义结构体来接收请求参数
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:confirmPassword`
	}
	var req SignUpReq
	// Bind会根据http的Content-type来处理，如果请求是json格式，content-type就是application/json,gin就会使用json反序列化
	// 如果格式不正确，Bind会自动返回一个400的错误码
	if err := cxt.Bind(&req); err != nil {
		return
	}
	// 校验参数

	//isEmail, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	//用预编译之后就要换一种写法了
	isEmail, err := h.emailRegex.MatchString(req.Email)
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}

	if !isEmail {
		cxt.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	isPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}

	if !isPassword {
		cxt.String(http.StatusOK, "密码必须包含字母，数字和特殊字符，并且长度不能小于8")
		return
	}

	if req.Password != req.ConfirmPassword {
		cxt.String(http.StatusOK, "两次输入密码不一致")
		return
	}
	err = h.svc.SignUp(cxt, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	cxt.String(http.StatusOK, "注册成功")
}

// 登录用户
func (h *UserHandler) Login(cxt *gin.Context) {

}

// 查看用户信息
func (h *UserHandler) Profile(cxt *gin.Context) {

}

// 修改用户信息
func (h *UserHandler) Edit(cxt *gin.Context) {

}
