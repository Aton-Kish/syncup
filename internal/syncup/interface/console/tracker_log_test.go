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

package console

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trackerRepositoryForLog_InProgress(t *testing.T) {
	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				msg: "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			r := &trackerRepositoryForLog{
				logger: slog.New(slog.NewJSONHandler(w, nil)).With(slog.String("app", "syncup"), slog.String("category", "tracker")),
			}

			// Act
			r.InProgress(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}

func Test_trackerRepositoryForLog_Failed(t *testing.T) {
	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				msg: "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			r := &trackerRepositoryForLog{
				logger: slog.New(slog.NewJSONHandler(w, nil)).With(slog.String("app", "syncup"), slog.String("category", "tracker")),
			}

			// Act
			r.Failed(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}

func Test_trackerRepositoryForLog_Success(t *testing.T) {
	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				msg: "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			r := &trackerRepositoryForLog{
				logger: slog.New(slog.NewJSONHandler(w, nil)).With(slog.String("app", "syncup"), slog.String("category", "tracker")),
			}

			// Act
			r.Success(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}
