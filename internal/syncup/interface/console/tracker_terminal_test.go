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
	"fmt"
	"testing"

	mock_console "github.com/Aton-Kish/syncup/internal/syncup/interface/console/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_trackerRepositoryForTerminal_InProgress(t *testing.T) {
	type args struct {
		msg string
	}

	type mockSpinnerSetSuffixReturn struct {
	}
	type mockSpinnerSetSuffix struct {
		calls   int
		returns []mockSpinnerSetSuffixReturn
	}

	type mockSpinnerStartReturn struct {
	}
	type mockSpinnerStart struct {
		calls   int
		returns []mockSpinnerStartReturn
	}

	tests := []struct {
		name                 string
		args                 args
		mockSpinnerSetSuffix mockSpinnerSetSuffix
		mockSpinnerStart     mockSpinnerStart
	}{
		{
			name: "happy path",
			args: args{
				msg: "hello",
			},
			mockSpinnerSetSuffix: mockSpinnerSetSuffix{
				returns: []mockSpinnerSetSuffixReturn{
					{},
				},
			},
			mockSpinnerStart: mockSpinnerStart{
				returns: []mockSpinnerStartReturn{
					{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpinner := mock_console.NewMockispinner(ctrl)

			mockSpinner.
				EXPECT().
				SetSuffix(gomock.Any()).
				DoAndReturn(func(suffix string) {
					tt.mockSpinnerSetSuffix.calls++
				}).
				Times(len(tt.mockSpinnerSetSuffix.returns))

			mockSpinner.
				EXPECT().
				Start().
				DoAndReturn(func() {
					tt.mockSpinnerStart.calls++
					fmt.Fprint(w, fmt.Sprintln("spinner", tt.args.msg))
				}).
				Times(len(tt.mockSpinnerStart.returns))

			r := &trackerRepositoryForTerminal{
				writer:  w,
				spinner: mockSpinner,
			}

			// Act
			r.InProgress(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}

func Test_trackerRepositoryForTerminal_Failed(t *testing.T) {
	type args struct {
		msg string
	}

	type mockSpinnerActiveReturn struct {
		res bool
	}
	type mockSpinnerActive struct {
		calls   int
		returns []mockSpinnerActiveReturn
	}

	type mockSpinnerSetFinalMsgReturn struct {
	}
	type mockSpinnerSetFinalMsg struct {
		calls   int
		returns []mockSpinnerSetFinalMsgReturn
	}

	type mockSpinnerStopReturn struct {
	}
	type mockSpinnerStop struct {
		calls   int
		returns []mockSpinnerStopReturn
	}

	tests := []struct {
		name                   string
		args                   args
		mockSpinnerActive      mockSpinnerActive
		mockSpinnerSetFinalMsg mockSpinnerSetFinalMsg
		mockSpinnerStop        mockSpinnerStop
	}{
		{
			name: "happy path: active spinner",
			args: args{
				msg: "hello",
			},
			mockSpinnerActive: mockSpinnerActive{
				returns: []mockSpinnerActiveReturn{
					{
						res: true,
					},
				},
			},
			mockSpinnerSetFinalMsg: mockSpinnerSetFinalMsg{
				returns: []mockSpinnerSetFinalMsgReturn{
					{},
				},
			},
			mockSpinnerStop: mockSpinnerStop{
				returns: []mockSpinnerStopReturn{
					{},
				},
			},
		},
		{
			name: "happy path: inactive spinner",
			args: args{
				msg: "hello",
			},
			mockSpinnerActive: mockSpinnerActive{
				returns: []mockSpinnerActiveReturn{
					{
						res: false,
					},
				},
			},
			mockSpinnerSetFinalMsg: mockSpinnerSetFinalMsg{
				returns: []mockSpinnerSetFinalMsgReturn{},
			},
			mockSpinnerStop: mockSpinnerStop{
				returns: []mockSpinnerStopReturn{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpinner := mock_console.NewMockispinner(ctrl)

			mockSpinner.
				EXPECT().
				Active().
				DoAndReturn(func() bool {
					r := tt.mockSpinnerActive.returns[tt.mockSpinnerActive.calls]
					tt.mockSpinnerActive.calls++
					return r.res
				}).
				Times(len(tt.mockSpinnerActive.returns))

			mockSpinner.
				EXPECT().
				SetFinalMsg(gomock.Any()).
				DoAndReturn(func(suffix string) {
					tt.mockSpinnerSetFinalMsg.calls++
				}).
				Times(len(tt.mockSpinnerSetFinalMsg.returns))

			mockSpinner.
				EXPECT().
				Stop().
				DoAndReturn(func() {
					tt.mockSpinnerStop.calls++
					fmt.Fprint(w, fmt.Sprintln("spinner", tt.args.msg))
				}).
				Times(len(tt.mockSpinnerStop.returns))

			r := &trackerRepositoryForTerminal{
				writer:  w,
				spinner: mockSpinner,
			}

			// Act
			r.Failed(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}

func Test_trackerRepositoryForTerminal_Success(t *testing.T) {
	type args struct {
		msg string
	}

	type mockSpinnerActiveReturn struct {
		res bool
	}
	type mockSpinnerActive struct {
		calls   int
		returns []mockSpinnerActiveReturn
	}

	type mockSpinnerSetFinalMsgReturn struct {
	}
	type mockSpinnerSetFinalMsg struct {
		calls   int
		returns []mockSpinnerSetFinalMsgReturn
	}

	type mockSpinnerStopReturn struct {
	}
	type mockSpinnerStop struct {
		calls   int
		returns []mockSpinnerStopReturn
	}

	tests := []struct {
		name                   string
		args                   args
		mockSpinnerActive      mockSpinnerActive
		mockSpinnerSetFinalMsg mockSpinnerSetFinalMsg
		mockSpinnerStop        mockSpinnerStop
	}{
		{
			name: "happy path: active spinner",
			args: args{
				msg: "hello",
			},
			mockSpinnerActive: mockSpinnerActive{
				returns: []mockSpinnerActiveReturn{
					{
						res: true,
					},
				},
			},
			mockSpinnerSetFinalMsg: mockSpinnerSetFinalMsg{
				returns: []mockSpinnerSetFinalMsgReturn{
					{},
				},
			},
			mockSpinnerStop: mockSpinnerStop{
				returns: []mockSpinnerStopReturn{
					{},
				},
			},
		},
		{
			name: "happy path: inactive spinner",
			args: args{
				msg: "hello",
			},
			mockSpinnerActive: mockSpinnerActive{
				returns: []mockSpinnerActiveReturn{
					{
						res: false,
					},
				},
			},
			mockSpinnerSetFinalMsg: mockSpinnerSetFinalMsg{
				returns: []mockSpinnerSetFinalMsgReturn{},
			},
			mockSpinnerStop: mockSpinnerStop{
				returns: []mockSpinnerStopReturn{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			w := new(bytes.Buffer)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpinner := mock_console.NewMockispinner(ctrl)

			mockSpinner.
				EXPECT().
				Active().
				DoAndReturn(func() bool {
					r := tt.mockSpinnerActive.returns[tt.mockSpinnerActive.calls]
					tt.mockSpinnerActive.calls++
					return r.res
				}).
				Times(len(tt.mockSpinnerActive.returns))

			mockSpinner.
				EXPECT().
				SetFinalMsg(gomock.Any()).
				DoAndReturn(func(suffix string) {
					tt.mockSpinnerSetFinalMsg.calls++
				}).
				Times(len(tt.mockSpinnerSetFinalMsg.returns))

			mockSpinner.
				EXPECT().
				Stop().
				DoAndReturn(func() {
					tt.mockSpinnerStop.calls++
					fmt.Fprint(w, fmt.Sprintln("spinner", tt.args.msg))
				}).
				Times(len(tt.mockSpinnerStop.returns))

			r := &trackerRepositoryForTerminal{
				writer:  w,
				spinner: mockSpinner,
			}

			// Act
			r.Success(ctx, tt.args.msg)

			// Assert
			assert.Greater(t, w.Len(), 0)
		})
	}
}
