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
	"sync"

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

type customLevel struct {
	key   customKeyGen
	event reflector.Interface
	lock  tools.RefreshingLock
	// key is object id
	customWatch map[string]context.CancelFunc
	customLock  sync.Mutex
}

func (m *customLevel) Run() error {

	if err := m.runMainlineTopology(); err != nil {
		return fmt.Errorf("watch mainline topology association failed, err: %v", err)
	}

	if err := m.runCustomLevelInstance(); err != nil {
		return fmt.Errorf("run mainline instance watch failed, err: %v", err)
	}

	return nil
}

// to watch mainline topology change for it's cache update.
func (m *customLevel) runMainlineTopology() error {

	_, err := redis.Client().Get(context.Background(), m.key.mainlineListDoneKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			blog.Errorf("get biz list done redis key failed, err: %v", err)
			return fmt.Errorf("get biz list done redis key failed, err: %v", err)
		}
		cap := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister:     m.onChange,
				OnAdd:        m.onChange,
				OnUpdate:     m.onChange,
				OnListerDone: m.onMainlineTopologyListDone,
				OnDelete:     m.onChange,
			},
		}
		page := 500
		opts := &types.ListWatchOptions{
			Options: types.Options{
				EventStruct: new(map[string]interface{}),
				Collection:  common.BKTableNameObjAsst,
				Filter: mapstr.MapStr{
					common.AssociationKindIDField: common.AssociationKindMainline,
				},
			},
			PageSize: &page,
		}
		blog.Info("do mainline topology cache with list watch.")
		return m.event.ListWatcher(context.Background(), opts, cap)
	}

	watchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnAdd:    m.onChange,
			OnUpdate: m.onChange,
			OnDelete: m.onChange,
		},
	}
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: new(map[string]interface{}),
			Collection:  common.BKTableNameObjAsst,
			Filter: mapstr.MapStr{
				common.AssociationKindIDField: common.AssociationKindMainline,
			},
		},
	}
	blog.Info("do mainline topology cache with only watch.")
	return m.event.Watcher(context.Background(), watchOpts, watchCap)
}

// to watch each custom level object instance's change for cache update.
func (m *customLevel) runCustomLevelInstance() error {

	relations, err := m.getMainlineTopology()
	if err != nil {
		blog.Errorf("received mainline topology change event, but get it from mongodb failed, err: %v", err)
		return err
	}
	// rank start from biz to host
	rank := m.rankMainlineTopology(relations)

	// reconcile the custom watch
	m.reconcileCustomWatch(rank)

	for _, objID := range rank {
		if objID == "biz" || objID == "module" || objID == "set" || objID == "host" {
			// skip system embed object
			continue
		}

		if err := m.runCustomWatch(objID); err != nil {
			blog.Errorf("run biz custom level %s watch failed. err: %v", err)
			return err
		}
	}
	return nil
}

func (m *customLevel) runCustomWatch(objID string) error {

	opts := types.Options{
		EventStruct: new(map[string]interface{}),
		Collection:  common.BKTableNameBaseInst,
		Filter: mapstr.MapStr{
			common.BKObjIDField: objID,
		},
	}

	_, err := redis.Client().Get(context.Background(), m.key.listDoneKey(objID)).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			blog.Errorf("get biz mainline key: %s list done redis key failed, err: %v", m.key.listDoneKey(objID), err)
			return fmt.Errorf("get biz mainline key: %s list done redis key failed, err: %v", m.key.listDoneKey(objID), err)
		}
		listCap := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister: m.onUpsertCustomInstance,
				OnAdd:    m.onUpsertCustomInstance,
				OnUpdate: m.onUpsertCustomInstance,
				OnListerDone: func() {
					m.onCustomInstanceListDone(objID)
				},
				OnDelete: m.onDeleteCustomInstance,
			},
		}
		// do with list watch
		page := 500
		listOpts := &types.ListWatchOptions{
			Options:  opts,
			PageSize: &page,
		}

		ctx, cancel := context.WithCancel(context.Background())
		m.customLock.Lock()
		m.customWatch[objID] = cancel
		m.customLock.Unlock()
		blog.Infof("do custom level object: %s instance sync cache with list watch.", objID)
		return m.event.ListWatcher(ctx, listOpts, listCap)

	}

	watchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnAdd:    m.onUpsertCustomInstance,
			OnUpdate: m.onUpsertCustomInstance,
			OnDelete: m.onDeleteCustomInstance,
		},
	}

	watchOpts := &types.WatchOptions{
		Options: opts,
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.customLock.Lock()
	m.customWatch[objID] = cancel
	m.customLock.Unlock()
	blog.Infof("do custom level object: %s instance sync cache with only watch.", objID)
	return m.event.Watcher(ctx, watchOpts, watchCap)
}

