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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	ptr "github.com/Aton-Kish/goptr"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/xfilepath"
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

func (r *functionRepositoryForFS) List(ctx context.Context, apiID string) ([]model.Function, error) {
	es, err := os.ReadDir(filepath.Join(r.BaseDir(ctx), dirNameFunctions))
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	fns := make([]model.Function, 0)
	errs := make([]error, 0)

	for _, e := range es {
		if !e.IsDir() {
			continue
		}

		functionID := e.Name()
		wg.Add(1)
		go func() {
			defer wg.Done()

			fn, err := r.Get(ctx, apiID, functionID)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			fns = append(fns, *fn)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return fns, nil
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

func (r *functionRepositoryForFS) Save(ctx context.Context, apiID string, function *model.Function) (*model.Function, error) {
	if function == nil || function.FunctionId == nil {
		return nil, &model.LibError{Err: model.ErrNilValue}
	}

	dir := filepath.Join(r.BaseDir(ctx), dirNameFunctions, *function.FunctionId)
	if !xfilepath.Exist(dir) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, &model.LibError{Err: err}
		}
	}

	metadata, err := json.MarshalIndent(function, "", "  ")
	if err != nil {
		return nil, &model.LibError{Err: err}
	}

	if err := os.WriteFile(filepath.Join(dir, fileNameFunctionMetadata), metadata, 0o644); err != nil {
		return nil, &model.LibError{Err: err}
	}

	switch {
	case function.Runtime == nil:
		// VTL runtime
		if err := os.WriteFile(filepath.Join(dir, fileNameFunctionVTLRequestMappingTemplate), []byte(ptr.ToValue(function.RequestMappingTemplate)), 0o644); err != nil {
			return nil, &model.LibError{Err: err}
		}

		if err := os.WriteFile(filepath.Join(dir, fileNameFunctionVTLResponseMappingTemplate), []byte(ptr.ToValue(function.ResponseMappingTemplate)), 0o644); err != nil {
			return nil, &model.LibError{Err: err}
		}
	case function.Runtime.Name == model.RuntimeNameAppsyncJs:
		// AppSync JS runtime
		if err := os.WriteFile(filepath.Join(dir, fileNameFunctionAppSyncJSCode), []byte(ptr.ToValue(function.Code)), 0o644); err != nil {
			return nil, &model.LibError{Err: err}
		}
	default:
		// invalid runtime
		return nil, &model.LibError{Err: fmt.Errorf("%w: runtime %s", model.ErrInvalidValue, function.Runtime.Name)}
	}

	return function, nil
}

func (r *functionRepositoryForFS) Delete(ctx context.Context, apiID string, functionID string) error {
	if err := os.RemoveAll(filepath.Join(r.BaseDir(ctx), dirNameFunctions, functionID)); err != nil {
		return &model.LibError{Err: err}
	}

	return nil
}
