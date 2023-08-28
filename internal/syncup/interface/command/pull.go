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
	"sync"

	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/repository"
	"github.com/Aton-Kish/syncup/internal/syncup/usecase"
	"github.com/spf13/cobra"
)

type pullFlags struct {
	region  string
	profile string

	apiID   string
	baseDir string
}

type PullCommand interface {
	Command
}

type pullCommand struct {
	options *options

	useCase                    usecase.PullUseCase
	awsActivator               repository.AWSActivator
	baseDirProvider            repository.BaseDirProvider
	mfaTokenProviderRepository repository.MFATokenProviderRepository

	cmd   *cobra.Command
	flags *pullFlags
	once  sync.Once
}

func NewPullCommand(repo repository.Repository, optFns ...func(o *options)) PullCommand {
	return &pullCommand{
		options: newOptions(optFns...),

		useCase:                    usecase.NewPullUseCase(repo),
		awsActivator:               repo,
		baseDirProvider:            repo,
		mfaTokenProviderRepository: repo.MFATokenProviderRepository(),
	}
}

func (c *pullCommand) Execute(ctx context.Context, args ...string) error {
	cmd := c.command()
	cmd.SetArgs(args)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return &commandError{Err: err}
	}

	return nil
}

func (c *pullCommand) RegisterSubCommands(cmds ...Command) {
	subs := make([]*cobra.Command, 0, len(cmds))
	for _, cmd := range cmds {
		subs = append(subs, cmd.command())
	}

	cmd := c.command()
	cmd.AddCommand(subs...)
}

func (c *pullCommand) command() *cobra.Command {
	c.once.Do(func() {
		c.flags = new(pullFlags)

		c.cmd = &cobra.Command{
			Use:   "pull",
			Short: "Pull resources from AWS AppSync",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

				if err := c.awsActivator.ActivateAWS(
					ctx,
					model.AWSOptionsWithRegion(c.flags.region),
					model.AWSOptionsWithProfile(c.flags.profile),
					model.AWSOptionsWithMFATokenProvider(c.mfaTokenProviderRepository.Get(ctx)),
				); err != nil {
					return err
				}

				c.baseDirProvider.SetBaseDir(ctx, c.flags.baseDir)

				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

				if _, err := c.useCase.Execute(
					ctx,
					&usecase.PullInput{
						APIID: c.flags.apiID,
					},
				); err != nil {
					return err
				}

				return nil
			},
			SilenceUsage: true,
		}

		c.cmd.Flags().StringVar(&c.flags.profile, "profile", "", "Use a specific profile from your AWS credential file.")
		c.cmd.Flags().StringVar(&c.flags.region, "region", "", "The AWS region to use. Overrides config/env settings.")

		c.cmd.Flags().StringVar(&c.flags.apiID, "api-id", "", "The API ID of AWS AppSync.")
		_ = c.cmd.MarkFlagRequired("api-id")
		c.cmd.Flags().StringVar(&c.flags.baseDir, "dir", "", "The directory in which the resources will be saved (instead of current directory).")

		c.cmd.SetIn(c.options.stdio.in)
		c.cmd.SetOut(c.options.stdio.out)
		c.cmd.SetErr(c.options.stdio.err)
	})

	return c.cmd
}
