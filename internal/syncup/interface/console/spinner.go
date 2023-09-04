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

package console

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

type ispinner interface {
	Active() bool
	Color(colors ...string) error
	Disable()
	Enable()
	Enabled() bool
	Lock()
	Restart()
	Reverse()
	Start()
	Stop()
	Unlock()
	UpdateCharSet(cs []string)
	UpdateSpeed(d time.Duration)

	SetSuffix(suffix string)
	SetFinalMsg(msg string)
}

type xspinner struct {
	mu sync.Mutex
	*spinner.Spinner
}

func newSpinner(cs []string, d time.Duration, options ...spinner.Option) ispinner {
	return &xspinner{
		Spinner: spinner.New(cs, d, options...),
	}
}

func (s *xspinner) SetSuffix(suffix string) {
	s.mu.Lock()
	s.Suffix = suffix
	s.mu.Unlock()
}

func (s *xspinner) SetFinalMsg(msg string) {
	s.mu.Lock()
	s.FinalMSG = msg
	s.mu.Unlock()
}
