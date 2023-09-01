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

type Schema string

type Function struct {
	FunctionId              *string     `json:"functionId,omitempty"`
	FunctionArn             *string     `json:"-"`
	Name                    *string     `json:"name,omitempty"`
	Description             *string     `json:"description,omitempty"`
	DataSourceName          *string     `json:"dataSourceName,omitempty"`
	RequestMappingTemplate  *string     `json:"-"`
	ResponseMappingTemplate *string     `json:"-"`
	FunctionVersion         *string     `json:"functionVersion,omitempty"`
	SyncConfig              *SyncConfig `json:"syncConfig,omitempty"`
	MaxBatchSize            int32       `json:"maxBatchSize"`
	Runtime                 *Runtime    `json:"runtime,omitempty"`
	Code                    *string     `json:"-"`
}

type SyncConfig struct {
	ConflictHandler             ConflictHandlerType          `json:"conflictHandler,omitempty"`
	ConflictDetection           ConflictDetectionType        `json:"conflictDetection,omitempty"`
	LambdaConflictHandlerConfig *LambdaConflictHandlerConfig `json:"lambdaConflictHandlerConfig,omitempty"`
}

type ConflictDetectionType string

type ConflictHandlerType string

type LambdaConflictHandlerConfig struct {
	LambdaConflictHandlerArn *string `json:"lambdaConflictHandlerArn,omitempty"`
}

type Runtime struct {
	Name           RuntimeName `json:"name,omitempty"`
	RuntimeVersion *string     `json:"runtimeVersion,omitempty"`
}

type RuntimeName string
