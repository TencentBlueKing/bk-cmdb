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

package y3_14_202603111314

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Migration logic:
//
//	For all existing documents where `is_hidden` is missing
//	or null, set the default value to false.
//
// Mongo command:
//
//	db.cc_ObjAttDes.updateMany(
//	  { $or: [{ is_hidden: { $exists: false } }, { is_hidden: null }] },
//	  { $set: { is_hidden: false } }
//	)
//
// Impact:
//	Only updates legacy documents.
//	New documents created after this migration will contain the field.

// Idempotent:
//
//	Running multiple times will not change existing correct data.
func upsertObjAttIsHidden(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	updateCond := bson.M{
		common.BKDBOR: bson.A{
			bson.M{common.BKIsHidden: bson.M{
				common.BKDBExists: false,
			}},
			bson.M{common.BKIsHidden: nil},
		},
	}
	updateData := map[string]interface{}{common.BKIsHidden: false}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, updateCond, updateData); err != nil {
		blog.Errorf("upsert attribute failed, err: %v, cond: %v, updateData: %v", err, updateCond,
			updateData)
		return err
	}
	return nil
}
