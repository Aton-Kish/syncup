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

package command

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

const (
	fileNameReadme = "README.md"
)

var (
	templateFuncMap = template.FuncMap{
		"replace": strings.ReplaceAll,
		"now":     time.Now,
	}

	//go:embed template/README.md.gotmpl
	readmeGoTemplate string
	readmeTemplate   = template.Must(
		template.
			New("readme").
			Funcs(templateFuncMap).
			Parse(readmeGoTemplate),
	)

	//go:embed template/reference.md.gotmpl
	referenceGoTemplate string
	referenceTemplate   = template.Must(
		template.
			New("reference").
			Funcs(templateFuncMap).
			Parse(referenceGoTemplate),
	)
)

type xcommand struct {
	*cobra.Command
}

func newCommand(cmd *cobra.Command) *xcommand {
	return &xcommand{
		Command: cmd,
	}
}

func (c *xcommand) GenerateReadme(dir string) error {
	c.InitDefaultHelpCmd()
	c.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	if err := readmeTemplate.Execute(buf, c); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, fileNameReadme), buf.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

func (c *xcommand) GenerateReference(dir string) error {
	c.InitDefaultHelpCmd()
	c.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	if err := referenceTemplate.Execute(buf, c); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	name := fmt.Sprintf("%s.md", strings.ReplaceAll(c.CommandPath(), " ", "_"))
	if err := os.WriteFile(filepath.Join(dir, name), buf.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

func (c *xcommand) GenerateReferences(dir string) error {
	if err := c.GenerateReference(dir); err != nil {
		return err
	}

	for _, sub := range c.Commands() {
		if !sub.IsAvailableCommand() || sub.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := newCommand(sub).GenerateReferences(dir); err != nil {
			return err
		}
	}

	return nil
}
