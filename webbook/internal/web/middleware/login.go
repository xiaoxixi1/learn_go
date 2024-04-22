package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddleware struct {
}

func (lm *LoginMiddleware) CheckLoginBuild() gin.HandlerFunc {
	return func(cxt *gin.Context) {
		// 注册一下这个类型
		gob.Register(time.Now())
		url := cxt.Request.URL.Path
		if url == "/users/signup" || url == "/users/login" {
			// 不需要校验
			return
		}
		sess := sessions.Default(cxt)
		userId := sess.Get("userid")
		if userId == nil { // 说明没有登录,中断，不再执行后面的业务
			cxt.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		/**
		  加入刷新策略是1分钟1次，则需要定义一个字段记录刷新时间
		*/
		const updateTimeKey = "UpdateTime"
		updateTime := sess.Get(updateTimeKey)
		lastUpdateTime, ok := updateTime.(time.Time)
		if updateTime == nil || !ok || time.Now().Sub(lastUpdateTime) > time.Second*10 {
			//第一次进来，或者时间大于1分钟时，都重新刷新时间
			sess.Set(updateTimeKey, time.Now())
			sess.Set("userid", userId) // 因为set之后会重新生成一个，所以这里得重新set
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
