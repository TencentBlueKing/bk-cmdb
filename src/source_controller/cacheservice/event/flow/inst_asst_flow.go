/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package flow

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"
)

func newInstAsstFlow(ctx context.Context, opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}
	instAsstFlow := InstAsstFlow{
		Flow: flow,
	}

	return instAsstFlow.RunFlow(ctx)
}

// InstAsstFlow instance association event watch flow
type InstAsstFlow struct {
	Flow
}

// RunFlow run instance association event watch flow
func (f *InstAsstFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.metrics)

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name: f.key.Namespace(),
			CollOpts: &types.WatchCollOptions{
				CollectionOptions: types.CollectionOptions{
					CollectionFilter: &types.CollectionFilter{
						Regex: fmt.Sprintf("_%s", common.BKObjectInstAsstShardingTablePrefix),
					},
					EventStruct: f.EventStruct,
				},
			},
			TokenHandler: f.tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.TaskBatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	err := f.task.AddLoopBatchTask(opts)
	if err != nil {
		blog.Errorf("run %s flow, but add loop batch task failed, err: %v", f.key.Namespace(), err)
		return err
	}

	return nil
}
