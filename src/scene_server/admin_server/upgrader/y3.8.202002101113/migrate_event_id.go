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

package y3_8_202002101113

import (
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
)

func migrateEventIDToMongo(ctx context.Context, db dal.RDB, cache redis.Client, conf *upgrader.Config) error {
	sid, err := cache.Get(ctx, common.EventCacheEventIDKey).Result()
	if redis.IsNilErr(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}

	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return err
	}

	docs := map[string]interface{}{
		"_id":        common.EventCacheEventIDKey,
		"SequenceID": id,
	}

	filter := map[string]interface{}{
		"_id": common.EventCacheEventIDKey,
	}

	err = db.Table(common.BKTableNameIDgenerator).Upsert(ctx, filter, docs)
	if err != nil {
		return err
	}

	return cache.Del(ctx, common.EventCacheEventIDKey).Err()
}
