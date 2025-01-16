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

package task

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/stream/types"
)

// genWatchDBRelationMap generate db uuid to watch db uuid map
func genWatchDBRelationMap(db dal.Dal) (map[string]string, error) {
	ctx := context.Background()
	masterDB := db.Shard(sharding.NewShardOpts().WithIgnoreTenant())

	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	shardingConf := new(sharding.ShardingDBConf)
	if err := masterDB.Table(common.BKTableNameSystem).Find(cond).One(ctx, &shardingConf); err != nil {
		blog.Errorf("get sharding db config failed, err: %v", err)
		return nil, err
	}

	watchDBRelation := make(map[string]string)
	relations := make([]sharding.WatchDBRelation, 0)
	if err := masterDB.Table(common.BKTableNameWatchDBRelation).Find(nil).All(ctx, &relations); err != nil {
		return nil, fmt.Errorf("get db and watch db relation failed, err: %v", err)
	}

	for _, relation := range relations {
		watchDBRelation[relation.DB] = watchDBRelation[relation.WatchDB]
	}
	return watchDBRelation, nil
}

// compareToken compare event with token, returns if event is greater than the token
func compareToken(event *types.Event, token *types.TokenInfo) bool {
	if event.ClusterTime.Sec > token.StartAtTime.Sec {
		return true
	}

	if event.ClusterTime.Sec == token.StartAtTime.Sec {
		if event.ClusterTime.Nano > token.StartAtTime.Nano {
			return true
		}
		return event.ClusterTime.Nano == token.StartAtTime.Nano && event.Token.Data > token.Token
	}
	return false
}
