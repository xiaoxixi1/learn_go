package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"project_go/webbook/integration/startup"
	"project_go/webbook/internal/web"
	"testing"
	"time"
)

/**
集成测试单独放在一个包里面
*/
// 可以少输出很多debug日志
func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandle_sendCode(t *testing.T) {
	rdb := startup.InitRedis()
	service := startup.InitWebServer()
	testCases := []struct {
		name string

		phone string

		wantCode int
		wantBody web.Result
		before   func(t *testing.T)
		after    func(t *testing.T)
	}{
		{
			name: "发送成功的用例",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:userLogin:12345"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "12345",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name:  "发送消息过于频繁",
			phone: "12345",
			before: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				//提前准备一条过期的数据
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				// 删除数据
				_, err := rdb.Del(ctx, key).Result()
				assert.NoError(t, err)

			},
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 400,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			name:  "系统错误",
			phone: "12345",
			before: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				//提前准备一条过期的数据
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
				//dur, err := rdb.TTL(ctx, key).Result()

			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				// 删除数据
				_, err := rdb.Del(ctx, key).Result()
				assert.NoError(t, err)

			},
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// 准备请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			//执行请求
			service.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			var respons web.Result
			json.NewDecoder(recorder.Body).Decode(&respons)
			assert.Equal(t, tc.wantBody, respons)
		})
	}
}
