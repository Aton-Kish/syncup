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
	"strings"
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository/mock"
	"github.com/Aton-Kish/syncup/internal/syncup/usecase"
	mock_usecase "github.com/Aton-Kish/syncup/internal/syncup/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_pushCommand_Execute(t *testing.T) {
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

	type mockPushUseCaseExecuteReturn struct {
		res *usecase.PushOutput
		err error
	}
	type mockPushUseCaseExecute struct {
		calls   int
		returns []mockPushUseCaseExecuteReturn
	}

	type expected struct {
		errIs error
	}

	tests := []struct {
		name                              string
		args                              args
		mockMFATokenProviderRepositoryGet mockMFATokenProviderRepositoryGet
		mockAWSActivatorActivateAWS       mockAWSActivatorActivateAWS
		mockBaseDirProviderSetBaseDir     mockBaseDirProviderSetBaseDir
		mockPushUseCaseExecute            mockPushUseCaseExecute
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
			mockPushUseCaseExecute: mockPushUseCaseExecute{
				returns: []mockPushUseCaseExecuteReturn{
					{
						res: &usecase.PushOutput{},
						err: nil,
					},
				},
			},
			expected: expected{
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
			mockPushUseCaseExecute: mockPushUseCaseExecute{
				returns: []mockPushUseCaseExecuteReturn{},
			},
			expected: expected{
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
			mockPushUseCaseExecute: mockPushUseCaseExecute{
				returns: []mockPushUseCaseExecuteReturn{},
			},
			expected: expected{
				errIs: nil,
			},
		},
		{
			name: "edge path: PushUseCase.Execute() error",
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
			mockPushUseCaseExecute: mockPushUseCaseExecute{
				returns: []mockPushUseCaseExecuteReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
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

			mockPushUseCase := mock_usecase.NewMockPushUseCase(ctrl)
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

			mockPushUseCase.
				EXPECT().
				Execute(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, params *usecase.PushInput) (*usecase.PushOutput, error) {
					r := tt.mockPushUseCaseExecute.returns[tt.mockPushUseCaseExecute.calls]
					tt.mockPushUseCaseExecute.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockPushUseCaseExecute.returns))

			stdin := new(bytes.Reader)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			c := &pushCommand{
				options:                    newOptions(WithStdio(stdin, stdout, stderr)),
				useCase:                    mockPushUseCase,
				awsActivator:               mockAWSActivator,
				baseDirProvider:            mockBaseDirProvider,
				mfaTokenProviderRepository: mockMFATokenProviderRepository,
			}

			// Act
			err := c.Execute(ctx, tt.args.args...)

			// Assert
			if strings.HasPrefix(tt.name, "happy") {
				assert.NoError(t, err)

				assert.Equal(t, 0, stdin.Len())
				assert.Equal(t, 0, stdout.Len())
				assert.Equal(t, 0, stderr.Len())
			} else {
				var ce *commandError
				assert.ErrorAs(t, err, &ce)

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
