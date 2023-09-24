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

package syncup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	type args struct {
		ctx context.Context //nolint:containedctx
	}

	type expected struct {
		res string
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: request id was found",
			args: args{
				ctx: context.WithValue(context.Background(), contextKeyRequestID, "RequestID"),
			},
			expected: expected{
				res: "RequestID",
			},
		},
		{
			name: "happy path: request id was not found",
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				res: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			actual := RequestID(tt.args.ctx)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func TestWithRequestID(t *testing.T) {
	type args struct {
		id string
	}

	type expected struct {
		res context.Context //nolint:containedctx
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path",
			args: args{
				id: "RequestID",
			},
			expected: expected{
				res: context.WithValue(context.Background(), contextKeyRequestID, "RequestID"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			// Act
			actual := WithRequestID(ctx, tt.args.id)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}
