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
	"testing"
	"time"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure/mapper"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/stretchr/testify/assert"
)

func Test_functionRepositoryForAppSync_List(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.FunctionId = ptr.Pointer("FunctionId")
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.FunctionId = ptr.Pointer("FunctionId")
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID string
	}

	type mockAppSyncClientListFunctionsReturn struct {
		out *appsync.ListFunctionsOutput
		err error
	}
	type mockAppSyncClientListFunctions struct {
		calls   int
		returns []mockAppSyncClientListFunctionsReturn
	}

	type expected struct {
		out   []model.Function
		errAs error
		errIs error
	}

	tests := []struct {
		name                           string
		args                           args
		mockAppSyncClientListFunctions mockAppSyncClientListFunctions
		expected                       expected
	}{
		{
			name: "happy path",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out: []model.Function{
					functionVTL_2018_05_29,
					functionAPPSYNC_JS_1_0_0,
				},
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListFunctions() error",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListFunctions":
									defer func() { tt.mockAppSyncClientListFunctions.calls++ }()
									r := tt.mockAppSyncClientListFunctions.returns[tt.mockAppSyncClientListFunctions.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
				config.WithRetryer(func() aws.Retryer {
					return retry.AddWithMaxBackoffDelay(retry.NewStandard(), duration)
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &functionRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.List(ctx, tt.args.apiID)

			// Assert
			assert.Equal(t, tt.expected.out, actual)

			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}

func Test_functionRepositoryForAppSync_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.FunctionId = ptr.Pointer("FunctionId")
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.FunctionId = ptr.Pointer("FunctionId")
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID string
		name  string
	}

	type mockAppSyncClientListFunctionsReturn struct {
		out *appsync.ListFunctionsOutput
		err error
	}
	type mockAppSyncClientListFunctions struct {
		calls   int
		returns []mockAppSyncClientListFunctionsReturn
	}

	type expected struct {
		out   *model.Function
		errAs error
		errIs error
	}

	tests := []struct {
		name                           string
		args                           args
		mockAppSyncClientListFunctions mockAppSyncClientListFunctions
		expected                       expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   &functionVTL_2018_05_29,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListFunctions() error",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: function not found",
			args: args{
				apiID: "apiID",
				name:  "notExistName",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: model.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			cfg, err := config.LoadDefaultConfig(
				ctx,
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListFunctions":
									defer func() { tt.mockAppSyncClientListFunctions.calls++ }()
									r := tt.mockAppSyncClientListFunctions.returns[tt.mockAppSyncClientListFunctions.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
				config.WithRetryer(func() aws.Retryer {
					return retry.AddWithMaxBackoffDelay(retry.NewStandard(), duration)
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &functionRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID, tt.args.name)

			// Assert
			assert.Equal(t, tt.expected.out, actual)

			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}

func Test_functionRepositoryForAppSync_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.FunctionId = ptr.Pointer("FunctionId")
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.FunctionId = ptr.Pointer("FunctionId")
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID    string
		function *model.Function
	}

	type mockAppSyncClientListFunctionsReturn struct {
		out *appsync.ListFunctionsOutput
		err error
	}
	type mockAppSyncClientListFunctions struct {
		calls   int
		returns []mockAppSyncClientListFunctionsReturn
	}

	type mockAppSyncClientCreateFunctionReturn struct {
		out *appsync.CreateFunctionOutput
		err error
	}
	type mockAppSyncClientCreateFunction struct {
		calls   int
		returns []mockAppSyncClientCreateFunctionReturn
	}

	type mockAppSyncClientUpdateFunctionReturn struct {
		out *appsync.UpdateFunctionOutput
		err error
	}
	type mockAppSyncClientUpdateFunction struct {
		calls   int
		returns []mockAppSyncClientUpdateFunctionReturn
	}

	type expected struct {
		out   *model.Function
		errAs error
		errIs error
	}

	tests := []struct {
		name                            string
		args                            args
		mockAppSyncClientListFunctions  mockAppSyncClientListFunctions
		mockAppSyncClientCreateFunction mockAppSyncClientCreateFunction
		mockAppSyncClientUpdateFunction mockAppSyncClientUpdateFunction
		expected                        expected
	}{
		{
			name: "happy path: create",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{
					{
						out: &appsync.CreateFunctionOutput{
							FunctionConfiguration: mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   &functionVTL_2018_05_29,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "happy path: update",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{
					{
						out: &appsync.UpdateFunctionOutput{
							FunctionConfiguration: mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   &functionVTL_2018_05_29,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "happy path: create - retries on ConcurrentModificationException",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: &appsync.CreateFunctionOutput{
							FunctionConfiguration: mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   &functionVTL_2018_05_29,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "happy path: update - retries on ConcurrentModificationException",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: &appsync.UpdateFunctionOutput{
							FunctionConfiguration: mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   &functionVTL_2018_05_29,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: create - exceeds max retry count",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
				},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: update - exceeds max retry count",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
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
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: nil function",
			args: args{
				apiID:    "apiID",
				function: nil,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil name",
			args: args{
				apiID:    "apiID",
				function: &model.Function{},
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: appsync.ListFunctions() error",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.CreateFunction() error",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.UpdateFunction() error",
			args: args{
				apiID:    "apiID",
				function: &functionVTL_2018_05_29,
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateFunction: mockAppSyncClientCreateFunction{
				returns: []mockAppSyncClientCreateFunctionReturn{},
			},
			mockAppSyncClientUpdateFunction: mockAppSyncClientUpdateFunction{
				returns: []mockAppSyncClientUpdateFunctionReturn{
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListFunctions":
									defer func() { tt.mockAppSyncClientListFunctions.calls++ }()
									r := tt.mockAppSyncClientListFunctions.returns[tt.mockAppSyncClientListFunctions.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								case "CreateFunction":
									defer func() { tt.mockAppSyncClientCreateFunction.calls++ }()
									r := tt.mockAppSyncClientCreateFunction.returns[tt.mockAppSyncClientCreateFunction.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								case "UpdateFunction":
									defer func() { tt.mockAppSyncClientUpdateFunction.calls++ }()
									r := tt.mockAppSyncClientUpdateFunction.returns[tt.mockAppSyncClientUpdateFunction.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
				config.WithRetryer(func() aws.Retryer {
					return retry.AddWithMaxBackoffDelay(retry.NewStandard(), duration)
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &functionRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.function)

			// Assert
			assert.Equal(t, tt.expected.out, actual)

			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}

func Test_functionRepositoryForAppSync_Delete(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.FunctionId = ptr.Pointer("FunctionId")
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.FunctionId = ptr.Pointer("FunctionId")
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID string
		name  string
	}

	type mockAppSyncClientListFunctionsReturn struct {
		out *appsync.ListFunctionsOutput
		err error
	}
	type mockAppSyncClientListFunctions struct {
		calls   int
		returns []mockAppSyncClientListFunctionsReturn
	}

	type mockAppSyncClientDeleteFunctionReturn struct {
		out *appsync.DeleteFunctionOutput
		err error
	}
	type mockAppSyncClientDeleteFunction struct {
		calls   int
		returns []mockAppSyncClientDeleteFunctionReturn
	}

	type expected struct {
		errAs error
		errIs error
	}

	tests := []struct {
		name                            string
		args                            args
		mockAppSyncClientListFunctions  mockAppSyncClientListFunctions
		mockAppSyncClientDeleteFunction mockAppSyncClientDeleteFunction
		expected                        expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: &appsync.DeleteFunctionOutput{},
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
			name: "happy path: retries on ConcurrentModificationException",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: &appsync.DeleteFunctionOutput{},
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
			name: "edge path: exceeds max retry count",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 409}},
								Err:      &types.ConcurrentModificationException{},
							},
						},
					},
					{
						out: nil,
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
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListFunctions() error",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: function not found",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				errAs: &model.LibError{},
				errIs: model.ErrNotFound,
			},
		},
		{
			name: "edge path: appsync.DeleteFunction() error",
			args: args{
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			mockAppSyncClientListFunctions: mockAppSyncClientListFunctions{
				returns: []mockAppSyncClientListFunctionsReturn{
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionVTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						out: &appsync.ListFunctionsOutput{
							Functions: []types.FunctionConfiguration{
								*mapper.NewFunctionMapper().FromModel(context.Background(), &functionAPPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientDeleteFunction: mockAppSyncClientDeleteFunction{
				returns: []mockAppSyncClientDeleteFunctionReturn{
					{
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				errAs: &model.LibError{},
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListFunctions":
									defer func() { tt.mockAppSyncClientListFunctions.calls++ }()
									r := tt.mockAppSyncClientListFunctions.returns[tt.mockAppSyncClientListFunctions.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								case "DeleteFunction":
									defer func() { tt.mockAppSyncClientDeleteFunction.calls++ }()
									r := tt.mockAppSyncClientDeleteFunction.returns[tt.mockAppSyncClientDeleteFunction.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
								default:
									t.Fatal("unexpected operation")
									return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
								}
							}), smithymiddleware.After,
						)
					},
				}),
				config.WithRetryer(func() aws.Retryer {
					return retry.AddWithMaxBackoffDelay(retry.NewStandard(), duration)
				}),
			)
			assert.NoError(t, err)

			mockAppSyncClient := appsync.NewFromConfig(cfg)

			r := &functionRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			err = r.Delete(ctx, tt.args.apiID, tt.args.name)

			// Assert
			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}
			}
		})
	}
}
