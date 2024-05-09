package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project_go/webbook/internal/service"
	"project_go/webbook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	JwtHandler
	svc     wechat.Service
	userSvc service.UserService
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (o *OAuth2WechatHandler) RegisterRoute(server *gin.Engine) {
	wg := server.Group("/wechat/oauth")
	wg.GET("/authurl", o.AuthURL)
	wg.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "构造跳转URL失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
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
	o.setToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}
