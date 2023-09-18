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
	"strings"
	"sync"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_repository "github.com/Aton-Kish/syncup/internal/syncup/domain/repository/mock"
	mock_service "github.com/Aton-Kish/syncup/internal/syncup/domain/service/mock"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_pullUseCase_Execute(t *testing.T) {
	testdataBaseDir := "../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.FunctionId = ptr.Pointer("VTL_2018-05-29")
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.FunctionId = ptr.Pointer("APPSYNC_JS_1.0.0")
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))
	resolverPIPELINE_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/metadata.json")))
	resolverPIPELINE_VTL_2018_05_29.PipelineConfig.FunctionNames = nil
	resolverPIPELINE_VTL_2018_05_29.PipelineConfig.Functions = []string{"VTL_2018-05-29", "APPSYNC_JS_1.0.0"}
	resolverPIPELINE_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/request.vtl"))))
	resolverPIPELINE_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/response.vtl"))))
	resolverPIPELINE_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/metadata.json")))
	resolverPIPELINE_APPSYNC_JS_1_0_0.PipelineConfig.FunctionNames = nil
	resolverPIPELINE_APPSYNC_JS_1_0_0.PipelineConfig.Functions = []string{"VTL_2018-05-29", "APPSYNC_JS_1.0.0"}
	resolverPIPELINE_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/code.js"))))

	type args struct {
		params *PullInput
	}

	type mockSchemaRepositoryForAppSyncGetReturn struct {
		res *model.Schema
		err error
	}
	type mockSchemaRepositoryForAppSyncGet struct {
		calls   int
		returns []mockSchemaRepositoryForAppSyncGetReturn
	}

	type mockSchemaRepositoryForFSSaveReturn struct {
		res *model.Schema
		err error
	}
	type mockSchemaRepositoryForFSSave struct {
		calls   int
		returns []mockSchemaRepositoryForFSSaveReturn
	}

	type mockFunctionRepositoryForAppSyncListReturn struct {
		res []model.Function
		err error
	}
	type mockFunctionRepositoryForAppSyncList struct {
		calls   int
		returns []mockFunctionRepositoryForAppSyncListReturn
	}

	type mockFunctionRepositoryForFSSaveReturn struct {
		res *model.Function
		err error
	}
	type mockFunctionRepositoryForFSSave struct {
		calls   int
		returns []mockFunctionRepositoryForFSSaveReturn
	}

	type mockResolverRepositoryForAppSyncListReturn struct {
		res []model.Resolver
		err error
	}
	type mockResolverRepositoryForAppSyncList struct {
		calls   int
		returns []mockResolverRepositoryForAppSyncListReturn
	}

	type mockResolverServiceResolvePipelineConfigFunctionNamesReturn struct {
		err error
	}
	type mockResolverServiceResolvePipelineConfigFunctionNames struct {
		calls   int
		returns []mockResolverServiceResolvePipelineConfigFunctionNamesReturn
	}

	type mockResolverRepositoryForFSSaveReturn struct {
		res *model.Resolver
		err error
	}
	type mockResolverRepositoryForFSSave struct {
		calls   int
		returns []mockResolverRepositoryForFSSaveReturn
	}

	type expected struct {
		res   *PullOutput
		errIs error
	}

	tests := []struct {
		name                                                  string
		args                                                  args
		mockSchemaRepositoryForAppSyncGet                     mockSchemaRepositoryForAppSyncGet
		mockSchemaRepositoryForFSSave                         mockSchemaRepositoryForFSSave
		mockFunctionRepositoryForAppSyncList                  mockFunctionRepositoryForAppSyncList
		mockFunctionRepositoryForFSSave                       mockFunctionRepositoryForFSSave
		mockResolverRepositoryForAppSyncList                  mockResolverRepositoryForAppSyncList
		mockResolverServiceResolvePipelineConfigFunctionNames mockResolverServiceResolvePipelineConfigFunctionNames
		mockResolverRepositoryForFSSave                       mockResolverRepositoryForFSSave
		expected                                              expected
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
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
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
						res: &functionVTL_2018_05_29,
						err: nil,
					},
					{
						res: &functionAPPSYNC_JS_1_0_0,
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{
					{
						res: []model.Resolver{
							resolverUNIT_VTL_2018_05_29,
							resolverUNIT_APPSYNC_JS_1_0_0,
							resolverPIPELINE_VTL_2018_05_29,
							resolverPIPELINE_APPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{
					{
						err: nil,
					},
					{
						err: nil,
					},
					{
						err: nil,
					},
					{
						err: nil,
					},
				},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{
					{
						res: &resolverUNIT_VTL_2018_05_29,
						err: nil,
					},
					{
						res: &resolverUNIT_APPSYNC_JS_1_0_0,
						err: nil,
					},
					{
						res: &resolverPIPELINE_VTL_2018_05_29,
						err: nil,
					},
					{
						res: &resolverPIPELINE_APPSYNC_JS_1_0_0,
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &PullOutput{},
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
						res: nil,
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
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: SchemaRepositoryForFS.Save() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: nil,
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
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForFS.List() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionRepositoryForFSSave: mockFunctionRepositoryForFSSave{
				returns: []mockFunctionRepositoryForFSSaveReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForFS.Save() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
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
						res: nil,
						err: &model.LibError{},
					},
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForAppSync.List() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
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
						res: &functionVTL_2018_05_29,
						err: nil,
					},
					{
						res: &functionAPPSYNC_JS_1_0_0,
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverService.ResolvePipelineConfigFunctionNames() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
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
						res: &functionVTL_2018_05_29,
						err: nil,
					},
					{
						res: &functionAPPSYNC_JS_1_0_0,
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{
					{
						res: []model.Resolver{
							resolverUNIT_VTL_2018_05_29,
							resolverUNIT_APPSYNC_JS_1_0_0,
							resolverPIPELINE_VTL_2018_05_29,
							resolverPIPELINE_APPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{
					{
						err: &model.LibError{},
					},
					{
						err: &model.LibError{},
					},
					{
						err: &model.LibError{},
					},
					{
						err: &model.LibError{},
					},
				},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForFS.Save() error",
			args: args{
				params: &PullInput{
					APIID: "APIID",
				},
			},
			mockSchemaRepositoryForAppSyncGet: mockSchemaRepositoryForAppSyncGet{
				returns: []mockSchemaRepositoryForAppSyncGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForFSSave: mockSchemaRepositoryForFSSave{
				returns: []mockSchemaRepositoryForFSSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
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
						res: &functionVTL_2018_05_29,
						err: nil,
					},
					{
						res: &functionAPPSYNC_JS_1_0_0,
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{
					{
						res: []model.Resolver{
							resolverUNIT_VTL_2018_05_29,
							resolverUNIT_APPSYNC_JS_1_0_0,
							resolverPIPELINE_VTL_2018_05_29,
							resolverPIPELINE_APPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockResolverServiceResolvePipelineConfigFunctionNames: mockResolverServiceResolvePipelineConfigFunctionNames{
				returns: []mockResolverServiceResolvePipelineConfigFunctionNamesReturn{
					{
						err: nil,
					},
					{
						err: nil,
					},
					{
						err: nil,
					},
					{
						err: nil,
					},
				},
			},
			mockResolverRepositoryForFSSave: mockResolverRepositoryForFSSave{
				returns: []mockResolverRepositoryForFSSaveReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
					{
						res: nil,
						err: &model.LibError{},
					},
					{
						res: nil,
						err: &model.LibError{},
					},
					{
						res: nil,
						err: &model.LibError{},
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mu sync.Mutex

			mockResolverService := mock_service.NewMockResolverService(ctrl)
			mockTrackerRepository := mock_repository.NewMockTrackerRepository(ctrl)
			mockSchemaRepositoryForAppSync := mock_repository.NewMockSchemaRepository(ctrl)
			mockSchemaRepositoryForFS := mock_repository.NewMockSchemaRepository(ctrl)
			mockFunctionRepositoryForAppSync := mock_repository.NewMockFunctionRepository(ctrl)
			mockFunctionRepositoryForFS := mock_repository.NewMockFunctionRepository(ctrl)
			mockResolverRepositoryForAppSync := mock_repository.NewMockResolverRepository(ctrl)
			mockResolverRepositoryForFS := mock_repository.NewMockResolverRepository(ctrl)

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
					r := tt.mockSchemaRepositoryForAppSyncGet.returns[tt.mockSchemaRepositoryForAppSyncGet.calls]
					tt.mockSchemaRepositoryForAppSyncGet.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForAppSyncGet.returns))

			mockSchemaRepositoryForFS.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
					r := tt.mockSchemaRepositoryForFSSave.returns[tt.mockSchemaRepositoryForFSSave.calls]
					tt.mockSchemaRepositoryForFSSave.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForFSSave.returns))

			mockFunctionRepositoryForAppSync.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Function, error) {
					r := tt.mockFunctionRepositoryForAppSyncList.returns[tt.mockFunctionRepositoryForAppSyncList.calls]
					tt.mockFunctionRepositoryForAppSyncList.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockFunctionRepositoryForAppSyncList.returns))

			mockFunctionRepositoryForFS.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
					mu.Lock()
					r := tt.mockFunctionRepositoryForFSSave.returns[tt.mockFunctionRepositoryForFSSave.calls]
					tt.mockFunctionRepositoryForFSSave.calls++
					mu.Unlock()
					return r.res, r.err
				}).
				MaxTimes(len(tt.mockFunctionRepositoryForFSSave.returns))

			mockResolverRepositoryForAppSync.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Resolver, error) {
					r := tt.mockResolverRepositoryForAppSyncList.returns[tt.mockResolverRepositoryForAppSyncList.calls]
					tt.mockResolverRepositoryForAppSyncList.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockResolverRepositoryForAppSyncList.returns))

			mockResolverService.
				EXPECT().
				ResolvePipelineConfigFunctionNames(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, resolver *model.Resolver, functions []model.Function) error {
					mu.Lock()
					r := tt.mockResolverServiceResolvePipelineConfigFunctionNames.returns[tt.mockResolverServiceResolvePipelineConfigFunctionNames.calls]
					tt.mockResolverServiceResolvePipelineConfigFunctionNames.calls++
					mu.Unlock()
					return r.err
				}).
				MaxTimes(len(tt.mockResolverServiceResolvePipelineConfigFunctionNames.returns))

			mockResolverRepositoryForFS.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, resolver *model.Resolver) (*model.Resolver, error) {
					mu.Lock()
					r := tt.mockResolverRepositoryForFSSave.returns[tt.mockResolverRepositoryForFSSave.calls]
					tt.mockResolverRepositoryForFSSave.calls++
					mu.Unlock()
					return r.res, r.err
				}).
				MaxTimes(len(tt.mockResolverRepositoryForFSSave.returns))

			uc := &pullUseCase{
				resolverService:              mockResolverService,
				trackerRepository:            mockTrackerRepository,
				schemaRepositoryForAppSync:   mockSchemaRepositoryForAppSync,
				schemaRepositoryForFS:        mockSchemaRepositoryForFS,
				functionRepositoryForAppSync: mockFunctionRepositoryForAppSync,
				functionRepositoryForFS:      mockFunctionRepositoryForFS,
				resolverRepositoryForAppSync: mockResolverRepositoryForAppSync,
				resolverRepositoryForFS:      mockResolverRepositoryForFS,
			}

			// Act
			actual, err := uc.Execute(ctx, tt.args.params)

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
