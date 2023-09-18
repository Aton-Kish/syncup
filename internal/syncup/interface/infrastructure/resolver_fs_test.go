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
	"path/filepath"
	"strings"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_resolverRepositoryForFS_List(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))
	resolverPIPELINE_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/metadata.json")))
	resolverPIPELINE_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/request.vtl"))))
	resolverPIPELINE_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/VTL_2018-05-29/response.vtl"))))
	resolverPIPELINE_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/metadata.json")))
	resolverPIPELINE_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/PIPELINE/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
	}

	type expected struct {
		res   []model.Resolver
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path",
			fields: fields{
				baseDir: testdataBaseDir,
			},
			args: args{
				apiID: "apiID",
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
			name: "edge path: non-existing dir",
			fields: fields{
				baseDir: filepath.Join(testdataBaseDir, "invalidBaseDir"),
			},
			args: args{
				apiID: "apiID",
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

			r := &resolverRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.List(ctx, tt.args.apiID)

			// Assert
			assert.ElementsMatch(t, tt.expected.res, actual)

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

func Test_resolverRepositoryForFS_ListByTypeName(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID    string
		typeName string
	}

	type expected struct {
		res   []model.Resolver
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path",
			fields: fields{
				baseDir: testdataBaseDir,
			},
			args: args{
				apiID:    "apiID",
				typeName: "UNIT",
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
			name: "edge path: non-existing dir",
			fields: fields{
				baseDir: filepath.Join(testdataBaseDir, "invalidBaseDir"),
			},
			args: args{
				apiID:    "apiID",
				typeName: "UNIT",
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

			r := &resolverRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.ListByTypeName(ctx, tt.args.apiID, tt.args.typeName)

			// Assert
			assert.ElementsMatch(t, tt.expected.res, actual)

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

func Test_resolverRepositoryForFS_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID     string
		typeName  string
		fieldName string
	}

	type expected struct {
		res   *model.Resolver
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path: VTL runtime",
			fields: fields{
				baseDir: testdataBaseDir,
			},
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "VTL_2018-05-29",
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: AppSync JS runtime",
			fields: fields{
				baseDir: testdataBaseDir,
			},
			args: args{
				apiID:     "apiID",
				typeName:  "UNIT",
				fieldName: "APPSYNC_JS_1.0.0",
			},
			expected: expected{
				res:   &resolverUNIT_APPSYNC_JS_1_0_0,
				errIs: nil,
			},
		},
		{
			name: "edge path: non-existing dir",
			fields: fields{
				baseDir: filepath.Join(testdataBaseDir, "invalidBaseDir"),
			},
			args: args{
				apiID:     "apiID",
				typeName:  "invalidTypeName",
				fieldName: "invalidFieldName",
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

			r := &resolverRepositoryForFS{
				baseDir: tt.fields.baseDir,
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

func Test_resolverRepositoryForFS_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID    string
		resolver *model.Resolver
	}

	type expected struct {
		res   *model.Resolver
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path: VTL runtime - existing dir",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: AppSync JS runtime - existing dir",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_APPSYNC_JS_1_0_0,
			},
			expected: expected{
				res:   &resolverUNIT_APPSYNC_JS_1_0_0,
				errIs: nil,
			},
		},
		{
			name: "happy path: VTL runtime - non-existing dir",
			fields: fields{
				baseDir: filepath.Join(t.TempDir(), "notExist"),
			},
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_VTL_2018_05_29,
			},
			expected: expected{
				res:   &resolverUNIT_VTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: AppSync JS runtime - non-existing dir",
			fields: fields{
				baseDir: filepath.Join(t.TempDir(), "notExist"),
			},
			args: args{
				apiID:    "apiID",
				resolver: &resolverUNIT_APPSYNC_JS_1_0_0,
			},
			expected: expected{
				res:   &resolverUNIT_APPSYNC_JS_1_0_0,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil resolver",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:    "apiID",
				resolver: nil,
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil type name",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID: "apiID",
				resolver: &model.Resolver{
					FieldName: ptr.Pointer("FieldName"),
				},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil field name",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID: "apiID",
				resolver: &model.Resolver{
					TypeName: ptr.Pointer("TypeName"),
				},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: invalid runtime",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID: "apiID",
				resolver: &model.Resolver{
					TypeName:  ptr.Pointer("TypeName"),
					FieldName: ptr.Pointer("FieldName"),
					Runtime: &model.Runtime{
						Name: "INVALID",
					},
				},
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrInvalidValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &resolverRepositoryForFS{
				baseDir: tt.fields.baseDir,
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

func Test_resolverRepositoryForFS_Delete(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	resolverUNIT_VTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/metadata.json")))
	resolverUNIT_VTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/request.vtl"))))
	resolverUNIT_VTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/VTL_2018-05-29/response.vtl"))))
	resolverUNIT_APPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Resolver](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/metadata.json")))
	resolverUNIT_APPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "resolvers/UNIT/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID     string
		typeName  string
		fieldName string
	}

	type expected struct {
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path: VTL runtime",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:     "apiID",
				typeName:  *resolverUNIT_VTL_2018_05_29.TypeName,
				fieldName: *resolverUNIT_VTL_2018_05_29.FieldName,
			},
			expected: expected{
				errIs: nil,
			},
		},
		{
			name: "happy path: AppSync JS runtime",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:     "apiID",
				typeName:  *resolverUNIT_APPSYNC_JS_1_0_0.TypeName,
				fieldName: *resolverUNIT_APPSYNC_JS_1_0_0.FieldName,
			},
			expected: expected{
				errIs: nil,
			},
		},
		{
			name: "happy path: non-existing dir",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:     "apiID",
				typeName:  "notExistName",
				fieldName: "notExistName",
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

			r := &resolverRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			_, err := r.Save(ctx, tt.args.apiID, &resolverUNIT_VTL_2018_05_29)
			assert.NoError(t, err)

			_, err = r.Save(ctx, tt.args.apiID, &resolverUNIT_APPSYNC_JS_1_0_0)
			assert.NoError(t, err)

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
