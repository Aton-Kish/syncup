// Code generated by MockGen. DO NOT EDIT.
// Source: schema.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockSchemaRepository is a mock of SchemaRepository interface.
type MockSchemaRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSchemaRepositoryMockRecorder
}

// MockSchemaRepositoryMockRecorder is the mock recorder for MockSchemaRepository.
type MockSchemaRepositoryMockRecorder struct {
	mock *MockSchemaRepository
}

// NewMockSchemaRepository creates a new mock instance.
func NewMockSchemaRepository(ctrl *gomock.Controller) *MockSchemaRepository {
	mock := &MockSchemaRepository{ctrl: ctrl}
	mock.recorder = &MockSchemaRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchemaRepository) EXPECT() *MockSchemaRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockSchemaRepository) Get(ctx context.Context, apiID string) (*model.Schema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, apiID)
	ret0, _ := ret[0].(*model.Schema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSchemaRepositoryMockRecorder) Get(ctx, apiID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSchemaRepository)(nil).Get), ctx, apiID)
}

// Save mocks base method.
func (m *MockSchemaRepository) Save(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, apiID, schema)
	ret0, _ := ret[0].(*model.Schema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockSchemaRepositoryMockRecorder) Save(ctx, apiID, schema interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSchemaRepository)(nil).Save), ctx, apiID, schema)
}