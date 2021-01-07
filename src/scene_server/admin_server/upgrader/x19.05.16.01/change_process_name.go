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

package x19_05_16_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func changeProcessName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	mapping := map[string]string{
		common.BKProcessNameField: "进程别名",
		common.BKFuncName:         "进程名称",
		common.BKDescriptionField: "备注",
	}

	for key, newname := range mapping {
		cond := condition.CreateCondition().
			Field(common.BKObjIDField).Eq(common.BKInnerObjIDProc).
			Field(common.BKPropertyIDField).Eq(key).
			ToMapStr()
		if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, mapstr.MapStr{common.BKPropertyNameField: newname}); err != nil {
			return err
		}
	}

	cond := condition.CreateCondition().
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDProc).
		Field(common.BKPropertyIDField).Eq(common.BKProcessNameField).
		ToMapStr()
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, mapstr.MapStr{common.BKIsRequiredField: false}); err != nil {
		return err
	}

	procNameAtt := metadata.ObjAttDes{}
	if err := db.Table(common.BKTableNameObjAttDes).Find(condition.CreateCondition().
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDProc).
		Field(common.BKPropertyIDField).Eq(common.BKProcessNameField).
		ToMapStr()).One(ctx, &procNameAtt); err != nil {
		return err
	}
	bizIDAtt := metadata.ObjAttDes{}
	if err := db.Table(common.BKTableNameObjAttDes).Find(condition.CreateCondition().
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDProc).
		Field(common.BKPropertyIDField).Eq(common.BKAppIDField).
		ToMapStr()).One(ctx, &bizIDAtt); err != nil {
		return err
	}

	uniques := make([]metadata.ObjectUnique, 0)
	if err := db.Table(common.BKTableNameObjUnique).Find(condition.CreateCondition().
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDProc).
		ToMapStr()).All(ctx, &uniques); err != nil {
		return err
	}

	shouldDeleteUnique := metadata.ObjectUnique{
		Keys: []metadata.UniqueKey{
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   uint64(bizIDAtt.ID),
			},
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   uint64(procNameAtt.ID),
			},
		},
	}.KeysHash()
	for _, unique := range uniques {
		if unique.KeysHash() == shouldDeleteUnique {
			if err := db.Table(common.BKTableNameObjUnique).Delete(ctx, condition.CreateCondition().Field(common.BKFieldID).Eq(unique.ID).ToMapStr()); err != nil {
				return err
			}
		}
	}

	return nil
}
