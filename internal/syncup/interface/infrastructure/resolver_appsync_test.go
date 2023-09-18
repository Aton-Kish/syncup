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

func Test_resolverRepositoryForAppSync_List(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))
	resolverPIPELINE_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/metadata.json")))
	resolverPIPELINE_VTL_2018_05_29.PipelineConfig.FunctionNames = nil
	resolverPIPELINE_VTL_2018_05_29.PipelineConfig.Functions = []string{"FunctionId1", "FunctionId2"}
	resolverPIPELINE_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/request.vtl"))))
	resolverPIPELINE_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/response.vtl"))))
	resolverPIPELINE_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/metadata.json")))
	resolverPIPELINE_APPSYNC_JS_1_0_0.PipelineConfig.FunctionNames = nil
	resolverPIPELINE_APPSYNC_JS_1_0_0.PipelineConfig.Functions = []string{"FunctionId1", "FunctionId2"}
	resolverPIPELINE_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/code.js"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID string
	}

	type mockAppSyncClientListTypesReturn struct {
		res *appsync.ListTypesOutput
		err error
	}
	type mockAppSyncClientListTypes struct {
		calls   int
		returns []mockAppSyncClientListTypesReturn
	}

	type mockAppSyncClientListResolversReturn struct {
		res *appsync.ListResolversOutput
		err error
	}
	type mockAppSyncClientListResolvers struct {
		calls   int
		returns []mockAppSyncClientListResolversReturn
	}

	type expected struct {
		res   []model.Resolver
		errIs error
	}

	tests := []struct {
		name                           string
		args                           args
		mockAppSyncClientListTypes     mockAppSyncClientListTypes
		mockAppSyncClientListResolvers mockAppSyncClientListResolvers
		expected                       expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("PIPELINE"),
								},
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverPIPELINE_VTL_2018_05_29),
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverPIPELINE_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
					resolverPIPELINE_VTL_2018_05_29,
					resolverPIPELINE_APPSYNC_JS_1_0_0,
				},
				errIs: nil,
			},
		},
		{
			name: "happy path: retries on ConcurrentModificationException",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
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
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("PIPELINE"),
								},
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverPIPELINE_VTL_2018_05_29),
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverPIPELINE_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
					resolverPIPELINE_VTL_2018_05_29,
					resolverPIPELINE_APPSYNC_JS_1_0_0,
				},
				errIs: nil,
			},
		},
		{
			name: "edge path: exceeds max retry count",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
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
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListTypes() error",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil type name",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: nil,
								},
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListResolvers() error",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListTypes: mockAppSyncClientListTypes{
				returns: []mockAppSyncClientListTypesReturn{
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("UNIT"),
								},
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						res: &appsync.ListTypesOutput{
							Types: []types.Type{
								{
									Name: aws.String("PIPELINE"),
								},
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListTypes":
									r := tt.mockAppSyncClientListTypes.returns[tt.mockAppSyncClientListTypes.calls]
									tt.mockAppSyncClientListTypes.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								case "ListResolvers":
									r := tt.mockAppSyncClientListResolvers.returns[tt.mockAppSyncClientListResolvers.calls]
									tt.mockAppSyncClientListResolvers.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
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

			r := &resolverRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.List(ctx, tt.args.apiID)

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

func Test_resolverRepositoryForAppSync_ListByTypeName(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))

	type args struct {
		apiID    string
		typeName string
	}

	type mockAppSyncClientListResolversReturn struct {
		res *appsync.ListResolversOutput
		err error
	}
	type mockAppSyncClientListResolvers struct {
		calls   int
		returns []mockAppSyncClientListResolversReturn
	}

	type expected struct {
		res   []model.Resolver
		errIs error
	}

	tests := []struct {
		name                           string
		args                           args
		mockAppSyncClientListResolvers mockAppSyncClientListResolvers
		expected                       expected
	}{
		{
			name: "happy path",
			args: args{
				apiID:    "apiID",
				typeName: "UNIT",
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_APPSYNC_JS_1_0_0),
							},
							NextToken: nil,
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.ListResolvers() error",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientListResolvers: mockAppSyncClientListResolvers{
				returns: []mockAppSyncClientListResolversReturn{
					{
						res: &appsync.ListResolversOutput{
							Resolvers: []types.Resolver{
								*mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
							},
							NextToken: aws.String("NextToken"),
						},
						err: nil,
					},
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "ListResolvers":
									r := tt.mockAppSyncClientListResolvers.returns[tt.mockAppSyncClientListResolvers.calls]
									tt.mockAppSyncClientListResolvers.calls++
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

			r := &resolverRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.ListByTypeName(ctx, tt.args.apiID, tt.args.typeName)

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

func Test_resolverRepositoryForAppSync_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID     string
		typeName  string
		fieldName string
	}

	type mockAppSyncClientGetResolverReturn struct {
		res *appsync.GetResolverOutput
		err error
	}
	type mockAppSyncClientGetResolver struct {
		calls   int
		returns []mockAppSyncClientGetResolverReturn
	}

	type expected struct {
		res   *model.Resolver
		errIs error
	}

	tests := []struct {
		name                         string
		args                         args
		mockAppSyncClientGetResolver mockAppSyncClientGetResolver
		expected                     expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: retries on ConcurrentModificationException",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
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
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "edge path: exceeds max retry count",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
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
			name: "edge path: appsync.GetResolver() NotFoundException",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}},
								Err:      &types.NotFoundException{},
							},
						},
					},
				},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNotFound,
			},
		},
		{
			name: "edge path: appsync.GetResolver() except NotFoundException",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "GetResolver":
									r := tt.mockAppSyncClientGetResolver.returns[tt.mockAppSyncClientGetResolver.calls]
									tt.mockAppSyncClientGetResolver.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
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

			r := &resolverRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID, tt.args.typeName, tt.args.fieldName)

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

