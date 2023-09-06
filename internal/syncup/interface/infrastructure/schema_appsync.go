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
	"time"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

const (
	defaultDuration = time.Duration(1) * time.Second
)

type schemaRepositoryForAppSync struct {
	appsyncClient appsyncClient

	duration time.Duration
}

var (
	_ interface {
		repository.AWSActivator
	} = (*schemaRepositoryForAppSync)(nil)
)

func NewSchemaRepositoryForAppSync() repository.SchemaRepository {
	return &schemaRepositoryForAppSync{
		duration: defaultDuration,
	}
}

func (r *schemaRepositoryForAppSync) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) error {
	c, err := activatedAWSClients(ctx, optFns...)
	if err != nil {
		return err
	}

	r.appsyncClient = c.appsyncClient

	return nil
}

func (r *schemaRepositoryForAppSync) Get(ctx context.Context, apiID string) (*model.Schema, error) {
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

func (r *schemaRepositoryForAppSync) Save(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
	if schema == nil {
		return nil, &model.LibError{Err: fmt.Errorf("%w: missing arguments in save schema method", model.ErrNilValue)}
	}

	if err := r.startCreation(ctx, apiID, []byte(*schema)); err != nil {
		return nil, err
	}

	ticker := time.NewTicker(r.duration)
	defer ticker.Stop()

	ch := make(chan error)
	go func() {
		defer close(ch)

		for {
			select {
			case <-ticker.C:
				isCreated, err := r.isCreated(ctx, apiID)
				if err != nil {
					ch <- err
					return
				}

				if isCreated {
					ch <- nil
					return
				}
			case <-ctx.Done():
				ch <- &model.LibError{Err: ctx.Err()}
				return
			}
		}
	}()

	if err := <-ch; err != nil {
		return nil, err
	}

	s, err := r.Get(ctx, apiID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *schemaRepositoryForAppSync) startCreation(ctx context.Context, apiID string, definition []byte) error {
	if _, err := r.appsyncClient.StartSchemaCreation(
		ctx,
		&appsync.StartSchemaCreationInput{
			ApiId:      &apiID,
			Definition: definition,
		},
		func(o *appsync.Options) {
			o.Retryer = retry.AddWithErrorCodes(o.Retryer, (*types.ConcurrentModificationException)(nil).ErrorCode())
		},
	); err != nil {
		return &model.LibError{Err: err}
	}

	return nil
}

func (r *schemaRepositoryForAppSync) isCreated(ctx context.Context, apiID string) (bool, error) {
	out, err := r.appsyncClient.GetSchemaCreationStatus(
		ctx,
		&appsync.GetSchemaCreationStatusInput{
			ApiId: &apiID,
		},
	)
	if err != nil {
		return false, &model.LibError{Err: err}
	}

	switch out.Status {
	case types.SchemaStatusActive, types.SchemaStatusSuccess:
		return true, nil
	case types.SchemaStatusProcessing:
		return false, nil
	case types.SchemaStatusFailed:
		return false, &model.LibError{Err: fmt.Errorf("%w: schema status %s", model.ErrCreateFailed, out.Status)}
	default:
		return false, &model.LibError{Err: fmt.Errorf("%w: schema status %s", model.ErrInvalidValue, out.Status)}
	}
}
