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

package x18_11_19_01

// import (
// 	"context"

// 	"github.com/rs/xid"

// 	"configcenter/src/common"
// 	"configcenter/src/common/condition"
// 	"configcenter/src/common/mapstr"
// 	"configcenter/src/scene_server/admin_server/upgrader"
// 	"configcenter/src/storage/dal"
// )

// func reconcilAsstID(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
// 	start, limit := uint64(0), uint64(100)

// 	type HostInst struct {
// 		HostID  uint64 `bson:"bk_host_id"`
// 		AssetID string `bson:"bk_asset_id"`
// 		OwnerID string `bson:"bk_supplier_account"`
// 	}

// 	cond := condition.CreateCondition()
// 	or := cond.NewOR()
// 	or.Item(mapstr.MapStr{common.BKAssetIDField: nil})
// 	or.Item(mapstr.MapStr{common.BKAssetIDField: ""})

// 	for {
// 		hosts := []HostInst{}

// 		err := db.Table(common.BKTableNameBaseHost).Find(cond.ToMapStr()).
// 			Sort(common.BKHostIDField).
// 			Start(start).Limit(limit).All(ctx, &hosts)
// 		if err != nil {
// 			return err
// 		}

// 		for index := range hosts {
// 			if hosts[index].AssetID == "" {
// 				data := mapstr.MapStr{
// 					common.BKAssetIDField: xid.New().String(),
// 				}
// 				updateCond := condition.CreateCondition()
// 				updateCond.Field(common.BKHostIDField).Eq(hosts[index].HostID)
// 				updateCond.Field(common.BKOwnerIDField).Eq(hosts[index].OwnerID)

// 				if err := db.Table(common.BKTableNameBaseHost).
// 					Update(ctx, updateCond.ToMapStr(), data); err != nil {
// 					return err
// 				}
// 			}
// 		}

// 		if uint64(len(hosts)) < limit {
// 			break
// 		}

// 		start += limit
// 	}

// 	return nil
// }
