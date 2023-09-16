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

package mapper

import (
	"context"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

type ResolverMapper interface {
	ToModel(ctx context.Context, v *types.Resolver) *model.Resolver
	FromModel(ctx context.Context, v *model.Resolver) *types.Resolver
}

type resolverMapper struct{}

func NewResolverMapper() ResolverMapper {
	return (*resolverMapper)(nil)
}

func (*resolverMapper) ToModel(ctx context.Context, v *types.Resolver) *model.Resolver {
	if v == nil {
		return nil
	}

	return &model.Resolver{
		TypeName:                v.TypeName,
		FieldName:               v.FieldName,
		DataSourceName:          v.DataSourceName,
		ResolverArn:             v.ResolverArn,
		RequestMappingTemplate:  v.RequestMappingTemplate,
		ResponseMappingTemplate: v.ResponseMappingTemplate,
		Kind:                    model.ResolverKind(v.Kind),
		PipelineConfig:          (*pipelineConfigMapper)(nil).ToModel(ctx, v.PipelineConfig),
		SyncConfig:              (*syncConfigMapper)(nil).ToModel(ctx, v.SyncConfig),
		CachingConfig:           (*cachingConfigMapper)(nil).ToModel(ctx, v.CachingConfig),
		MaxBatchSize:            v.MaxBatchSize,
		Runtime:                 (*runtimeMapper)(nil).ToModel(ctx, v.Runtime),
		Code:                    v.Code,
	}
}

func (*resolverMapper) FromModel(ctx context.Context, v *model.Resolver) *types.Resolver {
	if v == nil {
		return nil
	}

	return &types.Resolver{
		TypeName:                v.TypeName,
		FieldName:               v.FieldName,
		DataSourceName:          v.DataSourceName,
		ResolverArn:             v.ResolverArn,
		RequestMappingTemplate:  v.RequestMappingTemplate,
		ResponseMappingTemplate: v.ResponseMappingTemplate,
		Kind:                    types.ResolverKind(v.Kind),
		PipelineConfig:          (*pipelineConfigMapper)(nil).FromModel(ctx, v.PipelineConfig),
		SyncConfig:              (*syncConfigMapper)(nil).FromModel(ctx, v.SyncConfig),
		CachingConfig:           (*cachingConfigMapper)(nil).FromModel(ctx, v.CachingConfig),
		MaxBatchSize:            v.MaxBatchSize,
		Runtime:                 (*runtimeMapper)(nil).FromModel(ctx, v.Runtime),
		Code:                    v.Code,
	}
}
