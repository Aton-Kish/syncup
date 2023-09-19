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

type ResolverService interface {
	Difference(ctx context.Context, resolvers1, resolvers2 []model.Resolver) ([]model.Resolver, error)

	ResolvePipelineConfigFunctionIDs(ctx context.Context, resolver *model.Resolver, functions []model.Function) error
	ResolvePipelineConfigFunctionNames(ctx context.Context, resolver *model.Resolver, functions []model.Function) error
}

type resolverService struct {
}

func NewResolverService(repo repository.Repository) ResolverService {
	return &resolverService{}
}

func (s *resolverService) Difference(ctx context.Context, resolvers1, resolvers2 []model.Resolver) (res []model.Resolver, err error) {
	defer wrap(&err)

	encountered := make(map[string]bool)
	for _, rslv := range resolvers2 {
		if rslv.TypeName == nil {
			return nil, fmt.Errorf("%w: missing type name", model.ErrNilValue)
		}

		if rslv.FieldName == nil {
			return nil, fmt.Errorf("%w: missing field name", model.ErrNilValue)
		}

		encountered[fmt.Sprintf("%s.%s", *rslv.TypeName, *rslv.FieldName)] = true
	}

	diff := make([]model.Resolver, 0)
	for _, rslv := range resolvers1 {
		if rslv.TypeName == nil {
			return nil, fmt.Errorf("%w: missing type name", model.ErrNilValue)
		}

		if rslv.FieldName == nil {
			return nil, fmt.Errorf("%w: missing field name", model.ErrNilValue)
		}

		if !encountered[fmt.Sprintf("%s.%s", *rslv.TypeName, *rslv.FieldName)] {
			diff = append(diff, rslv)
		}
	}

	return diff, nil
}

func (s *resolverService) ResolvePipelineConfigFunctionIDs(ctx context.Context, resolver *model.Resolver, functions []model.Function) (err error) {
	defer wrap(&err)

	if resolver == nil {
		return fmt.Errorf("%w: missing arguments in ResolverService.ResolvePipelineConfigFunctionIDs method", model.ErrNilValue)
	}

	if resolver.PipelineConfig == nil {
		return nil
	}

	nameToID := make(map[string]string)
	for _, fn := range functions {
		if fn.FunctionId == nil {
			return fmt.Errorf("%w: missing function id", model.ErrNilValue)
		}

		if fn.Name == nil {
			return fmt.Errorf("%w: missing name", model.ErrNilValue)
		}

		if _, ok := nameToID[*fn.Name]; ok {
			return fmt.Errorf("%w: function name %s", model.ErrDuplicateValue, *fn.Name)
		}

		nameToID[*fn.Name] = *fn.FunctionId
	}

	ids := make([]string, 0, len(resolver.PipelineConfig.FunctionNames))
	for _, name := range resolver.PipelineConfig.FunctionNames {
		id, ok := nameToID[name]
		if !ok {
			return fmt.Errorf("%w: name %s", model.ErrNotFound, name)
		}

		ids = append(ids, id)
	}

	resolver.PipelineConfig.Functions = ids

	return nil
}

func (s *resolverService) ResolvePipelineConfigFunctionNames(ctx context.Context, resolver *model.Resolver, functions []model.Function) (err error) {
	defer wrap(&err)

	if resolver == nil {
		return fmt.Errorf("%w: missing arguments in ResolverService.ResolvePipelineConfigFunctionNames method", model.ErrNilValue)
	}

	if resolver.PipelineConfig == nil {
		return nil
	}

	idToName := make(map[string]string)
	for _, fn := range functions {
		if fn.FunctionId == nil {
			return fmt.Errorf("%w: missing function id", model.ErrNilValue)
		}

		if fn.Name == nil {
			return fmt.Errorf("%w: missing name", model.ErrNilValue)
		}

		idToName[*fn.FunctionId] = *fn.Name
	}

	names := make([]string, 0, len(resolver.PipelineConfig.Functions))
	for _, id := range resolver.PipelineConfig.Functions {
		name, ok := idToName[id]
		if !ok {
			return fmt.Errorf("%w: id %s", model.ErrNotFound, id)
		}

		names = append(names, name)
	}

	resolver.PipelineConfig.FunctionNames = names

	return nil
}
