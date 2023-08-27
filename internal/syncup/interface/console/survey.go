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

package console

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type isurvey interface {
	InputString(ctx context.Context, prompt *survey.Input, opts ...survey.AskOpt) (string, error)
	InputInt(ctx context.Context, prompt *survey.Input, opts ...survey.AskOpt) (int, error)
	Password(ctx context.Context, prompt *survey.Password, opts ...survey.AskOpt) (string, error)
	Confirm(ctx context.Context, prompt *survey.Confirm, opts ...survey.AskOpt) (bool, error)
	Select(ctx context.Context, prompt *survey.Select, opts ...survey.AskOpt) (string, error)
	MultiSelect(ctx context.Context, prompt *survey.MultiSelect, opts ...survey.AskOpt) ([]string, error)
}

type xsurvey struct {
	options []survey.AskOpt
}

func newSurvey(opts ...survey.AskOpt) isurvey {
	return &xsurvey{
		options: opts,
	}
}

func (s *xsurvey) InputString(ctx context.Context, prompt *survey.Input, opts ...survey.AskOpt) (string, error) {
	return askOne[string](s, ctx, prompt, opts...)
}

func (s *xsurvey) InputInt(ctx context.Context, prompt *survey.Input, opts ...survey.AskOpt) (int, error) {
	return askOne[int](s, ctx, prompt, opts...)
}

func (s *xsurvey) Password(ctx context.Context, prompt *survey.Password, opts ...survey.AskOpt) (string, error) {
	return askOne[string](s, ctx, prompt, opts...)
}

func (s *xsurvey) Confirm(ctx context.Context, prompt *survey.Confirm, opts ...survey.AskOpt) (bool, error) {
	return askOne[bool](s, ctx, prompt, opts...)
}

func (s *xsurvey) Select(ctx context.Context, prompt *survey.Select, opts ...survey.AskOpt) (string, error) {
	return askOne[string](s, ctx, prompt, opts...)
}

func (s *xsurvey) MultiSelect(ctx context.Context, prompt *survey.MultiSelect, opts ...survey.AskOpt) ([]string, error) {
	return askOne[[]string](s, ctx, prompt, opts...)
}

func askOne[T any](s *xsurvey, ctx context.Context, prompt survey.Prompt, opts ...survey.AskOpt) (T, error) {
	optFns := make([]survey.AskOpt, 0, len(s.options)+len(opts))
	optFns = append(optFns, s.options...)
	optFns = append(optFns, opts...)

	var res T
	if err := survey.AskOne(prompt, &res, optFns...); err != nil {
		var zero T

		return zero, err
	}

	return res, nil
}
