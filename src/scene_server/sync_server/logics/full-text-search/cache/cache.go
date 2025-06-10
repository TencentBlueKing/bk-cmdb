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

// Package cache defines full-text search caching logics
package cache

import (
	"context"
	"sync"

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
)

// Cache is the full-text search cache
type Cache struct {
	cli cacheservice.Cache
	// resPoolBizIDMap is used to judge if biz/set... is in resource pool
	resPoolBizIDMap sync.Map
	// instEnumInfo is used to cache object enum options id to name
	instEnumInfo *instEnumIDToName
	// uuidObjIDMap is object uuid to objID map
	uuidObjIDMap sync.Map
}

type instEnumIDToName struct {
	// instEnumMap struct like: map[tenantID]map[obj]map[bk_property_id]map[option.id]option.name
	instEnumMap map[string]map[string]map[string]map[string]string
	rw          sync.RWMutex
}

// New creates a new cache
func New(cli cacheservice.Cache) (*Cache, error) {
	c := &Cache{
		cli:             cli,
		resPoolBizIDMap: sync.Map{},
		instEnumInfo: &instEnumIDToName{
			instEnumMap: make(map[string]map[string]map[string]map[string]string),
		},
		uuidObjIDMap: sync.Map{},
	}

	ctx := context.Background()

	// initialize resource pool biz info and instEnumMap
	err := tenant.ExecForAllTenants(func(tenantID string) error {
		_, err := c.getTenantResPoolBizInfo(ctx, tenantID)
		if err != nil {
			return err
		}

		if err = c.initTenantInstEnumInfo(ctx, tenantID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cache) getTenantResPoolBizInfo(ctx context.Context, tenantID string) (int64, error) {
	resPoolCond := mapstr.MapStr{common.BKDefaultField: common.DefaultAppFlag}

	biz := new(metadata.BizInst)
	err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameBaseApp).
		Find(resPoolCond).Fields(common.BKAppIDField).One(ctx, &biz)
	if err != nil {
		blog.Errorf("get resource pool biz for tenant: %s failed, err: %v", tenantID, err)
		if mongodb.IsNotFoundError(err) {
			return 0, nil
		}
		return 0, err
	}

	c.resPoolBizIDMap.Store(tenantID, biz.BizID)
	return biz.BizID, nil
}

// IsResourcePoolBiz check if biz id is resource pool biz id
func (c *Cache) IsResourcePoolBiz(kit *rest.Kit, bizID int64) bool {
	resPoolBizID, ok := c.resPoolBizIDMap.Load(kit.TenantID)
	if ok {
		return resPoolBizID == bizID
	}

	// get resource pool biz info for tenant from db
	var err error
	ferrors.FatalErrHandler(200, 100, func() error {
		resPoolBizID, err = c.getTenantResPoolBizInfo(kit.Ctx, kit.TenantID)
		return err
	})

	return resPoolBizID == bizID
}

func (c *Cache) initTenantInstEnumInfo(ctx context.Context, tenantID string) error {
	enumCond := mapstr.MapStr{
		common.BKPropertyTypeField: mapstr.MapStr{
			common.BKDBIN: []string{common.FieldTypeEnum, common.FieldTypeEnumMulti},
		},
	}

	attributes := make([]metadata.Attribute, 0)
	err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameObjAttDes).
		Find(enumCond).Fields(common.BKObjIDField, common.BKPropertyIDField, common.BKOptionField).All(ctx, &attributes)
	if err != nil {
		blog.Errorf("get %s object attributes enum info failed, err: %v", tenantID, err)
		return err
	}

	for _, attr := range attributes {
		tenantObjEnumMap, exists := c.instEnumInfo.instEnumMap[tenantID]
		if !exists {
			tenantObjEnumMap = make(map[string]map[string]map[string]string)
		}

		objEnumMap, exists := tenantObjEnumMap[attr.ObjectID]
		if !exists {
			objEnumMap = make(map[string]map[string]string)
		}

		option, err := metadata.ParseEnumOption(attr.Option)
		if err != nil {
			blog.Errorf("parse %v enum option failed, err: %v", attr, err)
			continue
		}

		enumMap := make(map[string]string)
		for _, opt := range option {
			enumMap[opt.ID] = opt.Name
		}
		objEnumMap[attr.PropertyID] = enumMap
		tenantObjEnumMap[attr.ObjectID] = objEnumMap
		c.instEnumInfo.instEnumMap[tenantID] = tenantObjEnumMap
	}

	return nil
}

// EnumIDToName change instance data enum id to enum name.
func (c *Cache) EnumIDToName(kit *rest.Kit, document mapstr.MapStr, objID string) mapstr.MapStr {
	c.instEnumInfo.rw.RLock()
	defer c.instEnumInfo.rw.RUnlock()

	tenantInstEnumInfo, ok := c.instEnumInfo.instEnumMap[kit.TenantID]
	if !ok {
		return document
	}

	objInstEnumInfo, ok := tenantInstEnumInfo[objID]
	if !ok {
		return document
	}

	for propertyId, enumInfo := range objInstEnumInfo {
		if _, ok := document[propertyId]; ok {
			if v, ok := document[propertyId].(string); ok {
				document[propertyId] = enumInfo[v]
			}
		}
	}

	return document
}

// SetObjEnumInfo set object enum cache info
func (c *Cache) SetObjEnumInfo(tenantID, objID string, attributes []mapstr.MapStr) {
	c.instEnumInfo.rw.Lock()
	defer c.instEnumInfo.rw.Unlock()

	if len(attributes) == 0 {
		if len(c.instEnumInfo.instEnumMap[tenantID]) > 0 {
			delete(c.instEnumInfo.instEnumMap[tenantID], objID)
		}
		return
	}

	tenantObjEnumMap, exists := c.instEnumInfo.instEnumMap[tenantID]
	if !exists {
		tenantObjEnumMap = make(map[string]map[string]map[string]string)
	}

	objEnumMap, exists := tenantObjEnumMap[objID]
	if !exists {
		objEnumMap = make(map[string]map[string]string)
	}

	for _, attr := range attributes {
		option, err := metadata.ParseEnumOption(attr[common.BKOptionField])
		if err != nil {
			blog.Errorf("parse %v enum option failed, err: %v", attr, err)
			continue
		}

		enumMap := make(map[string]string)
		for _, opt := range option {
			enumMap[opt.ID] = opt.Name
		}
		objEnumMap[util.GetStrByInterface(attr[common.BKPropertyIDField])] = enumMap
	}
	tenantObjEnumMap[objID] = objEnumMap
	c.instEnumInfo.instEnumMap[tenantID] = tenantObjEnumMap
}

// GetObjIDByUUID get object id by uuid
func (c *Cache) GetObjIDByUUID(kit *rest.Kit, uuid string) string {
	key := kit.TenantID + "_" + uuid
	objID, ok := c.uuidObjIDMap.Load(key)
	if ok {
		return util.GetStrByInterface(objID)
	}

	// get object id by uuid from db
	cond := mapstr.MapStr{metadata.ModelFieldObjUUID: uuid}
	var objectID string
	ferrors.FatalErrHandler(200, 100, func() error {
		object := new(metadata.Object)
		err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(cond).Fields(common.BKObjIDField).
			One(kit.Ctx, &object)
		if err != nil {
			if mongodb.IsNotFoundError(err) {
				return nil
			}
			blog.Errorf("get obj id failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
			return err
		}
		objectID = object.ObjectID
		return nil
	})

	c.uuidObjIDMap.Store(key, objectID)
	return objectID
}
