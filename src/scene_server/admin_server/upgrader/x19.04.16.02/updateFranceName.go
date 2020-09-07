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

package x19_04_16_02

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func updateFranceName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	cond.Field(common.BKPropertyIDField).Eq("bk_state_name")
	state := metadata.Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, &state)
	if err != nil {
		// bk_state_name already delete
		if err == types.ErrDocumentNotFound {
			return nil
		}
		return err
	}

	enums, err := metadata.ParseEnumOption(ctx, state.Option)
	if err != nil {
		return err
	}
	for index := range enums {
		if enums[index].ID == "FR" {
			enums[index].Name = "法国"
		}
	}

	return db.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), mapstr.MapStr{common.BKOptionField: enums})
}
