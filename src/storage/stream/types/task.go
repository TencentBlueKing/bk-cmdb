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

package types

import (
	"context"
	"errors"
	"time"

	"configcenter/src/storage/dal/mongo/local"
)

// NewTaskOptions is the new task options
type NewTaskOptions struct {
	StopNotifier <-chan struct{}
}

// Validate NewTaskOptions
func (o *NewTaskOptions) Validate() error {
	if o.StopNotifier == nil {
		// if not set, then set never stop loop as default
		o.StopNotifier = make(<-chan struct{})
	}
	return nil
}

// TaskTokenHandler is the token handler for db watch task
type TaskTokenHandler interface {
	SetLastWatchToken(ctx context.Context, uuid string, watchDB local.DB, token *TokenInfo) error
	GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*TokenInfo, error)
}

// WatchTaskOptions is the common options for watch task
type WatchTaskOptions struct {
	Name              string
	CollOpts          *WatchCollOptions
	TokenHandler      TaskTokenHandler
	RetryOptions      *RetryOptions
	MajorityCommitted *bool
	MaxAwaitTime      *time.Duration
}

// Validate WatchTaskOptions
func (o *WatchTaskOptions) Validate() error {
	if len(o.Name) == 0 {
		return errors.New("watch task name is not set")
	}

	if o.CollOpts == nil {
		return errors.New("watch task coll options is not set")
	}

	if err := o.CollOpts.Validate(); err != nil {
		return err
	}

	if o.TokenHandler == nil {
		return errors.New("token handler is not set")
	}

	if o.TokenHandler.SetLastWatchToken == nil || o.TokenHandler.GetStartWatchToken == nil {
		return errors.New("some token handler functions is not set")
	}

	if o.RetryOptions != nil {
		if o.RetryOptions.MaxRetryCount <= 0 {
			o.RetryOptions.MaxRetryCount = DefaultRetryCount
		}

		if o.RetryOptions.RetryDuration == 0 {
			o.RetryOptions.RetryDuration = DefaultRetryDuration
		}

		if o.RetryOptions.RetryDuration < 500*time.Millisecond {
			return errors.New("invalid retry duration, can not less than 500ms")
		}
	} else {
		o.RetryOptions = &RetryOptions{
			MaxRetryCount: DefaultRetryCount,
			RetryDuration: DefaultRetryDuration,
		}
	}

	return nil
}

// LoopOneTaskOptions is the options for loop watch events one by one operation of one task
type LoopOneTaskOptions struct {
	*WatchTaskOptions
	EventHandler *TaskOneHandler
}

// Validate LoopOneTaskOptions
func (o *LoopOneTaskOptions) Validate() error {
	if o.WatchTaskOptions == nil {
		return errors.New("common watch task options is not set")
	}

	if err := o.WatchTaskOptions.Validate(); err != nil {
		return err
	}

	if o.EventHandler == nil {
		return errors.New("event handler is not set")
	}

	if o.EventHandler.DoAdd == nil || o.EventHandler.DoUpdate == nil || o.EventHandler.DoDelete == nil {
		return errors.New("some event handler functions is not set")
	}
	return nil
}

// TaskOneHandler is the watch task's event handler that handles events one by one
type TaskOneHandler struct {
	DoAdd    func(dbInfo *DBInfo, event *Event) (retry bool)
	DoUpdate func(dbInfo *DBInfo, event *Event) (retry bool)
	DoDelete func(dbInfo *DBInfo, event *Event) (retry bool)
}

// LoopBatchTaskOptions is the options for loop watch batch events operation of one task
type LoopBatchTaskOptions struct {
	*WatchTaskOptions
	BatchSize    int
	EventHandler *TaskBatchHandler
}

// Validate LoopBatchTaskOptions
func (o *LoopBatchTaskOptions) Validate() error {
	if o.WatchTaskOptions == nil {
		return errors.New("common watch task options is not set")
	}

	if err := o.WatchTaskOptions.Validate(); err != nil {
		return err
	}

	if o.BatchSize <= 0 {
		return errors.New("batch size is invalid")
	}

	if o.EventHandler == nil {
		return errors.New("event handler is not set")
	}

	if o.EventHandler.DoBatch == nil {
		return errors.New("event handler DoBatch function is not set")
	}
	return nil
}

// TaskBatchHandler is the watch task's batch events handler
type TaskBatchHandler struct {
	DoBatch func(dbInfo *DBInfo, es []*Event) bool
}

// DBInfo is the db info for watch task
type DBInfo struct {
	// UUID is the cc db uuid
	UUID    string
	WatchDB *local.Mongo
	CcDB    local.DB
}
