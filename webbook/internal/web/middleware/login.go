package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleware struct {
}

func (lm *LoginMiddleware) CheckLoginBuild() gin.HandlerFunc {
	return func(cxt *gin.Context) {
		url := cxt.Request.URL.Path
		if url == "/users/signup" || url == "/users/login" {
			// 不需要校验
			return
		}
		sess := sessions.Default(cxt)
		if sess.Get("userid") == nil {
			// 说明没有登录,中断，不再执行后面的业务
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
