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
	"fmt"
	"time"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
)

// newCacheWithID new general cache whose data uses id as id key
func newCacheWithIDAndSubRes[T any](key *general.Key, idField string, subResFields []string,
	getTable func(kit *rest.Kit, filter *types.BasicFilter) (string, error),
	parseData func(data dataWithTenant[T]) (*basicInfo, error)) *Cache {

	cache := NewCache()
	cache.key = key
	cache.expireSeconds = 30 * 60 * time.Second
	cache.expireRangeSeconds = [2]int{-600, 600}
	cache.needCacheAll = false
	cache.parseData = parseDataWithIDAndSubRes[T](parseData)
	cache.getDataByID = getDataByIDAndSubRes[T](idField, getTable)
	cache.listData = listDataWithIDAndSubRes[T](idField, subResFields, getTable)
	return cache
}

type dataWithTenant[T any] struct {
	TenantID string `json:"-" bson:"-"`
	Data     T      `json:",inline" bson:",inline"`
}

// MarshalJSON marshal json
func (data dataWithTenant[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.Data)
}

func parseDataWithIDAndSubRes[T any](parser func(data dataWithTenant[T]) (*basicInfo, error)) dataParser {
	return func(data any) (*basicInfo, error) {
		var info *basicInfo
		var err error

		switch val := data.(type) {
		case dataWithTenant[T]:
			info, err = parser(val)
		case types.WatchEventData:
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

		if len(info.subRes) == 0 {
			return nil, fmt.Errorf("sub resource is empty")
		}

		return info, nil
	}
}

func getDataByIDAndSubRes[T any](idField string,
	getTable func(kit *rest.Kit, filter *types.BasicFilter) (string, error)) dataGetterByKeys {

	return func(kit *rest.Kit, opt *getDataByKeysOpt) ([]any, error) {
		table, err := getTable(kit, opt.BasicFilter)
		if err != nil {
			blog.Errorf("get table by basic filter(%+v) failed, err: %v, rid: %s", opt.BasicFilter, err, kit.Rid)
			return nil, err
		}

		dataArr, err := getDBDataByID[T](kit, opt, table, idField)
		if err != nil {
			return nil, err
		}

		allData := make([]interface{}, 0)
		for _, data := range dataArr {
			allData = append(allData, dataWithTenant[T]{
				TenantID: kit.TenantID,
				Data:     data,
			})
		}
		return allData, nil
	}
}

func listDataWithIDAndSubRes[T any](idField string, subResFields []string,
	getTable func(kit *rest.Kit, filter *types.BasicFilter) (string, error)) dataLister {

	return func(kit *rest.Kit, opt *listDataOpt) (*listDataRes, error) {
		kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

		if rawErr := opt.Validate(true); rawErr.ErrCode != 0 {
			blog.Errorf("list general data option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
			return nil, fmt.Errorf("list data option is invalid")
		}

		table, err := getTable(kit, opt.BasicFilter)
		if err != nil {
			blog.Errorf("get table by basic filter(%+v) failed, err: %v, rid: %s", opt.BasicFilter, err, kit.Rid)
			return nil, err
		}

		if opt.OnlyListID {
			opt.Fields = append(subResFields, idField)
		}

		cnt, dataArr, err := listDBDataWithID[T](kit, opt, table, idField)
		if err != nil {
			return nil, err
		}

		allData := make([]interface{}, 0)
		for _, data := range dataArr {
			allData = append(allData, dataWithTenant[T]{
				TenantID: kit.TenantID,
				Data:     data,
			})
		}

		return &listDataRes{Count: cnt, Data: allData}, nil
	}
}
