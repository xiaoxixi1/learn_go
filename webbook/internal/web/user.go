package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/service"
	"time"
	//"regexp" // 官方正则表达式不支持?=的写法
	regexp "github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
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
	nameRegex     *regexp.Regexp
	svc           service.UserService
	codeSvc       service.CodeService
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	nameRegexPattern     = `^[a-zA-Z0-9_]{4,16}$` //这里因为还限制了特殊字符，就暂时不用数据库字段长度来限制了
	userBiz              = "userLogin"
)

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegex:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		nameRegex:     regexp.MustCompile(nameRegexPattern, regexp.None),
		svc:           svc,
		codeSvc:       codeSvc,
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
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
	/**
	  验证码登录增加2个接口:
	    1 发送验证码
	    2 校验验证码登录
	*/
	ug.POST("/login_sms/code/send", h.SmsSendLoginCode)
	ug.POST("/login_sms", h.LoginSms)

}

// 注册用户
func (h *UserHandler) SignUp(cxt *gin.Context) {
	// 在内部定义结构体来接收请求参数
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
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
	// 解决邮箱冲突
	err = h.svc.SignUp(cxt, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		cxt.String(http.StatusOK, "注册成功")
		return
	case service.UserDuplicateError:
		cxt.String(http.StatusOK, "邮箱冲突")
		return
	default:
		cxt.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginJWT(cxt *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := cxt.Bind(&req); err != nil {
		return
	}
	us, err := h.svc.Login(cxt, req.Email, req.Password)
	switch err {
	case nil:
		h.setToken(cxt, us.Id)
		return
	case service.InvalidPasswordOrUser:
		cxt.String(http.StatusOK, "账号或者密码不正确")
		return
	default:
		cxt.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) setToken(cxt *gin.Context, userId int64) {
	us := UserClaim{
		userid:    userId,
		UserAgent: cxt.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 30S后过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, us)
	tokenStr, err := token.SignedString([]byte(JWTKEY))
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	// 自定义头部传输token，前端配合接收
	cxt.Header("x-jwt-token", tokenStr)
	cxt.String(http.StatusOK, "登录成功")
}

// 登录用户
func (h *UserHandler) Login(cxt *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := cxt.Bind(&req); err != nil {
		return
	}
	us, err := h.svc.Login(cxt, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(cxt)
		sess.Set("userid", us.Id)
		sess.Options(sessions.Options{
			// session。Options里面的参数，除了MaxAge有多层含义之后，其他参数就是cookie中那些参数的含义
			// 可以理解位session用这些选项来初始化cookie
			// MaxAge一方面用来控制cookie，也用来控制session中key，value过期时间
			// 有效时间15分钟
			MaxAge: 60,
		})
		err = sess.Save()
		if err != nil {
			cxt.String(http.StatusOK, "系统错误")
			return
		}
		cxt.String(http.StatusOK, "登录成功")
		return
	case service.InvalidPasswordOrUser:
		cxt.String(http.StatusOK, "账号或者密码不正确")
		return
	default:
		cxt.String(http.StatusOK, "系统错误")
	}

}

// 查看用户信息
func (h *UserHandler) Profile(cxt *gin.Context) {
	type ProfileResponse struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		Name            string `json:"name"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"PersonalProfile"`
	}
	/**	sess := sessions.Default(cxt)
	userId, ok := sess.Get("userid").(int64)
	if !ok {
		cxt.String(http.StatusOK, "系统错误")
		return
	}*/
	// 使用jwt
	uc, ok := cxt.MustGet("user").(UserClaim)
	if !ok {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	userId := uc.userid
	user, err := h.svc.Profile(cxt, userId)
	switch err {
	case service.UserNotFoundError:
		cxt.String(http.StatusOK, "该用户不存在")
		return
	case nil:
		cxt.JSON(http.StatusOK, ProfileResponse{
			Email:           user.Email,
			Name:            user.Name,
			Birthday:        user.Birthday.Format(time.DateOnly),
			Password:        user.Password,
			PersonalProfile: user.PersonalProfile,
		})
		return
	default:
		cxt.String(http.StatusOK, "系统错误")
	}
}

// 修改用户信息
func (h *UserHandler) Edit(cxt *gin.Context) {
	sess := sessions.Default(cxt)
	userId, ok := sess.Get("userid").(int64)
	if !ok {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	type EditReq struct {
		Name            string `json:"name"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"PersonalProfile"`
	}
	var req EditReq
	if err := cxt.Bind(&req); err != nil {
		return
	}
	isName, err := h.nameRegex.MatchString(req.Name)
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	if !isName {
		cxt.String(http.StatusOK, "昵称格式不正确")
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		cxt.String(http.StatusOK, "生日格式不正确")
		return
	}
	err = h.svc.Edit(cxt, domain.User{
		Id:              userId,
		Name:            req.Name,
		Birthday:        birthday,
		PersonalProfile: req.PersonalProfile,
	})
	if err != nil {
		cxt.String(http.StatusOK, "更新失败")
		return
	}
	cxt.String(http.StatusOK, "更新成功")

}

// 发送验证码
func (h *UserHandler) SmsSendLoginCode(cxt *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := cxt.Bind(&req); err != nil {
		return
	}
	err := h.codeSvc.SendCode(cxt, userBiz, req.Phone)
	switch err {
	case nil:
		cxt.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
		//cxt.String(http.StatusOK, "发送成功")
	case service.ErrCodeSendTooMany:
		cxt.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		cxt.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		// 补日志的
	}
}

// 验证码登录
func (h *UserHandler) LoginSms(cxt *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := cxt.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		cxt.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "请输入手机号",
		})
		return
	}
	ok, err := h.codeSvc.Verify(cxt, userBiz, req.Phone, req.Code)
	if err != nil {
		cxt.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		cxt.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	/**
	  验证码正确，登录或者用户不存在则创建一个用户
	  但是这里存在并发问题，有可能一个号码，创建2个用户，所以要将phone设为uniqe
	*/
	user, err := h.svc.FindOrCreate(cxt, req.Phone)
	if err != nil {
		cxt.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	h.setToken(cxt, user.Id)
	cxt.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

type UserClaim struct {
	jwt.RegisteredClaims
	userid    int64
	UserAgent string
}

var JWTKEY = []byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8v4")
