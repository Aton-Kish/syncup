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

package usecase

import (
	"context"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

type PullInput struct {
	APIID string
}

type PullOutput struct {
}

type PullUseCase interface {
	Execute(ctx context.Context, params *PullInput) (*PullOutput, error)
}

type pullUseCase struct {
	trackerRepository          repository.TrackerRepository
	schemaRepositoryForAppSync repository.SchemaRepository
	schemaRepositoryForFS      repository.SchemaRepository
}

func NewPullUseCase(repo repository.Repository) PullUseCase {
	return &pullUseCase{
		trackerRepository:          repo.TrackerRepository(),
		schemaRepositoryForAppSync: repo.SchemaRepositoryForAppSync(),
		schemaRepositoryForFS:      repo.SchemaRepositoryForFS(),
	}
}

func (uc *pullUseCase) Execute(ctx context.Context, params *PullInput) (*PullOutput, error) {
	apiID := params.APIID

	uc.trackerRepository.Doing(ctx, model.TrackerStatusInfo, "fetching schema")

	schema, err := uc.schemaRepositoryForAppSync.Get(ctx, apiID)
	if err != nil {
		uc.trackerRepository.Done(ctx, model.TrackerStatusDanger, "failed to fetch schema")
		return nil, err
	}

	if _, err := uc.schemaRepositoryForFS.Save(ctx, apiID, schema); err != nil {
		uc.trackerRepository.Done(ctx, model.TrackerStatusDanger, "failed to save schema")
		return nil, err
	}

	uc.trackerRepository.Done(ctx, model.TrackerStatusSuccess, "saved schema")

	return &PullOutput{}, nil
}
