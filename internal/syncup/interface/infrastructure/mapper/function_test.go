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
	"testing"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/stretchr/testify/assert"
)

func Test_functionMapper_ToModel(t *testing.T) {
	type args struct {
		v *types.FunctionConfiguration
	}

	type expected struct {
		res *model.Function
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: nil",
			args: args{
				v: nil,
			},
			expected: expected{
				res: nil,
			},
		},
		{
			name: "happy path: not nil",
			args: args{
				v: &types.FunctionConfiguration{
					FunctionId:              aws.String("FunctionId"),
					FunctionArn:             aws.String("FunctionArn"),
					Name:                    aws.String("Name"),
					Description:             aws.String("Description"),
					DataSourceName:          aws.String("DataSourceName"),
					RequestMappingTemplate:  aws.String("RequestMappingTemplate"),
					ResponseMappingTemplate: aws.String("ResponseMappingTemplate"),
					FunctionVersion:         aws.String("FunctionVersion"),
					SyncConfig: &types.SyncConfig{
						ConflictDetection: types.ConflictDetectionType("ConflictDetection"),
						ConflictHandler:   types.ConflictHandlerType("ConflictHandler"),
						LambdaConflictHandlerConfig: &types.LambdaConflictHandlerConfig{
							LambdaConflictHandlerArn: aws.String("LambdaConflictHandlerArn"),
						},
					},
					MaxBatchSize: 0,
					Runtime: &types.AppSyncRuntime{
						Name:           types.RuntimeName("Name"),
						RuntimeVersion: aws.String("RuntimeVersion"),
					},
					Code: aws.String("Code"),
				},
			},
			expected: expected{
				res: &model.Function{
					FunctionId:              ptr.Pointer("FunctionId"),
					FunctionArn:             ptr.Pointer("FunctionArn"),
					Name:                    ptr.Pointer("Name"),
					Description:             ptr.Pointer("Description"),
					DataSourceName:          ptr.Pointer("DataSourceName"),
					RequestMappingTemplate:  ptr.Pointer("RequestMappingTemplate"),
					ResponseMappingTemplate: ptr.Pointer("ResponseMappingTemplate"),
					FunctionVersion:         ptr.Pointer("FunctionVersion"),
					SyncConfig: &model.SyncConfig{
						ConflictHandler:   model.ConflictHandlerType("ConflictHandler"),
						ConflictDetection: model.ConflictDetectionType("ConflictDetection"),
						LambdaConflictHandlerConfig: &model.LambdaConflictHandlerConfig{
							LambdaConflictHandlerArn: ptr.Pointer("LambdaConflictHandlerArn"),
						},
					},
					MaxBatchSize: 0,
					Runtime: &model.Runtime{
						Name:           model.RuntimeName("Name"),
						RuntimeVersion: ptr.Pointer("RuntimeVersion"),
					},
					Code: ptr.Pointer("Code"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			// Act
			actual := (*functionMapper)(nil).ToModel(ctx, tt.args.v)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func Test_functionMapper_FromModel(t *testing.T) {
	type args struct {
		v *model.Function
	}

	type expected struct {
		res *types.FunctionConfiguration
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "happy path: nil",
			args: args{
				v: nil,
			},
			expected: expected{
				res: nil,
			},
		},
		{
			name: "happy path: not nil",
			args: args{
				v: &model.Function{
					FunctionId:              ptr.Pointer("FunctionId"),
					FunctionArn:             ptr.Pointer("FunctionArn"),
					Name:                    ptr.Pointer("Name"),
					Description:             ptr.Pointer("Description"),
					DataSourceName:          ptr.Pointer("DataSourceName"),
					RequestMappingTemplate:  ptr.Pointer("RequestMappingTemplate"),
					ResponseMappingTemplate: ptr.Pointer("ResponseMappingTemplate"),
					FunctionVersion:         ptr.Pointer("FunctionVersion"),
					SyncConfig: &model.SyncConfig{
						ConflictHandler:   model.ConflictHandlerType("ConflictHandler"),
						ConflictDetection: model.ConflictDetectionType("ConflictDetection"),
						LambdaConflictHandlerConfig: &model.LambdaConflictHandlerConfig{
							LambdaConflictHandlerArn: ptr.Pointer("LambdaConflictHandlerArn"),
						},
					},
					MaxBatchSize: 0,
					Runtime: &model.Runtime{
						Name:           model.RuntimeName("Name"),
						RuntimeVersion: ptr.Pointer("RuntimeVersion"),
					},
					Code: ptr.Pointer("Code"),
				},
			},
			expected: expected{
				res: &types.FunctionConfiguration{
					FunctionId:              aws.String("FunctionId"),
					FunctionArn:             aws.String("FunctionArn"),
					Name:                    aws.String("Name"),
					Description:             aws.String("Description"),
					DataSourceName:          aws.String("DataSourceName"),
					RequestMappingTemplate:  aws.String("RequestMappingTemplate"),
					ResponseMappingTemplate: aws.String("ResponseMappingTemplate"),
					FunctionVersion:         aws.String("FunctionVersion"),
					SyncConfig: &types.SyncConfig{
						ConflictDetection: types.ConflictDetectionType("ConflictDetection"),
						ConflictHandler:   types.ConflictHandlerType("ConflictHandler"),
						LambdaConflictHandlerConfig: &types.LambdaConflictHandlerConfig{
							LambdaConflictHandlerArn: aws.String("LambdaConflictHandlerArn"),
						},
					},
					MaxBatchSize: 0,
					Runtime: &types.AppSyncRuntime{
						Name:           types.RuntimeName("Name"),
						RuntimeVersion: aws.String("RuntimeVersion"),
					},
					Code: aws.String("Code"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			// Act
			actual := (*functionMapper)(nil).FromModel(ctx, tt.args.v)

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}
