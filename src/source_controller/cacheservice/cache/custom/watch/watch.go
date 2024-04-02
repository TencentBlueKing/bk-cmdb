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

// Package watch defines the custom resource cache data watch logics
package watch

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
)

// Watcher defines mongodb event watcher for custom resource
type Watcher struct {
	loopW    stream.LoopInterface
	cacheSet *cache.CacheSet
}

// Init custom resource mongodb event watcher
func Init(loopW stream.LoopInterface, cacheSet *cache.CacheSet) error {
	watcher := &Watcher{
		loopW:    loopW,
		cacheSet: cacheSet,
	}

	if err := watcher.watchPodLabel(); err != nil {
		return err
	}

	if err := watcher.watchSharedNsRel(); err != nil {
		return err
	}

	return nil
}

type watchOptions struct {
	watchType WatchType
	watchOpts *types.WatchOptions
	doBatch   func(es []*types.Event) (retry bool)
}

// WatchType is the custom resource watch type
type WatchType string

const (
	// PodLabelWatchType is the kube pod label watch type
	PodLabelWatchType WatchType = "pod_label"
	// SharedNsRelWatchType is the shared namespace relation watch type
	SharedNsRelWatchType WatchType = "shared_namespace_relation"
)

// watchCustomResource watch custom resource
func (w *Watcher) watchCustomResource(opt *watchOptions) (bool, error) {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	name := fmt.Sprintf("%s:%s", cache.Namespace, opt.watchType)

	tokenHandler := tokenhandler.NewSingleTokenHandler(name, mongodb.Client())

	exists, err := tokenHandler.IsTokenExists(ctx)
	if err != nil {
		blog.Errorf("check if custom resource %s watch token exists failed, err: %v", name, err)
		return false, err
	}

	if exists {
		startAtTime, err := tokenHandler.GetStartWatchTime(ctx)
		if err != nil {
			blog.Errorf("get custom resource %s start watch time failed, err: %v", name, err)
			return false, err
		}
		opt.watchOpts.StartAtTime = startAtTime
	} else {
		opt.watchOpts.StartAtTime = new(types.TimeStamp)
	}

	opt.watchOpts.WatchFatalErrorCallback = tokenHandler.ResetWatchToken

	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         name,
			WatchOpt:     opt.watchOpts,
			TokenHandler: tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 3,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: opt.doBatch,
		},
		BatchSize: 200,
	}

	if err = w.loopW.WithBatch(loopOptions); err != nil {
		blog.Errorf("watch custom resource %s failed, err: %v", name, err)
		return false, err
	}

	return exists, nil
}
