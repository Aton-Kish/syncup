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

package registry

import (
	"context"
	"os"
	"strconv"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/console"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/infrastructure"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
	goVersion = "unknown"
	goOS      = "unknown"
	goArch    = "unknown"
	buildTime = "unknown"
)

type repo struct {
	version *model.Version

	trackerRepository repository.TrackerRepository

	mfaTokenProviderRepository repository.MFATokenProviderRepository

	schemaRepositoryForAppSync repository.SchemaRepository
	schemaRepositoryForFS      repository.SchemaRepository

	functionRepositoryForAppSync repository.FunctionRepository
	functionRepositoryForFS      repository.FunctionRepository
}

func NewRepository() repository.Repository {
	version := &model.Version{
		Version:   version,
		GitCommit: gitCommit,
		GoVersion: goVersion,
		OS:        goOS,
		Arch:      goArch,
		BuildTime: buildTime,
	}

	var trackerRepository repository.TrackerRepository
	if isCI, _ := strconv.ParseBool(os.Getenv("CI")); isCI {
		trackerRepository = console.NewTrackerRepositoryForLog(os.Stderr)
	} else {
		trackerRepository = console.NewTrackerRepositoryForTerminal(os.Stderr)
	}

	mfaTokenProviderRepository := console.NewMFATokenProviderRepository()

	schemaRepositoryForAppSync := infrastructure.NewSchemaRepositoryForAppSync()
	schemaRepositoryForFS := infrastructure.NewSchemaRepositoryForFS()

	functionRepositoryForAppSync := infrastructure.NewFunctionRepositoryForAppSync()
	functionRepositoryForFS := infrastructure.NewFunctionRepositoryForFS()

	return &repo{
		version: version,

		trackerRepository: trackerRepository,

		mfaTokenProviderRepository: mfaTokenProviderRepository,

		schemaRepositoryForAppSync: schemaRepositoryForAppSync,
		schemaRepositoryForFS:      schemaRepositoryForFS,

		functionRepositoryForAppSync: functionRepositoryForAppSync,
		functionRepositoryForFS:      functionRepositoryForFS,
	}
}

func (r *repo) repositories() []any {
	return []any{
		r.TrackerRepository(),

		r.MFATokenProviderRepository(),

		r.SchemaRepositoryForAppSync(),
		r.SchemaRepositoryForFS(),

		r.FunctionRepositoryForAppSync(),
		r.FunctionRepositoryForFS(),
	}
}

func (r *repo) ActivateAWS(ctx context.Context, optFns ...func(o *model.AWSOptions)) error {
	for _, rr := range r.repositories() {
		if activator, ok := rr.(repository.AWSActivator); ok {
			if err := activator.ActivateAWS(ctx, optFns...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *repo) BaseDir(ctx context.Context) string {
	for _, rr := range r.repositories() {
		if baseDirProvider, ok := rr.(repository.BaseDirProvider); ok {
			return baseDirProvider.BaseDir(ctx)
		}
	}

	return ""
}

func (r *repo) SetBaseDir(ctx context.Context, dir string) {
	for _, rr := range r.repositories() {
		if baseDirProvider, ok := rr.(repository.BaseDirProvider); ok {
			baseDirProvider.SetBaseDir(ctx, dir)
		}
	}
}

func (r *repo) Version() *model.Version {
	return r.version
}

func (r *repo) TrackerRepository() repository.TrackerRepository {
	return r.trackerRepository
}

func (r *repo) MFATokenProviderRepository() repository.MFATokenProviderRepository {
	return r.mfaTokenProviderRepository
}

func (r *repo) SchemaRepositoryForAppSync() repository.SchemaRepository {
	return r.schemaRepositoryForAppSync
}

func (r *repo) SchemaRepositoryForFS() repository.SchemaRepository {
	return r.schemaRepositoryForFS
}

func (r *repo) FunctionRepositoryForAppSync() repository.FunctionRepository {
	return r.functionRepositoryForAppSync
}

func (r *repo) FunctionRepositoryForFS() repository.FunctionRepository {
	return r.functionRepositoryForFS
}
