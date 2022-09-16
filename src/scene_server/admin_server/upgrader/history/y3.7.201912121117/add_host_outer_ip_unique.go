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

package y3_7_201912121117

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addHostOuterIPUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// find host outer ip property's id
	attrCond := condition.CreateCondition()
	attrCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	attrCond.Field(common.BKPropertyIDField).Eq(common.BKHostOuterIPField)
	attrCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	attr := new(metadata.Attribute)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrCond.ToMapStr()).One(ctx, attr)
	if err != nil {
		return err
	}
	id := uint64(attr.ID)

	// check if this unique exists
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	uniqueCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	existUniques := make([]objectUnique, 0)
	err = db.Table(common.BKTableNameObjUnique).Find(uniqueCond.ToMapStr()).All(ctx, &existUniques)
	if err != nil {
		return err
	}
	for _, u := range existUniques {
		if len(u.Keys) == 1 && u.Keys[0].ID == id {
			return nil
		}
	}

	// insert host outer ip unique
	unique := objectUnique{
		ObjID:     common.BKInnerObjIDHost,
		MustCheck: false,
		Keys: []UniqueKey{
			{
				Kind: UniqueKeyKindProperty,
				ID:   id,
			},
		},
		Ispre:    false,
		OwnerID:  conf.OwnerID,
		LastTime: Now(),
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

type objectUnique struct {
	ID        uint64      `json:"id" bson:"id"`
	ObjID     string      `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool        `json:"must_check" bson:"must_check"`
	Keys      []UniqueKey `json:"keys" bson:"keys"`
	Ispre     bool        `json:"ispre" bson:"ispre"`
	OwnerID   string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  time.Time   `json:"last_time" bson:"last_time"`
}

// Now TODO
func Now() time.Time {
	return time.Now().UTC()
}

// UniqueKey TODO
type UniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}

const (
	// UniqueKeyKindProperty TODO
	UniqueKeyKindProperty = "property"
)
