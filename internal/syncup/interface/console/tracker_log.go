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

package console

import (
	"context"
	"io"
	"log/slog"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
)

type trackerRepositoryForLog struct {
	logger *slog.Logger
}

func NewTrackerRepositoryForLog(w io.Writer) repository.TrackerRepository {
	return &trackerRepositoryForLog{
		logger: slog.New(slog.NewJSONHandler(w, nil)).With(slog.String("app", "syncup"), slog.String("category", "tracker")),
	}
}

func (r *trackerRepositoryForLog) InProgress(ctx context.Context, msg string) {
	r.logContext(ctx, model.TrackerStatusInProgress, msg)
}

func (r *trackerRepositoryForLog) Failed(ctx context.Context, msg string) {
	r.logContext(ctx, model.TrackerStatusFailed, msg)
}

func (r *trackerRepositoryForLog) Success(ctx context.Context, msg string) {
	r.logContext(ctx, model.TrackerStatusSuccess, msg)
}

func (r *trackerRepositoryForLog) logContext(ctx context.Context, status model.TrackerStatus, msg string) {
	l := r.logger.With(slog.String("status", string(status)))
	switch status {
	case model.TrackerStatusInProgress, model.TrackerStatusSuccess:
		l.InfoContext(ctx, msg)
	case model.TrackerStatusFailed:
		l.ErrorContext(ctx, msg)
	default:
		l.DebugContext(ctx, msg)
	}
}
