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

package cache

import (
	"context"
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/redis"

	"github.com/tidwall/gjson"
)

// ListDetailByIDs list general resource detail cache by ids
// if a resource does not exist in the cache and cannot be found in mongodb, it will not be returned
// therefore the length of the returned array may not be equal to the requested ids
func (c *Cache) ListDetailByIDs(kit *rest.Kit, opt *types.ListDetailByIDsOpt) ([]string, error) {
	if rawErr := opt.Validate(c.key.HasSubRes()); rawErr.ErrCode != 0 {
		blog.Errorf("list detail by ids option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return nil, rawErr.ToCCError(kit.CCError)
	}

	idDetailMap, err := c.listDetailByIDs(kit, opt)
	if err != nil {
		blog.Errorf("list detail cache by id keys(%+v) failed, err: %v, rid: %s", opt.IDKeys, err, kit.Rid)
		return nil, err
	}

	details := make([]string, 0)
	for _, id := range opt.IDKeys {
		detail, exists := idDetailMap[id]
		if exists {
			details = append(details, detail)
		}
	}

	return details, nil
}

// listDetailByIDs list general resource detail cache by ids, returns the id to detail map
func (c *Cache) listDetailByIDs(kit *rest.Kit, opt *types.ListDetailByIDsOpt) (map[string]string, error) {
	// get detail by ids from cache
	idKeys := util.StrArrayUnique(opt.IDKeys)
	detailKeys := make([]string, len(idKeys))
	for i, id := range idKeys {
		if opt.SubRes == "" {
			detailKeys[i] = c.key.DetailKey(id)
		} else {
			detailKeys[i] = c.key.DetailKey(id, opt.SubRes)
		}
	}

	results, err := redis.Client().MGet(kit.Ctx, detailKeys...).Result()
	if err != nil {
		blog.Errorf("list %s ids %+v detail cache failed, err: %v, rid: %s", c.key.Resource(), idKeys, err, kit.Rid)
		return nil, err
	}

	// generate id detail map and find the ids that need to be refreshed
	idDetailMap := make(map[string]string)
	needRefreshIDs, needRefreshKeys := make([]string, 0), make([]string, 0)
	for idx, res := range results {
		if res == nil {
			needRefreshIDs = append(needRefreshIDs, idKeys[idx])
			needRefreshKeys = append(needRefreshKeys, detailKeys[idx])
			continue
		}

		detail, ok := res.(string)
		if !ok {
			blog.Errorf("%s %s detail(%+v) is invalid, rid: %s", c.key.Resource(), idKeys[idx], res, kit.Rid)
			continue
		}

		if detail == "" {
			blog.Errorf("%s %s detail is empty, rid: %s", c.key.Resource(), idKeys[idx], kit.Rid)
			continue
		}

		if !opt.IsSystem && kit.SupplierAccount != common.BKSuperOwnerID {
			supplierAccount := gjson.Get(detail, common.BkSupplierAccount).String()
			if supplierAccount != common.BKDefaultOwnerID && supplierAccount != kit.SupplierAccount {
				continue
			}
		}

		if len(opt.Fields) != 0 {
			idDetailMap[idKeys[idx]] = *json.CutJsonDataWithFields(&detail, opt.Fields)
		} else {
			idDetailMap[idKeys[idx]] = detail
		}
	}

	if len(needRefreshIDs) == 0 {
		return idDetailMap, nil
	}

	// can not find detail in cache, need refresh the cache
	getDataOpt := &getDataByKeysOpt{
		BasicFilter: &types.BasicFilter{
			SubRes:          opt.SubRes,
			SupplierAccount: kit.SupplierAccount,
			IsSystem:        opt.IsSystem,
		},
		Keys: needRefreshIDs,
	}
	dbData, err := c.getDataByID(kit.Ctx, getDataOpt, kit.Rid)
	if err != nil {
		return nil, err
	}

	c.tryRefreshDetail(&tryRefreshDetailOpt{toRefreshKeys: needRefreshKeys, dbData: dbData, fields: opt.Fields,
		idDetailMap: idDetailMap}, kit.Rid)

	return idDetailMap, nil
}

type tryRefreshDetailOpt struct {
	toRefreshKeys []string
	dbData        []any
	fields        []string
	idDetailMap   map[string]string
	uniqueKeyType general.UniqueKeyType
	keyDetailMap  map[string]string
}

// tryRefreshDetail try refresh the general resource detail cache if it's not locked
func (c *Cache) tryRefreshDetail(opt *tryRefreshDetailOpt, rid string) {
	toRefreshKeyMap := make(map[string]struct{})
	for _, key := range opt.toRefreshKeys {
		toRefreshKeyMap[key] = struct{}{}
	}

	for _, data := range opt.dbData {
		// generate id detail map
		info, err := c.parseData(data)
		if err != nil {
			blog.Errorf("parse %s data: %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, rid)
			continue
		}

		idKey, _ := c.key.IDKey(info.id, info.oid)

		detailJs, err := json.Marshal(data)
		if err != nil {
			blog.Errorf("marshal %s mongo data %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, rid)
			continue
		}
		detailStr := string(detailJs)

		if len(opt.fields) != 0 {
			opt.idDetailMap[idKey] = *json.CutJsonDataWithFields(&detailStr, opt.fields)
		} else {
			opt.idDetailMap[idKey] = detailStr
		}

		lgc, exists := c.uniqueKeyLogics[opt.uniqueKeyType]
		if exists {
			redisKeys, _ := lgc.genKey(data, info)
			for _, redisKey := range redisKeys {
				opt.keyDetailMap[redisKey] = opt.idDetailMap[idKey]
				delete(toRefreshKeyMap, c.key.UniqueKey(string(opt.uniqueKeyType), redisKey))
			}
		}

		// refresh the general resource detail cache when we had the lock
		detailKey := c.key.DetailKey(idKey, info.subRes...)
		delete(toRefreshKeyMap, detailKey)
		if !c.refreshingLock.CanRefresh(detailKey) {
			continue
		}
		c.refreshingLock.SetRefreshing(detailKey)

		go func(data any) {
			defer c.refreshingLock.SetUnRefreshing(detailKey)

			pipeline := redis.Client().Pipeline()
			ttl := c.withRandomExpireSeconds(c.expireSeconds)

			// upsert all related unique key cache
			for typ, lgc := range c.uniqueKeyLogics {
				redisKeys, err := lgc.genKey(data, info)
				if err != nil {
					blog.Errorf("generate %s %s key failed, err: %v, data: %+v, rid: %s", c.key.Resource(), typ, err,
						data, rid)
					continue
				}

				for _, redisKey := range redisKeys {
					pipeline.SetNX(c.key.UniqueKey(string(typ), redisKey), idKey, ttl)
				}
			}

			// upsert general resource detail cache
			pipeline.SetNX(detailKey, detailStr, c.key.WithRandomExpireSeconds())

			_, err = pipeline.Exec()
			if err != nil {
				blog.Errorf("refresh %s cache failed, err: %v, data: %s, rid: %s", c.key.Resource(), idKey, err,
					detailStr, rid)
				return
			}

			blog.V(4).Infof("refresh %s cache success, id: %s, rid: %s", c.key.Resource(), idKey, rid)
		}(data)
	}

	go c.handleNotExistKey(toRefreshKeyMap, rid)
}

// handleNotExistKey set not exist refresh key cache to empty string to avoid cache penetration
func (c *Cache) handleNotExistKey(notExistKeyMap map[string]struct{}, rid string) error {
	if len(notExistKeyMap) == 0 {
		return nil
	}

	pipeline := redis.Client().Pipeline()
	for notExistKey := range notExistKeyMap {
		pipeline.SetNX(notExistKey, "", c.key.WithRandomExpireSeconds())
	}

	if _, err := pipeline.Exec(); err != nil {
		blog.Errorf("refresh not exist %s cache failed, err: %v, key info: %+v, rid: %s", c.key.Resource(), err,
			notExistKeyMap, rid)
		return err
	}

	blog.V(4).Infof("refresh not exist %s cache success, key info: %+v, rid: %s", c.key.Resource(), notExistKeyMap, rid)
	return nil
}

// ListDetailByUniqueKey list general resource detail cache by unique keys
// if a resource does not exist in the cache and cannot be found in mongodb, it will not be returned
// therefore the length of the returned array may not be equal to the requested keys
func (c *Cache) ListDetailByUniqueKey(kit *rest.Kit, opt *types.ListDetailByUniqueKeyOpt) ([]string, error) {
	if rawErr := opt.Validate(c.key.HasSubRes()); rawErr.ErrCode != 0 {
		blog.Errorf("list detail by ids option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return nil, rawErr.ToCCError(kit.CCError)
	}

	keyDetailMap, err := c.listDetailByUniqueKey(kit, opt)
	if err != nil {
		blog.Errorf("list detail cache by unique keys(%+v) failed, err: %v, rid: %s", opt.Keys, err, kit.Rid)
		return nil, err
	}

	details := make([]string, 0)
	for _, key := range opt.Keys {
		detail, exists := keyDetailMap[key]
		if exists {
			details = append(details, detail)
		}
	}

	return details, nil
}

// listDetailByUniqueKey list general resource detail cache by unique keys, returns the unique key to detail map
func (c *Cache) listDetailByUniqueKey(kit *rest.Kit, opt *types.ListDetailByUniqueKeyOpt) (map[string]string,
	error) {

	uniqueKeyLgc, exists := c.uniqueKeyLogics[opt.Type]
	if !exists {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "type")
	}

	// get id keys by unique keys from cache
	keys := util.StrArrayUnique(opt.Keys)
	uniqueKeys := make([]string, len(keys))
	for i, key := range keys {
		uniqueKeys[i] = c.key.UniqueKey(string(opt.Type), key)
	}

	results, err := redis.Client().MGet(kit.Ctx, uniqueKeys...).Result()
	if err != nil {
		blog.Errorf("list %s unique keys %+v cache failed, err: %v, rid: %s", c.key.Resource(), keys, err, kit.Rid)
		return nil, err
	}

	// generate unique key to id map and find the unique keys that need to be refreshed
	keyIDKeyMap := make(map[string]string)
	idKeys := make([]string, 0)
	needRefreshKeys, needRefreshRedisKeys := make([]string, 0), make([]string, 0)
	for idx, res := range results {
		if res == nil {
			needRefreshKeys = append(needRefreshKeys, keys[idx])
			needRefreshRedisKeys = append(needRefreshRedisKeys, uniqueKeys[idx])
			continue
		}

		idKey, ok := res.(string)
		if !ok {
			blog.Errorf("%s unique key %s id key(%+v) is invalid, rid: %s", c.key.Resource(), keys[idx], res, kit.Rid)
			continue
		}

		if idKey == "" {
			blog.Errorf("%s unique key %s id key is empty, rid: %s", c.key.Resource(), keys[idx], kit.Rid)
			continue
		}

		keyIDKeyMap[keys[idx]] = idKey
		idKeys = append(idKeys, idKey)
	}

	// list detail by ids from redis
	idDetailMap := make(map[string]string)
	keyDetailMap := make(map[string]string)
	if len(idKeys) > 0 {
		listByIDOpt := &types.ListDetailByIDsOpt{SubRes: opt.SubRes, IsSystem: opt.IsSystem, IDKeys: idKeys,
			Fields: opt.Fields}
		idDetailMap, err = c.listDetailByIDs(kit, listByIDOpt)
		if err != nil {
			blog.Errorf("list detail by ids(%+v) failed, err: %v, rid: %s", listByIDOpt, err, kit.Rid)
			return nil, err
		}
	}

	for key, idKey := range keyIDKeyMap {
		keyDetailMap[key] = idDetailMap[idKey]
	}

	if len(needRefreshKeys) == 0 {
		return keyDetailMap, nil
	}

	// can not find detail in cache, need refresh the cache
	getDataOpt := &getDataByKeysOpt{
		BasicFilter: &types.BasicFilter{SubRes: opt.SubRes, SupplierAccount: kit.SupplierAccount,
			IsSystem: opt.IsSystem},
		Keys: needRefreshKeys,
	}
	dbData, err := uniqueKeyLgc.getData(kit.Ctx, getDataOpt, kit.Rid)
	if err != nil {
		return nil, err
	}

	c.tryRefreshDetail(&tryRefreshDetailOpt{toRefreshKeys: needRefreshRedisKeys, dbData: dbData, fields: opt.Fields,
		idDetailMap: idDetailMap, uniqueKeyType: opt.Type, keyDetailMap: keyDetailMap}, kit.Rid)

	return keyDetailMap, nil
}

// ListDetail list general resource detail using id list
func (c *Cache) ListDetail(kit *rest.Kit, opt *types.ListDetailOpt) ([]string, error) {
	if rawErr := opt.Validate(c.key.HasSubRes()); rawErr.ErrCode != 0 {
		blog.Errorf("list detail option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return nil, rawErr.ToCCError(kit.CCError)
	}

	idListTTL, err := c.validateIDList(opt.IDListFilter)
	if err != nil {
		blog.Errorf("id list filter option is invalid, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	refreshOpt := &refreshIDListOpt{
		filterOpt: opt.IDListFilter,
		ttl:       idListTTL,
	}
	notExists, _, err := c.tryRefreshIDListIfNeeded(kit.Ctx, refreshOpt, kit.Rid)
	if err != nil {
		blog.Errorf("try refresh id list failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	// id list not exists, get detail from db
	if notExists {
		dbRes, err := c.listDataFromDB(kit.Ctx, opt, kit.Rid)
		if err != nil {
			return nil, err
		}

		details := make([]string, len(dbRes.Data))
		for i, data := range dbRes.Data {
			detailJs, err := json.Marshal(data)
			if err != nil {
				blog.Errorf("marshal %s mongo data %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, kit.Rid)
				continue
			}
			details[i] = string(detailJs)
		}
		return details, nil
	}

	// id list exists, get id list and detail from redis
	idKeys, err := c.listIDsFromRedis(kit.Ctx, opt.IDListFilter.IDListKey, opt.Page, kit.Rid)
	if err != nil {
		return nil, err
	}

	if len(idKeys) == 0 {
		return make([]string, 0), nil
	}

	listByIDsOpt := &types.ListDetailByIDsOpt{
		SubRes:   opt.IDListFilter.SubRes,
		IsSystem: opt.IDListFilter.IsSystem,
		IDKeys:   idKeys,
		Fields:   opt.Fields,
	}
	return c.ListDetailByIDs(kit, listByIDsOpt)
}

// listDataFromDB list detail from db
func (c *Cache) listDataFromDB(ctx context.Context, opt *types.ListDetailOpt, rid string) (*listDataRes, error) {
	if rawErr := opt.Validate(c.key.HasSubRes()); rawErr.ErrCode != 0 {
		blog.Errorf("list %s detail option is invalid, err: %v, opt: %+v, rid: %s", c.key.Resource(), rawErr, opt, rid)
		return nil, fmt.Errorf("list detail option is invalid")
	}

	listOpt := &listDataOpt{
		BasicFilter: opt.IDListFilter.BasicFilter,
		Fields:      opt.Fields,
		OnlyListID:  opt.OnlyListID,
		Page:        opt.Page,
	}

	var err error
	if !opt.IDListFilter.IsAll {
		listOpt.Cond, err = opt.IDListFilter.Cond.ToMgo()
		if err != nil {
			blog.Errorf("parse list %s detail cond(%s) failed, err: %v, rid: %s", c.key.Resource(),
				opt.IDListFilter.Cond, err, rid)
			return nil, err
		}
	}

	return c.listData(ctx, listOpt, rid)
}

// RefreshDetailByIDs refresh general resource detail cache by ids
func (c *Cache) RefreshDetailByIDs(kit *rest.Kit, opt *types.RefreshDetailByIDsOpt) error {
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("refresh detail by ids option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return rawErr.ToCCError(kit.CCError)
	}

	getDataOpt := &getDataByKeysOpt{
		BasicFilter: &types.BasicFilter{
			SubRes:          opt.SubResource,
			SupplierAccount: kit.SupplierAccount,
			IsSystem:        true,
		},
		Keys: opt.IDKeys,
	}
	dbData, err := c.getDataByID(kit.Ctx, getDataOpt, kit.Rid)
	if err != nil {
		return err
	}

	c.tryRefreshDetail(&tryRefreshDetailOpt{dbData: dbData, idDetailMap: make(map[string]string)}, kit.Rid)
	return nil
}

// CountData count general resource data
func (c *Cache) CountData(kit *rest.Kit, opt *types.ListDetailOpt) (int64, error) {
	if rawErr := opt.Validate(c.key.HasSubRes()); rawErr.ErrCode != 0 {
		blog.Errorf("list ids from start option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return 0, rawErr.ToCCError(kit.CCError)
	}

	_, err := c.validateIDList(opt.IDListFilter)
	if err != nil {
		blog.Errorf("id list filter option is invalid, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return 0, err
	}

	exists, err := isIDListExists(kit.Ctx, opt.IDListFilter.IDListKey, kit.Rid)
	if err != nil {
		return 0, err
	}

	// id list not exists, get data count from db
	if !exists {
		dbRes, err := c.listDataFromDB(kit.Ctx, opt, kit.Rid)
		if err != nil {
			return 0, err
		}
		return int64(dbRes.Count), nil
	}

	// id list exists, get id list count from redis
	cnt, err := c.countIDsFromRedis(kit.Ctx, opt.IDListFilter.IDListKey, kit.Rid)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
