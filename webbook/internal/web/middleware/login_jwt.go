package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"project_go/webbook/internal/web"
	"strings"
	"time"
)

type LoginJWTMiddleware struct {
}

func (lm *LoginJWTMiddleware) CheckLoginJWTBuild() gin.HandlerFunc {
	return func(cxt *gin.Context) {
		url := cxt.Request.URL.Path
		if url == "/users/signup" || url == "/users/login" {
			// 不需要校验
			return
		}
		authcode := cxt.GetHeader("Authorization")
		if authcode == "" {
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authSeg := strings.Split(authcode, " ")
		if len(authSeg) != 2 {
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := authSeg[1]
		uc := web.UserClaim{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKEY, nil
		})
		if err != nil {
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			// token解析处理是非法的或者过期的
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != cxt.GetHeader("User-Agent") {
			// 这个地方后面要买点，因为能够进来这个分支的，大概率是攻击者
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 刷新token
		expireTime := uc.ExpiresAt
		// 如果过期时间小于20s则进行刷新
		if expireTime.Sub(time.Now()) > time.Second*20 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 30))
			tokenStr, err = token.SignedString(web.JWTKEY)
			if err != nil {
				//只是刷新时间时间，登录校验成功，所以不中断，只打印日志
				log.Println(err)
			}
			cxt.Header("x-jwt-token", tokenStr)
		}
		// 将user的信息缓存下来，便于profile或者edit接口使用
		cxt.Set("user", uc)

	}
}
