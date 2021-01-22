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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/stream/types"
)

func (t *Topology) watchSet() error {
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: new(setBase),
			Collection:  common.BKTableNameBaseSet,
			Filter:      mapstr.MapStr{},
		},
	}

	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         "topology cache with set",
			WatchOpt:     watchOpts,
			TokenHandler: newTokenHandler("set"),
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: t.onSetChange,
		},
		BatchSize: 50,
	}

	return t.loopW.WithBatch(loopOptions)
}

func (t *Topology) onSetChange(es []*types.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	rid := es[0].ID()
	bizList := make([]int64, 0)
	for idx := range es {
		one := es[idx]

		var set *setBase
		switch one.OperationType {
		case types.Insert:
			set = one.Document.(*setBase)

		case types.Update:
			// when a custom level is deleted from mainline topology, then
			// we will change it's children's parent id, we will received
			// it's parent's update event, and bk_parent_id is changed.
			if _, exists := one.ChangeDesc.UpdatedFields[common.BKParentIDField]; !exists {
				// only handle bk_parent_id changed events
				continue
			}

			set = one.Document.(*setBase)

		case types.Delete:
			filter := mapstr.MapStr{
				"oid":  one.Oid,
				"coll": common.BKTableNameBaseSet,
			}
			archive := new(setArchive)
			err := t.db.Table(common.BKTableNameDelArchive).Find(filter).One(context.TODO(), archive)
			if err != nil {
				blog.Errorf("topology cache, get deleted set %s failed, err: %v, rid: %s", one.Oid, err, rid)
				if t.db.IsNotFoundError(err) {
					blog.Errorf("can not find deleted set %s detail, skip, rid: %s", one.Oid, rid)
					continue
				} else {
					return true
				}
			}

			set = archive.Detail

		default:
			// only handle insert and delete event.
			continue
		}

		blog.Infof("topology cache, received biz: %d, set: %d/%s, op-time: %s, changed event, rid: %s",
			set.Business, set.ID, set.Name, one.ClusterTime.String(), rid)

		bizList = append(bizList, set.Business)
	}

	bizList = util.IntArrayUnique(bizList)

	err := t.refreshBatch(bizList, rid)
	if err != nil {
		return true
	}

	return false
}

func (t *Topology) watchModule() error {
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: new(moduleBase),
			Collection:  common.BKTableNameBaseModule,
			Filter:      mapstr.MapStr{},
		},
	}

	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         "topology cache with module",
			WatchOpt:     watchOpts,
			TokenHandler: newTokenHandler("module"),
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: t.onModuleChange,
		},
		BatchSize: 100,
	}

	return t.loopW.WithBatch(loopOptions)
}

func (t *Topology) onModuleChange(es []*types.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	rid := es[0].ID()
	bizList := make([]int64, 0)
	for idx := range es {
		one := es[idx]

		var module *moduleBase
		switch one.OperationType {
		case types.Insert:
			module = one.Document.(*moduleBase)

		case types.Delete:
			filter := mapstr.MapStr{
				"oid":  one.Oid,
				"coll": common.BKTableNameBaseModule,
			}
			archive := new(moduleArchive)
			err := t.db.Table(common.BKTableNameDelArchive).Find(filter).One(context.TODO(), archive)
			if err != nil {
				blog.Errorf("topology cache, get deleted module %s failed, err: %v, rid: %s", one.Oid, err, rid)
				if t.db.IsNotFoundError(err) {
					blog.Errorf("can not find deleted module %s detail, skip, rid: %s", one.Oid, rid)
					continue
				} else {
					return true
				}
			}

			module = archive.Detail

		default:
			// only handle insert and delete event.
			continue
		}

		blog.Infof("topology cache, received biz: %d, module: %d/%s, op-time: %s, changed event, rid: %s",
			module.Business, module.ID, module.Name, one.ClusterTime.String(), rid)

		bizList = append(bizList, module.Business)
	}

	bizList = util.IntArrayUnique(bizList)

	err := t.refreshBatch(bizList, rid)
	if err != nil {
		return true
	}

	return false
}

// watchCustom watch business custom change event
func (t *Topology) watchCustom() error {
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: new(customBase),
			Collection:  common.BKTableNameBaseInst,
			Filter: mapstr.MapStr{
				common.BKAppIDField: mapstr.MapStr{
					common.BKDBGT: 0,
				},
				common.BKInstIDField: mapstr.MapStr{
					common.BKDBGT: 0,
				},
				common.BKParentIDField: mapstr.MapStr{
					common.BKDBGT: 0,
				},
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBExists: true,
				},
			},
		},
	}

	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         "topology cache with custom level",
			WatchOpt:     watchOpts,
			TokenHandler: newTokenHandler("custom_level"),
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: t.onCustomChange,
		},
		BatchSize: 20,
	}

	return t.loopW.WithBatch(loopOptions)
}

func (t *Topology) onCustomChange(es []*types.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	rid := es[0].ID()
	bizList := make([]int64, 0)
	for idx := range es {
		one := es[idx]

		var custom *customBase
		switch one.OperationType {
		case types.Insert:
			custom = one.Document.(*customBase)

		case types.Update:
			// when a custom level is deleted from mainline topology, then
			// we will change it's children's parent id, we will received
			// it's parent's update event, and bk_parent_id is changed.
			if _, exists := one.ChangeDesc.UpdatedFields[common.BKParentIDField]; !exists {
				// only handle bk_parent_id changed events
				continue
			}

			custom = one.Document.(*customBase)

		case types.Delete:
			filter := mapstr.MapStr{
				"oid":  one.Oid,
				"coll": common.BKTableNameBaseInst,
			}
			archive := new(customArchive)
			err := t.db.Table(common.BKTableNameDelArchive).Find(filter).One(context.TODO(), archive)
			if err != nil {
				blog.Errorf("topology cache, get deleted custom level %s failed, err: %v, rid: %s", one.Oid, err, rid)
				if t.db.IsNotFoundError(err) {
					blog.Errorf("can not find deleted custom level %s detail, skip, rid: %s", one.Oid, rid)
					continue
				} else {
					return true
				}
			}

			custom = archive.Detail

		default:
			// only handle insert, update bk_parent_id , delete event, drop the other event.
			continue
		}

		blog.Infof("topology cache, received biz: %d, custom level %s: %d/%s, op-time: %s, changed event, rid: %s",
			custom.Business, custom.Object, custom.ID, custom.Name, one.ClusterTime.String(), rid)

		bizList = append(bizList, custom.Business)
	}

	bizList = util.IntArrayUnique(bizList)

	err := t.refreshBatch(bizList, rid)
	if err != nil {
		return true
	}

	return false
}
