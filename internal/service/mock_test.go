// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mirhijinam/pg-project/internal/model"
	repository "github.com/mirhijinam/pg-project/internal/repository"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepo) Create(arg0 context.Context, arg1 *model.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepoMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), arg0, arg1)
}

// GetList mocks base method.
func (m *MockRepo) GetList(arg0 context.Context) ([]model.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", arg0)
	ret0, _ := ret[0].([]model.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockRepoMockRecorder) GetList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockRepo)(nil).GetList), arg0)
}

// GetOne mocks base method.
func (m *MockRepo) GetOne(arg0 context.Context, arg1 int) (model.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", arg0, arg1)
	ret0, _ := ret[0].(model.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockRepoMockRecorder) GetOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockRepo)(nil).GetOne), arg0, arg1)
}

// SetError mocks base method.
func (m *MockRepo) SetError(arg0 context.Context, arg1 int, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetError", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetError indicates an expected call of SetError.
func (mr *MockRepoMockRecorder) SetError(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetError", reflect.TypeOf((*MockRepo)(nil).SetError), arg0, arg1, arg2)
}

// SetSuccess mocks base method.
func (m *MockRepo) SetSuccess(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSuccess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSuccess indicates an expected call of SetSuccess.
func (mr *MockRepoMockRecorder) SetSuccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSuccess", reflect.TypeOf((*MockRepo)(nil).SetSuccess), arg0, arg1)
}

// Writer mocks base method.
func (m *MockRepo) Writer(arg0 context.Context, arg1 int) repository.WriterFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer", arg0, arg1)
	ret0, _ := ret[0].(repository.WriterFunc)
	return ret0
}

// Writer indicates an expected call of Writer.
func (mr *MockRepoMockRecorder) Writer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockRepo)(nil).Writer), arg0, arg1)
}

// MockPool is a mock of Pool interface.
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool.
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance.
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// Go mocks base method.
func (m *MockPool) Go(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Go", arg0)
}

// Go indicates an expected call of Go.
func (mr *MockPoolMockRecorder) Go(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Go", reflect.TypeOf((*MockPool)(nil).Go), arg0)
}

// activate mocks base method.
func (m *MockPool) activate() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "activate")
}

// activate indicates an expected call of activate.
func (mr *MockPoolMockRecorder) activate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "activate", reflect.TypeOf((*MockPool)(nil).activate))
}
