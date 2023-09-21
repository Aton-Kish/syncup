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

package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure/mapper"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

type functionRepositoryForAppSync struct {
	appsyncClient appsyncClient
}

var (
	_ interface {
		repository.AWSActivator
	} = (*functionRepositoryForAppSync)(nil)
)

func NewFunctionRepositoryForAppSync() repository.FunctionRepository {
	return &functionRepositoryForAppSync{}
}

func (r *functionRepositoryForAppSync) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) (err error) {
	defer wrap(&err)

	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *functionRepositoryForAppSync) List(ctx context.Context, apiID string) (res []model.Function, err error) {
	defer wrap(&err)

	fns := make([]model.Function, 0)

	var token *string
	for {
		out, err := r.appsyncClient.ListFunctions(
			ctx,
			&appsync.ListFunctionsInput{
				ApiId:     &apiID,
				NextToken: token,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, fn := range out.Functions {
			fns = append(fns, *mapper.NewFunctionMapper().ToModel(ctx, &fn))
		}

		token = out.NextToken
		if token == nil {
			break
		}
	}

	encountered := make(map[string]bool)
	for _, fn := range fns {
		if fn.Name == nil {
			return nil, fmt.Errorf("%w: missing name", model.ErrNilValue)
		}

		name := *fn.Name
		if encountered[name] {
			return nil, fmt.Errorf("%w: function name %s", model.ErrDuplicateValue, name)
		}

		encountered[name] = true
	}

	return fns, nil
}

func (r *functionRepositoryForAppSync) Get(ctx context.Context, apiID string, name string) (res *model.Function, err error) {
	defer wrap(&err)

	fns, err := r.List(ctx, apiID)
	if err != nil {
		return nil, err
	}

	for _, fn := range fns {
		if *fn.Name == name {
			return &fn, nil
		}
	}

	return nil, model.ErrNotFound
}

func (r *functionRepositoryForAppSync) Save(ctx context.Context, apiID string, function *model.Function) (res *model.Function, err error) {
	defer wrap(&err)

	if function == nil {
		return nil, fmt.Errorf("%w: missing arguments in save function method", model.ErrNilValue)
	}

	if function.Name == nil {
		return nil, fmt.Errorf("%w: missing name", model.ErrNilValue)
	}

	save := r.update
	fnToSave := function
	if fn, err := r.Get(ctx, apiID, *function.Name); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			save = r.create
		} else {
			return nil, err
		}
	} else {
		fnToSave = fn
	}

	fn, err := save(ctx, apiID, fnToSave)
	if err != nil {
		return nil, err
	}

	return fn, nil
}

func (r *functionRepositoryForAppSync) create(ctx context.Context, apiID string, function *model.Function) (res *model.Function, err error) {
	defer wrap(&err)

	if function == nil {
		return nil, fmt.Errorf("%w: missing arguments in create function method", model.ErrNilValue)
	}

	f := mapper.NewFunctionMapper().FromModel(ctx, function)
	out, err := r.appsyncClient.CreateFunction(
		ctx,
		&appsync.CreateFunctionInput{
			ApiId:                   &apiID,
			Name:                    f.Name,
			Description:             f.Description,
			DataSourceName:          f.DataSourceName,
			RequestMappingTemplate:  f.RequestMappingTemplate,
			ResponseMappingTemplate: f.ResponseMappingTemplate,
			FunctionVersion:         f.FunctionVersion,
			SyncConfig:              f.SyncConfig,
			MaxBatchSize:            f.MaxBatchSize,
			Runtime:                 f.Runtime,
			Code:                    f.Code,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, err
	}

	fn := mapper.NewFunctionMapper().ToModel(ctx, out.FunctionConfiguration)
	if fn == nil {
		return nil, fmt.Errorf("%w: missing function in AppSync CreateFunction API response", model.ErrNilValue)
	}

	return fn, nil
}

func (r *functionRepositoryForAppSync) update(ctx context.Context, apiID string, function *model.Function) (res *model.Function, err error) {
	defer wrap(&err)

	if function == nil {
		return nil, fmt.Errorf("%w: missing arguments in update function method", model.ErrNilValue)
	}

	f := mapper.NewFunctionMapper().FromModel(ctx, function)
	out, err := r.appsyncClient.UpdateFunction(
		ctx,
		&appsync.UpdateFunctionInput{
			ApiId:                   &apiID,
			FunctionId:              f.FunctionId,
			Name:                    f.Name,
			Description:             f.Description,
			DataSourceName:          f.DataSourceName,
			RequestMappingTemplate:  f.RequestMappingTemplate,
			ResponseMappingTemplate: f.ResponseMappingTemplate,
			FunctionVersion:         f.FunctionVersion,
			SyncConfig:              f.SyncConfig,
			MaxBatchSize:            f.MaxBatchSize,
			Runtime:                 f.Runtime,
			Code:                    f.Code,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, err
	}

	fn := mapper.NewFunctionMapper().ToModel(ctx, out.FunctionConfiguration)
	if fn == nil {
		return nil, fmt.Errorf("%w: missing function in AppSync UpdateFunction API response", model.ErrNilValue)
	}

	return fn, nil
}

func (r *functionRepositoryForAppSync) Delete(ctx context.Context, apiID string, name string) (err error) {
	defer wrap(&err)

	fn, err := r.Get(ctx, apiID, name)
	if err != nil {
		return err
	}

	if _, err := r.appsyncClient.DeleteFunction(
		ctx,
		&appsync.DeleteFunctionInput{
			ApiId:      &apiID,
			FunctionId: fn.FunctionId,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	); err != nil {
		return err
	}

	return nil
}
