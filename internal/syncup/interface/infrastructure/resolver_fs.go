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
	dirNameResolvers                           = "resolvers"
	fileNameResolverMetadata                   = "metadata.json"
	fileNameResolverVTLRequestMappingTemplate  = "request.vtl"
	fileNameResolverVTLResponseMappingTemplate = "response.vtl"
	fileNameResolverAppSyncJSCode              = "code.js"
)

type resolverRepositoryForFS struct {
	baseDir string
}

var (
	_ interface {
		repository.BaseDirProvider
	} = (*resolverRepositoryForFS)(nil)
)

func NewResolverRepositoryForFS() repository.ResolverRepository {
	return &resolverRepositoryForFS{}
}

func (r *resolverRepositoryForFS) BaseDir(ctx context.Context) string {
	return r.baseDir
}

func (r *resolverRepositoryForFS) SetBaseDir(ctx context.Context, dir string) {
	r.baseDir = dir
}

func (r *resolverRepositoryForFS) List(ctx context.Context, apiID string) (res []model.Resolver, err error) {
	defer wrap(&err)

	es, err := os.ReadDir(filepath.Join(r.BaseDir(ctx), dirNameResolvers))
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	rslvs := make([]model.Resolver, 0)
	errs := make([]error, 0)

	for _, e := range es {
		if !e.IsDir() {
			continue
		}

		typeName := e.Name()
		wg.Add(1)
		go func() {
			defer wg.Done()

			resolvers, err := r.ListByTypeName(ctx, apiID, typeName)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			rslvs = append(rslvs, resolvers...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return rslvs, nil
}

func (r *resolverRepositoryForFS) ListByTypeName(ctx context.Context, apiID string, typeName string) (res []model.Resolver, err error) {
	defer wrap(&err)

	es, err := os.ReadDir(filepath.Join(r.BaseDir(ctx), dirNameResolvers, typeName))
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	rslvs := make([]model.Resolver, 0)
	errs := make([]error, 0)

	for _, e := range es {
		if !e.IsDir() {
			continue
		}

		fieldName := e.Name()
		wg.Add(1)
		go func() {
			defer wg.Done()

			rslv, err := r.Get(ctx, apiID, typeName, fieldName)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			rslvs = append(rslvs, *rslv)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return rslvs, nil
}

func (r *resolverRepositoryForFS) Get(ctx context.Context, apiID string, typeName string, fieldName string) (res *model.Resolver, err error) {
	defer wrap(&err)

	dir := filepath.Join(r.BaseDir(ctx), dirNameResolvers, typeName, fieldName)
	metadata, err := os.ReadFile(filepath.Join(dir, fileNameResolverMetadata))
	if err != nil {
		return nil, err
	}

	rslv := new(model.Resolver)
	if err := json.Unmarshal(metadata, rslv); err != nil {
		return nil, err
	}

	switch {
	case rslv.Runtime == nil:
		// VTL runtime
		requestMappingTemplate, err := os.ReadFile(filepath.Join(dir, fileNameResolverVTLRequestMappingTemplate))
		if err != nil {
			return nil, err
		}

		rslv.RequestMappingTemplate = ptr.Pointer(string(requestMappingTemplate))

		responseMappingTemplate, err := os.ReadFile(filepath.Join(dir, fileNameResolverVTLResponseMappingTemplate))
		if err != nil {
			return nil, err
		}

		rslv.ResponseMappingTemplate = ptr.Pointer(string(responseMappingTemplate))
	case rslv.Runtime.Name == model.RuntimeNameAppsyncJs:
		// AppSync JS runtime
		code, err := os.ReadFile(filepath.Join(dir, fileNameResolverAppSyncJSCode))
		if err != nil {
			return nil, err
		}

		rslv.Code = ptr.Pointer(string(code))
	default:
		// invalid runtime
		return nil, fmt.Errorf("%w: runtime %s", model.ErrInvalidValue, rslv.Runtime.Name)
	}

	return rslv, nil
}

func (r *resolverRepositoryForFS) Save(ctx context.Context, apiID string, resolver *model.Resolver) (res *model.Resolver, err error) {
	defer wrap(&err)

	if resolver == nil {
		return nil, fmt.Errorf("%w: missing arguments in save resolver method", model.ErrNilValue)
	}

	if resolver.TypeName == nil {
		return nil, fmt.Errorf("%w: missing type name", model.ErrNilValue)
	}

	if resolver.FieldName == nil {
		return nil, fmt.Errorf("%w: missing field name", model.ErrNilValue)
	}

	dir := filepath.Join(r.BaseDir(ctx), dirNameResolvers, *resolver.TypeName, *resolver.FieldName)
	if !xfilepath.Exist(dir) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	metadata, err := json.MarshalIndent(resolver, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(dir, fileNameResolverMetadata), metadata, 0o644); err != nil {
		return nil, err
	}

	switch {
	case resolver.Runtime == nil:
		// VTL runtime
		if err := os.WriteFile(filepath.Join(dir, fileNameResolverVTLRequestMappingTemplate), []byte(ptr.ToValue(resolver.RequestMappingTemplate)), 0o644); err != nil {
			return nil, err
		}

		if err := os.WriteFile(filepath.Join(dir, fileNameResolverVTLResponseMappingTemplate), []byte(ptr.ToValue(resolver.ResponseMappingTemplate)), 0o644); err != nil {
			return nil, err
		}
	case resolver.Runtime.Name == model.RuntimeNameAppsyncJs:
		// AppSync JS runtime
		if err := os.WriteFile(filepath.Join(dir, fileNameResolverAppSyncJSCode), []byte(ptr.ToValue(resolver.Code)), 0o644); err != nil {
			return nil, err
		}
	default:
		// invalid runtime
		return nil, fmt.Errorf("%w: runtime %s", model.ErrInvalidValue, resolver.Runtime.Name)
	}

	return resolver, nil
}

func (r *resolverRepositoryForFS) Delete(ctx context.Context, apiID string, typeName string, fieldName string) (err error) {
	defer wrap(&err)

	if err := os.RemoveAll(filepath.Join(r.BaseDir(ctx), dirNameResolvers, typeName, fieldName)); err != nil {
		return err
	}

	return nil
}