// reconcileCustomWatch to check if the custom watch is coordinate with
// watch or not.
func (m *customLevel) reconcileCustomWatch(rank []string) {
	rankMap := make(map[string]bool)
	for _, objID := range rank {
		if objID == "biz" || objID == "module" || objID == "set" || objID == "host" {
			// skip system embed object
			continue
		}
		rankMap[objID] = true
	}

	// check if new watch is need.
	for objID := range rankMap {
		m.customLock.Lock()
		_, exist := m.customWatch[objID]
		m.customLock.Unlock()
		if exist {
			continue
		}
		// not exist, need to add a new watch immediately.
		if err := m.runCustomWatch(objID); err != nil {
			blog.Errorf("reconcile custom watch and need new watch with object %s, but run watch it failed, err: %v", objID, err)
			continue
		}
	}

	// check if the watch need to be canceled.
	m.customLock.Lock()
	for obj, cancel := range m.customWatch {
		if _, exist := rankMap[obj]; !exist {
			blog.Warnf("reconcile custom watch, find a redundant one with object %s, cancel it now.", obj)
			// cancel the watch immediately.
			cancel()
			// Obviously, we need to delete list keys, instance key and expire key belongs to this object.
			// Normally, we do not need to do this, because of if this object has instances, then it can not be deleted.
			// so if the object has already be delete, then the instances is no longer exist at the same time.
			// if we do not delete instance failed unfortunately, the cache may stay, but we have the trigger to
			// refresh it, it will be deleted later.
		}
	}
	m.customLock.Unlock()
	blog.V(4).Infof("reconcile custom watch with rank: %v finished.", rank)
}

// onChange is to reform the topology with the info from mongodb
// we do not use the event data, because it's complicated to change
// a mainline topology and a change is always associated with add and
// delete operation.
func (m *customLevel) onChange(_ *types.Event) {
	// read information from mongodb
	relations, err := m.getMainlineTopology()
	if err != nil {
		blog.Errorf("received mainline topology change event, but get it from mongodb failed, err: %v", err)
		return
	}
	// rank start from biz to host
	rank := m.rankMainlineTopology(relations)

	// reconcile the custom watch
	m.reconcileCustomWatch(rank)

	// then set the rank to cache
	redis.Client().Set(context.Background(), m.key.topologyKey(), m.key.topologyValue(rank), 0)
}

func (m *customLevel) onMainlineTopologyListDone() {
	if err := redis.Client().Set(context.Background(), m.key.mainlineListDoneKey(), "done", 0).Err(); err != nil {
		blog.Errorf("list business mainline topology to cache and list done, but set list done key: %s failed, err: %v",
			m.key.mainlineListDoneKey(), err)
		return
	}
	blog.Info("list business mainline topology to cache and list done")
}

// onUpsertCustomInstance is to upsert the custom object instance cache.
// do operation like:
// 1. upsert the object instance list cache.
// 2. upsert the object instance detail cache.
// 3. upsert the object instance oid relation for delete it usage.
func (m *customLevel) onUpsertCustomInstance(e *types.Event) {
	blog.V(4).Infof("received biz custom level instance upsert event, detail: %s", e.String())

	fields := gjson.GetManyBytes(e.DocBytes, "bk_obj_id", "bk_inst_id", "bk_inst_name", "metadata.label.bk_biz_id", "bk_parent_id")
	objID := fields[0].String()
	if len(objID) == 0 {
		blog.Errorf("received biz custom level instance upsert event, but parse object id failed, doc: %s", e.String())
		return
	}

	instID := fields[1].Int()
	if instID == 0 {
		blog.Errorf("received biz custom level instance upsert event, but parse object instance id failed, doc: %s", e.String())
		return
	}

	instName := fields[2].String()
	if len(instName) == 0 {
		blog.Errorf("received biz custom level instance upsert event, but parse object instance name failed, doc: %s", e.String())
		return
	}

	biz := fields[3].String()
	bizID, err := strconv.ParseInt(biz, 10, 64)
	if err != nil {
		blog.Errorf("received biz custom level instance upsert event, but parse business id failed failed, doc: %s, err: %v", e.String(), err)
		return
	}

	parentID := fields[4].Int()
	if parentID == 0 {
		blog.Errorf("received biz custom level instance upsert event, but parse parent id failed, doc: %s", e.String())
		return
	}

	forUpsert := &forUpsertCache{
		instID:            instID,
		parentID:          parentID,
		name:              instName,
		doc:               e.DocBytes,
		listKey:           customKey.objListKeyWithBiz(objID, bizID),
		listExpireKey:     customKey.objListExpireKeyWithBiz(objID, bizID),
		detailKey:         customKey.detailKey(objID, instID),
		detailExpireKey:   customKey.detailExpireKey(objID, bizID),
		parseListKeyValue: customKey.parseListKeyValue,
		genListKeyValue:   customKey.genListKeyValue,
		getInstName:       m.getCustomObjInstName,
	}
	// update the cache
	upsertListCache(forUpsert)

	// record the object id relation
	m.upsertOid(objID, bizID, instID, e.Oid)
}

