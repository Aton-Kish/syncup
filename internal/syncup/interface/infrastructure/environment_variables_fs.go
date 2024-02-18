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

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/xfilepath"
)

const (
	fileNameEnvironmentVariables = "env.json"
)

type environmentVariablesRepositoryForFS struct {
	baseDir string
}

var (
	_ interface {
		repository.BaseDirProvider
	} = (*environmentVariablesRepositoryForFS)(nil)
)

func NewEnvironmentVariablesRepositoryForFS() repository.EnvironmentVariablesRepository {
	return &environmentVariablesRepositoryForFS{}
}

func (r *environmentVariablesRepositoryForFS) BaseDir(ctx context.Context) string {
	return r.baseDir
}

func (r *environmentVariablesRepositoryForFS) SetBaseDir(ctx context.Context, dir string) {
	r.baseDir = dir
}

func (r *environmentVariablesRepositoryForFS) Get(ctx context.Context, apiID string) (res model.EnvironmentVariables, err error) {
	defer wrap(&err)

	data, err := os.ReadFile(filepath.Join(r.BaseDir(ctx), fileNameEnvironmentVariables))
	if err != nil {
		return nil, err
	}

	vs := make(model.EnvironmentVariables)
	if err := json.Unmarshal(data, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (r *environmentVariablesRepositoryForFS) Save(ctx context.Context, apiID string, variables model.EnvironmentVariables) (res model.EnvironmentVariables, err error) {
	defer wrap(&err)

	if variables == nil {
		return nil, fmt.Errorf("%w: missing arguments in save environment variables method", model.ErrNilValue)
	}

	dir := r.BaseDir(ctx)
	if !xfilepath.Exist(dir) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	data, err := json.MarshalIndent(variables, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(dir, fileNameEnvironmentVariables), data, 0o644); err != nil {
		return nil, err
	}

	return variables, nil
}
