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
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

var (
	mfaTokenPattern = regexp.MustCompile(`\d{6}`)
)

type mfaTokenProviderRepository struct {
	survey isurvey
}

func NewMFATokenProviderRepository() repository.MFATokenProviderRepository {
	return &mfaTokenProviderRepository{
		survey: newSurvey(survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)),
	}
}

func (r *mfaTokenProviderRepository) Get(ctx context.Context) model.MFATokenProvider {
	return func() (res string, err error) {
		defer wrap(&err)

		token, err := r.survey.Password(
			ctx,
			&survey.Password{
				Message: "Enter MFA token code:",
			},
			survey.WithValidator(r.tokenValidator(ctx)),
		)
		if err != nil {
			return "", err
		}

		return token, nil
	}
}

func (r *mfaTokenProviderRepository) tokenValidator(ctx context.Context) survey.Validator {
	return func(val any) error {
		if str, ok := val.(string); ok {
			if !mfaTokenPattern.MatchString(str) {
				return fmt.Errorf("value must satisfy regular expression pattern: %s", mfaTokenPattern)
			}
		} else {
			return fmt.Errorf("cannot cast value of type %v to string", reflect.TypeOf(val).Name())
		}

		return nil
	}
}
