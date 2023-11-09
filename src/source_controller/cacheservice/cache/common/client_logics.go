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

package common

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/common/key"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// listWithIDs list common resource info from redis with ids
// if a common resource does not exist in the cache and cannot be found in mongodb, it will not be returned
// therefore the length and sequence of the returned array may not be equal to the requested ids
func (c *Client) listWithIDs(kit *rest.Kit, generator key.KeyGenerator, ids, fields []string) ([]string, error) {
	if len(ids) == 0 {
		return make([]string, 0), nil
	}

	ids = util.StrArrayUnique(ids)
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = generator.DetailKey(id)
	}

	results, err := redis.Client().MGet(kit.Ctx, keys...).Result()
	if err != nil {
		blog.Errorf("get %s ids %+v details from redis failed, err: %v, rid: %s", generator.Type(), ids, err, kit.Rid)
		return nil, err
	}

	needRefreshIdx := make([]int, 0)
	details := make([]string, 0)
	for idx, res := range results {
		if res == nil {
			needRefreshIdx = append(needRefreshIdx, idx)
			continue
		}

		detail, ok := res.(string)
		if !ok {
			blog.Errorf("%s %s detail(%+v) is invalid, rid: %s", generator.Type(), ids[idx], res, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "detail")
		}

		if len(fields) != 0 {
			details = append(details, *json.CutJsonDataWithFields(&detail, fields))
		} else {
			details = append(details, detail)
		}
	}

	if len(needRefreshIdx) == 0 {
		return details, nil
	}

	// can not find detail in cache, need refresh the cache
	refreshIDs := make([]string, len(needRefreshIdx))
	for i, idx := range needRefreshIdx {
		refreshIDs[i] = ids[idx]
	}

	refreshDetails, err := c.listWithRefreshKeys(kit, generator, key.IDKind, refreshIDs, fields)
	if err != nil {
		return nil, err
	}

	return append(details, refreshDetails...), nil
}

// listWithRefreshKeys list common resource info from mongo with refresh keys
func (c *Client) listWithRefreshKeys(kit *rest.Kit, generator key.KeyGenerator, kind key.KeyKind,
	keys, fields []string) ([]string, error) {

	if len(keys) == 0 {
		return make([]string, 0), nil
	}

	mongoData, err := generator.GetMongoData(kind, mongodb.Client(), keys...)
	if err != nil {
		blog.Errorf("get %s ids %+v mongo data failed, err: %v, rid: %s", generator.Type(), keys, err, kit.Rid)
		return nil, err
	}

	details := make([]string, 0)
	for _, data := range mongoData {
		c.tryRefreshDetail(generator, data, kit.Rid)

		detailJs, err := json.Marshal(data)
		if err != nil {
			blog.Errorf("marshal %s mongo data %+v failed, err: %v, rid: %s", generator.Type(), data, err, kit.Rid)
			return nil, err
		}
		detailStr := string(detailJs)

		if len(fields) != 0 {
			details = append(details, *json.CutJsonDataWithFields(&detailStr, fields))
		} else {
			details = append(details, detailStr)
		}
	}
	return details, nil
}

func (c *Client) tryRefreshDetail(generator key.KeyGenerator, data interface{}, rid string) {
	idKey, _, err := generator.GetIDKey(data)
	if err != nil {
		blog.Errorf("generate %s try refresh key failed, err: %v, data: %+v, rid: %s", generator.Type(), err, data, rid)
		return
	}

	detailKey := generator.DetailKey(idKey)
	if !c.lock.CanRefresh(detailKey) {
		return
	}

	// set refreshing status
	c.lock.SetRefreshing(detailKey)

	// check if we can refresh the common resource detail cache
	go func() {
		defer c.lock.SetUnRefreshing(detailKey)

		refreshCache(generator, data, rid)
	}()
}
