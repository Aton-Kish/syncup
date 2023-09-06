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

type syncConfigMapper struct{}

var (
	_ interface {
		ToModel(ctx context.Context, v *types.SyncConfig) *model.SyncConfig
		FromModel(ctx context.Context, v *model.SyncConfig) *types.SyncConfig
	} = (*syncConfigMapper)(nil)
)

func (*syncConfigMapper) ToModel(ctx context.Context, v *types.SyncConfig) *model.SyncConfig {
	if v == nil {
		return nil
	}

	return &model.SyncConfig{
		ConflictHandler:             model.ConflictHandlerType(v.ConflictHandler),
		ConflictDetection:           model.ConflictDetectionType(v.ConflictDetection),
		LambdaConflictHandlerConfig: (*lambdaConflictHandlerConfigMapper)(nil).ToModel(ctx, v.LambdaConflictHandlerConfig),
	}
}

func (*syncConfigMapper) FromModel(ctx context.Context, v *model.SyncConfig) *types.SyncConfig {
	if v == nil {
		return nil
	}

	return &types.SyncConfig{
		ConflictHandler:             types.ConflictHandlerType(v.ConflictHandler),
		ConflictDetection:           types.ConflictDetectionType(v.ConflictDetection),
		LambdaConflictHandlerConfig: (*lambdaConflictHandlerConfigMapper)(nil).FromModel(ctx, v.LambdaConflictHandlerConfig),
	}
}

type lambdaConflictHandlerConfigMapper struct{}

var (
	_ interface {
		ToModel(ctx context.Context, v *types.LambdaConflictHandlerConfig) *model.LambdaConflictHandlerConfig
		FromModel(ctx context.Context, v *model.LambdaConflictHandlerConfig) *types.LambdaConflictHandlerConfig
	} = (*lambdaConflictHandlerConfigMapper)(nil)
)

func (*lambdaConflictHandlerConfigMapper) ToModel(ctx context.Context, v *types.LambdaConflictHandlerConfig) *model.LambdaConflictHandlerConfig {
	if v == nil {
		return nil
	}

	return &model.LambdaConflictHandlerConfig{
		LambdaConflictHandlerArn: v.LambdaConflictHandlerArn,
	}
}

func (*lambdaConflictHandlerConfigMapper) FromModel(ctx context.Context, v *model.LambdaConflictHandlerConfig) *types.LambdaConflictHandlerConfig {
	if v == nil {
		return nil
	}

	return &types.LambdaConflictHandlerConfig{
		LambdaConflictHandlerArn: v.LambdaConflictHandlerArn,
	}
}
