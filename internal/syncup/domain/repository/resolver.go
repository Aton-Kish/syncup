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

package repository

import (
	"context"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
)

type ResolverRepository interface {
	List(ctx context.Context, apiID string) ([]model.Resolver, error)
	ListByTypeName(ctx context.Context, apiID string, typeName string) ([]model.Resolver, error)
	Get(ctx context.Context, apiID string, typeName string, fieldName string) (*model.Resolver, error)
	Save(ctx context.Context, apiID string, resolver *model.Resolver) (*model.Resolver, error)
	Delete(ctx context.Context, apiID string, typeName string, fieldName string) error
}
