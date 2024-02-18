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
	"fmt"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

type environmentVariablesRepositoryForAppSync struct {
	appsyncClient appsyncClient
}

var (
	_ interface {
		repository.AWSActivator
	} = (*environmentVariablesRepositoryForAppSync)(nil)
)

func NewEnvironmentVariablesRepositoryForAppSync() repository.EnvironmentVariablesRepository {
	return &environmentVariablesRepositoryForAppSync{}
}

func (r *environmentVariablesRepositoryForAppSync) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) (err error) {
	defer wrap(&err)

	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *environmentVariablesRepositoryForAppSync) Get(ctx context.Context, apiID string) (res model.EnvironmentVariables, err error) {
	defer wrap(&err)

	out, err := r.appsyncClient.GetGraphqlApiEnvironmentVariables(
		ctx,
		&appsync.GetGraphqlApiEnvironmentVariablesInput{
			ApiId: &apiID,
		},
	)
	if err != nil {
		return nil, err
	}

	vs := make(model.EnvironmentVariables)
	if out.EnvironmentVariables != nil {
		vs = out.EnvironmentVariables
	}

	return vs, nil
}

func (r *environmentVariablesRepositoryForAppSync) Save(ctx context.Context, apiID string, variables model.EnvironmentVariables) (res model.EnvironmentVariables, err error) {
	defer wrap(&err)

	if variables == nil {
		return nil, fmt.Errorf("%w: missing arguments in save environment variables method", model.ErrNilValue)
	}

	out, err := r.appsyncClient.PutGraphqlApiEnvironmentVariables(
		ctx,
		&appsync.PutGraphqlApiEnvironmentVariablesInput{
			ApiId:                &apiID,
			EnvironmentVariables: variables,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	)
	if err != nil {
		return nil, err
	}

	vs := model.EnvironmentVariables(out.EnvironmentVariables)
	return vs, nil
}
