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
	"fmt"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

type FunctionService interface {
	Difference(ctx context.Context, functions1, functions2 []model.Function) ([]model.Function, error)
}

type functionService struct {
}

func NewFunctionService(repo repository.Repository) FunctionService {
	return &functionService{}
}

func (s *functionService) Difference(ctx context.Context, functions1, functions2 []model.Function) (res []model.Function, err error) {
	defer wrap(&err)

	encountered := make(map[string]bool)
	for _, fn := range functions2 {
		if fn.Name == nil {
			return nil, fmt.Errorf("%w: missing name", model.ErrNilValue)
		}

		encountered[*fn.Name] = true
	}

	diff := make([]model.Function, 0)
	for _, fn := range functions1 {
		if fn.Name == nil {
			return nil, fmt.Errorf("%w: missing name", model.ErrNilValue)
		}

		if !encountered[*fn.Name] {
			diff = append(diff, fn)
		}
	}

	return diff, nil
}
