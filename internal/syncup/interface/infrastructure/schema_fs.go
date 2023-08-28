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
	"os"
	"path/filepath"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

const (
	fileNameSchema = "schema.graphqls"
)

type schemaRepositoryForFS struct {
	baseDir string
}

var (
	_ interface {
		repository.BaseDirProvider
	} = (*schemaRepositoryForFS)(nil)
)

func NewSchemaRepositoryForFS() repository.SchemaRepository {
	return &schemaRepositoryForFS{}
}

func (r *schemaRepositoryForFS) BaseDir(ctx context.Context) string {
	return r.baseDir
}

func (r *schemaRepositoryForFS) SetBaseDir(ctx context.Context, dir string) {
	r.baseDir = dir
}

func (r *schemaRepositoryForFS) Get(ctx context.Context, apiID string) (*model.Schema, error) {
	data, err := os.ReadFile(filepath.Join(r.BaseDir(ctx), fileNameSchema))
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	s := model.Schema(data)
	return &s, nil
}

func (r *schemaRepositoryForFS) Save(ctx context.Context, apiID string, schema *model.Schema) (*model.Schema, error) {
	if schema == nil {
		return nil, &model.LibError{Err: model.ErrNilValue}
	}

	panic("unimplemented")
}
