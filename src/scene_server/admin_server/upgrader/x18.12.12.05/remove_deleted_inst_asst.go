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

package x18_12_12_05

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func removeDeletedInstAsst(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	assts := []metadata.InstAsst{}
	err := db.Table(common.BKTableNameInstAsst).Find(nil).All(ctx, &assts)
	if err != nil {
		return err
	}

	for _, asst := range assts {
		count, err := getInst(ctx, db, asst.ObjectID, asst.InstID, asst.OwnerID)
		if err != nil {
			return err
		}
		if count <= 0 {
			err := db.Table(common.BKTableNameInstAsst).Delete(ctx, mapstr.MapStr{common.BKFieldID: asst.ID})
			if err != nil {
				return err
			}
			continue
		}

		count, err = getInst(ctx, db, asst.AsstObjectID, asst.AsstInstID, asst.OwnerID)
		if err != nil {
			return err
		}
		if count <= 0 {
			err := db.Table(common.BKTableNameInstAsst).Delete(ctx, mapstr.MapStr{common.BKFieldID: asst.ID})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func getInst(ctx context.Context, db dal.RDB, objID string, instID int64, ownerID string) (uint64, error) {
	insttable := common.GetInstTableName(objID)
	cond := condition.CreateCondition()
	cond.Field(common.GetInstIDField(objID)).Eq(instID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)
	if insttable == common.BKTableNameBaseInst {
		cond.Field(common.BKObjIDField).Eq(objID)
	}
	return db.Table(insttable).Find(cond.ToMapStr()).Count(ctx)
}
