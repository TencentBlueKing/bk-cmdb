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
package x18_12_12_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func setTCPDefault(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(common.BKDefaultOwnerID)
	cond.Field(common.BKPropertyIDField).Eq(common.BKProtocol)

	ostypeProperty := Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, &ostypeProperty)
	if err != nil {
		return err
	}

	enumOpts, err := metadata.ParseEnumOption(ctx, ostypeProperty.Option)
	if err != nil {
		return err
	}
	for index := range enumOpts {
		if enumOpts[index].Name == "TCP" {
			enumOpts[index].IsDefault = true
		}
	}

	data := mapstr.MapStr{
		common.BKOptionField: enumOpts,
	}

	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		return err
	}

	procCond := condition.CreateCondition()
	procCond.Field(common.BKProtocol).Eq(nil)

	procData := mapstr.MapStr{
		common.BKProtocol: "1",
	}
	err = db.Table(common.BKTableNameBaseProcess).Update(ctx, procCond.ToMapStr(), procData)
	if err != nil {
		return err
	}

	return nil

}
