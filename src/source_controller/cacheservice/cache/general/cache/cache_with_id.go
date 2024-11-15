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
	"strconv"
	"time"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/mongodb"
)

func init() {
	addCache(newMapStrCacheWithID(general.BizKey, false, common.BKTableNameBaseApp, common.BKAppIDField))
	addCache(newMapStrCacheWithID(general.SetKey, false, common.BKTableNameBaseSet, common.BKSetIDField))
	addCache(newMapStrCacheWithID(general.ModuleKey, false, common.BKTableNameBaseModule, common.BKModuleIDField))
	addCache(newMapStrCacheWithID(general.BizSetKey, false, common.BKTableNameBaseBizSet, common.BKBizSetIDField))
	addCache(newMapStrCacheWithID(general.PlatKey, false, common.BKTableNameBasePlat, common.BKCloudIDField))
}

// newMapStrCacheWithID new general cache whose data is of mapstr type and uses id as id key
func newMapStrCacheWithID(key *general.Key, needCacheAll bool, table, idField string) *Cache {
	return newCacheWithID[mapstr.MapStr](key, needCacheAll, table, idField,
		func(data mapstr.MapStr, idField string) (*basicInfo, error) {
			id, err := util.GetInt64ByInterface(data[idField])
			if err != nil {
				return nil, fmt.Errorf("parse id %+v failed, err: %v", data[idField], err)
			}
			return &basicInfo{
				id:     id,
				tenant: util.GetStrByInterface(data[common.TenantID]),
			}, nil
		})
}

// newCacheWithID new general cache whose data uses id as id key
func newCacheWithID[T any](key *general.Key, needCacheAll bool, table, idField string,
	parser func(data T, idField string) (*basicInfo, error)) *Cache {

	cache := NewCache()
	cache.key = key
	cache.expireSeconds = 30 * 60 * time.Second
	cache.expireRangeSeconds = [2]int{-600, 600}
	cache.needCacheAll = needCacheAll
	cache.parseData = parseDataWithID[T](idField, parser)
	cache.getDataByID = getDataByID[T](table, idField)
	cache.listData = listDataWithID[T](table, idField)
	return cache
}

// parseDataWithID returns the dataParser for resource that uses id as id key
func parseDataWithID[T any](idField string, parser func(data T, idField string) (*basicInfo, error)) dataParser {
	return func(data any) (*basicInfo, error) {
		var info *basicInfo
		var err error

		switch val := data.(type) {
		case T:
			// parse db data
			info, err = parser(val, idField)
		case types.WatchEventData:
			// parse event watch data
			info, err = parseWatchChainNode(val.ChainNode)
		default:
			return nil, fmt.Errorf("data type %T is invalid", data)
		}

		if err != nil {
			return nil, err
		}

		if info.id == 0 {
			return nil, fmt.Errorf("id is zero")
		}

		return info, nil
	}
}

func getDataByID[T any](table, idField string) dataGetterByKeys {
	return func(ctx context.Context, opt *getDataByKeysOpt, rid string) ([]any, error) {
		ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

		dataArr, err := getDBDataByID[T](ctx, opt, table, idField, rid)
		if err != nil {
			return nil, err
		}

		return convertToAnyArr(dataArr), nil
	}
}

func getDBDataByID[T any](ctx context.Context, opt *getDataByKeysOpt, table, idField string, rid string) ([]T, error) {
	if len(opt.Keys) == 0 {
		return make([]T, 0), nil
	}

	ids := make([]int64, len(opt.Keys))

	for i, key := range opt.Keys {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			blog.Errorf("parse id (index: %d, key: %s) failed, err: %v, rid: %s", i, key, err, rid)
			return nil, err
		}
		ids[i] = id
	}

	cond := mapstr.MapStr{
		idField: mapstr.MapStr{common.BKDBIN: ids},
	}

	dataArr := make([]T, 0)
	if err := mongodb.Client().Table(table).Find(cond).All(ctx, &dataArr); err != nil {
		blog.Errorf("get %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, rid)
		return nil, err
	}
	return dataArr, nil
}

func listDataWithID[T any](table, idField string) dataLister {
	return func(ctx context.Context, opt *listDataOpt, rid string) (*listDataRes, error) {
		ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

		if rawErr := opt.Validate(false); rawErr.ErrCode != 0 {
			blog.Errorf("list general data option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, rid)
			return nil, fmt.Errorf("list data option is invalid")
		}

		if opt.OnlyListID {
			opt.Fields = []string{idField}
		}

		cnt, dataArr, err := listDBDataWithID[T](ctx, opt, table, idField, rid)
		if err != nil {
			return nil, err
		}

		return &listDataRes{Count: cnt, Data: convertToAnyArr(dataArr)}, nil
	}
}

func listDBDataWithID[T any](ctx context.Context, opt *listDataOpt, table, idField string, rid string) (uint64, []T,
	error) {

	cond := opt.Cond
	if cond == nil {
		cond = make(mapstr.MapStr)
	}

	if opt.Page.EnableCount {
		cnt, err := mongodb.Client().Table(table).Find(cond).Count(ctx)
		if err != nil {
			blog.Errorf("count %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, rid)
			return 0, nil, err
		}

		return cnt, make([]T, 0), nil
	}

	if opt.Page.StartID != 0 {
		_, exists := cond[idField]
		if exists {
			cond = mapstr.MapStr{
				common.BKDBAND: []mapstr.MapStr{
					{idField: mapstr.MapStr{common.BKDBGT: opt.Page.StartID}}, cond,
				},
			}
		} else {
			cond[idField] = mapstr.MapStr{common.BKDBGT: opt.Page.StartID}
		}
	}

	if len(opt.Fields) > 0 {
		opt.Fields = append(opt.Fields, idField)
	}

	dataArr := make([]T, 0)
	err := mongodb.Client().Table(table).Find(cond).Sort(idField).Start(uint64(opt.Page.StartIndex)).
		Limit(uint64(opt.Page.Limit)).Fields(opt.Fields...).All(ctx, &dataArr)
	if err != nil {
		blog.Errorf("list %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, rid)
		return 0, nil, err
	}

	return 0, dataArr, nil
}
