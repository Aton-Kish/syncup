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

package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/stretchr/testify/assert"
)

func Test_rootCommand_Execute(t *testing.T) {
	type args struct {
		args []string
	}

	type expected struct {
		errAs error
		errIs error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: default",
			args: args{
				args: []string{},
			},
			expected: expected{
				errAs: nil,
				errIs: nil,
			},
		},
		{
			name: "happy path: version flag",
			args: args{
				args: []string{"-v"},
			},
			expected: expected{
				errAs: nil,
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			stdin := new(bytes.Reader)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			c := &rootCommand{
				options: newOptions(WithStdio(stdin, stdout, stderr)),
				version: &model.Version{
					Version:   "Version",
					GitCommit: "GitCommit",
					GoVersion: "GoVersion",
					OS:        "OS",
					Arch:      "Arch",
					BuildTime: "BuildTime",
				},
			}

			// Act
			err := c.Execute(ctx, tt.args.args...)

			// Assert
			if tt.expected.errAs == nil && tt.expected.errIs == nil {
				assert.NoError(t, err)

				assert.Equal(t, 0, stdin.Len())
				assert.Equal(t, 0, stdout.Len())
				assert.Greater(t, stderr.Len(), 0)
			} else {
				if tt.expected.errAs != nil {
					assert.ErrorAs(t, err, &tt.expected.errAs)
				}

				if tt.expected.errIs != nil {
					assert.ErrorIs(t, err, tt.expected.errIs)
				}

				assert.Equal(t, 0, stdin.Len())
				assert.Equal(t, 0, stdout.Len())
				assert.Greater(t, stderr.Len(), 0)
			}
		})
	}
}
