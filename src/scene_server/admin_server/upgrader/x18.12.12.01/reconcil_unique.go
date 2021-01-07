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

func reconcilUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	oldAttributes := []Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDProc}).All(ctx, &oldAttributes)
	if err != nil {
		return err
	}
	var obj2IsOnlyProperty = map[string][]Attribute{}
	var propertyIDToProperty = map[string]Attribute{}

	var keyfunc = func(a, b string) string { return a + ":" + b }
	for _, oldAttr := range oldAttributes {
		if oldAttr.IsOnly {
			obj2IsOnlyProperty[oldAttr.ObjectID] = append(obj2IsOnlyProperty[oldAttr.ObjectID], oldAttr)
		}
		propertyIDToProperty[keyfunc(oldAttr.ObjectID, oldAttr.PropertyID)] = oldAttr
	}

	shouldCheck := []string{
		keyfunc(common.BKInnerObjIDProc, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDProc, "bk_start_param_regex"),
		keyfunc(common.BKInnerObjIDProc, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDProc, common.BKFuncName),
	}

	checkKeysShouldExists(propertyIDToProperty, shouldCheck)

	uniques := []metadata.ObjectUnique{
		// process
		{
			ObjID:     common.BKInnerObjIDProc,
			MustCheck: false,
			Keys: []metadata.UniqueKey{
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKAppIDField)].ID),
				},
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, "bk_start_param_regex")].ID),
				},
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKFuncName)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
	}

	uniqueIDs, err := db.NextSequences(ctx, common.BKTableNameObjUnique, len(uniques))
	if err != nil {
		return err
	}

	for index, unique := range uniques {
		exists, err := isUniqueExists(ctx, db, conf, unique)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		unique.ID = uniqueIDs[index]
		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	return nil
}

func checkKeysShouldExists(m map[string]Attribute, shouldExistKeys []string) []string {
	notValidKeys := []string{}
	for _, k := range shouldExistKeys {
		if _, ok := m[k]; !ok {
			notValidKeys = append(notValidKeys, k)
		}
	}
	return notValidKeys
}

func isUniqueExists(ctx context.Context, db dal.RDB, conf *upgrader.Config, unique metadata.ObjectUnique) (bool, error) {
	keyhash := unique.KeysHash()
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(unique.ObjID)
	uniqueCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	existUniques := []metadata.ObjectUnique{}

	err := db.Table(common.BKTableNameObjUnique).Find(uniqueCond.ToMapStr()).All(ctx, &existUniques)
	if err != nil {
		return false, err
	}

	for _, uni := range existUniques {
		if uni.KeysHash() == keyhash {
			return true, nil
		}
	}
	return false, nil

}
