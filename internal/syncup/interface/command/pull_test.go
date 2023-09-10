// Copyright (c) 2023 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package command

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository/mock"
	"github.com/Aton-Kish/syncup/internal/syncup/usecase"
	mock_usecase "github.com/Aton-Kish/syncup/internal/syncup/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_pullCommand_Execute(t *testing.T) {
	type args struct {
		args []string
	}

	type mockMFATokenProviderRepositoryGetReturn struct {
		res model.MFATokenProvider
	}
	type mockMFATokenProviderRepositoryGet struct {
		calls   int
		returns []mockMFATokenProviderRepositoryGetReturn
	}

	type mockAWSActivatorActivateAWSReturn struct {
		err error
	}
	type mockAWSActivatorActivateAWS struct {
		calls   int
		returns []mockAWSActivatorActivateAWSReturn
	}

	type mockBaseDirProviderSetBaseDirReturn struct {
	}
	type mockBaseDirProviderSetBaseDir struct {
		calls   int
		returns []mockBaseDirProviderSetBaseDirReturn
	}

	type mockPullUseCaseExecuteReturn struct {
		res *usecase.PullOutput
		err error
	}
	type mockPullUseCaseExecute struct {
		calls   int
		returns []mockPullUseCaseExecuteReturn
	}

	type expected struct {
		errAs error
		errIs error
	}

	tests := []struct {
		name                              string
		args                              args
		mockMFATokenProviderRepositoryGet mockMFATokenProviderRepositoryGet
		mockAWSActivatorActivateAWS       mockAWSActivatorActivateAWS
		mockBaseDirProviderSetBaseDir     mockBaseDirProviderSetBaseDir
		mockPullUseCaseExecute            mockPullUseCaseExecute
		expected                          expected
	}{
		{
			name: "happy path",
			args: args{
				args: []string{"--api-id", "apiID"},
			},
			mockMFATokenProviderRepositoryGet: mockMFATokenProviderRepositoryGet{
				returns: []mockMFATokenProviderRepositoryGetReturn{
					{
						res: func() (string, error) {
							return "123456", nil
						},
					},
				},
			},
			mockAWSActivatorActivateAWS: mockAWSActivatorActivateAWS{
				returns: []mockAWSActivatorActivateAWSReturn{
					{
						err: nil,
					},
				},
			},
			mockBaseDirProviderSetBaseDir: mockBaseDirProviderSetBaseDir{
				returns: []mockBaseDirProviderSetBaseDirReturn{
					{},
				},
			},
			mockPullUseCaseExecute: mockPullUseCaseExecute{
				returns: []mockPullUseCaseExecuteReturn{
					{
						res: &usecase.PullOutput{},
						err: nil,
					},
				},
			},
			expected: expected{
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: missing --api-id flag",
			args: args{
				args: []string{},
			},
			mockMFATokenProviderRepositoryGet: mockMFATokenProviderRepositoryGet{
				returns: []mockMFATokenProviderRepositoryGetReturn{
					{
						res: func() (string, error) {
							return "123456", nil
						},
					},
				},
			},
			mockAWSActivatorActivateAWS: mockAWSActivatorActivateAWS{
				returns: []mockAWSActivatorActivateAWSReturn{
					{
						err: nil,
					},
				},
			},
			mockBaseDirProviderSetBaseDir: mockBaseDirProviderSetBaseDir{
				returns: []mockBaseDirProviderSetBaseDirReturn{
					{},
				},
			},
			mockPullUseCaseExecute: mockPullUseCaseExecute{
				returns: []mockPullUseCaseExecuteReturn{},
			},
			expected: expected{
				errAs: &commandError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: AWSActivator.ActivateAWS() error",
			args: args{
				args: []string{"--api-id", "apiID"},
			},
			mockMFATokenProviderRepositoryGet: mockMFATokenProviderRepositoryGet{
				returns: []mockMFATokenProviderRepositoryGetReturn{
					{
						res: func() (string, error) {
							return "123456", nil
						},
					},
				},
			},
			mockAWSActivatorActivateAWS: mockAWSActivatorActivateAWS{
				returns: []mockAWSActivatorActivateAWSReturn{
					{
						err: errors.New("error"),
					},
				},
			},
			mockBaseDirProviderSetBaseDir: mockBaseDirProviderSetBaseDir{
				returns: []mockBaseDirProviderSetBaseDirReturn{},
			},
			mockPullUseCaseExecute: mockPullUseCaseExecute{
				returns: []mockPullUseCaseExecuteReturn{},
			},
			expected: expected{
				errAs: &commandError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: PullUseCase.Execute() error",
			args: args{
				args: []string{"--api-id", "apiID"},
			},
			mockMFATokenProviderRepositoryGet: mockMFATokenProviderRepositoryGet{
				returns: []mockMFATokenProviderRepositoryGetReturn{
					{
						res: func() (string, error) {
							return "123456", nil
						},
					},
				},
			},
			mockAWSActivatorActivateAWS: mockAWSActivatorActivateAWS{
				returns: []mockAWSActivatorActivateAWSReturn{
					{
						err: nil,
					},
				},
			},
			mockBaseDirProviderSetBaseDir: mockBaseDirProviderSetBaseDir{
				returns: []mockBaseDirProviderSetBaseDirReturn{
					{},
				},
			},
			mockPullUseCaseExecute: mockPullUseCaseExecute{
				returns: []mockPullUseCaseExecuteReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				errAs: &commandError{},
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPullUseCase := mock_usecase.NewMockPullUseCase(ctrl)
			mockAWSActivator := mock_repository.NewMockAWSActivator(ctrl)
			mockBaseDirProvider := mock_repository.NewMockBaseDirProvider(ctrl)
			mockMFATokenProviderRepository := mock_repository.NewMockMFATokenProviderRepository(ctrl)

			mockMFATokenProviderRepository.
				EXPECT().
				Get(ctx).
				DoAndReturn(func(ctx context.Context) model.MFATokenProvider {
					r := tt.mockMFATokenProviderRepositoryGet.returns[tt.mockMFATokenProviderRepositoryGet.calls]
					tt.mockMFATokenProviderRepositoryGet.calls++
					return r.res
				}).
				Times(len(tt.mockMFATokenProviderRepositoryGet.returns))

			mockAWSActivator.
				EXPECT().
				ActivateAWS(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, optFns ...func(o *model.AWSOptions)) error {
					r := tt.mockAWSActivatorActivateAWS.returns[tt.mockAWSActivatorActivateAWS.calls]
					tt.mockAWSActivatorActivateAWS.calls++
					return r.err
				}).
				Times(len(tt.mockAWSActivatorActivateAWS.returns))

			mockBaseDirProvider.
				EXPECT().
				SetBaseDir(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, dir string) {
					tt.mockBaseDirProviderSetBaseDir.calls++
				}).
				Times(len(tt.mockBaseDirProviderSetBaseDir.returns))

			mockPullUseCase.
				EXPECT().
				Execute(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, params *usecase.PullInput) (*usecase.PullOutput, error) {
					r := tt.mockPullUseCaseExecute.returns[tt.mockPullUseCaseExecute.calls]
					tt.mockPullUseCaseExecute.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockPullUseCaseExecute.returns))

			stdin := new(bytes.Reader)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			c := &pullCommand{
				options:                    newOptions(WithStdio(stdin, stdout, stderr)),
				useCase:                    mockPullUseCase,
				awsActivator:               mockAWSActivator,
				baseDirProvider:            mockBaseDirProvider,
				mfaTokenProviderRepository: mockMFATokenProviderRepository,
			}

			// Act
			err := c.Execute(ctx, tt.args.args...)

			// Assert
			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)

				assert.Equal(t, 0, stdin.Len())
				assert.Equal(t, 0, stdout.Len())
				assert.Equal(t, 0, stderr.Len())
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}

				assert.Equal(t, 0, stdin.Len())
				assert.Equal(t, 0, stdout.Len())
				assert.Greater(t, stderr.Len(), 0)
			}
		})
	}
}
