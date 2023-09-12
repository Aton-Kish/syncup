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

package mapper

import (
	"context"
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/stretchr/testify/assert"
)

func Test_runtimeMapper_ToModel(t *testing.T) {
	type args struct {
		v *types.AppSyncRuntime
	}

	type expected struct {
		res *model.Runtime
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: nil",
			args: args{
				v: nil,
			},
			expected: expected{
				res: nil,
			},
		},
		{
			name: "happy path: not nil",
			args: args{
				v: &types.AppSyncRuntime{
					Name:           types.RuntimeName("Name"),
					RuntimeVersion: aws.String("RuntimeVersion"),
				},
			},
			expected: expected{
				res: &model.Runtime{
					Name:           model.RuntimeName("Name"),
					RuntimeVersion: ptr.Pointer("RuntimeVersion"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			// Act
			actual := (*runtimeMapper)(nil).ToModel(ctx, tt.args.v)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func Test_runtimeMapper_FromModel(t *testing.T) {
	type args struct {
		v *model.Runtime
	}

	type expected struct {
		res *types.AppSyncRuntime
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: nil",
			args: args{
				v: nil,
			},
			expected: expected{
				res: nil,
			},
		},
		{
			name: "happy path: not nil",
			args: args{
				v: &model.Runtime{
					Name:           model.RuntimeName("Name"),
					RuntimeVersion: ptr.Pointer("RuntimeVersion"),
				},
			},
			expected: expected{
				res: &types.AppSyncRuntime{
					Name:           types.RuntimeName("Name"),
					RuntimeVersion: aws.String("RuntimeVersion"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			// Act
			actual := (*runtimeMapper)(nil).FromModel(ctx, tt.args.v)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}
