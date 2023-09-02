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
	"path/filepath"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure/mapper"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

func Test_functionRepositoryForAppSync_List(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustJSONUnmarshal[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustJSONUnmarshal[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

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
	functionVTL_2018_05_29 := testhelpers.MustJSONUnmarshal[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))

	type args struct {
		apiID      string
		functionID string
	}

	type mockAppSyncClientGetFunctionReturn struct {
		out *appsync.GetFunctionOutput
		err error
	}
	type mockAppSyncClientGetFunction struct {
		calls   int
		returns []mockAppSyncClientGetFunctionReturn
	}

	type expected struct {
		out   *model.Function
		errAs error
		errIs error
	}

	tests := []struct {
		name                         string
		args                         args
		mockAppSyncClientGetFunction mockAppSyncClientGetFunction
		expected                     expected
	}{
		{
			name: "happy path",
			args: args{
				apiID:      "apiID",
				functionID: "functionID",
			},
			mockAppSyncClientGetFunction: mockAppSyncClientGetFunction{
				returns: []mockAppSyncClientGetFunctionReturn{
					{
						out: &appsync.GetFunctionOutput{
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
			name: "edge path: appsync.GetFunction() error",
			args: args{
				apiID:      "apiID",
				functionID: "functionID",
			},
			mockAppSyncClientGetFunction: mockAppSyncClientGetFunction{
				returns: []mockAppSyncClientGetFunctionReturn{
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
			name: "edge path: appsync.GetFunction() returns nil function",
			args: args{
				apiID:      "apiID",
				functionID: "functionID",
			},
			mockAppSyncClientGetFunction: mockAppSyncClientGetFunction{
				returns: []mockAppSyncClientGetFunctionReturn{
					{
						out: &appsync.GetFunctionOutput{
							FunctionConfiguration: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: model.ErrNilValue,
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
								case "GetFunction":
									defer func() { tt.mockAppSyncClientGetFunction.calls++ }()
									r := tt.mockAppSyncClientGetFunction.returns[tt.mockAppSyncClientGetFunction.calls]
									return smithymiddleware.FinalizeOutput{Result: r.out}, smithymiddleware.Metadata{}, r.err
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

			r := &functionRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID, tt.args.functionID)

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
