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
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_schemaRepositoryForFS_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
	}

	type expected struct {
		res   *model.Schema
		errAs error
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
				baseDir: filepath.Join(testdataBaseDir, "schema"),
			},
			args: args{
				apiID: "apiID",
			},
			expected: expected{
				res:   &schema,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path",
			fields: fields{
				baseDir: filepath.Join(testdataBaseDir, "invalidBaseDir"),
			},
			args: args{
				apiID: "apiID",
			},
			expected: expected{
				res:   nil,
				errAs: &model.LibError{},
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &schemaRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.Get(ctx, tt.args.apiID)

			// Assert
			assert.Equal(t, tt.expected.res, actual)

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

func Test_schemaRepositoryForFS_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	schema := model.Schema(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "schema/schema.graphqls")))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID  string
		schema *model.Schema
	}

	type expected struct {
		res   *model.Schema
		errAs error
		errIs error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "happy path: existing dir",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			expected: expected{
				res:   &schema,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "happy path: non-existing dir",
			fields: fields{
				baseDir: filepath.Join(t.TempDir(), "notExist"),
			},
			args: args{
				apiID:  "apiID",
				schema: &schema,
			},
			expected: expected{
				res:   &schema,
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil schema",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:  "apiID",
				schema: nil,
			},
			expected: expected{
				res:   nil,
				errAs: &model.LibError{},
				errIs: model.ErrNilValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &schemaRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.schema)

			// Assert
			assert.Equal(t, tt.expected.res, actual)

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
