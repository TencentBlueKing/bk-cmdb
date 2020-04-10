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

package y3_9_202002131522

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addPlatUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// find bk_vpc_id property's id
	attrCond := condition.CreateCondition()
	attrCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDPlat)
	attrCond.Field(common.BKPropertyIDField).Eq(common.BKVpcID)
	attrCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	attr := new(metadata.Attribute)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrCond.ToMapStr()).One(ctx, attr)
	if err != nil {
		return err
	}
	id := uint64(attr.ID)

	// check if this unique exists
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDPlat)
	uniqueCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	existUniques := make([]metadata.ObjectUnique, 0)
	err = db.Table(common.BKTableNameObjUnique).Find(uniqueCond.ToMapStr()).All(ctx, &existUniques)
	if err != nil {
		return err
	}
	for _, u := range existUniques {
		if len(u.Keys) == 1 && u.Keys[0].ID == id {
			return nil
		}
	}

	// insert bk_vpc_id unique
	unique := metadata.ObjectUnique{
		ObjID:     common.BKInnerObjIDPlat,
		MustCheck: false,
		Keys: []metadata.UniqueKey{
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   id,
			},
		},
		Ispre:    false,
		OwnerID:  conf.OwnerID,
		LastTime: metadata.Now(),
	}
	uid, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
	if err != nil {
		return err
	}
	unique.ID = uid
	if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
		return err
	}
	return nil
}
