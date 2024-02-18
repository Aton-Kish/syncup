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

type PushInput struct {
	APIID                     string
	DeleteExtraneousResources bool
}

type PushOutput struct {
}

type PushUseCase interface {
	Execute(ctx context.Context, params *PushInput) (*PushOutput, error)
}

type pushUseCase struct {
	functionService                          service.FunctionService
	resolverService                          service.ResolverService
	trackerRepository                        repository.TrackerRepository
	environmentVariablesRepositoryForAppSync repository.EnvironmentVariablesRepository
	environmentVariablesRepositoryForFS      repository.EnvironmentVariablesRepository
	schemaRepositoryForAppSync               repository.SchemaRepository
	schemaRepositoryForFS                    repository.SchemaRepository
	functionRepositoryForAppSync             repository.FunctionRepository
	functionRepositoryForFS                  repository.FunctionRepository
	resolverRepositoryForAppSync             repository.ResolverRepository
	resolverRepositoryForFS                  repository.ResolverRepository
}

func NewPushUseCase(repo repository.Repository) PushUseCase {
	return &pushUseCase{
		functionService:                          service.NewFunctionService(repo),
		resolverService:                          service.NewResolverService(repo),
		trackerRepository:                        repo.TrackerRepository(),
		environmentVariablesRepositoryForAppSync: repo.EnvironmentVariablesRepositoryForAppSync(),
		environmentVariablesRepositoryForFS:      repo.EnvironmentVariablesRepositoryForFS(),
		schemaRepositoryForAppSync:               repo.SchemaRepositoryForAppSync(),
		schemaRepositoryForFS:                    repo.SchemaRepositoryForFS(),
		functionRepositoryForAppSync:             repo.FunctionRepositoryForAppSync(),
		functionRepositoryForFS:                  repo.FunctionRepositoryForFS(),
		resolverRepositoryForAppSync:             repo.ResolverRepositoryForAppSync(),
		resolverRepositoryForFS:                  repo.ResolverRepositoryForFS(),
	}
}

func (uc *pushUseCase) Execute(ctx context.Context, params *PushInput) (res *PushOutput, err error) {
	defer wrap(&err)

	if _, err := uc.pushEnvironmentVariables(ctx, params.APIID); err != nil {
		return nil, err
	}

	if _, err := uc.pushSchema(ctx, params.APIID); err != nil {
		return nil, err
	}

	fns, err := uc.pushFunctions(ctx, params.APIID)
	if err != nil {
		return nil, err
	}

	rslvs, err := uc.pushResolvers(ctx, params.APIID, fns)
	if err != nil {
		return nil, err
	}

	if params.DeleteExtraneousResources {
		if err := uc.deleteExtraneousFunctions(ctx, params.APIID, fns); err != nil {
			return nil, err
		}

		if err := uc.deleteExtraneousResolvers(ctx, params.APIID, rslvs); err != nil {
			return nil, err
		}
	}

	return &PushOutput{}, nil
}

func (uc *pushUseCase) pushEnvironmentVariables(ctx context.Context, apiID string) (res model.EnvironmentVariables, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "loading environment variables")

	vs, err := uc.environmentVariablesRepositoryForFS.Get(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to load environment variables")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "pushing environment variables")

	varibales, err := uc.environmentVariablesRepositoryForAppSync.Save(ctx, apiID, vs)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to push environment variables")
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "pushed environment variables")

	return varibales, nil
}

func (uc *pushUseCase) pushSchema(ctx context.Context, apiID string) (res *model.Schema, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "loading schema")

	s, err := uc.schemaRepositoryForFS.Get(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to load schema")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "pushing schema")

	schema, err := uc.schemaRepositoryForAppSync.Save(ctx, apiID, s)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to push schema")
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "pushed schema")

	return schema, nil
}

