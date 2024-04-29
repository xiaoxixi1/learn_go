// Code generated by MockGen. DO NOT EDIT.
// Source: ./webbook/internal/service/user.go
//
// Generated by this command:
//
//	mockgen -source=./webbook/internal/service/user.go -package=svcmocks -destination=./webbook/internal/service/mocks/user.mock.go
//

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "project_go/webbook/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Edit mocks base method.
func (m *MockUserService) Edit(cxt context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", cxt, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit.
func (mr *MockUserServiceMockRecorder) Edit(cxt, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockUserService)(nil).Edit), cxt, user)
}

// FindOrCreate mocks base method.
func (m *MockUserService) FindOrCreate(cxt context.Context, phone string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", cxt, phone)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockUserServiceMockRecorder) FindOrCreate(cxt, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockUserService)(nil).FindOrCreate), cxt, phone)
}

// Login mocks base method.
func (m *MockUserService) Login(cxt context.Context, email, password string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", cxt, email, password)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(cxt, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), cxt, email, password)
}

// Profile mocks base method.
func (m *MockUserService) Profile(cxt context.Context, userid int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", cxt, userid)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockUserServiceMockRecorder) Profile(cxt, userid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockUserService)(nil).Profile), cxt, userid)
}

// SignUp mocks base method.
func (m *MockUserService) SignUp(cxt context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", cxt, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceMockRecorder) SignUp(cxt, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserService)(nil).SignUp), cxt, user)
}
