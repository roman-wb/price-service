// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/roman-wb/price-service/internal/servers (interfaces: Logger,Parser,PriceRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	models "github.com/roman-wb/price-service/internal/models"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Infof mocks base method.
func (m *MockLogger) Infof(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// MockParser is a mock of Parser interface.
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance.
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockParser) Fetch(arg0 string) ([]models.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].([]models.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockParserMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockParser)(nil).Fetch), arg0)
}

// MockPriceRepo is a mock of PriceRepo interface.
type MockPriceRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPriceRepoMockRecorder
}

// MockPriceRepoMockRecorder is the mock recorder for MockPriceRepo.
type MockPriceRepoMockRecorder struct {
	mock *MockPriceRepo
}

// NewMockPriceRepo creates a new mock instance.
func NewMockPriceRepo(ctrl *gomock.Controller) *MockPriceRepo {
	mock := &MockPriceRepo{ctrl: ctrl}
	mock.recorder = &MockPriceRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPriceRepo) EXPECT() *MockPriceRepoMockRecorder {
	return m.recorder
}

// Import mocks base method.
func (m *MockPriceRepo) Import(arg0 time.Time, arg1 []models.Price) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Import indicates an expected call of Import.
func (mr *MockPriceRepoMockRecorder) Import(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockPriceRepo)(nil).Import), arg0, arg1)
}

// List mocks base method.
func (m *MockPriceRepo) List(arg0, arg1 int, arg2 string, arg3 int32) ([]models.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]models.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPriceRepoMockRecorder) List(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPriceRepo)(nil).List), arg0, arg1, arg2, arg3)
}
