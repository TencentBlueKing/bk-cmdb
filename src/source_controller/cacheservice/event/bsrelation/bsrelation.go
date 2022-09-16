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

// Package bsrelation TODO
package bsrelation

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/event"
	mixevent "configcenter/src/source_controller/cacheservice/event/mix-event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream"
)

const (
	bizSetRelationLockKey = common.BKCacheKeyV3Prefix + "biz_set_relation:event_lock"
	bizSetRelationLockTTL = 1 * time.Minute
)

// NewBizSetRelation init and run biz set relation event watch
func NewBizSetRelation(watch stream.LoopInterface, watchDB *local.Mongo, ccDB dal.DB) error {
	base := mixevent.MixEventFlowOptions{
		MixKey:       event.BizSetRelationKey,
		Watch:        watch,
		WatchDB:      watchDB,
		CcDB:         ccDB,
		EventLockKey: bizSetRelationLockKey,
		EventLockTTL: bizSetRelationLockTTL,
	}

	// watch biz set event
	bizSet := base
	bizSet.Key = event.BizSetKey
	bizSet.WatchFields = []string{common.BKBizSetIDField, common.BKBizSetScopeField}
	if err := newBizSetRelation(context.Background(), bizSet); err != nil {
		blog.Errorf("watch biz set event for biz set relation failed, err: %v", err)
		return err
	}
	blog.Info("watch biz set relation events, watch biz set success")

	// watch biz event
	biz := base
	biz.Key = event.BizKey
	if err := newBizSetRelation(context.Background(), biz); err != nil {
		blog.Errorf("watch biz event for biz set relation failed, err: %v", err)
		return err
	}
	blog.Info("watch biz set relation events, watch biz success")

	return nil
}