func (uc *pushUseCase) pushFunctions(ctx context.Context, apiID string) (res []model.Function, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "loading functions")

	fns, err := uc.functionRepositoryForFS.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to load functions")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "pushing functions")

	var mu sync.Mutex
	var wg sync.WaitGroup
	functions := make([]model.Function, 0, len(fns))
	errs := make([]error, 0)

	for _, fn := range fns {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()

			function, err := uc.functionRepositoryForAppSync.Save(ctx, apiID, &fn)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to push function %s", ptr.ToValue(fn.Name)))
				return
			}

			mu.Lock()
			functions = append(functions, *function)
			mu.Unlock()

			uc.trackerRepository.Success(ctx, fmt.Sprintf("pushed function %s", ptr.ToValue(function.Name)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "pushed all functions")

	return functions, nil
}

func (uc *pushUseCase) pushResolvers(ctx context.Context, apiID string, functions []model.Function) (res []model.Resolver, err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "loading resolvers")

	rslvs, err := uc.resolverRepositoryForFS.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to load resolvers")
		return nil, err
	}

	uc.trackerRepository.InProgress(ctx, "pushing resolvers")

	var mu sync.Mutex
	var wg sync.WaitGroup
	resolvers := make([]model.Resolver, 0, len(rslvs))
	errs := make([]error, 0)

	for _, rslv := range rslvs {
		rslv := rslv
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := uc.resolverService.ResolvePipelineConfigFunctionIDs(ctx, &rslv, functions); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to push resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
				return
			}

			resovler, err := uc.resolverRepositoryForAppSync.Save(ctx, apiID, &rslv)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to push resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
				return
			}

			mu.Lock()
			resolvers = append(resolvers, *resovler)
			mu.Unlock()

			uc.trackerRepository.Success(ctx, fmt.Sprintf("pushed resolver %s.%s", ptr.ToValue(resovler.TypeName), ptr.ToValue(resovler.FieldName)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	uc.trackerRepository.Success(ctx, "pushed all resolvers")

	return resolvers, nil
}

func (uc *pushUseCase) deleteExtraneousFunctions(ctx context.Context, apiID string, functions []model.Function) (err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "fetching functions")

	fns, err := uc.functionRepositoryForAppSync.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to fetch functions")
		return err
	}

	extraneousFns, err := uc.functionService.Difference(ctx, fns, functions)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to retrieve extraneous functions")
		return err
	}

	if len(extraneousFns) == 0 {
		uc.trackerRepository.Success(ctx, "there were no extraneous functions")
		return nil
	}

	uc.trackerRepository.InProgress(ctx, "deleting extraneous functions")

	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)

	for _, fn := range extraneousFns {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := uc.functionRepositoryForAppSync.Delete(ctx, apiID, *fn.Name); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to delete extraneous function %s", ptr.ToValue(fn.Name)))
				return
			}

			uc.trackerRepository.Success(ctx, fmt.Sprintf("deleted extraneous function %s", ptr.ToValue(fn.Name)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return err
	}

	uc.trackerRepository.Success(ctx, "deleted all extraneous functions")

	return nil
}

func (uc *pushUseCase) deleteExtraneousResolvers(ctx context.Context, apiID string, resolvers []model.Resolver) (err error) {
	defer wrap(&err)

	uc.trackerRepository.InProgress(ctx, "fetching resolvers")

	rslvs, err := uc.resolverRepositoryForAppSync.List(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to fetch resolvers")
		return err
	}

	extraneousRslvs, err := uc.resolverService.Difference(ctx, rslvs, resolvers)
	if err != nil {
		uc.trackerRepository.Failed(ctx, "failed to retrieve extraneous resolvers")
		return err
	}

	if len(extraneousRslvs) == 0 {
		uc.trackerRepository.Success(ctx, "there were no extraneous resolvers")
		return nil
	}

	uc.trackerRepository.InProgress(ctx, "deleting extraneous resolvers")

	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)

	for _, rslv := range extraneousRslvs {
		rslv := rslv
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := uc.resolverRepositoryForAppSync.Delete(ctx, apiID, *rslv.TypeName, *rslv.FieldName); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()

				uc.trackerRepository.Failed(ctx, fmt.Sprintf("failed to delete extraneous resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
				return
			}

			uc.trackerRepository.Success(ctx, fmt.Sprintf("deleted extraneous resolver %s.%s", ptr.ToValue(rslv.TypeName), ptr.ToValue(rslv.FieldName)))
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return err
	}

	uc.trackerRepository.Success(ctx, "deleted all extraneous resolvers")

	return nil
}