func Test_resolverRepositoryForAppSync_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID    string
		resolver *model.Resolver
	}

	type mockAppSyncClientGetResolverReturn struct {
		res *appsync.GetResolverOutput
		err error
	}
	type mockAppSyncClientGetResolver struct {
		calls   int
		returns []mockAppSyncClientGetResolverReturn
	}

	type mockAppSyncClientCreateResolverReturn struct {
		res *appsync.CreateResolverOutput
		err error
	}
	type mockAppSyncClientCreateResolver struct {
		calls   int
		returns []mockAppSyncClientCreateResolverReturn
	}

	type mockAppSyncClientUpdateResolverReturn struct {
		res *appsync.UpdateResolverOutput
		err error
	}
	type mockAppSyncClientUpdateResolver struct {
		calls   int
		returns []mockAppSyncClientUpdateResolverReturn
	}

	type expected struct {
		res   *model.Resolver
		errIs error
	}

	tests := []struct {
		name                            string
		args                            args
		mockAppSyncClientGetResolver    mockAppSyncClientGetResolver
		mockAppSyncClientCreateResolver mockAppSyncClientCreateResolver
		mockAppSyncClientUpdateResolver mockAppSyncClientUpdateResolver
		expected                        expected
	}{
		{
			name: "happy path: create",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}},
								Err:      &types.NotFoundException{},
							},
						},
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{
					{
						res: &appsync.CreateResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: update",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{
					{
						res: &appsync.UpdateResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: create - retries on ConcurrentModificationException",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}},
								Err:      &types.NotFoundException{},
							},
						},
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{
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
						res: &appsync.CreateResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: update - retries on ConcurrentModificationException",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{
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
						res: &appsync.UpdateResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "edge path: create - exceeds max retry count",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}},
								Err:      &types.NotFoundException{},
							},
						},
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{
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
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: update - exceeds max retry count",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{
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
			name: "edge path: nil resolver",
			args: args{
				apiID:    "apiID",
				resolver: nil,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil type name",
			args: args{
				apiID: "apiID",
				resolver: &model.Resolver{
					FieldName: ptr.Pointer("FieldName"),
				},
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil field name",
			args: args{
				apiID: "apiID",
				resolver: &model.Resolver{
					TypeName: ptr.Pointer("TypeName"),
				},
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: appsync.GetResolver() error",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.CreateResolver() error",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: nil,
						err: &awshttp.ResponseError{
							ResponseError: &smithyhttp.ResponseError{
								Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 404}},
								Err:      &types.NotFoundException{},
							},
						},
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.UpdateResolver() error",
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			mockAppSyncClientGetResolver: mockAppSyncClientGetResolver{
				returns: []mockAppSyncClientGetResolverReturn{
					{
						res: &appsync.GetResolverOutput{
							Resolver: mapper.NewResolverMapper().FromModel(context.Background(), &resolverUNIT_VTL_2018_05_29),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientCreateResolver: mockAppSyncClientCreateResolver{
				returns: []mockAppSyncClientCreateResolverReturn{},
			},
			mockAppSyncClientUpdateResolver: mockAppSyncClientUpdateResolver{
				returns: []mockAppSyncClientUpdateResolverReturn{
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
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "GetResolver":
									r := tt.mockAppSyncClientGetResolver.returns[tt.mockAppSyncClientGetResolver.calls]
									tt.mockAppSyncClientGetResolver.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								case "CreateResolver":
									r := tt.mockAppSyncClientCreateResolver.returns[tt.mockAppSyncClientCreateResolver.calls]
									tt.mockAppSyncClientCreateResolver.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								case "UpdateResolver":
									r := tt.mockAppSyncClientUpdateResolver.returns[tt.mockAppSyncClientUpdateResolver.calls]
									tt.mockAppSyncClientUpdateResolver.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
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

			r := &resolverRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.resolver)

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

func Test_resolverRepositoryForAppSync_Delete(t *testing.T) {
	duration := time.Duration(1) * time.Millisecond

	type args struct {
		apiID     string
		typeName  string
		fieldName string
	}

	type mockAppSyncClientDeleteResolverReturn struct {
		res *appsync.DeleteResolverOutput
		err error
	}
	type mockAppSyncClientDeleteResolver struct {
		calls   int
		returns []mockAppSyncClientDeleteResolverReturn
	}

	type expected struct {
		errIs error
	}

	tests := []struct {
		name                            string
		args                            args
		mockAppSyncClientDeleteResolver mockAppSyncClientDeleteResolver
		expected                        expected
	}{
		{
			name: "happy path: default",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientDeleteResolver: mockAppSyncClientDeleteResolver{
				returns: []mockAppSyncClientDeleteResolverReturn{
					{
						res: &appsync.DeleteResolverOutput{},
						err: nil,
					},
				},
			},
			expected: expected{
				errIs: nil,
			},
		},
		{
			name: "happy path: retries on ConcurrentModificationException",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientDeleteResolver: mockAppSyncClientDeleteResolver{
				returns: []mockAppSyncClientDeleteResolverReturn{
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
						res: &appsync.DeleteResolverOutput{},
						err: nil,
					},
				},
			},
			expected: expected{
				errIs: nil,
			},
		},
		{
			name: "edge path: exceeds max retry count",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientDeleteResolver: mockAppSyncClientDeleteResolver{
				returns: []mockAppSyncClientDeleteResolverReturn{
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
				errIs: nil,
			},
		},
		{
			name: "edge path: appsync.DeleteResolver()",
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018_05_29",
			},
			mockAppSyncClientDeleteResolver: mockAppSyncClientDeleteResolver{
				returns: []mockAppSyncClientDeleteResolverReturn{
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

			cfg, err := config.LoadDefaultConfig(
				ctx,
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "DeleteResolver":
									r := tt.mockAppSyncClientDeleteResolver.returns[tt.mockAppSyncClientDeleteResolver.calls]
									tt.mockAppSyncClientDeleteResolver.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
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

			r := &resolverRepositoryForAppSync{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			err = r.Delete(ctx, tt.args.apiID, tt.args.typeName, tt.args.fieldName)

			// Assert
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
