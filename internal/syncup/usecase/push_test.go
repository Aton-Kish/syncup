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

func Test_pushUseCase_Execute(t *testing.T) {
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
		params *PushInput
	}

	type mockSchemaRepositoryForFSGetReturn struct {
		res *model.Schema
		err error
	}
	type mockSchemaRepositoryForFSGet struct {
		calls   int
		returns []mockSchemaRepositoryForFSGetReturn
	}

	type mockSchemaRepositoryForAppSyncSaveReturn struct {
		res *model.Schema
		err error
	}
	type mockSchemaRepositoryForAppSyncSave struct {
		calls   int
		returns []mockSchemaRepositoryForAppSyncSaveReturn
	}

	type mockFunctionRepositoryForFSListReturn struct {
		res []model.Function
		err error
	}
	type mockFunctionRepositoryForFSList struct {
		calls   int
		returns []mockFunctionRepositoryForFSListReturn
	}

	type mockFunctionRepositoryForAppSyncSaveReturn struct {
		res *model.Function
		err error
	}
	type mockFunctionRepositoryForAppSyncSave struct {
		calls   int
		returns []mockFunctionRepositoryForAppSyncSaveReturn
	}

	type mockResolverRepositoryForFSListReturn struct {
		res []model.Resolver
		err error
	}
	type mockResolverRepositoryForFSList struct {
		calls   int
		returns []mockResolverRepositoryForFSListReturn
	}

	type mockResolverServiceResolvePipelineConfigFunctionIDsReturn struct {
		err error
	}
	type mockResolverServiceResolvePipelineConfigFunctionIDs struct {
		calls   int
		returns []mockResolverServiceResolvePipelineConfigFunctionIDsReturn
	}

	type mockResolverRepositoryForAppSyncSaveReturn struct {
		res *model.Resolver
		err error
	}
	type mockResolverRepositoryForAppSyncSave struct {
		calls   int
		returns []mockResolverRepositoryForAppSyncSaveReturn
	}

	type mockFunctionRepositoryForAppSyncListReturn struct {
		res []model.Function
		err error
	}
	type mockFunctionRepositoryForAppSyncList struct {
		calls   int
		returns []mockFunctionRepositoryForAppSyncListReturn
	}

	type mockFunctionServiceDifferenceReturn struct {
		res []model.Function
		err error
	}
	type mockFunctionServiceDifference struct {
		calls   int
		returns []mockFunctionServiceDifferenceReturn
	}

	type mockFunctionRepositoryForAppSyncDeleteReturn struct {
		err error
	}
	type mockFunctionRepositoryForAppSyncDelete struct {
		calls   int
		returns []mockFunctionRepositoryForAppSyncDeleteReturn
	}

	type mockResolverRepositoryForAppSyncListReturn struct {
		res []model.Resolver
		err error
	}
	type mockResolverRepositoryForAppSyncList struct {
		calls   int
		returns []mockResolverRepositoryForAppSyncListReturn
	}

	type mockResolverServiceDifferenceReturn struct {
		res []model.Resolver
		err error
	}
	type mockResolverServiceDifference struct {
		calls   int
		returns []mockResolverServiceDifferenceReturn
	}

	type mockResolverRepositoryForAppSyncDeleteReturn struct {
		err error
	}
	type mockResolverRepositoryForAppSyncDelete struct {
		calls   int
		returns []mockResolverRepositoryForAppSyncDeleteReturn
	}

	type expected struct {
		res   *PushOutput
		errIs error
	}

	tests := []struct {
		name                                                string
		args                                                args
		mockSchemaRepositoryForFSGet                        mockSchemaRepositoryForFSGet
		mockSchemaRepositoryForAppSyncSave                  mockSchemaRepositoryForAppSyncSave
		mockFunctionRepositoryForFSList                     mockFunctionRepositoryForFSList
		mockFunctionRepositoryForAppSyncSave                mockFunctionRepositoryForAppSyncSave
		mockResolverRepositoryForFSList                     mockResolverRepositoryForFSList
		mockResolverServiceResolvePipelineConfigFunctionIDs mockResolverServiceResolvePipelineConfigFunctionIDs
		mockResolverRepositoryForAppSyncSave                mockResolverRepositoryForAppSyncSave
		mockFunctionRepositoryForAppSyncList                mockFunctionRepositoryForAppSyncList
		mockFunctionServiceDifference                       mockFunctionServiceDifference
		mockFunctionRepositoryForAppSyncDelete              mockFunctionRepositoryForAppSyncDelete
		mockResolverRepositoryForAppSyncList                mockResolverRepositoryForAppSyncList
		mockResolverServiceDifference                       mockResolverServiceDifference
		mockResolverRepositoryForAppSyncDelete              mockResolverRepositoryForAppSyncDelete
		expected                                            expected
	}{
		{
			name: "happy path: exist no extraneous files",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
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
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{
					{
						res: []model.Resolver{},
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   &PushOutput{},
				errIs: nil,
			},
		},
		{
			name: "happy path: delete extraneous files",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{
					{
						err: nil,
					},
					{
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
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName1"), FieldName: ptr.Pointer("ExtraneousResolverFieldName1")},
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName2"), FieldName: ptr.Pointer("ExtraneousResolverFieldName2")},
						},
						err: nil,
					},
				},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{
					{
						res: []model.Resolver{
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName1"), FieldName: ptr.Pointer("ExtraneousResolverFieldName1")},
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName2"), FieldName: ptr.Pointer("ExtraneousResolverFieldName2")},
						},
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{
					{
						err: nil,
					},
					{
						err: nil,
					},
				},
			},
			expected: expected{
				res:   &PushOutput{},
				errIs: nil,
			},
		},
		{
			name: "happy path: skip deleting extraneous files",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: false,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   &PushOutput{},
				errIs: nil,
			},
		},
		{
			name: "edge path: SchemaRepositoryForFS.Get() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{},
			},
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{},
			},
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: SchemaRepositoryForAppSync.Save() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{},
			},
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{},
			},
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForAppSync.List() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{},
			},
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{},
			},
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForAppSync.Save() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{},
			},
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{},
			},
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForFS.List() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{},
			},
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverService.ResolvePipelineConfigFunctionIDs() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{},
			},
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForAppSync.Save() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForAppSync.List() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionService.Difference() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: FunctionRepositoryForAppSync.Delete() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{
					{
						err: &model.LibError{},
					},
					{
						err: &model.LibError{},
					},
				},
			},
			mockResolverRepositoryForAppSyncList: mockResolverRepositoryForAppSyncList{
				returns: []mockResolverRepositoryForAppSyncListReturn{},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForAppSync.List() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{
					{
						err: nil,
					},
					{
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
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverService.Difference() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{
					{
						err: nil,
					},
					{
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
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName1"), FieldName: ptr.Pointer("ExtraneousResolverFieldName1")},
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName2"), FieldName: ptr.Pointer("ExtraneousResolverFieldName2")},
						},
						err: nil,
					},
				},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{
					{
						res: nil,
						err: &model.LibError{},
					},
				},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{},
			},
			expected: expected{
				res:   nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: ResolverRepositoryForAppSync.Delete() error",
			args: args{
				params: &PushInput{
					APIID:                     "APIID",
					DeleteExtraneousResources: true,
				},
			},
			mockSchemaRepositoryForFSGet: mockSchemaRepositoryForFSGet{
				returns: []mockSchemaRepositoryForFSGetReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockSchemaRepositoryForAppSyncSave: mockSchemaRepositoryForAppSyncSave{
				returns: []mockSchemaRepositoryForAppSyncSaveReturn{
					{
						res: &schema,
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForFSList: mockFunctionRepositoryForFSList{
				returns: []mockFunctionRepositoryForFSListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncSave: mockFunctionRepositoryForAppSyncSave{
				returns: []mockFunctionRepositoryForAppSyncSaveReturn{
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
			mockResolverRepositoryForFSList: mockResolverRepositoryForFSList{
				returns: []mockResolverRepositoryForFSListReturn{
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
			mockResolverServiceResolvePipelineConfigFunctionIDs: mockResolverServiceResolvePipelineConfigFunctionIDs{
				returns: []mockResolverServiceResolvePipelineConfigFunctionIDsReturn{
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
			mockResolverRepositoryForAppSyncSave: mockResolverRepositoryForAppSyncSave{
				returns: []mockResolverRepositoryForAppSyncSaveReturn{
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
			mockFunctionRepositoryForAppSyncList: mockFunctionRepositoryForAppSyncList{
				returns: []mockFunctionRepositoryForAppSyncListReturn{
					{
						res: []model.Function{
							functionVTL_2018_05_29,
							functionAPPSYNC_JS_1_0_0,
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionServiceDifference: mockFunctionServiceDifference{
				returns: []mockFunctionServiceDifferenceReturn{
					{
						res: []model.Function{
							{Name: ptr.Pointer("ExtraneousFunction1")},
							{Name: ptr.Pointer("ExtraneousFunction2")},
						},
						err: nil,
					},
				},
			},
			mockFunctionRepositoryForAppSyncDelete: mockFunctionRepositoryForAppSyncDelete{
				returns: []mockFunctionRepositoryForAppSyncDeleteReturn{
					{
						err: nil,
					},
					{
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
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName1"), FieldName: ptr.Pointer("ExtraneousResolverFieldName1")},
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName2"), FieldName: ptr.Pointer("ExtraneousResolverFieldName2")},
						},
						err: nil,
					},
				},
			},
			mockResolverServiceDifference: mockResolverServiceDifference{
				returns: []mockResolverServiceDifferenceReturn{
					{
						res: []model.Resolver{
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName1"), FieldName: ptr.Pointer("ExtraneousResolverFieldName1")},
							{TypeName: ptr.Pointer("ExtraneousResolverTypeName2"), FieldName: ptr.Pointer("ExtraneousResolverFieldName2")},
						},
						err: nil,
					},
				},
			},
			mockResolverRepositoryForAppSyncDelete: mockResolverRepositoryForAppSyncDelete{
				returns: []mockResolverRepositoryForAppSyncDeleteReturn{
					{
						err: &model.LibError{},
					},
					{
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

			mockFunctionService := mock_service.NewMockFunctionService(ctrl)
			mockResolverService := mock_service.NewMockResolverService(ctrl)
			mockTrackerRepository := mock_repository.NewMockTrackerRepository(ctrl)
			mockSchemaRepositoryForFS := mock_repository.NewMockSchemaRepository(ctrl)
			mockSchemaRepositoryForAppSync := mock_repository.NewMockSchemaRepository(ctrl)
			mockFunctionRepositoryForFS := mock_repository.NewMockFunctionRepository(ctrl)
			mockFunctionRepositoryForAppSync := mock_repository.NewMockFunctionRepository(ctrl)
			mockResolverRepositoryForFS := mock_repository.NewMockResolverRepository(ctrl)
			mockResolverRepositoryForAppSync := mock_repository.NewMockResolverRepository(ctrl)

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

			mockSchemaRepositoryForFS.
				EXPECT().
				Get(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) (*model.Schema, error) {
					r := tt.mockSchemaRepositoryForFSGet.returns[tt.mockSchemaRepositoryForFSGet.calls]
					tt.mockSchemaRepositoryForFSGet.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForFSGet.returns))

			mockSchemaRepositoryForAppSync.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
					r := tt.mockSchemaRepositoryForAppSyncSave.returns[tt.mockSchemaRepositoryForAppSyncSave.calls]
					tt.mockSchemaRepositoryForAppSyncSave.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockSchemaRepositoryForAppSyncSave.returns))

			mockFunctionRepositoryForFS.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Function, error) {
					r := tt.mockFunctionRepositoryForFSList.returns[tt.mockFunctionRepositoryForFSList.calls]
					tt.mockFunctionRepositoryForFSList.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockFunctionRepositoryForFSList.returns))

			mockFunctionRepositoryForAppSync.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
					mu.Lock()
					r := tt.mockFunctionRepositoryForAppSyncSave.returns[tt.mockFunctionRepositoryForAppSyncSave.calls]
					tt.mockFunctionRepositoryForAppSyncSave.calls++
					mu.Unlock()
					return r.res, r.err
				}).
				MaxTimes(len(tt.mockFunctionRepositoryForAppSyncSave.returns))

			mockFunctionRepositoryForAppSync.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Function, error) {
					r := tt.mockFunctionRepositoryForAppSyncList.returns[tt.mockFunctionRepositoryForAppSyncList.calls]
					tt.mockFunctionRepositoryForAppSyncList.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockFunctionRepositoryForAppSyncList.returns))

			mockFunctionService.
				EXPECT().
				Difference(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, functions1 []model.Function, functions2 []model.Function) ([]model.Function, error) {
					r := tt.mockFunctionServiceDifference.returns[tt.mockFunctionServiceDifference.calls]
					tt.mockFunctionServiceDifference.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockFunctionServiceDifference.returns))

			mockFunctionRepositoryForAppSync.
				EXPECT().
				Delete(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, name string) error {
					mu.Lock()
					r := tt.mockFunctionRepositoryForAppSyncDelete.returns[tt.mockFunctionRepositoryForAppSyncDelete.calls]
					tt.mockFunctionRepositoryForAppSyncDelete.calls++
					mu.Unlock()
					return r.err
				}).
				MaxTimes(len(tt.mockFunctionRepositoryForAppSyncDelete.returns))

			mockResolverRepositoryForFS.
				EXPECT().
				List(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string) ([]model.Resolver, error) {
					r := tt.mockResolverRepositoryForFSList.returns[tt.mockResolverRepositoryForFSList.calls]
					tt.mockResolverRepositoryForFSList.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockResolverRepositoryForFSList.returns))

			mockResolverService.
				EXPECT().
				ResolvePipelineConfigFunctionIDs(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, resolver *model.Resolver, functions []model.Function) error {
					mu.Lock()
					r := tt.mockResolverServiceResolvePipelineConfigFunctionIDs.returns[tt.mockResolverServiceResolvePipelineConfigFunctionIDs.calls]
					tt.mockResolverServiceResolvePipelineConfigFunctionIDs.calls++
					mu.Unlock()
					return r.err
				}).
				MaxTimes(len(tt.mockResolverServiceResolvePipelineConfigFunctionIDs.returns))

			mockResolverRepositoryForAppSync.
				EXPECT().
				Save(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, resolver *model.Resolver) (*model.Resolver, error) {
					mu.Lock()
					r := tt.mockResolverRepositoryForAppSyncSave.returns[tt.mockResolverRepositoryForAppSyncSave.calls]
					tt.mockResolverRepositoryForAppSyncSave.calls++
					mu.Unlock()
					return r.res, r.err
				}).
				MaxTimes(len(tt.mockResolverRepositoryForAppSyncSave.returns))

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
				Difference(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, resolvers1 []model.Resolver, resolvers2 []model.Resolver) ([]model.Resolver, error) {
					r := tt.mockResolverServiceDifference.returns[tt.mockResolverServiceDifference.calls]
					tt.mockResolverServiceDifference.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockResolverServiceDifference.returns))

			mockResolverRepositoryForAppSync.
				EXPECT().
				Delete(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, apiID string, typeName string, fieldName string) error {
					mu.Lock()
					r := tt.mockResolverRepositoryForAppSyncDelete.returns[tt.mockResolverRepositoryForAppSyncDelete.calls]
					tt.mockResolverRepositoryForAppSyncDelete.calls++
					mu.Unlock()
					return r.err
				}).
				MaxTimes(len(tt.mockResolverRepositoryForAppSyncDelete.returns))

			uc := &pushUseCase{
				functionService:              mockFunctionService,
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
