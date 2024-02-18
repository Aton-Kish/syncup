// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=./mock/mock_repository.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockAWSActivator is a mock of AWSActivator interface.
type MockAWSActivator struct {
	ctrl     *gomock.Controller
	recorder *MockAWSActivatorMockRecorder
}

// MockAWSActivatorMockRecorder is the mock recorder for MockAWSActivator.
type MockAWSActivatorMockRecorder struct {
	mock *MockAWSActivator
}

// NewMockAWSActivator creates a new mock instance.
func NewMockAWSActivator(ctrl *gomock.Controller) *MockAWSActivator {
	mock := &MockAWSActivator{ctrl: ctrl}
	mock.recorder = &MockAWSActivatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAWSActivator) EXPECT() *MockAWSActivatorMockRecorder {
	return m.recorder
}

// ActivateAWS mocks base method.
func (m *MockAWSActivator) ActivateAWS(ctx context.Context, optFns ...func(*model.AWSOptions)) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ActivateAWS", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateAWS indicates an expected call of ActivateAWS.
func (mr *MockAWSActivatorMockRecorder) ActivateAWS(ctx any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateAWS", reflect.TypeOf((*MockAWSActivator)(nil).ActivateAWS), varargs...)
}

// MockBaseDirProvider is a mock of BaseDirProvider interface.
type MockBaseDirProvider struct {
	ctrl     *gomock.Controller
	recorder *MockBaseDirProviderMockRecorder
}

// MockBaseDirProviderMockRecorder is the mock recorder for MockBaseDirProvider.
type MockBaseDirProviderMockRecorder struct {
	mock *MockBaseDirProvider
}

// NewMockBaseDirProvider creates a new mock instance.
func NewMockBaseDirProvider(ctrl *gomock.Controller) *MockBaseDirProvider {
	mock := &MockBaseDirProvider{ctrl: ctrl}
	mock.recorder = &MockBaseDirProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBaseDirProvider) EXPECT() *MockBaseDirProviderMockRecorder {
	return m.recorder
}

// BaseDir mocks base method.
func (m *MockBaseDirProvider) BaseDir(ctx context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseDir", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// BaseDir indicates an expected call of BaseDir.
func (mr *MockBaseDirProviderMockRecorder) BaseDir(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseDir", reflect.TypeOf((*MockBaseDirProvider)(nil).BaseDir), ctx)
}

// SetBaseDir mocks base method.
func (m *MockBaseDirProvider) SetBaseDir(ctx context.Context, dir string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBaseDir", ctx, dir)
}

// SetBaseDir indicates an expected call of SetBaseDir.
func (mr *MockBaseDirProviderMockRecorder) SetBaseDir(ctx, dir any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaseDir", reflect.TypeOf((*MockBaseDirProvider)(nil).SetBaseDir), ctx, dir)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// ActivateAWS mocks base method.
func (m *MockRepository) ActivateAWS(ctx context.Context, optFns ...func(*model.AWSOptions)) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ActivateAWS", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateAWS indicates an expected call of ActivateAWS.
func (mr *MockRepositoryMockRecorder) ActivateAWS(ctx any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateAWS", reflect.TypeOf((*MockRepository)(nil).ActivateAWS), varargs...)
}

// BaseDir mocks base method.
func (m *MockRepository) BaseDir(ctx context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseDir", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// BaseDir indicates an expected call of BaseDir.
func (mr *MockRepositoryMockRecorder) BaseDir(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseDir", reflect.TypeOf((*MockRepository)(nil).BaseDir), ctx)
}

// FunctionRepositoryForAppSync mocks base method.
func (m *MockRepository) FunctionRepositoryForAppSync() repository.FunctionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FunctionRepositoryForAppSync")
	ret0, _ := ret[0].(repository.FunctionRepository)
	return ret0
}

// FunctionRepositoryForAppSync indicates an expected call of FunctionRepositoryForAppSync.
func (mr *MockRepositoryMockRecorder) FunctionRepositoryForAppSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FunctionRepositoryForAppSync", reflect.TypeOf((*MockRepository)(nil).FunctionRepositoryForAppSync))
}

// FunctionRepositoryForFS mocks base method.
func (m *MockRepository) FunctionRepositoryForFS() repository.FunctionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FunctionRepositoryForFS")
	ret0, _ := ret[0].(repository.FunctionRepository)
	return ret0
}

