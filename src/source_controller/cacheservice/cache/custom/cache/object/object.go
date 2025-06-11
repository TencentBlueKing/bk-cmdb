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

// Package object is the object cache
package object

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache/logics"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

var (
	objUUIDCache = logics.NewStrCache(logics.NewKey(types.ObjUUIDType, 6*time.Hour))
)

// ObjectCache is object cache
type ObjectCache struct {
	isMaster     discovery.ServiceManageInterface
	objUUIDCache *logics.StrCache
}

// NewObjectCache new object cache
func NewObjectCache(isMaster discovery.ServiceManageInterface) *ObjectCache {
	return &ObjectCache{
		isMaster:     isMaster,
		objUUIDCache: objUUIDCache,
	}
}

// GetUUIDByObj get object uuid by objID
func (c *ObjectCache) GetUUIDByObj(kit *rest.Kit, objID string) (string, error) {
	return GetUUIDByObj(kit, objID)
}

// GetUUIDByObj get object uuid by objID
func GetUUIDByObj(kit *rest.Kit, objID string) (string, error) {
	// get object uuid from cache
	objUUIDMap, err := objUUIDCache.List(kit, []string{objID})
	if err != nil {
		return "", err
	}

	uuid, exists := objUUIDMap[objID]
	if exists {
		return uuid, nil
	}

	// get object uuid from db and refresh cache
	cond := mapstr.MapStr{common.BKObjIDField: objID}
	object := new(metadata.Object)
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(cond).Fields(metadata.ModelFieldObjUUID).
		One(kit.Ctx, &object)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("get %s object %s failed, err: %v, rid: %v", kit.TenantID, objID, err, kit.Rid)
		return "", err
	}

	if err = objUUIDCache.BatchUpdate(kit, map[string]interface{}{objID: object.UUID}); err != nil {
		blog.Errorf("update %s object %s uuid %s cache failed, err: %v, data: %+v, rid: %s", kit.TenantID, objID,
			object.UUID, err, kit.Rid)
	}

	return object.UUID, nil
}

// UpdateObjUUIDCache update objID to uuid cache by map[objID]uuid
func (c *ObjectCache) UpdateObjUUIDCache(kit *rest.Kit, objUUIDMap map[string]string) error {
	redisDataMap := make(map[string]interface{})
	for objID, uuid := range objUUIDMap {
		redisDataMap[objID] = uuid
	}

	if err := c.objUUIDCache.BatchUpdate(kit, redisDataMap); err != nil {
		blog.Errorf("update objID to uuid cache failed, err: %v, data: %+v, rid: %s", err, objUUIDMap, kit.Rid)
		return err
	}
	return nil
}

// DeleteObjUUIDCache delete objID to uuid cache by objIDs
func (c *ObjectCache) DeleteObjUUIDCache(kit *rest.Kit, objIDs []string) error {
	if err := c.objUUIDCache.BatchDelete(kit, objIDs); err != nil {
		blog.Errorf("delete objID to uuid cache failed, err: %v, keys: %+v, rid: %s", err, objIDs, kit.Rid)
		return err
	}
	return nil
}

// RefreshCache refresh object cache
func (c *ObjectCache) RefreshCache(rid string) error {
	// lock refresh object cache operation, returns error if it is already locked
	lockKey := fmt.Sprintf("%s:shared_ns_rel_refresh:lock", logics.Namespace)

	locker := lock.NewLocker(redis.Client())
	locked, err := locker.Lock(lock.StrFormat(lockKey), 10*time.Minute)
	defer locker.Unlock()
	if err != nil {
		blog.Errorf("get %s lock failed, err: %v, rid: %s", lockKey, err, rid)
		return err
	}

	if !locked {
		blog.Errorf("%s task is already lock, rid: %s", lockKey, rid)
		return errors.New("there's a same refreshing task running, please retry later")
	}

	kit := rest.NewKit().WithCtx(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)).
		WithRid(rid)
	err = tenant.ExecForAllTenants(func(tenantID string) error {
		kit = kit.WithTenant(tenantID)

		objects := make([]metadata.Object, 0)
		err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
			metadata.ModelFieldObjUUID).All(kit.Ctx, &objects)
		if err != nil {
			blog.Errorf("list %s objects failed, err: %v, rid: %v", kit.TenantID, err, kit.Rid)
			return err
		}

		objUUIDMap := make(map[string]interface{})
		for _, obj := range objects {
			objUUIDMap[obj.ObjectID] = obj.UUID
		}

		// refresh label key and value count cache
		err = c.objUUIDCache.Refresh(kit, "*", objUUIDMap)
		if err != nil {
			blog.Errorf("refresh objID to uuid cache failed, err: %v, data: %+v, rid: %s", err, objUUIDMap, kit.Rid)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// LoopRefreshCache loop refresh object key and value cache every day at 3am
func (c *ObjectCache) LoopRefreshCache() {
	for {
		time.Sleep(2 * time.Hour)

		if !c.isMaster.IsMaster() {
			blog.V(4).Infof("loop refresh object cache, but not master, skip.")
			time.Sleep(time.Minute)
			continue
		}

		rid := util.GenerateRID()

		blog.Infof("start refresh object cache task, rid: %s", rid)
		err := c.RefreshCache(rid)
		if err != nil {
			blog.Errorf("refresh object cache failed, err: %v, rid: %s", err, rid)
			continue
		}
		blog.Infof("refresh object cache successfully, rid: %s", rid)
	}
}

// GetInstTableNameByObjID get object instance table name by objID
func GetInstTableNameByObjID(kit *rest.Kit, objID string) (string, error) {
	if common.IsInnerModel(objID) {
		return common.GetInnerInstTableName(objID), nil
	}
	uuid, err := GetUUIDByObj(kit, objID)
	if err != nil {
		return "", err
	}

	return common.GetObjInstTableName(uuid), nil
}
