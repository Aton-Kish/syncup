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
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	mock_console "github.com/Aton-Kish/syncup/internal/syncup/interface/console/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_mfaTokenProviderRepository_Get(t *testing.T) {
	type mockSurveyPasswordReturn struct {
		res string
		err error
	}
	type mockSurveyPassword struct {
		calls   int
		returns []mockSurveyPasswordReturn
	}

	type expected struct {
		res   string
		errIs error
	}

	tests := []struct {
		name               string
		mockSurveyPassword mockSurveyPassword
		expected           expected
	}{
		{
			name: "happy path",
			mockSurveyPassword: mockSurveyPassword{
				returns: []mockSurveyPasswordReturn{
					{
						res: "123456",
						err: nil,
					},
				},
			},
			expected: expected{
				res:   "123456",
				errIs: nil,
			},
		},
		{
			name: "edge path",
			mockSurveyPassword: mockSurveyPassword{
				returns: []mockSurveyPasswordReturn{
					{
						res: "",
						err: errors.New("error"),
					},
				},
			},
			expected: expected{
				res:   "",
				errIs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSurvey := mock_console.NewMockisurvey(ctrl)

			mockSurvey.
				EXPECT().
				Password(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, prompt *survey.Password, opts ...survey.AskOpt) (string, error) {
					r := tt.mockSurveyPassword.returns[tt.mockSurveyPassword.calls]
					tt.mockSurveyPassword.calls++
					return r.res, r.err
				}).
				Times(len(tt.mockSurveyPassword.returns))

			r := &mfaTokenProviderRepository{
				survey: mockSurvey,
			}
			provider := r.Get(ctx)

			// Act
			actual, err := provider()

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

func Test_mfaTokenProviderRepository_tokenValidator(t *testing.T) {
	type args struct {
		val any
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path",
			args: args{
				val: "123456",
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "edge path: not match pattern",
			args: args{
				val: "abc",
			},
			expected: expected{
				err: errors.New("value must satisfy regular expression pattern: \\d{6}"),
			},
		},
		{
			name: "edge path: not string",
			args: args{
				val: 123456,
			},
			expected: expected{
				err: errors.New("cannot cast value of type int to string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			r := &mfaTokenProviderRepository{}
			validator := r.tokenValidator(ctx)

			// Act
			err := validator(tt.args.val)

			// Assert
			if tt.expected.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}
