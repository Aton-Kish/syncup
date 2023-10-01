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
	"fmt"
	"sync"

	"github.com/Aton-Kish/syncup/internal/syncup"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type RootCommand interface {
	Command
}

type rootCommand struct {
	options *options

	version *model.Version

	cmd  *xcommand
	once sync.Once
}

func NewRootCommand(repo repository.Repository, optFns ...func(o *options)) RootCommand {
	return &rootCommand{
		options: newOptions(optFns...),

		version: repo.Version(),
	}
}

func (c *rootCommand) Execute(ctx context.Context, args ...string) (err error) {
	defer wrap(&err)

	cmd := c.command()
	cmd.SetArgs(args)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}

func (c *rootCommand) RegisterSubCommands(cmds ...Command) {
	subs := make([]*cobra.Command, 0, len(cmds))
	for _, cmd := range cmds {
		subs = append(subs, cmd.command().Command)
	}

	cmd := c.command()
	cmd.AddCommand(subs...)
}

func (c *rootCommand) GenerateReferences(ctx context.Context, dir string) (err error) {
	defer wrap(&err)

	cmd := c.command()
	cmd.InitDefaultVersionFlag()
	cmd.InitDefaultCompletionCmd()

	return cmd.GenerateReferences(dir)
}

func (c *rootCommand) command() *xcommand {
	c.once.Do(func() {
		c.cmd = newCommand(&cobra.Command{
			Use:     "syncup",
			Short:   "Sync up with AWS AppSync",
			Version: fmt.Sprintf("%s, build %s (%s/%s)", c.version.Version, c.version.GitCommit, c.version.OS, c.version.Arch),
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				ctx = syncup.WithRequestID(ctx, uuid.NewString())
				cmd.SetContext(ctx)
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				defer wrap(&err)

				if err := cmd.Help(); err != nil {
					return err
				}

				return nil
			},
			SilenceUsage: true,
		})

		c.cmd.SetIn(c.options.stdio.in)
		c.cmd.SetOutput(c.options.stdio.err)
	})

	return c.cmd
}
