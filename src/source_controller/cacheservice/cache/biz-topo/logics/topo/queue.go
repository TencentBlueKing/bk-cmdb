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

// Package topo defines business topology related common logics
package topo

import (
	"sync"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

type bizRefreshQueue struct {
	sync.Mutex
	topoKey  key.Key
	bizIDs   []int64
	bizIDMap map[int64]struct{}
}

func newBizRefreshQueue(topoType types.TopoType) *bizRefreshQueue {
	queue := &bizRefreshQueue{
		topoKey:  key.TopoKeyMap[topoType],
		bizIDs:   make([]int64, 0),
		bizIDMap: make(map[int64]struct{}),
	}

	return queue
}

// Run refreshing biz topo cache task
func (q *bizRefreshQueue) Run() {
	for {
		bizID, exists := q.Pop()
		if !exists {
			time.Sleep(time.Millisecond * 50)
			continue
		}

		rid := util.GenerateRID()
		err := TryRefreshBizTopoByCache(q.topoKey, bizID, rid)
		if err != nil {
			blog.Errorf("try refresh biz %d %s topo failed, err: %v, rid: %s", bizID, q.topoKey.Type(), err, rid)
			time.Sleep(time.Millisecond * 100)
			continue
		}
	}
}

// Push some need refresh bizs
func (q *bizRefreshQueue) Push(bizIDs ...int64) {
	q.Lock()
	defer q.Unlock()

	for _, bizID := range bizIDs {
		_, exists := q.bizIDMap[bizID]
		if !exists {
			q.bizIDs = append(q.bizIDs, bizID)
			q.bizIDMap[bizID] = struct{}{}
		}
	}
}

// Pop one need refresh biz
func (q *bizRefreshQueue) Pop() (int64, bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.bizIDs) == 0 {
		return 0, false
	}

	bizID := q.bizIDs[0]
	q.bizIDs = q.bizIDs[1:]
	delete(q.bizIDMap, bizID)

	return bizID, true
}

// Remove one need refresh biz
func (q *bizRefreshQueue) Remove(bizID int64) {
	q.Lock()
	defer q.Unlock()

	_, exists := q.bizIDMap[bizID]
	if !exists {
		return
	}

	for idx, id := range q.bizIDs {
		if id == bizID {
			q.bizIDs = append(q.bizIDs[:idx], q.bizIDs[idx+1:]...)
			break
		}
	}

	delete(q.bizIDMap, bizID)
}

var bizRefreshQueuePool = make(map[types.TopoType]*bizRefreshQueue)

func init() {
	refreshQueueTypes := []types.TopoType{types.KubeType}
	for _, queueType := range refreshQueueTypes {
		queue := newBizRefreshQueue(queueType)
		bizRefreshQueuePool[queueType] = queue
		go queue.Run()
	}
}

// AddRefreshBizTopoTask add refresh biz topo cache task
func AddRefreshBizTopoTask(topoType types.TopoType, bizIDs []int64, rid string) {
	queue, exists := bizRefreshQueuePool[topoType]
	if !exists {
		blog.Errorf("topo type %s has no biz refresh queue, rid: %s", topoType, rid)
		return
	}

	bizIDs = util.IntArrayUnique(bizIDs)
	queue.Push(bizIDs...)
}
