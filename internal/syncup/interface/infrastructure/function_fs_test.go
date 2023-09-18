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

func Test_functionRepositoryForFS_List(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
	}

	type expected struct {
		res   []model.Function
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
				res: []model.Function{
					functionVTL_2018_05_29,
					functionAPPSYNC_JS_1_0_0,
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

			r := &functionRepositoryForFS{
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

func Test_functionRepositoryForFS_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
		name  string
	}

	type expected struct {
		res   *model.Function
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
				apiID: "apiID",
				name:  "VTL_2018-05-29",
			},
			expected: expected{
				res:   &functionVTL_2018_05_29,
				errIs: nil,
			},
		},
		{
			name: "happy path: AppSync JS runtime",
			fields: fields{
				baseDir: testdataBaseDir,
			},
			args: args{
				apiID: "apiID",
				name:  "APPSYNC_JS_1.0.0",
			},
			expected: expected{
				res:   &functionAPPSYNC_JS_1_0_0,
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
				name:  "invalidName",
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

			r := &functionRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID, tt.args.name)

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

func Test_functionRepositoryForFS_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID    string
		function *model.Function
	}

	type expected struct {
		res   *model.Function
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
				function: &functionVTL_2018_05_29,
			},
			expected: expected{
				res:   &functionVTL_2018_05_29,
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
				function: &functionAPPSYNC_JS_1_0_0,
			},
			expected: expected{
				res:   &functionAPPSYNC_JS_1_0_0,
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
				function: &functionVTL_2018_05_29,
			},
			expected: expected{
				res:   &functionVTL_2018_05_29,
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
				function: &functionAPPSYNC_JS_1_0_0,
			},
			expected: expected{
				res:   &functionAPPSYNC_JS_1_0_0,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil function",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:    "apiID",
				function: nil,
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
		{
			name: "edge path: nil name",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:    "apiID",
				function: &model.Function{},
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
				function: &model.Function{
					Name: ptr.Pointer("Name"),
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

			r := &functionRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.function)

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

func Test_functionRepositoryForFS_Delete(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
		name  string
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
				apiID: "apiID",
				name:  *functionVTL_2018_05_29.Name,
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
				apiID: "apiID",
				name:  *functionAPPSYNC_JS_1_0_0.Name,
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
				apiID: "apiID",
				name:  "notExistName",
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

			r := &functionRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			_, err := r.Save(ctx, tt.args.apiID, &functionVTL_2018_05_29)
			assert.NoError(t, err)

			_, err = r.Save(ctx, tt.args.apiID, &functionAPPSYNC_JS_1_0_0)
			assert.NoError(t, err)

			// Act
			err = r.Delete(ctx, tt.args.apiID, tt.args.name)

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
