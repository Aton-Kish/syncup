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

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_environmentVariablesRepositoryForFS_Get(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	variables := testhelpers.MustUnmarshalJSON[model.EnvironmentVariables](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "environment_variables/env.json")))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID string
	}

	type expected struct {
		res   model.EnvironmentVariables
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
				baseDir: filepath.Join(testdataBaseDir, "environment_variables"),
			},
			args: args{
				apiID: "apiID",
			},
			expected: expected{
				res:   variables,
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
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &environmentVariablesRepositoryForFS{
				baseDir: tt.fields.baseDir,
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

func Test_environmentVariablesRepositoryForFS_Save(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	variables := testhelpers.MustUnmarshalJSON[model.EnvironmentVariables](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "environment_variables/env.json")))

	type fields struct {
		baseDir string
	}

	type args struct {
		apiID     string
		variables model.EnvironmentVariables
	}

	type expected struct {
		res   model.EnvironmentVariables
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
				apiID:     "apiID",
				variables: variables,
			},
			expected: expected{
				res:   variables,
				errIs: nil,
			},
		},
		{
			name: "happy path: non-existing dir",
			fields: fields{
				baseDir: filepath.Join(t.TempDir(), "notExist"),
			},
			args: args{
				apiID:     "apiID",
				variables: variables,
			},
			expected: expected{
				res:   variables,
				errIs: nil,
			},
		},
		{
			name: "edge path: nil environment variables",
			fields: fields{
				baseDir: t.TempDir(),
			},
			args: args{
				apiID:     "apiID",
				variables: nil,
			},
			expected: expected{
				res:   nil,
				errIs: model.ErrNilValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &environmentVariablesRepositoryForFS{
				baseDir: tt.fields.baseDir,
			}

			// Act
			actual, err := r.Save(ctx, tt.args.apiID, tt.args.variables)

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
