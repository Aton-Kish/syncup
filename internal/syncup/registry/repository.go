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

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/interface/console"
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

	mfaTokenProviderRepository repository.MFATokenProviderRepository
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

	mfaTokenProviderRepository := console.NewMFATokenProviderRepository()

	return &repo{
		version: version,

		mfaTokenProviderRepository: mfaTokenProviderRepository,
	}
}

func (r *repo) repositories() []any {
	return []any{
		r.MFATokenProviderRepository(),
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

func (r *repo) Version() *model.Version {
	return r.version
}

func (r *repo) MFATokenProviderRepository() repository.MFATokenProviderRepository {
	return r.mfaTokenProviderRepository
}
