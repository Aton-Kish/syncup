// Code generated by MockGen. DO NOT EDIT.
// Source: push.go
//
// Generated by this command:
//
//	mockgen -source=push.go -destination=./mock/mock_push.go
//
// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	usecase "github.com/Aton-Kish/syncup/internal/syncup/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockPushUseCase is a mock of PushUseCase interface.
type MockPushUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockPushUseCaseMockRecorder
}

// MockPushUseCaseMockRecorder is the mock recorder for MockPushUseCase.
type MockPushUseCaseMockRecorder struct {
	mock *MockPushUseCase
}

// NewMockPushUseCase creates a new mock instance.
func NewMockPushUseCase(ctrl *gomock.Controller) *MockPushUseCase {
	mock := &MockPushUseCase{ctrl: ctrl}
	mock.recorder = &MockPushUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPushUseCase) EXPECT() *MockPushUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockPushUseCase) Execute(ctx context.Context, params *usecase.PushInput) (*usecase.PushOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, params)
	ret0, _ := ret[0].(*usecase.PushOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockPushUseCaseMockRecorder) Execute(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockPushUseCase)(nil).Execute), ctx, params)
}
