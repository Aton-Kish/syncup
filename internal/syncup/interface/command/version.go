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
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/spf13/cobra"
)

type VersionCommand interface {
	Command
}

type versionCommand struct {
	options *options

	version *model.Version

	cmd  *cobra.Command
	once sync.Once
}

type versionOutput struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	GoVersion string `json:"go"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	BuildTime string `json:"built"`
}

func NewVersionCommand(repo repository.Repository, optFns ...func(o *options)) VersionCommand {
	return &versionCommand{
		options: newOptions(optFns...),

		version: repo.Version(),
	}
}

func (c *versionCommand) Execute(ctx context.Context, args ...string) error {
	cmd := c.command()
	cmd.SetArgs(args)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return &commandError{Err: err}
	}

	return nil
}

func (c *versionCommand) RegisterSubCommands(cmds ...Command) {
	subs := make([]*cobra.Command, 0, len(cmds))
	for _, cmd := range cmds {
		subs = append(subs, cmd.command())
	}

	cmd := c.command()
	cmd.AddCommand(subs...)
}

func (c *versionCommand) command() *cobra.Command {
	c.once.Do(func() {
		c.cmd = &cobra.Command{
			Use:   "version",
			Short: "Show the syncup version information",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := &versionOutput{
					Version:   c.version.Version,
					GitCommit: c.version.GitCommit,
					GoVersion: c.version.GoVersion,
					OS:        c.version.OS,
					Arch:      c.version.Arch,
					BuildTime: c.version.BuildTime,
				}

				data, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(data)); err != nil {
					return err
				}

				return nil
			},
			SilenceUsage: true,
		}

		c.cmd.SetIn(c.options.stdio.in)
		c.cmd.SetOut(c.options.stdio.out)
		c.cmd.SetErr(c.options.stdio.err)
	})

	return c.cmd
}
