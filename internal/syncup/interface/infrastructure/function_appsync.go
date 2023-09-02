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

func (r *functionRepositoryForAppSync) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) error {
	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *functionRepositoryForAppSync) List(ctx context.Context, apiID string) ([]model.Function, error) {
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
			return nil, &model.LibError{Err: err}
		}

		for _, fn := range out.Functions {
			fns = append(fns, *mapper.NewFunctionMapper().ToModel(ctx, &fn))
		}

		token = out.NextToken
		if token == nil {
			break
		}
	}

	return fns, nil
}

func (r *functionRepositoryForAppSync) Get(ctx context.Context, apiID string, functionID string) (*model.Function, error) {
	out, err := r.appsyncClient.GetFunction(
		ctx,
		&appsync.GetFunctionInput{
			ApiId:      &apiID,
			FunctionId: &functionID,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	fn := mapper.NewFunctionMapper().ToModel(ctx, out.FunctionConfiguration)
	if fn == nil {
		return nil, &model.LibError{Err: model.ErrNilValue}
	}

	return fn, nil
}

func (*functionRepositoryForAppSync) Save(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
	panic("unimplemented")
}

func (*functionRepositoryForAppSync) Delete(ctx context.Context, apiID string, functionID string) error {
	panic("unimplemented")
}
