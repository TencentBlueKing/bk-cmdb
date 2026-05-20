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

package y3_14_202604101200

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// addObjLabelField add field 'bk_labels' to object labels ,default:[]string
func addObjLabelField(ctx context.Context, db dal.RDB, _ *upgrader.Config) error {
	cond := bson.M{
		common.BKDBOR: bson.A{
			bson.M{metadata.ModelFieldLabels: bson.M{
				common.BKDBExists: false,
			}},
			bson.M{metadata.ModelFieldLabels: nil},
		},
	}
	doc := mapstr.MapStr{
		metadata.ModelFieldLabels: make([]string, 0),
	}
	if err := db.Table(common.BKTableNameObjDes).Update(ctx, cond, doc); err != nil {
		blog.ErrorJSON("failed to add object labels value of field %s to true, err: %s",
			metadata.ModelFieldLabels, err)
		return err
	}
	objLabelIndex := types.Index{
		Keys: bson.D{
			{metadata.ModelFieldLabels, 1},
		},
		Name:       common.CCLogicIndexNamePrefix + "obj_labels",
		Unique:     false,
		Background: true,
	}

	err := db.Table(common.BKTableNameObjDes).CreateIndex(ctx, objLabelIndex)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("create obj labels index %+v failed, err: %v", objLabelIndex, err)
		return err
	}

	return nil
}
