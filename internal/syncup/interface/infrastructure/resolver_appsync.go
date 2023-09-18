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
	"sync"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure/mapper"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

type resolverRepositoryForAppSync struct {
	appsyncClient appsyncClient
}

var (
	_ interface {
		repository.AWSActivator
	} = (*resolverRepositoryForAppSync)(nil)
)

func NewResolverRepositoryForAppSync() repository.ResolverRepository {
	return &resolverRepositoryForAppSync{}
}

func (r *resolverRepositoryForAppSync) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) (err error) {
	defer wrap(&err)

	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *resolverRepositoryForAppSync) List(ctx context.Context, apiID string) (res []model.Resolver, err error) {
	defer wrap(&err)

	ns, err := r.listTypeNames(ctx, apiID)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	rslvs := make([]model.Resolver, 0)
	errs := make([]error, 0)

	for _, n := range ns {
		n := n
		wg.Add(1)
		go func() {
			defer wg.Done()

			resolvers, err := r.ListByTypeName(ctx, apiID, n)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			rslvs = append(rslvs, resolvers...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return rslvs, nil
}

func (r *resolverRepositoryForAppSync) listTypeNames(ctx context.Context, apiID string) (res []string, err error) {
	defer wrap(&err)

	ns := make([]string, 0)

	var token *string
	for {
		out, err := r.appsyncClient.ListTypes(
			ctx,
			&appsync.ListTypesInput{
				ApiId:     &apiID,
				Format:    types.TypeDefinitionFormatSdl,
				NextToken: token,
			},
			func(o *appsync.Options) {
				o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
			},
		)
		if err != nil {
			return nil, err
		}

		for _, typ := range out.Types {
			if typ.Name == nil {
				return nil, fmt.Errorf("%w: missing type name in AppSync ListTypes API response", model.ErrNilValue)
			}

			ns = append(ns, *typ.Name)
		}

		token = out.NextToken
		if token == nil {
			break
		}
	}

	return ns, nil
}

func (r *resolverRepositoryForAppSync) ListByTypeName(ctx context.Context, apiID string, typeName string) (res []model.Resolver, err error) {
	defer wrap(&err)

	rslvs := make([]model.Resolver, 0)

	var token *string
	for {
		out, err := r.appsyncClient.ListResolvers(
			ctx,
			&appsync.ListResolversInput{
				ApiId:     &apiID,
				TypeName:  &typeName,
				NextToken: token,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, rslv := range out.Resolvers {
			rslvs = append(rslvs, *mapper.NewResolverMapper().ToModel(ctx, &rslv))
		}

		token = out.NextToken
		if token == nil {
			break
		}
	}

	return rslvs, nil
}

func (r *resolverRepositoryForAppSync) Get(ctx context.Context, apiID string, typeName string, fieldName string) (res *model.Resolver, err error) {
	defer wrap(&err)

	out, err := r.appsyncClient.GetResolver(
		ctx,
		&appsync.GetResolverInput{
			ApiId:     &apiID,
			TypeName:  &typeName,
			FieldName: &fieldName,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		if nfe := new(types.NotFoundException); errors.As(err, &nfe) {
			return nil, model.ErrNotFound
		}

		return nil, err
	}

	rslv := mapper.NewResolverMapper().ToModel(ctx, out.Resolver)
	if rslv == nil {
		return nil, fmt.Errorf("%w: missing resolver in AppSync GetResolver API response", model.ErrNilValue)
	}

	return rslv, nil
}

func (r *resolverRepositoryForAppSync) Save(ctx context.Context, apiID string, resolver *model.Resolver) (res *model.Resolver, err error) {
	defer wrap(&err)

	if resolver == nil {
		return nil, fmt.Errorf("%w: missing arguments in save resolver method", model.ErrNilValue)
	}

	if resolver.TypeName == nil {
		return nil, fmt.Errorf("%w: missing type name", model.ErrNilValue)
	}

	if resolver.FieldName == nil {
		return nil, fmt.Errorf("%w: missing field name", model.ErrNilValue)
	}

	save := r.update
	if _, err := r.Get(ctx, apiID, *resolver.TypeName, *resolver.FieldName); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			save = r.create
		} else {
			return nil, err
		}
	}

	rslv, err := save(ctx, apiID, resolver)
	if err != nil {
		return nil, err
	}

	return rslv, nil
}

func (r *resolverRepositoryForAppSync) create(ctx context.Context, apiID string, resolver *model.Resolver) (res *model.Resolver, err error) {
	defer wrap(&err)

	if resolver == nil {
		return nil, fmt.Errorf("%w: missing arguments in create resolver method", model.ErrNilValue)
	}

	rr := mapper.NewResolverMapper().FromModel(ctx, resolver)
	out, err := r.appsyncClient.CreateResolver(
		ctx,
		&appsync.CreateResolverInput{
			ApiId:                   &apiID,
			TypeName:                rr.TypeName,
			FieldName:               rr.FieldName,
			DataSourceName:          rr.DataSourceName,
			RequestMappingTemplate:  rr.RequestMappingTemplate,
			ResponseMappingTemplate: rr.ResponseMappingTemplate,
			Kind:                    rr.Kind,
			PipelineConfig:          rr.PipelineConfig,
			SyncConfig:              rr.SyncConfig,
			CachingConfig:           rr.CachingConfig,
			MaxBatchSize:            rr.MaxBatchSize,
			Runtime:                 rr.Runtime,
			Code:                    rr.Code,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, err
	}

	rslv := mapper.NewResolverMapper().ToModel(ctx, out.Resolver)
	if rslv == nil {
		return nil, fmt.Errorf("%w: missing resolver in AppSync CreateResolver API response", model.ErrNilValue)
	}

	return rslv, nil
}

func (r *resolverRepositoryForAppSync) update(ctx context.Context, apiID string, resolver *model.Resolver) (res *model.Resolver, err error) {
	defer wrap(&err)

	if resolver == nil {
		return nil, fmt.Errorf("%w: missing arguments in update resolver method", model.ErrNilValue)
	}

	rr := mapper.NewResolverMapper().FromModel(ctx, resolver)
	out, err := r.appsyncClient.UpdateResolver(
		ctx,
		&appsync.UpdateResolverInput{
			ApiId:                   &apiID,
			TypeName:                rr.TypeName,
			FieldName:               rr.FieldName,
			DataSourceName:          rr.DataSourceName,
			RequestMappingTemplate:  rr.RequestMappingTemplate,
			ResponseMappingTemplate: rr.ResponseMappingTemplate,
			Kind:                    rr.Kind,
			PipelineConfig:          rr.PipelineConfig,
			SyncConfig:              rr.SyncConfig,
			CachingConfig:           rr.CachingConfig,
			MaxBatchSize:            rr.MaxBatchSize,
			Runtime:                 rr.Runtime,
			Code:                    rr.Code,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, err
	}

	rslv := mapper.NewResolverMapper().ToModel(ctx, out.Resolver)
	if rslv == nil {
		return nil, fmt.Errorf("%w: missing resolver in AppSync UpdateResolver API response", model.ErrNilValue)
	}

	return rslv, nil
}

func (r *resolverRepositoryForAppSync) Delete(ctx context.Context, apiID string, typeName string, fieldName string) (err error) {
	defer wrap(&err)

	if _, err := r.appsyncClient.DeleteResolver(
		ctx,
		&appsync.DeleteResolverInput{
			ApiId:     &apiID,
			TypeName:  &typeName,
			FieldName: &fieldName,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	); err != nil {
		return err
	}

	return nil
}
