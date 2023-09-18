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

type appsyncClient interface {
	GetIntrospectionSchema(ctx context.Context, params *appsync.GetIntrospectionSchemaInput, optFns ...func(*appsync.Options)) (*appsync.GetIntrospectionSchemaOutput, error)
	StartSchemaCreation(ctx context.Context, params *appsync.StartSchemaCreationInput, optFns ...func(*appsync.Options)) (*appsync.StartSchemaCreationOutput, error)
	GetSchemaCreationStatus(ctx context.Context, params *appsync.GetSchemaCreationStatusInput, optFns ...func(*appsync.Options)) (*appsync.GetSchemaCreationStatusOutput, error)

	ListFunctions(ctx context.Context, params *appsync.ListFunctionsInput, optFns ...func(*appsync.Options)) (*appsync.ListFunctionsOutput, error)
	CreateFunction(ctx context.Context, params *appsync.CreateFunctionInput, optFns ...func(*appsync.Options)) (*appsync.CreateFunctionOutput, error)
	UpdateFunction(ctx context.Context, params *appsync.UpdateFunctionInput, optFns ...func(*appsync.Options)) (*appsync.UpdateFunctionOutput, error)
	DeleteFunction(ctx context.Context, params *appsync.DeleteFunctionInput, optFns ...func(*appsync.Options)) (*appsync.DeleteFunctionOutput, error)

	ListResolvers(ctx context.Context, params *appsync.ListResolversInput, optFns ...func(*appsync.Options)) (*appsync.ListResolversOutput, error)
	GetResolver(ctx context.Context, params *appsync.GetResolverInput, optFns ...func(*appsync.Options)) (*appsync.GetResolverOutput, error)
	CreateResolver(ctx context.Context, params *appsync.CreateResolverInput, optFns ...func(*appsync.Options)) (*appsync.CreateResolverOutput, error)
	UpdateResolver(ctx context.Context, params *appsync.UpdateResolverInput, optFns ...func(*appsync.Options)) (*appsync.UpdateResolverOutput, error)
	DeleteResolver(ctx context.Context, params *appsync.DeleteResolverInput, optFns ...func(*appsync.Options)) (*appsync.DeleteResolverOutput, error)

	ListTypes(ctx context.Context, params *appsync.ListTypesInput, optFns ...func(*appsync.Options)) (*appsync.ListTypesOutput, error)
}

type awsClients struct {
	appsyncClient appsyncClient
}

var (
	awsClts         *awsClients
	awsActivateOnce sync.Once
)

func activatedAWSClients(ctx context.Context, optFns ...func(o *model.AWSOptions)) (res *awsClients, err error) {
	defer wrap(&err)

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
				return
			}
		}

		if _, err = cfg.Credentials.Retrieve(ctx); err != nil {
			return
		}

		awsClts = &awsClients{
			appsyncClient: appsync.NewFromConfig(cfg),
		}
	})

	return awsClts, err
}
