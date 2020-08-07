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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccError "configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"gopkg.in/redis.v5"
)

func (c *Client) getBusinessFromMongo(bizID int64) (string, error) {
	biz := make(map[string]interface{})
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	err := c.db.Table(common.BKTableNameBaseApp).Find(filter).One(context.Background(), &biz)
	if err != nil {
		blog.Errorf("get business %d info from db, but failed, err: %v", bizID, err)
		return "", ccError.New(common.CCErrCommDBSelectFailed, err.Error())
	}
	js, _ := json.Marshal(biz)
	return string(js), nil
}

// if expireKey is "", then it means you can not use the list array, you need to get
// from the db directly.
func (c *Client) getBusinessBaseInfo() (list []BizBaseInfo, err error) {
	// get all keys which contains the biz id.
	keys, err := c.rds.SMembers(bizKey.listKeyWithBiz(0)).Result()
	if err != nil {
		return nil, fmt.Errorf("get bizlist keys %s falied. err: %v", bizKey.listKeyWithBiz(0), err)
	}
	for _, key := range keys {
		instID, _, name, err := bizKey.parseListKeyValue(key)
		if err != nil {
			// invalid key, delete immediately
			if err := c.rds.SRem(bizKey.listKeyWithBiz(0), key).Err(); err != nil {
				blog.Errorf("delete invalid biz hash %s key: %s failed, err: %v", bizKey.listKeyWithBiz(0), key, err)
			}
			return nil, fmt.Errorf("got invalid key %s", key)
		}
		list = append(list, BizBaseInfo{
			BusinessID:   instID,
			BusinessName: name,
		})
	}

	return list, nil
}

func (c *Client) genBusinessListKeys(_ int64) ([]string, error) {
	bizList := make([]BizBaseInfo, 0)
	err := c.db.Table(common.BKTableNameBaseApp).Find(nil).Fields(common.BKAppIDField, common.BKAppNameField).All(context.Background(), &bizList)
	if err != nil {
		blog.Errorf("get all biz list from mongodb failed, err: %v", err)
		return nil, err
	}
	list := make([]string, len(bizList))
	for idx := range bizList {
		list[idx] = bizKey.genListKeyValue(bizList[idx].BusinessID, 0, bizList[idx].BusinessName)
	}
	return list, nil
}

// it's a background goroutine to check and if it's need to refresh the list cache.
// if it's needed then refresh the list.
func (c *Client) tryRefreshBaseList(bizID int64, ref refreshList) {
	if !c.lock.CanRefresh(ref.mainKey) {
		return
	}
	// set refreshing status
	c.lock.SetRefreshing(ref.mainKey)

	// now, we check whether we can refresh the list
	go func() {
		defer func() {
			c.lock.SetUnRefreshing(ref.mainKey)
		}()

		// get expire key
		expireStart, err := c.rds.Get(ref.expireKey).Result()
		if err != nil {
			if err != redis.Nil {
				blog.Errorf("try to refresh list %s cache, but get expire key failed, er: %v", ref.expireKey, err)
				return
			}
			expireStart = "0"
		}

		at, err := strconv.ParseInt(expireStart, 10, 64)
		if err != nil {
			blog.Errorf("try to refresh list %s cache, but get expire time failed, er: %v", ref.expireKey, err)
		}

		// if the expire key is not exist, or has already expired, then refresh the cache.
		if at > 0 && (time.Now().Unix()-at <= int64(ref.expireDuration/time.Second)) {
			// do not need refresh
			return
		}

		// get the lock key, so that it can only be done by one instance.
		success, err := c.rds.SetNX(ref.lockKey, 1, 15*time.Second).Result()
		if err != nil {
			blog.Errorf("sync %s list to refresh cache, but got redis lock failed, err: %v", ref.mainKey, err)
			return
		}

		if !success {
			blog.V(4).Infof("sync %s list to refresh cache, but do not get the redis lock, give up.", ref.mainKey)
			return
		}

		defer func() {
			if err := c.rds.Del(ref.lockKey).Err(); err != nil {
				blog.Errorf("sync %s list to refresh cache, but delete redis lock failed, err: %v", ref.mainKey, err)
			}
		}()

		// already get lock, now we start to refresh.

		wantKeys, err := ref.getList(bizID)
		if err != nil {
			blog.Errorf("sync list to refresh cache, but get list from mongodb failed, key: %s failed, err: %v", ref.mainKey, err)
		}

		// get real keys for compare later.
		realKeys, err := c.rds.SMembers(ref.mainKey).Result()
		if err != nil {
			blog.Errorf("sync list to refresh cache, but get list from redis key: %s failed, err: %v", ref.mainKey, err)
			return
		}

		pipeline := c.rds.Pipeline()
		defer pipeline.Close()
		realMap := make(map[string]bool)
		wantMap := make(map[string]bool)
		for _, key := range realKeys {
			realMap[key] = true
		}

		for _, k := range wantKeys {
			wantMap[k] = true
			if _, exist := realMap[k]; exist {
				continue
			}
			// not exist, means a new key is found, need to add
			pipeline.SAdd(ref.mainKey, k)
		}

		for _, k := range realKeys {
			if _, exist := wantMap[k]; exist {
				continue
			}
			// not exit, means we got a redundant key, need to remove
			pipeline.SRem(ref.mainKey, k)
		}

		// reset the expire time
		pipeline.Set(ref.expireKey, time.Now().Unix(), 0)
		_, err = pipeline.Exec()
		if err != nil {
			blog.Errorf("sync biz list to refresh cache %s, but exec pipeline failed failed, err: %v", ref.mainKey, err)
			return
		}

		blog.V(4).Infof("refresh list %s cache success.", ref.mainKey)

	}()
}

