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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

const (
	dirNameFunctions                           = "functions"
	fileNameFunctionMetadata                   = "metadata.json"
	fileNameFunctionVTLRequestMappingTemplate  = "request.vtl"
	fileNameFunctionVTLResponseMappingTemplate = "response.vtl"
	fileNameFunctionAppSyncJSCode              = "code.js"
)

type functionRepositoryForFS struct {
	baseDir string
}

var (
	_ interface {
		repository.BaseDirProvider
	} = (*functionRepositoryForFS)(nil)
)

func NewFunctionRepositoryForFS() repository.FunctionRepository {
	return &functionRepositoryForFS{}
}

func (r *functionRepositoryForFS) BaseDir(ctx context.Context) string {
	return r.baseDir
}

func (r *functionRepositoryForFS) SetBaseDir(ctx context.Context, dir string) {
	r.baseDir = dir
}

func (*functionRepositoryForFS) List(ctx context.Context, apiID string) ([]model.Function, error) {
	panic("unimplemented")
}

func (r *functionRepositoryForFS) Get(ctx context.Context, apiID string, functionID string) (*model.Function, error) {
	dir := filepath.Join(r.BaseDir(ctx), dirNameFunctions, functionID)
	metadata, err := os.ReadFile(filepath.Join(dir, fileNameFunctionMetadata))
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	fn := new(model.Function)
	if err := json.Unmarshal(metadata, fn); err != nil {
		return nil, &model.LibError{Err: err}
	}

	switch {
	case fn.Runtime == nil:
		// VTL runtime
		requestMappingTemplate, err := os.ReadFile(filepath.Join(dir, fileNameFunctionVTLRequestMappingTemplate))
		if err != nil {
			return nil, &model.LibError{Err: err}
		}

		fn.RequestMappingTemplate = ptr.Pointer(string(requestMappingTemplate))

		responseMappingTemplate, err := os.ReadFile(filepath.Join(dir, fileNameFunctionVTLResponseMappingTemplate))
		if err != nil {
			return nil, &model.LibError{Err: err}
		}

		fn.ResponseMappingTemplate = ptr.Pointer(string(responseMappingTemplate))
	case fn.Runtime.Name == model.RuntimeNameAppsyncJs:
		// AppSync JS runtime
		code, err := os.ReadFile(filepath.Join(dir, fileNameFunctionAppSyncJSCode))
		if err != nil {
			return nil, &model.LibError{Err: err}
		}

		fn.Code = ptr.Pointer(string(code))
	default:
		// invalid runtime
		return nil, &model.LibError{Err: fmt.Errorf("%w: runtime %s", model.ErrInvalidValue, fn.Runtime.Name)}
	}

	return fn, nil
}

func (*functionRepositoryForFS) Save(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
	panic("unimplemented")
}

func (*functionRepositoryForFS) Delete(ctx context.Context, apiID string, functionID string) error {
	panic("unimplemented")
}
