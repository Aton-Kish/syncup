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

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
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

func Test_schemaRepositoryForAppSync_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))

	type args struct {
		apiID string
	}

	type mockAppSyncClientGetIntrospectionSchemaReturn struct {
		res *appsync.GetIntrospectionSchemaOutput
		err error
	}
	type mockAppSyncClientGetIntrospectionSchema struct {
		calls   int
		returns []mockAppSyncClientGetIntrospectionSchemaReturn
	}

	type expected struct {
		res   *model.Schema
		errIs error
	}

	tests := []struct {
		name                                    string
		args                                    args
		mockAppSyncClientGetIntrospectionSchema mockAppSyncClientGetIntrospectionSchema
		expected                                expected
	}{
		{
			name: "happy path",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{
					{
						res: &appsync.GetIntrospectionSchemaOutput{
							Schema: []byte(schema),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &schema,
				errIs: nil,
			},
		},
		{
			name: "edge path",
			args: args{
				apiID: "apiID",
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{
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
								case "GetIntrospectionSchema":
									r := tt.mockAppSyncClientGetIntrospectionSchema.returns[tt.mockAppSyncClientGetIntrospectionSchema.calls]
									tt.mockAppSyncClientGetIntrospectionSchema.calls++
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

			r := &schemaRepositoryForAppSync{
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

func Test_schemaRepositoryForAppSync_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))
	duration := time.Duration(1) * time.Millisecond

	type fields struct {
		pollingInterval time.Duration
	}

	type args struct {
		apiID  string
		schema *model.Schema
	}

	type mockAppSyncClientStartSchemaCreationReturn struct {
		res *appsync.StartSchemaCreationOutput
		err error
	}
	type mockAppSyncClientStartSchemaCreation struct {
		calls   int
		returns []mockAppSyncClientStartSchemaCreationReturn
	}

	type mockAppSyncClientGetSchemaCreationStatusReturn struct {
		res *appsync.GetSchemaCreationStatusOutput
		err error
	}
	type mockAppSyncClientGetSchemaCreationStatus struct {
		calls   int
		returns []mockAppSyncClientGetSchemaCreationStatusReturn
	}

	type mockAppSyncClientGetIntrospectionSchemaReturn struct {
		res *appsync.GetIntrospectionSchemaOutput
		err error
	}
	type mockAppSyncClientGetIntrospectionSchema struct {
		calls   int
		returns []mockAppSyncClientGetIntrospectionSchemaReturn
	}

	type expected struct {
		res   *model.Schema
		errIs error
	}

	tests := []struct {
		name    string
		timeout time.Duration
		fields  fields
		args    args
		mockAppSyncClientStartSchemaCreation
		mockAppSyncClientGetSchemaCreationStatus
		mockAppSyncClientGetIntrospectionSchema
		expected expected
	}{
		{
			name:    "happy path: default",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusSuccess,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{
					{
						res: &appsync.GetIntrospectionSchemaOutput{
							Schema: []byte(schema),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &schema,
				errIs: nil,
			},
		},
		{
			name:    "happy path: retries on ConcurrentModificationException",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
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
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusSuccess,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{
					{
						res: &appsync.GetIntrospectionSchemaOutput{
							Schema: []byte(schema),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &schema,
				errIs: nil,
			},
		},
		{
			name:    "edge path: timeout",
			timeout: duration, // timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: context.DeadlineExceeded,
			},
		},
		{
			name:    "edge path: exceeds max retry count",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
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
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name:    "edge path: nil schema",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: nil,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name:    "edge path: appsync.StartSchemaCreation() error",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name:    "edge path: failed to create",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusFailed,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrCreateFailed,
			},
		},
		{
			name:    "edge path: invalid schema status",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatus("invalid schema status"),
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrInvalidValue,
			},
		},
		{
			name:    "edge path: appsync.GetSchemaCreationStatus() error",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: nil,
						err: errors.New("error"),
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name:    "edge path: appsync.GetIntrospectionSchema() error",
			timeout: 100 * duration, // not timeout
			fields: fields{
				pollingInterval: duration,
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			mockAppSyncClientStartSchemaCreation: mockAppSyncClientStartSchemaCreation{
				returns: []mockAppSyncClientStartSchemaCreationReturn{
					{
						res: &appsync.StartSchemaCreationOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetSchemaCreationStatus: mockAppSyncClientGetSchemaCreationStatus{
				returns: []mockAppSyncClientGetSchemaCreationStatusReturn{
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusProcessing,
						},
						err: nil,
					},
					{
						res: &appsync.GetSchemaCreationStatusOutput{
							Status: types.SchemaStatusSuccess,
						},
						err: nil,
					},
				},
			},
			mockAppSyncClientGetIntrospectionSchema: mockAppSyncClientGetIntrospectionSchema{
				returns: []mockAppSyncClientGetIntrospectionSchemaReturn{
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
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			cfg, err := config.LoadDefaultConfig(
				ctx,
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session")),
				config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
					func(stack *smithymiddleware.Stack) error {
						return stack.Finalize.Add(
							smithymiddleware.FinalizeMiddlewareFunc("Mock", func(ctx context.Context, input smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
								switch awsmiddleware.GetOperationName(ctx) {
								case "StartSchemaCreation":
									r := tt.mockAppSyncClientStartSchemaCreation.returns[tt.mockAppSyncClientStartSchemaCreation.calls]
									tt.mockAppSyncClientStartSchemaCreation.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								case "GetSchemaCreationStatus":
									r := tt.mockAppSyncClientGetSchemaCreationStatus.returns[tt.mockAppSyncClientGetSchemaCreationStatus.calls]
									tt.mockAppSyncClientGetSchemaCreationStatus.calls++
									return smithymiddleware.FinalizeOutput{Result: r.res}, smithymiddleware.Metadata{}, r.err
								case "GetIntrospectionSchema":
									r := tt.mockAppSyncClientGetIntrospectionSchema.returns[tt.mockAppSyncClientGetIntrospectionSchema.calls]
									tt.mockAppSyncClientGetIntrospectionSchema.calls++
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

			r := &schemaRepositoryForAppSync{
				appsyncClient:   mockAppSyncClient,
				pollingInterval: tt.fields.pollingInterval,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.schema)

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
