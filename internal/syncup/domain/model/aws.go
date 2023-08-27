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

package model

type MFATokenProvider func() (string, error)

type AWSOptions struct {
	Region           string
	Profile          string
	MFATokenProvider MFATokenProvider
}

func NewAWSOptions(optFns ...func(o *AWSOptions)) *AWSOptions {
	o := new(AWSOptions)

	for _, fn := range optFns {
		fn(o)
	}

	return o
}

func AWSOptionsWithRegion(region string) func(o *AWSOptions) {
	return func(o *AWSOptions) {
		o.Region = region
	}
}

func AWSOptionsWithProfile(profile string) func(o *AWSOptions) {
	return func(o *AWSOptions) {
		o.Profile = profile
	}
}

func AWSOptionsWithMFATokenProvider(provider MFATokenProvider) func(o *AWSOptions) {
	return func(o *AWSOptions) {
		o.MFATokenProvider = provider
	}
}
