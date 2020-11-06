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

package business

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/cacheservice/cache/tools"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
)

type moduleSet struct {
	key        keyGenerator
	collection string
	event      reflector.Interface

	lock tools.RefreshingLock
}

func (ms *moduleSet) Run() error {

	opts := types.Options{
		EventStruct: new(map[string]interface{}),
		Collection:  ms.collection,
	}

	_, err := redis.Client().Get(context.Background(), ms.key.listDoneKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			blog.Errorf("get  %s list done redis key failed, err: %v", ms.collection, err)
			return fmt.Errorf("get %s list done redis key failed, err: %v", ms.collection, err)
		}
		// can not find list done key, do with list watch
		page := 500
		cap := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister:     ms.onUpsert,
				OnAdd:        ms.onUpsert,
				OnUpdate:     ms.onUpsert,
				OnListerDone: ms.onListDone,
				OnDelete:     ms.onDelete,
			},
		}

		listOpts := &types.ListWatchOptions{
			Options:  opts,
			PageSize: &page,
		}
		blog.Infof("do %s cache with list watch", ms.collection)
		return ms.event.ListWatcher(context.Background(), listOpts, cap)
	}

	watchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnAdd:    ms.onUpsert,
			OnUpdate: ms.onUpsert,
			OnDelete: ms.onDelete,
		},
	}
	watchOpts := &types.WatchOptions{
		Options: opts,
	}
	blog.Infof("do %s cache with only watch", ms.collection)
	return ms.event.Watcher(context.Background(), watchOpts, watchCap)

}

func (ms *moduleSet) onUpsert(e *types.Event) {
	blog.V(3).Infof("received %s upsert event, oid: %s, operate: %s, doc: %s", ms.collection, e.Oid, e.OperationType, e.DocBytes)

	var fields []string
	switch ms.collection {
	case common.BKTableNameBaseModule:
		fields = []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField}
	case common.BKTableNameBaseSet:
		fields = []string{common.BKSetIDField, common.BKSetNameField, common.BKAppIDField}
	}
	info := gjson.GetManyBytes(e.DocBytes, fields...)
	instID := info[0].Int()
	name := info[1].String()
	if instID == 0 {
		blog.Errorf("received %s upsert event, invalid instance id: %d, oid: %s", ms.collection, instID, e.Oid)
		return
	}

	if len(name) == 0 {
		blog.Errorf("received %s upsert event, invalid name: %s, oid: %s", ms.collection, name, e.Oid)
		return
	}
	bizID := info[2].Int()
	if bizID == 0 {
		blog.Errorf("received %s upsert event, got biz id is 0, oid: %s", ms.collection, e.Oid)
		return
	}

	// save the oid relation immediately.
	if err := ms.setOidRelation(e.Oid, instID); err != nil {
		blog.Errorf("received %s upsert event, but set oid relation failed, oid: %s, id: %s, err: %v", ms.collection, e.Oid, instID, err)
		// do not return, continue.
	}

	var key keyGenerator
	var parentID int64
	var nameFunc func(int64) (string, error)
	switch ms.collection {
	case common.BKTableNameBaseSet:
		key = setKey
		nameFunc = ms.getSetNameFromMongo
		parentID = gjson.GetBytes(e.DocBytes, "bk_parent_id").Int()
	case common.BKTableNameBaseModule:
		key = moduleKey
		nameFunc = ms.getModuleNameFromMongo
		parentID = gjson.GetBytes(e.DocBytes, "bk_set_id").Int()
	default:
		blog.Errorf("received %s upsert event, unsupported", ms.collection)
		return
	}

	if parentID == 0 {
		blog.Errorf("received %s upsert event, invalid parent id: %d, oid: %s", ms.collection, 0, e.Oid)
		return
	}

	forUpsert := forUpsertCache{
		instID:            instID,
		parentID:          parentID,
		name:              name,
		doc:               e.DocBytes,
		listKey:           key.listKeyWithBiz(bizID),
		listExpireKey:     key.listExpireKeyWithBiz(bizID),
		detailKey:         key.detailKey(instID),
		detailExpireKey:   key.detailExpireKey(instID),
		parseListKeyValue: key.parseListKeyValue,
		genListKeyValue:   key.genListKeyValue,
		getInstName:       nameFunc,
	}
	upsertListCache(&forUpsert)
}

func (ms *moduleSet) onDelete(e *types.Event) {
	blog.V(3).Infof("received %s delete event, oid: %s", ms.collection, e.Oid)
	instID, err := ms.getInstIDWithOid(e.Oid)
	if err != nil {
		blog.Errorf("received %s delete event, but get oid: %s relation failed, err: %v", ms.collection, e.Oid, err)
		return
	}

	pipeline := redis.Client().Pipeline()
	defer pipeline.Close()
	pipeline.HDel(ms.key.objectIDKey(), e.Oid)
	pipeline.Del(ms.key.detailKey(instID))
	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("received %s delete event, oid: %s, but delete data failed, err: %v", ms.collection, e.Oid, err)
		return
	}

	blog.V(4).Infof("received %s delete event, oid: %s, and delete data success.", ms.collection, e.Oid)
}

func (ms *moduleSet) onListDone() {
	if err := redis.Client().Set(context.Background(), ms.key.listDoneKey(), "done", 0).Err(); err != nil {
		blog.Errorf("list %s cache done, but set list done key failed, err: %v", ms.key.listDoneKey(), err)
		return
	}
	blog.Infof("list %s cache done", ms.key.listDoneKey())
}

func (ms *moduleSet) setOidRelation(oid string, instID int64) error {

	return redis.Client().HSet(context.Background(), ms.key.objectIDKey(), oid, instID).Err()
}

func (ms *moduleSet) getInstIDWithOid(oid string) (int64, error) {
	id, err := redis.Client().HGet(context.Background(), ms.key.objectIDKey(), oid).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(id, 10, 64)
}

// getModuleFromMongo return only module id and name.
func (ms *moduleSet) getModuleNameFromMongo(id int64) (string, error) {
	mod := new(ModuleBaseInfo)
	filter := mapstr.MapStr{
		common.BKModuleIDField: id,
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(filter).One(context.Background(), mod); err != nil {
		blog.Errorf("get module %d name from mongo failed, err: %v", id, err)
		return "", err
	}
	return mod.ModuleName, nil
}

func (ms *moduleSet) getSetNameFromMongo(id int64) (string, error) {
	mod := new(ModuleBaseInfo)
	filter := mapstr.MapStr{
		common.BKSetIDField: id,
	}

	if err := mongodb.Client().Table(common.BKTableNameBaseSet).Find(filter).One(context.Background(), mod); err != nil {
		blog.Errorf("get module %d name from mongo failed, err: %v", id, err)
		return "", err
	}
	return mod.ModuleName, nil
}