// tryRefreshInstanceDetail is a try to refresh the instance detail with the instance id given.
// it has a cache lock to avoid concurrent refresh try at local.
// and it has a redis lock to avoid refresh by multiple instance.
func (c *Client) tryRefreshInstanceDetail(instID int64, ref refreshInstance) {
	if !c.lock.CanRefresh(ref.mainKey) {
		return
	}
	// set refreshing status
	c.lock.SetRefreshing(ref.mainKey)

	go func() {
		defer func() {
			c.lock.SetUnRefreshing(ref.mainKey)
		}()

		// get expire key
		expireStart, err := c.rds.Get(ref.expireKey).Result()
		if err != nil {
			if err != redis.Nil {
				blog.Errorf("try to refresh instance %s cache, but get expire key failed, er: %v", ref.expireKey, err)
				return
			}
			expireStart = "0"
			return
		}

		at, err := strconv.ParseInt(expireStart, 10, 64)
		if err != nil {
			blog.Errorf("try to refresh instance %s cache, but get expire time failed, er: %v", ref.expireKey, err)
		}

		// if the expire key is not exist, or has already expired, then refresh the cache.
		if at > 0 && (time.Now().Unix()-at <= int64(ref.expireDuration/time.Second)) {
			// do not need refresh
			return
		}

		// get the lock key, so that it can only be done by one coreservice instance.
		success, err := c.rds.SetNX(ref.lockKey, 1, 15*time.Second).Result()
		if err != nil {
			blog.Errorf("sync %s instance to refresh cache, but got redis lock failed, err: %v", ref.mainKey, err)
			return
		}

		if !success {
			blog.V(4).Infof("sync %s instance to refresh cache, but do not get the redis lock, give up.", ref.mainKey)
			return
		}

		defer func() {
			if err := c.rds.Del(ref.lockKey).Err(); err != nil {
				blog.Errorf("sync %s instance to refresh cache, but delete redis lock failed, err: %v", ref.mainKey, err)
			}
		}()

		data, err := ref.getDetail(instID)
		if err != nil {
			blog.Errorf("refresh %s cache, but get instance data failed, err: %v", ref.mainKey, err)
			return
		}
		pipeline := c.rds.Pipeline()
		// refresh the data
		pipeline.Set(ref.mainKey, data, 0)
		// reset the expire time
		pipeline.Set(ref.expireKey, time.Now().Unix(), 0)
		_, err = pipeline.Exec()
		if err != nil {
			blog.Errorf("refresh cache %s, but exec pipeline failed failed, err: %v", ref.mainKey, err)
			return
		}

		blog.V(4).Infof("refresh %s cache success.", ref.mainKey)

	}()
}

