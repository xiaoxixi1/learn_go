package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type JwtHandler struct {
	refreshMethod jwt.SigningMethod
	refreshKey    []byte
}

func NewJwtHandler() *JwtHandler {
	return &JwtHandler{
		refreshMethod: jwt.SigningMethodHS256,
		refreshKey:    []byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8vh"),
	}
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
}

func (h *JwtHandler) setRefreshToken(cxt *gin.Context, userId int64) {
	rs := RefreshClaim{
		userId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			// 长token7天后过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(h.refreshMethod, rs)
	tokenStr, err := token.SignedString(h.refreshKey)
	if err != nil {
		cxt.String(http.StatusOK, "系统错误")
		return
	}
	// 自定义头部传输token，前端配合接收
	cxt.Header("x-jwt-refresh-token", tokenStr)
}

func ExtractToken(cxt *gin.Context) string {
	authcode := cxt.GetHeader("Authorization")
	if authcode == "" {
		cxt.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	authSeg := strings.Split(authcode, " ")
	if len(authSeg) != 2 {
		cxt.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return authSeg[1]
}

type UserClaim struct {
	jwt.RegisteredClaims
	userid    int64
	UserAgent string
}

type RefreshClaim struct {
	jwt.RegisteredClaims
	userId int64
}

var JWTKEY = []byte("Bhy3mfsThsmBvfpNwyCF2FEzS4GfR8v4")
