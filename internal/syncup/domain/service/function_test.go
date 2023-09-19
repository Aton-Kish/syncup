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

func Test_functionService_Difference(t *testing.T) {
	testdataBaseDir := "../../../../testdata"
	functionVTL_2018_05_29 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/metadata.json")))
	functionVTL_2018_05_29.RequestMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/request.vtl"))))
	functionVTL_2018_05_29.ResponseMappingTemplate = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/VTL_2018-05-29/response.vtl"))))
	functionAPPSYNC_JS_1_0_0 := testhelpers.MustUnmarshalJSON[model.Function](t, testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/metadata.json")))
	functionAPPSYNC_JS_1_0_0.Code = ptr.Pointer(string(testhelpers.MustReadFile(t, filepath.Join(testdataBaseDir, "functions/APPSYNC_JS_1.0.0/code.js"))))

	type args struct {
		functions1 []model.Function
		functions2 []model.Function
	}

	type expected struct {
		res   []model.Function
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
				functions1: []model.Function{
					functionVTL_2018_05_29,
					functionAPPSYNC_JS_1_0_0,
				},
				functions2: []model.Function{
					functionVTL_2018_05_29,
					functionAPPSYNC_JS_1_0_0,
				},
			},
			expected: expected{
				res:   []model.Function{},
				errIs: nil,
			},
		},
		{
			name: "happy path: some differences",
			args: args{
				functions1: []model.Function{
					functionVTL_2018_05_29,
					functionAPPSYNC_JS_1_0_0,
				},
				functions2: []model.Function{
					functionAPPSYNC_JS_1_0_0,
				},
			},
			expected: expected{
				res: []model.Function{
					functionVTL_2018_05_29,
				},
				errIs: nil,
			},
		},
		{
			name: "happy path: everything is different",
			args: args{
				functions1: []model.Function{
					functionVTL_2018_05_29,
				},
				functions2: []model.Function{
					functionAPPSYNC_JS_1_0_0,
				},
			},
			expected: expected{
				res: []model.Function{
					functionVTL_2018_05_29,
				},
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			s := &functionService{}

			// Act
			actual, err := s.Difference(ctx, tt.args.functions1, tt.args.functions2)

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
