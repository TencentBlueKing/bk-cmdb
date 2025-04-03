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

package scheduler

import (
	"context"
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/sharding"
)

// genWatchDBRelationInfo generate db uuid to watch db uuid map and default watch db uuid for db without watch db
func genWatchDBRelationInfo(db dal.Dal) (map[string]string, string, error) {
	ctx := context.Background()
	masterDB := db.Shard(sharding.NewShardOpts().WithIgnoreTenant())

	relations := make([]sharding.WatchDBRelation, 0)
	if err := masterDB.Table(common.BKTableNameWatchDBRelation).Find(nil).All(ctx, &relations); err != nil {
		return nil, "", fmt.Errorf("get db and watch db relation failed, err: %v", err)
	}

	watchDBRelation := make(map[string]string)
	for _, relation := range relations {
		watchDBRelation[relation.DB] = relation.WatchDB
	}

	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	conf := new(sharding.ShardingDBConf)
	err := masterDB.Table(common.BKTableNameSystem).Find(cond).One(ctx, &conf)
	if err != nil {
		return nil, "", fmt.Errorf("get sharding db conf failed, err: %v", err)
	}
	return watchDBRelation, conf.ForNewData, nil
}

type watchObserver struct {
	isMaster       discovery.ServiceManageInterface
	previousStatus bool
}

// canLoop describe whether we can still loop the next event or next batch events.
// this is a master slave service. we should re-watch the event from the previous
// event token, only when we do this, we can loop the continuous events later which
// is no events is skipped or duplicated.
func (o *watchObserver) canLoop() (reWatch bool, loop bool) {
	current := o.isMaster.IsMaster()

	if o.previousStatus == current {
		if !current {
			// not master, status not changed, and can not loop
			return false, false
		} else {
			// is master, status not changed, and can loop
			return false, true
		}
	}

	blog.Infof("loop watch, is master status changed from %v to %v.", o.previousStatus, current)

	// update status
	o.previousStatus = current

	// status already changed, and can not continue loop, need to re-watch again.
	return true, false
}
