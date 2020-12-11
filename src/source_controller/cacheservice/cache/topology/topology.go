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

package topology

import (
	"context"
	"sync"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
	drvRedis "configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
)

func NewTopology(isMaster discovery.ServiceManageInterface, loopW stream.LoopInterface) (*Topology, error) {

	t := &Topology{
		db:          mongodb.Client(),
		rds:         drvRedis.Client(),
		loopW:       loopW,
		checkMaster: isMaster,
		briefBizKey: newTopologyKey(),
	}

	if err := t.watchCustom(); err != nil {
		blog.Errorf("topology watch custom failed, err: %v", err)
		return nil, err
	}

	if err := t.watchSet(); err != nil {
		blog.Errorf("topology watch set failed, err: %v", err)
		return nil, err
	}

	if err := t.watchModule(); err != nil {
		blog.Errorf("topology watch module failed, err: %v", err)
		return nil, err
	}

	go t.loopBizBriefCache()

	return t, nil
}

type Topology struct {
	db          dal.DB
	rds         redis.Client
	loopW       stream.LoopInterface
	checkMaster discovery.ServiceManageInterface
	briefBizKey *cacheKey
}

// refreshBatch refresh business Topology with batch
func (t *Topology) refreshBatch(bizList []int64, rid string) error {
	blog.Infof("try to refresh biz: %v topology, rid: %s", bizList, rid)

	if len(bizList) == 0 {
		return nil
	}

	filter := mapstr.MapStr{
		common.BKAppIDField: mapstr.MapStr{
			common.BKDBIN: bizList,
		},
	}
	list := make([]*BizBase, 0)
	err := t.db.Table(common.BKTableNameBaseApp).Find(filter).Fields(bizBaseFields...).All(context.Background(), &list)
	if err != nil {
		blog.Errorf("list biz detail failed, err: %v, rid: %s", err, rid)
		return err
	}

	// set max goroutine number
	pipeline := make(chan struct{}, 5)
	wg := sync.WaitGroup{}
	var hitErr error
	for idx := range list {
		pipeline <- struct{}{}
		wg.Add(1)

		go func(biz *BizBase) {
			err := t.refreshBizTopology(biz, rid)
			if err != nil {
				hitErr = err
			}
			<-pipeline
			wg.Done()
		}(list[idx])
	}

	wg.Wait()
	if hitErr != nil {
		blog.Errorf("refresh biz list failed, err: %v, rid: %s", err, rid)
		return hitErr
	}

	blog.Infof("try to refresh biz topology success, rid: %s", rid)
	return nil
}

// refreshBizTopology construct a business Topology and update it to cache.
func (t *Topology) refreshBizTopology(biz *BizBase, rid string) error {
	ctx := context.WithValue(context.TODO(), common.ContextRequestIDField, rid)
	idle, common, err := t.getBusinessTopology(ctx, biz.ID)
	if err != nil {
		blog.Error("refresh biz %d/%s topology, but get topology failed, err: %v, rid: %s", biz.ID, biz.Name, err, rid)
		return err
	}

	topo := &BizBriefTopology{
		Biz:   biz,
		Idle:  idle,
		Nodes: common,
	}

	err = t.briefBizKey.updateTopology(ctx, topo)
	if err != nil {
		blog.Error("update biz %d/%s topology to cache failed, err: %v, rid: %s", biz.ID, biz.Name, err, rid)
		return err
	}

	return nil
}

// loopBizBriefCache launch the task to loop business's brief topology every interval minutes.
func (t *Topology) loopBizBriefCache() {
	blog.Infof("loop refresh biz brief topology task every %d minutes.", getBreifTopoCacheRefreshMinutes())
	for {

		if !t.checkMaster.IsMaster() {
			blog.V(4).Infof("loop biz brief cache, but not master, skip.")
			time.Sleep(time.Minute)
			continue
		}

		interval := getBreifTopoCacheRefreshMinutes()
		time.Sleep(time.Duration(interval) * time.Minute)
		// time.Sleep(30 * time.Second)
		rid := util.GenerateRID()

		blog.Infof("start loop refresh biz brief topology task, interval: %d, rid: %s", interval, rid)
		t.doLoopBizBriefTopologyToCache(rid)
		blog.Infof("finished loop refresh biz brief topology task, rid: %s", rid)
	}
}

func (t *Topology) doLoopBizBriefTopologyToCache(rid string) {
	// read from secondary in mongodb cluster.
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	all, err := t.listAllBusiness(ctx)
	if err != nil {
		blog.Errorf("loop biz brief topology, but list all business failed, err: %v, rid: %s", err, rid)
		return
	}

	for _, biz := range all {
		time.Sleep(50 * time.Millisecond)

		err := t.refreshBizTopology(biz, rid)
		if err != nil {
			blog.Errorf("loop refresh biz %d/%s topology failed, err: %v, rid: %s", biz.ID, biz.Name, err, rid)
		} else {
			blog.Infof("loop refresh biz %d/%s brief topology success, rid: %s", biz.ID, biz.Name, rid)
		}

	}

}
