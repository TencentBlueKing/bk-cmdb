/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package task defines event watch task logics
package task

import (
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"
)

// Task is the resource watch task
type Task struct {
	// Name is the watch task name that uniquely identifies the watch task
	Name string
	// CollOptions is the watch collection options
	CollOptions *types.WatchCollOptions
	// eventHandler is the batch event handler
	eventHandler *types.TaskBatchHandler
	// tokenHandler is the token handler
	tokenHandler types.TaskTokenHandler
	// retryOptions is the watch task retry options
	retryOptions *types.RetryOptions
	// NeedList defines whether to list all data before watch
	NeedList bool
	// BatchSize is the batch event size for one loop
	BatchSize         int
	MajorityCommitted *bool
	MaxAwaitTime      *time.Duration
}

// NewLoopOneTask create a loop watch task that handles one event at one time
func NewLoopOneTask(opts *types.LoopOneTaskOptions) (*Task, error) {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate loop one task options(%s) failed, err: %v", opts.Name, err)
		return nil, err
	}

	batchOpts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: opts.WatchTaskOptions,
		BatchSize:        1,
		EventHandler: &types.TaskBatchHandler{
			DoBatch: func(dbInfo *types.DBInfo, es []*types.Event) bool {
				for _, e := range es {
					var retry bool
					switch e.OperationType {
					case types.Insert:
						retry = opts.EventHandler.DoAdd(dbInfo, e)
					case types.Update, types.Replace:
						retry = opts.EventHandler.DoUpdate(dbInfo, e)
					case types.Delete:
						retry = opts.EventHandler.DoDelete(dbInfo, e)
					default:
						blog.Warnf("received unsupported operation type for %s job, doc: %s", opts.Name, e.DocBytes)
						continue
					}
					if retry {
						return retry
					}
				}
				return false
			},
		},
	}

	return newTask(batchOpts, false), nil
}

// NewLoopBatchTask add a loop watch task that handles batch events
func NewLoopBatchTask(opts *types.LoopBatchTaskOptions) (*Task, error) {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate loop batch task options(%s) failed, err: %v", opts.Name, err)
		return nil, err
	}

	return newTask(opts, false), nil
}

// NewListWatchTask add a list watch task
func NewListWatchTask(opts *types.LoopBatchTaskOptions) (*Task, error) {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate list watch task options(%s) failed, err: %v", opts.Name, err)
		return nil, err
	}

	return newTask(opts, true), nil
}

// newTask generate a new watch task
func newTask(opts *types.LoopBatchTaskOptions, needList bool) *Task {
	return &Task{
		Name:              opts.Name,
		CollOptions:       opts.CollOpts,
		eventHandler:      opts.EventHandler,
		tokenHandler:      opts.TokenHandler,
		retryOptions:      opts.RetryOptions,
		NeedList:          needList,
		BatchSize:         opts.BatchSize,
		MajorityCommitted: opts.MajorityCommitted,
		MaxAwaitTime:      opts.MaxAwaitTime,
	}
}
