package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/service"
	svcmocks "project_go/webbook/internal/service/mocks"
	"testing"
)

/**
单元测试：针对每一个方法进行测试，单独验证每一个方法的正确性
        讲究快速测试，快速修复
        测试该环节中的业务问题
        测试该环境中的技术问题
        单元测试从理论上讲，不能依赖任何第三方组件
集成测试：多个组件合并在一起的测试，验证各个方法，组件之间配合无误
go中编写单元测试：
  1 文件名以_test.go结尾
  2 测试方法以Test开头
  3 测试方法值接收一个参数： t *testing.T
同样可以用这个来写集成测试，冒烟测试，回归测试
Table Driven 模式：go里面惯用的组织测试的方式，都是用Table Driven
   Table Driven 的形式主要分成三个部分：
     测试用例的定义：每一个测试用例需要有什么
     具体的测试用例：设计的每一个测试用例都在这里
     执行测试用例：这里面还包括了对测试结果进行断言
mock工具，分成两个部分：
   mockgen：命令行工具
   测试中使用的控制mock对象的包
*/

func TestUserEmailPattern(t *testing.T) {
	testCases := []struct {
		// 测试用例定义
		name string // 用例的名字
		//预期输入，根据方法参数，接收器设计
		email string
		// 预期删除，根据方法返回值，接收器来设计
		match bool
		// mock 数据，在单元测试里面很常见,集成测试一般没有
		//mock func(ctrl *gmock.)
		// 测试用例准备环境，数据等
		//before func(t *testing.T)
		// 每个测试用里执行完后数据清理等
		//after func(t *testing.T)
	}{
		// 具体的测试用例
		{
			name:  "不带@",
			email: "123456",
			match: false,
		},
		{
			name:  "带@不带后缀",
			email: "123456@",
			match: false,
		},
		{
			name:  "正常邮箱",
			email: "123456@qq.com",
			match: true,
		},
	}
	// 执行测试用例
	h := NewUserHandler(nil, nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRegex.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

//func TestHTTP(t *testing.T) {
//	// 可以用过控制http方法，URL和传入的body
//	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("我的请求体")))
//	assert.NoError(t, err)
//	// 取响应体中的数据
//	recorder := httptest.NewRecorder()
//	assert.Equal(t, http.SameSiteLaxMode, recorder.Code)
//}

//

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 构造请求
		reqBuilder func(*testing.T) *http.Request

		// 期待输出
		wantcode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Jike150*",
				}).Return(nil) // 模拟db出错
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike150*","confirmPassword":"Jike150*"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "Bind出错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike150*"`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@","password":"Jike150*","confirmPassword":"Jike150*"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "邮箱格式不正确",
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike1500","confirmPassword":"Jike1500"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "密码必须包含字母，数字和特殊字符，并且长度不能小于8",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike150*","confirmPassword":"Jike1500"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "两次输入密码不一致",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Jike150*",
				}).Return(errors.New("db错误")) // 模拟db出错
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike150*","confirmPassword":"Jike150*"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "系统错误",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Jike150*",
				}).Return(service.UserDuplicateError) // 模拟db出错
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("{"+
					`"email":"123@qq.com","password":"Jike150*","confirmPassword":"Jike150*"}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantcode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			//效果等价于真的从网络里面收到了一个请求
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantcode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