// FunctionRepositoryForFS indicates an expected call of FunctionRepositoryForFS.
func (mr *MockRepositoryMockRecorder) FunctionRepositoryForFS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FunctionRepositoryForFS", reflect.TypeOf((*MockRepository)(nil).FunctionRepositoryForFS))
}

// MFATokenProviderRepository mocks base method.
func (m *MockRepository) MFATokenProviderRepository() repository.MFATokenProviderRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MFATokenProviderRepository")
	ret0, _ := ret[0].(repository.MFATokenProviderRepository)
	return ret0
}

// MFATokenProviderRepository indicates an expected call of MFATokenProviderRepository.
func (mr *MockRepositoryMockRecorder) MFATokenProviderRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MFATokenProviderRepository", reflect.TypeOf((*MockRepository)(nil).MFATokenProviderRepository))
}

// ResolverRepositoryForAppSync mocks base method.
func (m *MockRepository) ResolverRepositoryForAppSync() repository.ResolverRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolverRepositoryForAppSync")
	ret0, _ := ret[0].(repository.ResolverRepository)
	return ret0
}

// ResolverRepositoryForAppSync indicates an expected call of ResolverRepositoryForAppSync.
func (mr *MockRepositoryMockRecorder) ResolverRepositoryForAppSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolverRepositoryForAppSync", reflect.TypeOf((*MockRepository)(nil).ResolverRepositoryForAppSync))
}

// ResolverRepositoryForFS mocks base method.
func (m *MockRepository) ResolverRepositoryForFS() repository.ResolverRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolverRepositoryForFS")
	ret0, _ := ret[0].(repository.ResolverRepository)
	return ret0
}

// ResolverRepositoryForFS indicates an expected call of ResolverRepositoryForFS.
func (mr *MockRepositoryMockRecorder) ResolverRepositoryForFS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolverRepositoryForFS", reflect.TypeOf((*MockRepository)(nil).ResolverRepositoryForFS))
}

// SchemaRepositoryForAppSync mocks base method.
func (m *MockRepository) SchemaRepositoryForAppSync() repository.SchemaRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SchemaRepositoryForAppSync")
	ret0, _ := ret[0].(repository.SchemaRepository)
	return ret0
}

// SchemaRepositoryForAppSync indicates an expected call of SchemaRepositoryForAppSync.
func (mr *MockRepositoryMockRecorder) SchemaRepositoryForAppSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SchemaRepositoryForAppSync", reflect.TypeOf((*MockRepository)(nil).SchemaRepositoryForAppSync))
}

// SchemaRepositoryForFS mocks base method.
func (m *MockRepository) SchemaRepositoryForFS() repository.SchemaRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SchemaRepositoryForFS")
	ret0, _ := ret[0].(repository.SchemaRepository)
	return ret0
}

// SchemaRepositoryForFS indicates an expected call of SchemaRepositoryForFS.
func (mr *MockRepositoryMockRecorder) SchemaRepositoryForFS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SchemaRepositoryForFS", reflect.TypeOf((*MockRepository)(nil).SchemaRepositoryForFS))
}

// SetBaseDir mocks base method.
func (m *MockRepository) SetBaseDir(ctx context.Context, dir string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBaseDir", ctx, dir)
}

// SetBaseDir indicates an expected call of SetBaseDir.
func (mr *MockRepositoryMockRecorder) SetBaseDir(ctx, dir any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaseDir", reflect.TypeOf((*MockRepository)(nil).SetBaseDir), ctx, dir)
}

// TrackerRepository mocks base method.
func (m *MockRepository) TrackerRepository() repository.TrackerRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TrackerRepository")
	ret0, _ := ret[0].(repository.TrackerRepository)
	return ret0
}

// TrackerRepository indicates an expected call of TrackerRepository.
func (mr *MockRepositoryMockRecorder) TrackerRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackerRepository", reflect.TypeOf((*MockRepository)(nil).TrackerRepository))
}

// Version mocks base method.
func (m *MockRepository) Version() *model.Version {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(*model.Version)
	return ret0
}

// Version indicates an expected call of Version.
func (mr *MockRepositoryMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockRepository)(nil).Version))
}
