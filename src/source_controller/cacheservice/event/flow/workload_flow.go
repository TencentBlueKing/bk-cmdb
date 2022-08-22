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

package flow

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/stream/types"
)

func newWorkloadFlow(ctx context.Context, opts flowOptions, getDeleteEventDetails getDeleteEventDetailsFunc,
	parseEvent parseEventFunc) error {

	flow, err := NewFlow(opts, getDeleteEventDetails, parseEvent)
	if err != nil {
		return err
	}
	workloadFlow := WorkloadFlow{
		Flow: flow,
	}

	return workloadFlow.RunFlow(ctx)
}

// WorkloadFlow instance association event watch flow
type WorkloadFlow struct {
	Flow
}

// RunFlow run instance association event watch flow
func (f *WorkloadFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.watchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.key.Collection(), err)
		return err
	}

	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: f.EventStruct,
			// watch all kube workload tables
			CollectionFilter: map[string]interface{}{
				common.BKDBIN: []string{kubetypes.BKTableNameBaseDeployment, kubetypes.BKTableNameBaseStatefulSet,
					kubetypes.BKTableNameBaseDaemonSet, kubetypes.BKTableNameGameDeployment,
					kubetypes.BKTableNameGameStatefulSet, kubetypes.BKTableNameBaseCronJob,
					kubetypes.BKTableNameBaseJob, kubetypes.BKTableNameBasePodWorkload},
			},
			StartAtTime:             startAtTime,
			WatchFatalErrorCallback: f.tokenHandler.resetWatchToken,
		},
	}

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         f.key.Namespace(),
			WatchOpt:     watchOpts,
			TokenHandler: f.tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}
