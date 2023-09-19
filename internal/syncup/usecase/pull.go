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

package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/service"
)

type PullInput struct {
	APIID string
}

type PullOutput struct {
}

type PullUseCase interface {
	Execute(ctx context.Context, params *PullInput) (*PullOutput, error)
}

type pullUseCase struct {
	resolverService              service.ResolverService
	trackerRepository            repository.TrackerRepository
	schemaRepositoryForAppSync   repository.SchemaRepository
	schemaRepositoryForFS        repository.SchemaRepository
	functionRepositoryForAppSync repository.FunctionRepository
	functionRepositoryForFS      repository.FunctionRepository
	resolverRepositoryForAppSync repository.ResolverRepository
	resolverRepositoryForFS      repository.ResolverRepository
}

func NewPullUseCase(repo repository.Repository) PullUseCase {
	return &pullUseCase{
		resolverService:              service.NewResolverService(repo),
		trackerRepository:            repo.TrackerRepository(),
		schemaRepositoryForAppSync:   repo.SchemaRepositoryForAppSync(),
		schemaRepositoryForFS:        repo.SchemaRepositoryForFS(),
		functionRepositoryForAppSync: repo.FunctionRepositoryForAppSync(),
		functionRepositoryForFS:      repo.FunctionRepositoryForFS(),
		resolverRepositoryForAppSync: repo.ResolverRepositoryForAppSync(),
		resolverRepositoryForFS:      repo.ResolverRepositoryForFS(),
	}
}

func (uc *pullUseCase) Execute(ctx context.Context, params *PullInput) (res *PullOutput, err error) {
	defer wrap(&err)

	if _, err := uc.pullSchema(ctx, params.APIID); err != nil {
		return nil, err
	}

	fns, err := uc.pullFunctions(ctx, params.APIID)
	if err != nil {
		return nil, err
	}

	if _, err := uc.pullResolvers(ctx, params.APIID, fns); err != nil {
		return nil, err
	}

	return &PullOutput{}, nil
}

func (uc *pullUseCase) pullSchema(ctx context.Context, apiID string) (res *model.Schema, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "fetching schema")

	schema, err := uc.schemaRepositoryForAppSync.Get(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to fetch schema")
		return nil, err
	}

	if _, err := uc.schemaRepositoryForFS.Save(ctx, apiID, schema); err != nil {
		uc.trackerRepository.Failed(ctx, "failed to save schema")
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "saved schema")

	return schema, nil
}

func (uc *pullUseCase) pullFunctions(ctx context.Context, apiID string) (res []model.Function, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "fetching functions")

	functions, err := uc.functionRepositoryForAppSync.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to fetch functions")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "saving functions")

	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)

	for _, fn := range functions {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()

			if _, err := uc.functionRepositoryForFS.Save(ctx, apiID, &fn); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to save function %s", ptr.ToValue(fn.Name)))
				return
			}

			uc.trackerRepository.Success(ctx, fmt.Sprintf("saved function %s", ptr.ToValue(fn.Name)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "saved all functions")

	return functions, nil
}

func (uc *pullUseCase) pullResolvers(ctx context.Context, apiID string, functions []model.Function) (res []model.Resolver, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "fetching resolvers")

	rslvs, err := uc.resolverRepositoryForAppSync.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to fetch resolvers")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "saving resolvers")

	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)

	for _, rslv := range rslvs {
		rslv := rslv
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := uc.resolverService.ResolvePipelineConfigFunctionNames(ctx, &rslv, functions); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to save resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
				return
			}

			if _, err := uc.resolverRepositoryForFS.Save(ctx, apiID, &rslv); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to save resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
				return
			}

			uc.trackerRepository.Success(ctx, fmt.Sprintf("saved resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "saved all resolvers")

	return rslvs, nil
}
