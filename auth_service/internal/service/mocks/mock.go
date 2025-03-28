// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "auth_service/internal/model"
	service "auth_service/internal/service"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockAuthorization) Authenticate(user *model.Person, ctx context.Context) *service.AuthenticationServiceResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", user, ctx)
	ret0, _ := ret[0].(*service.AuthenticationServiceResponse)
	return ret0
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockAuthorizationMockRecorder) Authenticate(user, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuthorization)(nil).Authenticate), user, ctx)
}

// Authorizate mocks base method.
func (m *MockAuthorization) Authorizate(sessionID string) *service.AuthenticationServiceResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorizate", sessionID)
	ret0, _ := ret[0].(*service.AuthenticationServiceResponse)
	return ret0
}

// Authorizate indicates an expected call of Authorizate.
func (mr *MockAuthorizationMockRecorder) Authorizate(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorizate", reflect.TypeOf((*MockAuthorization)(nil).Authorizate), sessionID)
}

// GenerateSession mocks base method.
func (m *MockAuthorization) GenerateSession(userId uuid.UUID) (string, time.Time) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSession", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(time.Time)
	return ret0, ret1
}

// GenerateSession indicates an expected call of GenerateSession.
func (mr *MockAuthorizationMockRecorder) GenerateSession(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSession", reflect.TypeOf((*MockAuthorization)(nil).GenerateSession), userId)
}

// Registrate mocks base method.
func (m *MockAuthorization) Registrate(user *model.Person, ctx context.Context) *service.AuthenticationServiceResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registrate", user, ctx)
	ret0, _ := ret[0].(*service.AuthenticationServiceResponse)
	return ret0
}

// Registrate indicates an expected call of Registrate.
func (mr *MockAuthorizationMockRecorder) Registrate(user, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registrate", reflect.TypeOf((*MockAuthorization)(nil).Registrate), user, ctx)
}
