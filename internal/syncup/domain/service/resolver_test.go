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

package service

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_resolverService_Difference(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
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
		resolvers1 []model.Resolver
		resolvers2 []model.Resolver
	}

	type expected struct {
		res   []model.Resolver
		errIs error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: no differences",
			args: args{
				resolvers1: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
				resolvers2: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
			},
			expected: expected{
				res:   []model.Resolver{},
				errIs: nil,
			},
		},
		{
			name: "happy path: some differences",
			args: args{
				resolvers1: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
				resolvers2: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverPIPELINE_VTL_2018_05_29,
				},
			},
			expected: expected{
				res: []model.Resolver{
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
				errIs: nil,
			},
		},
		{
			name: "happy path: everything is different",
			args: args{
				resolvers1: []model.Resolver{
					resolverUNIT_VTL_2018_05_29,
					resolverUNIT_APPSYNC_JS_1_0_0,
				},
				resolvers2: []model.Resolver{
					resolverPIPELINE_VTL_2018_05_29,
					resolverPIPELINE_APPSYNC_JS_1_0_0,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			s := &resolverService{}

			// Act
			actual, err := s.Difference(ctx, tt.args.resolvers1, tt.args.resolvers2)

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

func Test_resolverService_ResolvePipelineConfigFunctionIDs(t *testing.T) {
	type args struct {
		resolver  *model.Resolver
		functions []model.Function
	}

	type expected struct {
		resolver *model.Resolver
		errIs    error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: default",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId1"),
						Name:       ptr.Pointer("FunctionName1"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId2"),
						Name:       ptr.Pointer("FunctionName2"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions:     []string{"FunctionId1", "FunctionId2"},
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: nil,
			},
		},
		{
			name: "happy path: nil pipeline config",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: nil,
				},
				functions: []model.Function{},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: nil,
				},
				errIs: nil,
			},
		},
		{
			name: "edge path: nil resolver",
			args: args{
				resolver:  nil,
				functions: []model.Function{},
			},
			expected: expected{
				resolver: nil,
				errIs:    model.ErrNilValue,
			},
		},
		{
			name: "edge path: missing function id",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName1"),
					},
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName2"),
					},
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: missing name",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId1"),
						Name:       nil,
					},
					{
						FunctionId: ptr.Pointer("FunctionId2"),
						Name:       nil,
					},
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       nil,
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: duplicate function name",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId1"),
						Name:       ptr.Pointer("FunctionName"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId2"),
						Name:       ptr.Pointer("FunctionName"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       ptr.Pointer("FunctionName"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: model.ErrDuplicateValue,
			},
		},
		{
			name: "edge path: function name not found",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: model.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			s := &resolverService{}

			// Act
			err := s.ResolvePipelineConfigFunctionIDs(ctx, tt.args.resolver, tt.args.functions)

			// Assert
			assert.Equal(t, tt.expected.resolver, tt.args.resolver)

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

func Test_resolverService_ResolvePipelineConfigFunctionNames(t *testing.T) {
	type args struct {
		resolver  *model.Resolver
		functions []model.Function
	}

	type expected struct {
		resolver *model.Resolver
		errIs    error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: default",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId1"),
						Name:       ptr.Pointer("FunctionName1"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId2"),
						Name:       ptr.Pointer("FunctionName2"),
					},
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions:     []string{"FunctionId1", "FunctionId2"},
						FunctionNames: []string{"FunctionName1", "FunctionName2"},
					},
				},
				errIs: nil,
			},
		},
		{
			name: "happy path: nil pipeline config",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: nil,
				},
				functions: []model.Function{},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: nil,
				},
				errIs: nil,
			},
		},
		{
			name: "edge path: nil resolver",
			args: args{
				resolver:  nil,
				functions: []model.Function{},
			},
			expected: expected{
				resolver: nil,
				errIs:    model.ErrNilValue,
			},
		},
		{
			name: "edge path: missing function id",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName1"),
					},
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName2"),
					},
					{
						FunctionId: nil,
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: missing name",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId1"),
						Name:       nil,
					},
					{
						FunctionId: ptr.Pointer("FunctionId2"),
						Name:       nil,
					},
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       nil,
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: function id not found",
			args: args{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				functions: []model.Function{
					{
						FunctionId: ptr.Pointer("FunctionId3"),
						Name:       ptr.Pointer("FunctionName3"),
					},
				},
			},
			expected: expected{
				resolver: &model.Resolver{
					PipelineConfig: &model.PipelineConfig{
						Functions: []string{"FunctionId1", "FunctionId2"},
					},
				},
				errIs: model.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			s := &resolverService{}

			// Act
			err := s.ResolvePipelineConfigFunctionNames(ctx, tt.args.resolver, tt.args.functions)

			// Assert
			assert.Equal(t, tt.expected.resolver, tt.args.resolver)

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
