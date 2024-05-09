package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	jwtHandler      *JwtHandler
	svc             wechat.Service
	userSvc         service.UserService
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
		jwtHandler:      NewJwtHandler(),
	}
}

func (o *OAuth2WechatHandler) RegisterRoute(server *gin.Engine) {
	wg := server.Group("/wechat/oauth")
	wg.GET("/authurl", o.AuthURL)
	wg.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "构造跳转URL失败",
		})
		return
	}
	err = o.setStateToJWTToken(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}

func (o *OAuth2WechatHandler) setStateToJWTToken(ctx *gin.Context, state string) error {
	stateClaim := StateClaim{
		State: state,
	}
	jwtState := jwt.NewWithClaims(jwt.SigningMethodHS256, stateClaim)
	tokenStr, err := jwtState.SignedString([]byte(o.key))
	if err != nil {

		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)
	return nil
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyStateFromJWTToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 400,
		})
		return
	}
	code := ctx.Query("code")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "微信授权失败",
			Code: 500,
		})
		return
	}
	user, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 500,
		})
		return
	}
	o.jwtHandler.setToken(ctx, user.Id)
	o.jwtHandler.setRefreshToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

func (o *OAuth2WechatHandler) verifyStateFromJWTToken(ctx *gin.Context) error {
	state := ctx.Query("state")
	cookie, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return err
	}
	var stateClaim StateClaim
	_, err = jwt.ParseWithClaims(cookie, &stateClaim, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return err
	}
	if state != stateClaim.State {
		return fmt.Errorf("state 不匹配")
	}
	return nil
}

type StateClaim struct {
	jwt.RegisteredClaims
	State string
}
