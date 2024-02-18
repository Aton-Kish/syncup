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

package infrastructure

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/stretchr/testify/assert"
)

func Test_environmentVariablesRepositoryForAppSync_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	variables := testhelpers.MustUnmarshalJSON[model.EnvironmentVariables](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "environment_variables/env.json")))

	type args struct {
		apiID string
	}

	type mockAppSyncClientGetGraphqlApiEnvironmentVariablesReturn struct {
		res *appsync.GetGraphqlApiEnvironmentVariablesOutput
		err error
	}
	type mockAppSyncClientGetGraphqlApiEnvironmentVariables struct {
		calls   int
		returns []mockAppSyncClientGetGraphqlApiEnvironmentVariablesReturn
	}

	type expected struct {
		res   model.EnvironmentVariables
		errIs error
	}

	tests := []struct {
		name                                               string
		args                                               args
		mockAppSyncClientGetGraphqlApiEnvironmentVariables mockAppSyncClientGetGraphqlApiEnvironmentVariables
		expected                                           expected
	}{
		{
			name: "happy path: no variables",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientGetGraphqlApiEnvironmentVariables: mockAppSyncClientGetGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientGetGraphqlApiEnvironmentVariablesReturn{
					{
						res: &appsync.GetGraphqlApiEnvironmentVariablesOutput{
							EnvironmentVariables: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   model.EnvironmentVariables{},
				errIs: nil,
			},
		},
		{
			name: "happy path: some variables",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientGetGraphqlApiEnvironmentVariables: mockAppSyncClientGetGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientGetGraphqlApiEnvironmentVariablesReturn{
					{
						res: &appsync.GetGraphqlApiEnvironmentVariablesOutput{
							EnvironmentVariables: variables,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   variables,
				errIs: nil,
			},
		},
		{
			name: "edge path",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientGetGraphqlApiEnvironmentVariables: mockAppSyncClientGetGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientGetGraphqlApiEnvironmentVariablesReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			cfg, err := config.LoadDefaultConfig(
				ctx,
				config.WithRegion("region"),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "GetGraphqlApiEnvironmentVariables":
									r := tt.mockAppSyncClientGetGraphqlApiEnvironmentVariables.returns[tt.mockAppSyncClientGetGraphqlApiEnvironmentVariables.calls]
									tt.mockAppSyncClientGetGraphqlApiEnvironmentVariables.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &environmentVariablesRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID)

			// Assert
			assert.Equal(t, tt.expected.res, actual)

			if strings.HasPrefix(tt.name, "happy") {
				assert.NoError(t, err)
			} else {
				var le *model.LibError
				assert.ErrorAs(t, err, &le)

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}

func Test_environmentVariablesRepositoryForAppSync_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	variables := testhelpers.MustUnmarshalJSON[model.EnvironmentVariables](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "environment_variables/env.json")))

	type args struct {
		apiID     string
		variables model.EnvironmentVariables
	}

	type mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn struct {
		res *appsync.PutGraphqlApiEnvironmentVariablesOutput
		err error
	}
	type mockAppSyncClientPutGraphqlApiEnvironmentVariables struct {
		calls   int
		returns []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn
	}

	type expected struct {
		res   model.EnvironmentVariables
		errIs error
	}

	tests := []struct {
		name                                               string
		args                                               args
		mockAppSyncClientPutGraphqlApiEnvironmentVariables mockAppSyncClientPutGraphqlApiEnvironmentVariables
		expected                                           expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID:     "apiID",
				variables: variables,
			},
			mockAppSyncClientPutGraphqlApiEnvironmentVariables: mockAppSyncClientPutGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn{
					{
						res: &appsync.PutGraphqlApiEnvironmentVariablesOutput{
							EnvironmentVariables: variables,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   variables,
				errIs: nil,
			},
		},
		{
			name: "happy path: retries on ConcurrentModificationException",
			args: args{
				apiID:     "apiID",
				variables: variables,
			},
			mockAppSyncClientPutGraphqlApiEnvironmentVariables: mockAppSyncClientPutGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						res: &appsync.PutGraphqlApiEnvironmentVariablesOutput{
							EnvironmentVariables: variables,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   variables,
				errIs: nil,
			},
		},
		{
			name: "edge path: exceeds max retry count",
			args: args{
				apiID:     "apiID",
				variables: variables,
			},
			mockAppSyncClientPutGraphqlApiEnvironmentVariables: mockAppSyncClientPutGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
				},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil environment variables",
			args: args{
				apiID:     "apiID",
				variables: nil,
			},
			mockAppSyncClientPutGraphqlApiEnvironmentVariables: mockAppSyncClientPutGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: appsync.PutGraphqlApiEnvironmentVariables() error",
			args: args{
				apiID:     "apiID",
				variables: variables,
			},
			mockAppSyncClientPutGraphqlApiEnvironmentVariables: mockAppSyncClientPutGraphqlApiEnvironmentVariables{
				returns: []mockAppSyncClientPutGraphqlApiEnvironmentVariablesReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			cfg, err := config.LoadDefaultConfig(
				ctx,
				config.WithRegion("region"),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "PutGraphqlApiEnvironmentVariables":
									r := tt.mockAppSyncClientPutGraphqlApiEnvironmentVariables.returns[tt.mockAppSyncClientPutGraphqlApiEnvironmentVariables.calls]
									tt.mockAppSyncClientPutGraphqlApiEnvironmentVariables.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &environmentVariablesRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.variables)

			// Assert
			assert.Equal(t, tt.expected.res, actual)

			if strings.HasPrefix(tt.name, "happy") {
				assert.NoError(t, err)
			} else {
				var le *model.LibError
				assert.ErrorAs(t, err, &le)

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}
