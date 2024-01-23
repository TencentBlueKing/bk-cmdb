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

// Package biztopo defines the business topology caching logics
package biztopo

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	topolgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/topo"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream"
)

// Topo defines the business topology caching logics
type Topo struct {
	isMaster discovery.ServiceManageInterface
	watcher  *watch.Watcher
}

// New Topo
func New(isMaster discovery.ServiceManageInterface, loopW stream.LoopInterface) (*Topo, error) {
	t := &Topo{
		isMaster: isMaster,
	}

	watcher, err := watch.New(loopW)
	if err != nil {
		return nil, fmt.Errorf("new watcher failed, err: %v", err)
	}

	t.watcher = watcher

	for _, topoKey := range key.TopoKeyMap {
		go t.loopBizTopoCache(topoKey)
	}
	return t, nil
}

// loopBizTopoCache launch the task to loop business's brief topology every interval minutes.
func (t *Topo) loopBizTopoCache(topoKey key.Key) {
	for {
		if !t.isMaster.IsMaster() {
			blog.V(4).Infof("loop %s biz brief cache, but not master, skip.", topoKey.Type())
			time.Sleep(time.Minute)
			continue
		}

		interval := topoKey.GetRefreshInterval()
		time.Sleep(interval)

		rid := util.GenerateRID()

		blog.Infof("start loop refresh %s biz topology task, interval: %s, rid: %s", topoKey.Type(), interval, rid)
		t.doLoopBizTopoToCache(topoKey, rid)
		blog.Infof("finished loop refresh %s biz topology task, rid: %s", topoKey.Type(), rid)
	}
}

func (t *Topo) doLoopBizTopoToCache(topoKey key.Key, rid string) {
	// read from secondary in mongodb cluster.
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	all, err := t.listAllBusiness(ctx)
	if err != nil {
		blog.Errorf("loop %s biz topology, but list all business failed, err: %v, rid: %s", topoKey.Type(), err, rid)
		return
	}

	for _, biz := range all {
		time.Sleep(50 * time.Millisecond)

		rid := fmt.Sprintf("%s:%d", rid, biz.BizID)

		err = topolgc.RefreshBizTopo(topoKey, biz.BizID, false, rid)
		if err != nil {
			blog.Errorf("loop refresh biz %d/%s %s topology failed, err: %v, rid: %s", biz.BizID, biz.BizName,
				topoKey.Type(), err, rid)
			continue
		}

		blog.Infof("loop refresh biz %d/%s %s topology success, rid: %s", biz.BizID, biz.BizName, topoKey.Type(), rid)
	}
}

const bizStep = 100

// listAllBusiness list all business brief info
func (t *Topo) listAllBusiness(ctx context.Context) ([]metadata.BizInst, error) {
	filter := mapstr.MapStr{}
	all := make([]metadata.BizInst, 0)

	for {
		oneStep := make([]metadata.BizInst, 0)
		err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(filter).Fields(common.BKAppIDField,
			common.BKAppNameField).Limit(bizStep).Sort(common.BKAppIDField).All(ctx, &oneStep)
		if err != nil {
			return nil, err
		}

		all = append(all, oneStep...)

		if len(oneStep) < bizStep {
			// we got all the data
			break
		}

		// update start position
		filter[common.BKAppIDField] = mapstr.MapStr{common.BKDBGT: oneStep[len(oneStep)-1].BizID}
	}

	return all, nil
}
