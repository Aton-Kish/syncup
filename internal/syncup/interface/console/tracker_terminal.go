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
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/briandowns/spinner"
	"github.com/mgutz/ansi"
)

type trackerRepositoryForTerminal struct {
	mu      sync.Mutex
	writer  io.Writer
	spinner ispinner
}

func NewTrackerRepositoryForTerminal(w io.Writer) repository.TrackerRepository {
	return &trackerRepositoryForTerminal{
		writer:  w,
		spinner: newSpinner(spinner.CharSets[11], time.Duration(120)*time.Millisecond, spinner.WithColor("fgCyan"), spinner.WithWriter(w)),
	}
}

func (r *trackerRepositoryForTerminal) InProgress(ctx context.Context, msg string) {
	r.spinner.SetSuffix(fmt.Sprintf(" %s", msg))
	r.spinner.Start()
}

func (r *trackerRepositoryForTerminal) Failed(ctx context.Context, msg string) {
	r.done(ctx, model.TrackerStatusFailed, msg)
}

func (r *trackerRepositoryForTerminal) Success(ctx context.Context, msg string) {
	r.done(ctx, model.TrackerStatusSuccess, msg)
}

func (r *trackerRepositoryForTerminal) done(ctx context.Context, status model.TrackerStatus, msg string) {
	var icon, iconStyle, msgStyle string
	switch status {
	case model.TrackerStatusFailed:
		icon = "X"
		iconStyle = "red"
		msgStyle = "red"
	case model.TrackerStatusSuccess:
		icon = "v"
		iconStyle = "green"
		msgStyle = "default+hb"
	default:
		icon = " "
		iconStyle = ""
		msgStyle = ""
	}

	cmsg := fmt.Sprintln(ansi.Color(icon, iconStyle), ansi.Color(msg, msgStyle))

	r.mu.Lock()
	if r.spinner.Active() {
		r.spinner.SetFinalMsg(cmsg)
		r.spinner.Stop()
	} else {
		fmt.Fprint(r.writer, cmsg)
	}
	r.mu.Unlock()
}
