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

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_infrastructure "github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure/mock"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_schemaAppSyncRepository_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))

	type args struct {
		apiID string
	}

	type mockAppSyncClientGetIntrospectionSchemaReturn struct {
		out *appsync.GetIntrospectionSchemaOutput
		err error
	}
	type mockAppSyncClientGetIntrospectionSchema struct {
		calls   int
		returns []mockAppSyncClientGetIntrospectionSchemaReturn
	}

	type expected struct {
		out   *model.Schema
		errAs error
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
						out: &appsync.GetIntrospectionSchemaOutput{
							Schema: []byte(schema),
						},
						err: nil,
					},
				},
			},
			expected: expected{
				out:   &schema,
				errAs: nil,
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
						out: nil,
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				out:   &schema,
				errAs: &model.LibError{},
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

			mockAppSyncClient := mock_infrastructure.NewMockappsyncClient(ctrl)

			mockAppSyncClient.
				EXPECT().
				GetIntrospectionSchema(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, params *appsync.GetIntrospectionSchemaInput, optFns ...func(*appsync.Options)) (*appsync.GetIntrospectionSchemaOutput, error) {
					defer func() { tt.mockAppSyncClientGetIntrospectionSchema.calls++ }()

					r := tt.mockAppSyncClientGetIntrospectionSchema.returns[tt.mockAppSyncClientGetIntrospectionSchema.calls]

					return r.out, r.err
				}).
				Times(len(tt.mockAppSyncClientGetIntrospectionSchema.returns))

			r := &schemaAppSyncRepository{
				appsyncClient: mockAppSyncClient,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID)

			// Assert
			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.Equal(t, tt.expected.out, actual)

				assert.NoError(t, err)
			} else {
				assert.Nil(t, actual)

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
