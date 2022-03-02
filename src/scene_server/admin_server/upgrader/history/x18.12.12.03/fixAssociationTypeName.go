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

package x18_12_12_03

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func fixAssociationTypeName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	nameKV := map[string]string{
		"run":         "运行",
		"group":       "组成",
		"default":     "默认关联",
		"cover":       "覆盖",
		"connect":     "上联",
		"bk_mainline": "拓扑组成",
		"belong":      "属于",
	}

	for id, name := range nameKV {
		cond := condition.CreateCondition()
		cond.Field(common.BKOwnerIDField).Eq(common.BKDefaultOwnerID)
		cond.Field(common.AssociationKindIDField).Eq(id)

		data := mapstr.MapStr{
			common.AssociationKindNameField: name,
		}

		err := db.Table(common.BKTableNameAsstDes).Update(ctx, cond.ToMapStr(), data)
		if err != nil {
			return err
		}
	}
	return nil
}
