// Code generated by MockGen. DO NOT EDIT.
// Source: mfa_token_provider.go
//
// Generated by this command:
//
//	mockgen -source=mfa_token_provider.go -destination=./mock/mock_mfa_token_provider.go
//
// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockMFATokenProviderRepository is a mock of MFATokenProviderRepository interface.
type MockMFATokenProviderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMFATokenProviderRepositoryMockRecorder
}

// MockMFATokenProviderRepositoryMockRecorder is the mock recorder for MockMFATokenProviderRepository.
type MockMFATokenProviderRepositoryMockRecorder struct {
	mock *MockMFATokenProviderRepository
}

// NewMockMFATokenProviderRepository creates a new mock instance.
func NewMockMFATokenProviderRepository(ctrl *gomock.Controller) *MockMFATokenProviderRepository {
	mock := &MockMFATokenProviderRepository{ctrl: ctrl}
	mock.recorder = &MockMFATokenProviderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMFATokenProviderRepository) EXPECT() *MockMFATokenProviderRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMFATokenProviderRepository) Get(ctx context.Context) model.MFATokenProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].(model.MFATokenProvider)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockMFATokenProviderRepositoryMockRecorder) Get(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMFATokenProviderRepository)(nil).Get), ctx)
}