func (m *customLevel) onCustomInstanceListDone(objID string) {
	if err := redis.Client().Set(context.Background(), m.key.listDoneKey(objID), "done", 0).Err(); err != nil {
		blog.Errorf("list business custom level %s to cache and list done, but set list done key: %s failed, err: %v",
			objID, m.key.listDoneKey(objID), err)
		return
	}
	blog.Infof("list business custom level %s to cache and list done", objID)
}

func (m *customLevel) onDeleteCustomInstance(e *types.Event) {
	blog.V(4).Infof("received biz custom level instance delete event, detail: %s", e.String())
	// get the instance id and business id from oid key
	v, err := m.getOidValue(e.Oid)
	if err != nil {
		blog.Errorf("received biz custom level instance delete event, but get oid relation failed, detail: %s, err: %v", e.String(), err)
		return
	}

	pipeline := redis.Client().Pipeline()
	pipeline.Del(customKey.detailKey(v.obj, v.instID))
	pipeline.Del(customKey.detailExpireKey(v.obj, v.instID))
	pipeline.HDel(customKey.objectIDKey(), e.Oid)
	pipeline.SRem(customKey.objListKeyWithBiz(v.obj, v.biz))
	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("received biz custom level instance delete event, but remove related keys failed, detail: %s, err: %v", e.String(), err)
		return
	}
	blog.V(4).Infof("received biz custom level instance delete event, detail: %s, delete related caches success", e.String())
}

func (m *customLevel) getMainlineTopology() ([]MainlineTopoAssociation, error) {
	relations := make([]MainlineTopoAssociation, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(filter).All(context.Background(), &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v", err)
		return nil, err
	}
	return relations, nil
}

// rankTopology is to rank the biz topology to a array, start from biz to host
func (m *customLevel) rankMainlineTopology(relations []MainlineTopoAssociation) []string {
	rank := make([]string, 0)
	next := "biz"
	rank = append(rank, next)
	for _, relation := range relations {
		if relation.AssociateTo == next {
			rank = append(rank, relation.ObjectID)
			next = relation.ObjectID
			continue
		} else {
			for _, rel := range relations {
				if rel.AssociateTo == next {
					rank = append(rank, rel.ObjectID)
					next = rel.ObjectID
					break
				}
			}
		}
	}
	return rank
}

func (m *customLevel) getCustomObjInstName(instID int64) (name string, err error) {
	instance := new(CustomInstanceBase)
	filter := mapstr.MapStr{
		common.BKInstIDField: instID,
	}
	err = mongodb.Client().Table(common.BKTableNameBaseInst).Find(filter).One(context.Background(), instance)
	if err != nil {
		blog.Errorf("find mainline custom level with instance: %d, failed, err: %v", instID, err)
		return "", err
	}
	return instance.InstanceName, nil
}

func (m *customLevel) upsertOid(objID string, bizID int64, instID int64, oid string) {
	value := customKey.genObjectIDKeyValue(bizID, instID, objID)
	if err := redis.Client().HSet(context.Background(), customKey.objectIDKey(), oid, value).Err(); err != nil {
		blog.Errorf("upsert business custom level object instance oid: %s relation failed, key: %s, value: %s, err: %v",
			oid, customKey.objectIDKey(), value, err)
	}
}

func (m *customLevel) delOid(oid string) {
	if err := redis.Client().HDel(context.Background(), customKey.objectIDKey(), oid).Err(); err != nil {
		blog.Errorf("delete business custom level object instance oid: %s relation failed, key: %s, err: %v",
			oid, customKey.objectIDKey(), err)
	}
}

func (m *customLevel) getOidValue(oid string) (*oidValue, error) {
	value, err := redis.Client().HGet(context.Background(), customKey.objectIDKey(), oid).Result()
	if err != nil {
		blog.Errorf("get business custom level object instance oid: %s relation failed, key: %s, err: %v",
			oid, customKey.objectIDKey(), err)
		return nil, err
	}

	return customKey.parseObjectIDKeyValue(value)
}