func (c *Client) getModuleBaseList(bizID int64) ([]ModuleBaseInfo, error) {
	// get all keys which contains the biz id.
	keys, err := c.rds.SMembers(moduleKey.listKeyWithBiz(bizID)).Result()
	if err != nil {
		return nil, fmt.Errorf("get module keys %s falied. err: %v", moduleKey.listKeyWithBiz(bizID), err)
	}

	list := make([]ModuleBaseInfo, 0)
	for _, key := range keys {
		moduleID, setID, name, err := moduleKey.parseListKeyValue(key)
		if err != nil {
			// invalid key, delete immediately
			if c.rds.SRem(moduleKey.listKeyWithBiz(bizID), key).Err() != nil {
				blog.Errorf("delete invalid module %s key: %s failed,", moduleKey.listKeyWithBiz(bizID), key)
			}
			return nil, fmt.Errorf("got invalid key %s", key)
		}
		list = append(list, ModuleBaseInfo{
			ModuleID:   moduleID,
			ModuleName: name,
			SetID:      setID,
		})
	}
	return list, nil
}

func (c *Client) getSetBaseList(bizID int64) ([]SetBaseInfo, error) {
	// get all keys which contains the biz id.
	keys, err := c.rds.SMembers(setKey.listKeyWithBiz(bizID)).Result()
	if err != nil {
		return nil, fmt.Errorf("get set keys %s falied. err: %v", setKey.listKeyWithBiz(bizID), err)
	}

	list := make([]SetBaseInfo, 0)
	for _, key := range keys {
		setID, parentID, name, err := setKey.parseListKeyValue(key)
		if err != nil {
			// invalid key, delete immediately
			if c.rds.SRem(setKey.listKeyWithBiz(bizID), key).Err() != nil {
				blog.Errorf("delete invalid set %s key: %s failed,", setKey.listKeyWithBiz(bizID), key)
			}
			return nil, fmt.Errorf("got invalid key %s", key)
		}
		list = append(list, SetBaseInfo{
			SetID:    setID,
			SetName:  name,
			ParentID: parentID,
		})
	}
	return list, nil
}

const step = 1000

func (c *Client) getAllModuleBase(bizID int64) ([]ModuleBaseInfo, error) {
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	cnt, err := c.db.Table(common.BKTableNameBaseModule).Find(filter).Count(context.Background())
	if err != nil {
		return nil, err
	}

	list := make([]ModuleBaseInfo, 0)
	for start := 0; start < int(cnt); start += step {
		modules := make([]ModuleBaseInfo, 0)
		err = c.db.Table(common.BKTableNameBaseModule).Find(filter).Fields(common.BKModuleIDField, common.BKModuleNameField).
			Start(uint64(start)).Limit(uint64(step)).All(context.Background(), &modules)
		if err != nil {
			blog.Errorf("get biz %d module list from mongodb failed, err: %v", bizID, err)
			return nil, err
		}
		list = append(list, modules...)
	}

	return list, nil
}

func (c *Client) genModuleListKeys(bizID int64) ([]string, error) {
	moduleBaseList, err := c.getAllModuleBase(bizID)
	if err != nil {
		blog.Errorf("sync list to refresh cache, but get biz: %d modules failed, err: %v", bizID, err)
		return nil, err
	}

	keys := make([]string, len(moduleBaseList))
	for idx, mod := range moduleBaseList {
		keys[idx] = moduleKey.genListKeyValue(mod.ModuleID, mod.SetID, mod.ModuleName)
	}
	return keys, nil
}

func (c *Client) getAllSetBase(bizID int64) ([]SetBaseInfo, error) {
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	cnt, err := c.db.Table(common.BKTableNameBaseSet).Find(nil).Count(context.Background())
	if err != nil {
		return nil, err
	}

	list := make([]SetBaseInfo, 0)
	for start := 0; start < int(cnt); start += step {
		modules := make([]SetBaseInfo, 0)
		err = c.db.Table(common.BKTableNameBaseSet).Find(filter).Fields(common.BKSetIDField, common.BKSetNameField).
			Start(uint64(start)).Limit(uint64(step)).All(context.Background(), &modules)
		if err != nil {
			blog.Errorf("get biz %d set list from mongodb failed, err: %v", bizID, err)
			return nil, err
		}
		list = append(list, modules...)
	}

	return list, nil
}

func (c *Client) genSetListKeys(bizID int64) ([]string, error) {
	setList, err := c.getAllSetBase(bizID)
	if err != nil {
		blog.Errorf("sync list to refresh cache, but get biz: %d sets failed, err: %v", bizID, err)
		return nil, err
	}

	keys := make([]string, len(setList))
	for idx, set := range setList {
		keys[idx] = moduleKey.genListKeyValue(set.SetID, set.ParentID, set.SetName)
	}
	return keys, nil
}

