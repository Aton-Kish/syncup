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

//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE

package usecase

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository/mock"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_pullUseCase_Execute(t *testing.T) {
	testdataBaseDir := "../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))

	type args struct {
		params *PullInput
	}

	type mockSchemaRepositoryForAppSyncGetReturn struct {
		out *model.Schema
		err error
	}
	type mockSchemaRepositoryForAppSyncGet struct {
		calls   int
		returns []mockSchemaRepositoryForAppSyncGetReturn
	}

	type mockSchemaRepositoryForFSSaveReturn struct {
		out *model.Schema
		err error
	}
	type mockSchemaRepositoryForFSSave struct {
		calls   int
		returns []mockSchemaRepositoryForFSSaveReturn
	}

	type expected struct {
		out   *PullOutput
		errAs error
		errIs error
	}

	tests := []struct {
		name                              string
		args                              args
		mockSchemaRepositoryForAppSyncGet mockSchemaRepositoryForAppSyncGet
		mockSchemaRepositoryForFSSave     mockSchemaRepositoryForFSSave
		expected                          expected
	}{
		{
			name: "happy path",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						out: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						out: &schema,
						err: nil,
					},
				},
			},
			expected: expected{
				out:   &PullOutput{},
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: SchemaRepositoryForAppSync.Get() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						out: nil,
						err: &model.LibError{},
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: SchemaRepositoryForFSSave.Save() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						out: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						out: nil,
						err: &model.LibError{},
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSchemaRepositoryForAppSync := mock_repository.NewMockSchemaRepository(ctrl)
			mockSchemaRepositoryForFS := mock_repository.NewMockSchemaRepository(ctrl)

			mockSchemaRepositoryForAppSync.
				EXPECT().
				Get(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) (*model.Schema, error) {
					defer func() { tt.mockSchemaRepositoryForAppSyncGet.calls++ }()
					r := tt.mockSchemaRepositoryForAppSyncGet.returns[tt.mockSchemaRepositoryForAppSyncGet.calls]
					return r.out, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForAppSyncGet.returns))

			mockSchemaRepositoryForFS.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
					defer func() { tt.mockSchemaRepositoryForFSSave.calls++ }()
					r := tt.mockSchemaRepositoryForFSSave.returns[tt.mockSchemaRepositoryForFSSave.calls]
					return r.out, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForFSSave.returns))

			uc := &pullUseCase{
				schemaRepositoryForAppSync: mockSchemaRepositoryForAppSync,
				schemaRepositoryForFS:      mockSchemaRepositoryForFS,
			}

			// Act
			actual, err := uc.Execute(ctx, tt.args.params)

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
