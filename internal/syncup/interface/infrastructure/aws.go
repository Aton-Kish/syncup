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

package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	credscache "github.com/Aton-Kish/aws-credscache-go/sdkv2"
	"github.com/Aton-Kish/syncup/internal/syncup/domain/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

type appsyncClient interface{}

type awsClients struct {
	appsyncClient appsyncClient
}

var (
	awsClts         *awsClients
	awsActivateOnce sync.Once
)

func activatedAWSClients(ctx context.Context, optFns ...func(o *model.AWSOptions)) (*awsClients, error) {
	var err error

	awsActivateOnce.Do(func() {
		o := model.NewAWSOptions(optFns...)

		var cfg aws.Config
		cfg, err = config.LoadDefaultConfig(
			ctx,
			config.WithRegion(o.Region),
			config.WithSharedConfigProfile(o.Profile),
			config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
				aro.TokenProvider = o.MFATokenProvider
			}),
		)
		if err != nil {
			err = &model.LibError{Err: err}
			return
		}

		home, _ := os.UserHomeDir()
		if len(home) > 0 {
			if _, err = credscache.InjectFileCacheProvider(
				&cfg,
				func(fco *credscache.FileCacheOptions) {
					fco.FileCacheDir = filepath.Join(home, ".aws", "cli", "cache")
				},
			); err != nil {
				err = &model.LibError{Err: err}
				return
			}
		}

		if _, err = cfg.Credentials.Retrieve(ctx); err != nil {
			err = &model.LibError{Err: err}
			return
		}

		awsClts = &awsClients{
			appsyncClient: appsync.NewFromConfig(cfg),
		}
	})

	return awsClts, err
}
