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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/coreservice/cache/tools"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	"gopkg.in/redis.v5"
)

type Client struct {
	rds   *redis.Client
	event reflector.Interface
	db    dal.DB
	lock  tools.RefreshingLock
}

func (c *Client) GetBizBaseList() ([]BizBaseInfo, error) {
	c.tryRefreshBaseList(0, refreshList{
		mainKey:        bizKey.listKeyWithBiz(0),
		lockKey:        bizKey.listLockKeyWithBiz(0),
		expireKey:      bizKey.listExpireKeyWithBiz(0),
		expireDuration: bizKey.listExpireDuration,
		getList:        c.genBusinessListKeys,
	})

	// get the business base info
	base, err := c.getBusinessBaseInfo()
	if err == nil {
		return base, nil
	}

	blog.Errorf("get biz base list from cache failed, will get from mongodb, err: %v", err)
	// get from db directly.
	list := make([]BizBaseInfo, 0)
	err = c.db.Table(common.BKTableNameBaseApp).Find(nil).Fields(common.BKAppIDField, common.BKAppNameField).All(context.Background(), &list)
	if err != nil {
		blog.Errorf("sync biz list to refresh cache, but get biz list from mongodb failed, err: %v", err)
		return nil, err
	}
	return list, err

}

// get a business's all info.
func (c *Client) GetBusiness(bizID int64) (string, error) {
	key := bizKey.detailKey(bizID)
	exist, err := c.rds.Exists(key).Result()
	if err != nil {
		blog.Warnf("get business info from cache,  biz: %d, but check exist failed, err: %v", bizID, err)
		// get from db directly.
		exist = false
	}

	// try to refresh cache.
	c.tryRefreshInstanceDetail(bizID, refreshInstance{
		mainKey:        bizKey.detailKey(bizID),
		lockKey:        bizKey.detailLockKey(bizID),
		expireKey:      bizKey.detailExpireKey(bizID),
		expireDuration: bizKey.detailExpireDuration,
		getDetail:      c.getModuleDetailFromMongo,
	})

	if exist {
		biz, err := c.rds.Get(key).Result()
		if err == nil {
			return biz, nil
		}
		// error occurs, get from db directly.
		// Note: this may cause high db query
		blog.Errorf("get business %d info from cache, but failed, will get from mongodb, err: %v", bizID, err)
	}

	// get from db
	return c.getBusinessFromMongo(bizID)
}

func (c *Client) GetModuleBaseList(bizID int64) ([]ModuleBaseInfo, error) {
	c.tryRefreshBaseList(bizID, refreshList{
		mainKey:        moduleKey.listKeyWithBiz(bizID),
		lockKey:        moduleKey.listLockKeyWithBiz(bizID),
		expireKey:      moduleKey.listExpireKeyWithBiz(bizID),
		expireDuration: moduleKey.listExpireDuration,
		getList:        c.genModuleListKeys,
	})
	base, err := c.getModuleBaseList(bizID)
	if err == nil {
		return base, nil
	}

	blog.Errorf("get biz %d module base list from cache failed, get from db now, err: %v", bizID, err)
	// do not return
	// get from db directly.
	list := make([]ModuleBaseInfo, 0)
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	err = c.db.Table(common.BKTableNameBaseModule).Find(filter).All(context.Background(), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Client) GetModule(moduleID int64) (string, error) {
	// try refresh the module cache
	c.tryRefreshInstanceDetail(moduleID, refreshInstance{
		mainKey:        moduleKey.detailKey(moduleID),
		lockKey:        moduleKey.detailLockKey(moduleID),
		expireKey:      moduleKey.detailExpireKey(moduleID),
		expireDuration: moduleKey.detailExpireDuration,
		getDetail:      c.getModuleDetailFromMongo,
	})

	mod, err := c.rds.Get(moduleKey.detailKey(moduleID)).Result()
	if err == nil && len(mod) != 0 {
		return mod, nil
	}
	blog.Errorf("get module: %d failed from redis, err: %v", err)
	// get from db directly.
	return c.getModuleDetailFromMongo(moduleID)
}

func (c *Client) GetSetBaseList(bizID int64) ([]SetBaseInfo, error) {
	c.tryRefreshBaseList(bizID, refreshList{
		mainKey:        setKey.listKeyWithBiz(bizID),
		lockKey:        setKey.listLockKeyWithBiz(bizID),
		expireKey:      setKey.listExpireKeyWithBiz(bizID),
		expireDuration: setKey.listExpireDuration,
		getList:        c.genSetListKeys,
	})
	base, err := c.getSetBaseList(bizID)
	if err == nil {
		return base, nil
	}

	blog.Errorf("get biz %d set base list from cache failed, get from db now, err: %v", bizID, err)
	// do not return
	// get from db directly.
	list := make([]SetBaseInfo, 0)
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	err = c.db.Table(common.BKTableNameBaseSet).Find(filter).All(context.Background(), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Client) GetSet(setID int64) (string, error) {
	// try refresh the module cache
	c.tryRefreshInstanceDetail(setID, refreshInstance{
		mainKey:        setKey.detailKey(setID),
		lockKey:        setKey.detailLockKey(setID),
		expireKey:      setKey.detailExpireKey(setID),
		expireDuration: setKey.detailExpireDuration,
		getDetail:      c.getSetDetailFromMongo,
	})

	set, err := c.rds.Get(moduleKey.detailKey(setID)).Result()
	if err == nil && len(set) != 0 {
		return set, nil
	}
	blog.Errorf("get set: %d failed from redis, err: %v", err)
	// get from db directly.
	return c.getSetDetailFromMongo(setID)
}

func (c *Client) GetCustomLevelBaseList(objectID string, bizID int64) ([]CustomInstanceBase, error) {
	c.tryRefreshBaseList(bizID, refreshList{
		mainKey:        customKey.objListKeyWithBiz(objectID, bizID),
		lockKey:        customKey.objListLockKeyWithBiz(objectID, bizID),
		expireKey:      customKey.objListExpireKeyWithBiz(objectID, bizID),
		expireDuration: customKey.listExpireDuration,
		getList: func(bizID int64) (strings []string, err error) {
			return c.genCustomLevelListKeys(objectID, bizID)
		},
	})

	list, err := c.getCustomLevelBaseList(objectID, bizID)
	if err == nil {
		return list, nil
	}
	blog.Errorf("get biz: %s, obj: %s custom level list keys from cache failed, will get from mongodb, err: %v",
		bizID, objectID, err)

	return c.getCustomLevelBaseFromMongodb(objectID, bizID)
}

func (c *Client) GetCustomLevelDetail(objID string, instID int64) (string, error) {
	c.tryRefreshInstanceDetail(instID, refreshInstance{
		mainKey:        customKey.detailKey(objID, instID),
		lockKey:        customKey.detailLockKey(objID, instID),
		expireKey:      customKey.detailExpireKey(objID, instID),
		expireDuration: customKey.detailExpireDuration,
		getDetail: func(instID int64) (s string, err error) {
			return c.getCustomLevelDetail(objID, instID)
		},
	})

	set, err := c.rds.Get(customKey.detailKey(objID, instID)).Result()
	if err == nil && len(set) != 0 {
		return set, nil
	}
	blog.Errorf("get biz custom level, obj:%s, inst: %d failed from redis, err: %v", objID, instID, err)
	// get from db directly.
	return c.getCustomLevelDetail(objID, instID)
}

func (c *Client) GetTopology() ([]string, error) {
	// TODO: try refresh the cache.
	key := customKey.topologyKey()
	rank, err := c.rds.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return strings.Split(rank, ","), nil
}
