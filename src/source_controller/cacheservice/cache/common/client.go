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
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/cacheservice/cache/common/key"
	"configcenter/src/source_controller/cacheservice/cache/tools"
	"configcenter/src/storage/driver/redis"
)

var client *Client
var clientOnce sync.Once

// NewClient new common resource client
func NewClient() *Client {
	if client != nil {
		return client
	}

	clientOnce.Do(func() {
		client = &Client{
			lock: tools.NewRefreshingLock(),
		}
	})

	return client
}

// Client is the common resource cache client
type Client struct {
	lock tools.RefreshingLock
}

// ListWithKey search common resource cache info with specified keys
func (c *Client) ListWithKey(kit *rest.Kit, cacheType string, opt *metadata.ListCommonCacheWithKeyOpt) (
	[]string, error) {

	if len(cacheType) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "type")
	}

	keyGenerator, err := key.GetKeyGenerator(key.KeyType(cacheType))
	if err != nil {
		blog.Errorf("get %s key generator failed, err: %v, rid: %s", cacheType, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "type")
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		return nil, rawErr.ToCCError(kit.CCError)
	}

	allIDs := make([]string, 0)
	needRefreshKeys := make([]string, 0)
	keyKind := key.KeyKind(opt.Kind)
	if keyKind == key.IDKind {
		allIDs = opt.Keys
	} else {
		for _, redisKey := range opt.Keys {
			existRes, err := redis.Client().Exists(kit.Ctx, redisKey).Result()
			if err != nil {
				blog.Errorf("check if %s key %s exists failed, err: %v, rid: %s", cacheType, redisKey, err, kit.Rid)
				return nil, err
			}

			if existRes != 1 {
				needRefreshKeys = append(needRefreshKeys, redisKey)
				continue
			}

			ids, err := redis.Client().SMembers(kit.Ctx, redisKey).Result()
			if err != nil {
				blog.Errorf("get %s ids by other key %s failed, err: %v, rid: %s", cacheType, redisKey, err, kit.Rid)
				return nil, err
			}
			allIDs = append(allIDs, ids...)
		}
	}

	refreshDetails, err := c.listWithRefreshKeys(kit, keyGenerator, keyKind, needRefreshKeys, opt.Fields)
	if err != nil {
		return nil, err
	}

	details, err := c.listWithIDs(kit, keyGenerator, allIDs, opt.Fields)
	if err != nil {
		blog.Errorf("list %s cache info by ids %+v failed, err: %v, rid: %s", cacheType, allIDs, err, kit.Rid)
		return nil, err
	}
	return append(refreshDetails, details...), nil
}
