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
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

type schemaAppSyncRepository struct {
	appsyncClient appsyncClient
}

var (
	_ interface {
		repository.AWSActivator
	} = (*schemaAppSyncRepository)(nil)
)

func NewSchemaAppSyncRepository() repository.SchemaRepository {
	return &schemaAppSyncRepository{}
}

func (r *schemaAppSyncRepository) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) error {
	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *schemaAppSyncRepository) Get(ctx context.Context, apiID string) (*model.Schema, error) {
	out, err := r.appsyncClient.GetIntrospectionSchema(
		ctx,
		&appsync.GetIntrospectionSchemaInput{
			ApiId:  &apiID,
			Format: types.OutputTypeSdl,
		},
	)
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	s := model.Schema(out.Schema)
	return &s, nil
}

func (r *schemaAppSyncRepository) Save(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
	panic("unimplemented")
}
