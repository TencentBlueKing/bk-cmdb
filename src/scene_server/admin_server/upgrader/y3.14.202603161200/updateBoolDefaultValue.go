/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_14_202603161200

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// upsertObjAttDesBoolDefaultValue
//
//option:true ->default:true
//option:false or default:null  ->default:false
func upsertObjAttDesBoolDefaultValue(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	updateCond := bson.M{
		common.BKPropertyTypeField: common.FieldTypeBool,
		common.BKDBOR: bson.A{
			bson.M{common.BKOptionField: false},
			bson.M{common.BKDefaultField: nil},
		},
	}
	updateData := map[string]interface{}{common.BKDefaultField: false}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, updateCond, updateData); err != nil {
		blog.Errorf("update bool attribute failed, err: %v, cond: %v, updateData: %v", err, updateCond,
			updateData)
		return err
	}

	updateCond = bson.M{common.BKOptionField: true}
	updateData = map[string]interface{}{common.BKDefaultField: true}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, updateCond, updateData); err != nil {
		blog.Errorf("update bool attribute failed, err: %v, cond: %v, updateData: %v", err, updateCond,
			updateData)
		return err
	}

	return nil
}
