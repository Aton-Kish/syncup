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

package usecase

import (
	"context"
	"path/filepath"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository/mock"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_pullUseCase_Execute(t *testing.T) {
	testdataBaseDir := "../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

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

	type mockFunctionRepositoryForAppSyncListReturn struct {
		out []model.Function
		err error
	}
	type mockFunctionRepositoryForAppSyncList struct {
		calls   int
		returns []mockFunctionRepositoryForAppSyncListReturn
	}

	type mockFunctionRepositoryForFSSaveReturn struct {
		out *model.Function
		err error
	}
	type mockFunctionRepositoryForFSSave struct {
		calls   int
		returns []mockFunctionRepositoryForFSSaveReturn
	}

	type expected struct {
		out   *PullOutput
		errAs error
		errIs error
	}

	tests := []struct {
		name                                 string
		args                                 args
		mockSchemaRepositoryForAppSyncGet    mockSchemaRepositoryForAppSyncGet
		mockSchemaRepositoryForFSSave        mockSchemaRepositoryForFSSave
		mockFunctionRepositoryForAppSyncList mockFunctionRepositoryForAppSyncList
		mockFunctionRepositoryForFSSave      mockFunctionRepositoryForFSSave
		expected                             expected
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						out: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{
					{
						out: &functionVTL_2018_05_29,
						err: nil,
					},
					{
						out: &functionAPPSYNC_JS_1_0_0,
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{},
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForFSSave.List() error",
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						out: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{},
			},
			expected: expected{
				out:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForFSSave.Save() error",
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						out: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{
					{
						out: nil,
						err: &model.LibError{},
					},
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

			mockTrackerRepository := mock_repository.NewMockTrackerRepository(ctrl)
			mockSchemaRepositoryForAppSync := mock_repository.NewMockSchemaRepository(ctrl)
			mockSchemaRepositoryForFS := mock_repository.NewMockSchemaRepository(ctrl)
			mockFunctionRepositoryForAppSync := mock_repository.NewMockFunctionRepository(ctrl)
			mockFunctionRepositoryForFS := mock_repository.NewMockFunctionRepository(ctrl)

			mockTrackerRepository.
				EXPECT().
				InProgress(ctx, gomock.Any()).
				AnyTimes()

			mockTrackerRepository.
				EXPECT().
				Success(ctx, gomock.Any()).
				AnyTimes()

			mockTrackerRepository.
				EXPECT().
				Failed(ctx, gomock.Any()).
				AnyTimes()

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

			mockFunctionRepositoryForAppSync.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Function, error) {
					defer func() { tt.mockFunctionRepositoryForAppSyncList.calls++ }()
					r := tt.mockFunctionRepositoryForAppSyncList.returns[tt.mockFunctionRepositoryForAppSyncList.calls]
					return r.out, r.err
				}).
				Times(len(tt.mockFunctionRepositoryForAppSyncList.returns))

			mockFunctionRepositoryForFS.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
					defer func() { tt.mockFunctionRepositoryForFSSave.calls++ }()
					r := tt.mockFunctionRepositoryForFSSave.returns[tt.mockFunctionRepositoryForFSSave.calls]
					return r.out, r.err
				}).
				Times(len(tt.mockFunctionRepositoryForFSSave.returns))

			uc := &pullUseCase{
				trackerRepository:            mockTrackerRepository,
				schemaRepositoryForAppSync:   mockSchemaRepositoryForAppSync,
				schemaRepositoryForFS:        mockSchemaRepositoryForFS,
				functionRepositoryForAppSync: mockFunctionRepositoryForAppSync,
				functionRepositoryForFS:      mockFunctionRepositoryForFS,
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
