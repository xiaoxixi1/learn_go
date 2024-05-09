package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JwtHandler struct {
}

func (h *JwtHandler) setToken(cxt *gin.Context, userId int64) {
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

type UserClaim struct {
	jwt.RegisteredClaims
	userid    int64
	UserAgent string
}

var JWTKEY = []byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8v4")