func (c *Client) getModuleDetailFromMongo(id int64) (string, error) {
	mod := make(map[string]interface{})
	filter := mapstr.MapStr{
		common.BKModuleIDField: id,
	}

	if err := c.db.Table(common.BKTableNameBaseModule).Find(filter).One(context.Background(), &mod); err != nil {
		blog.Errorf("get module %d update from mongo failed, err: %v", id, err)
		return "", err
	}
	js, _ := json.Marshal(mod)
	return string(js), nil
}

func (c *Client) getSetDetailFromMongo(id int64) (string, error) {
	set := make(map[string]interface{})
	filter := mapstr.MapStr{
		common.BKSetIDField: id,
	}

	if err := c.db.Table(common.BKTableNameBaseSet).Find(filter).One(context.Background(), &set); err != nil {
		blog.Errorf("get module %d update from mongo failed, err: %v", id, err)
		return "", err
	}
	js, _ := json.Marshal(set)
	return string(js), nil
}

func (c *Client) genCustomLevelListKeys(objID string, bizID int64) ([]string, error) {
	list, err := c.getCustomLevelBaseList(objID, bizID)
	if err != nil {
		blog.Errorf("get custom level list keys failed, err: %v", err)
		return nil, err
	}
	keys := make([]string, len(list))
	for _, inst := range list {
		keys = append(keys, customKey.genListKeyValue(inst.InstanceID, inst.ParentID, inst.InstanceName))
	}
	return keys, nil
}

func (c *Client) getCustomLevelBaseFromMongodb(objID string, bizID int64) ([]CustomInstanceBase, error) {
	filter := mapstr.MapStr{
		common.BKObjIDField:  objID,
		common.MetadataField: meta.NewMetadata(bizID),
	}
	// count for paging use.
	cnt, err := c.db.Table(common.BKTableNameBaseInst).Find(filter).Count(context.Background())
	if err != nil {
		blog.Errorf("get custom level object: %s, biz: %d, list keys, but count from mongodb failed, err: %v", objID, bizID, err)
		return nil, err
	}
	list := make([]CustomInstanceBase, 0)
	for start := 0; start < int(cnt); start += step {
		instances := make([]CustomInstanceBase, 0)
		err = c.db.Table(common.BKTableNameBaseInst).Find(filter).
			Start(uint64(start)).Limit(uint64(step)).All(context.Background(), &instances)
		if err != nil {
			blog.Errorf("get custom level object: %s, biz: %d, list keys, but get from mongodb failed, err: %v", objID, bizID, err)
			return nil, err
		}
		list = append(list, instances...)
	}
	return list, nil
}

func (c *Client) getCustomLevelBaseList(objectID string, bizID int64) ([]CustomInstanceBase, error) {
	// get all keys which contains the biz id.
	keys, err := c.rds.SMembers(customKey.objListKeyWithBiz(objectID, bizID)).Result()
	if err != nil {
		return nil, fmt.Errorf("get custom level keys %s falied. err: %v", customKey.objListKeyWithBiz(objectID, bizID), err)
	}

	list := make([]CustomInstanceBase, 0)
	for _, key := range keys {
		instID, parentID, name, err := customKey.parseListKeyValue(key)
		if err != nil {
			// invalid key, delete immediately
			if c.rds.SRem(customKey.objListKeyWithBiz(objectID, bizID), key).Err() != nil {
				blog.Errorf("delete invalid custom level %s key: %s failed,", customKey.objListKeyWithBiz(objectID, bizID), key)
			}
			return nil, fmt.Errorf("got invalid key %s", key)
		}
		list = append(list, CustomInstanceBase{
			ObjectID:     objectID,
			InstanceID:   instID,
			InstanceName: name,
			ParentID:     parentID,
		})
	}
	return list, nil
}

func (c *Client) getCustomLevelDetail(objID string, instID int64) (string, error) {
	filter := mapstr.MapStr{
		common.BKObjIDField:  objID,
		common.BKInstIDField: instID,
	}
	instance := make(map[string]interface{})
	err := c.db.Table(common.BKTableNameBaseInst).Find(filter).One(context.Background(), &instance)
	if err != nil {
		blog.Errorf("get custom level object: %s, inst: %d from mongodb failed, err: %v", objID, instID, err)
		return "", err
	}
	js, err := json.Marshal(instance)
	if err != nil {
		return "", err
	}
	return string(js), nil
}
