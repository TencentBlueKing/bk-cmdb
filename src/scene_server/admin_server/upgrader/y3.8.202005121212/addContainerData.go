/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_8_202005121212

import (
	"configcenter/src/common/mapstr"
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"gopkg.in/mgo.v2"
)

func addPresetObjects(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := addClassifications(ctx, db, conf); err != nil {
		return err
	}
	if err := addPropertyGroupData(ctx, db, conf); err != nil {
		return err
	}
	if err := addObjDesData(ctx, db, conf); err != nil {
		return err
	}
	if err := addObjAttDescData(ctx, db, conf); err != nil {
		return err
	}
	if err := addObjUnique(ctx, db, conf); err != nil {
		return err
	}
	return nil
}

var tables = map[string][]types.Index{
	common.BKTableNameBasePod: {
		types.Index{Name: "", Keys: map[string]int32{
			common.BKPodIDField: 1,
		}, Unique: true, Background: true},
		types.Index{Name: "", Keys: map[string]int32{
			common.BKPodUUIDField: 1,
		}, Unique: true, Background: true},
		types.Index{Name: "", Keys: map[string]int32{
			common.BKModuleIDField: 1,
		}, Background: true},
		types.Index{Name: "", Keys: map[string]int32{
			common.BKCloudIDField:     1,
			common.BKHostInnerIPField: 1,
		}, Background: true},
		types.Index{Name: "", Keys: map[string]int32{
			common.BKPodNameField:      1,
			common.BKPodNamespaceField: 1,
			common.BKPodClusterField:   1,
		}, Unique: true, Background: true},
	},
}

func createTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for tableName, indexes := range tables {
		exists, err := db.HasTable(ctx, tableName)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
				return err
			}
		}
		for index := range indexes {
			if err = db.Table(tableName).CreateIndex(ctx, indexes[index]); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

var classificationRows = []*metadata.Classification{
	&metadata.Classification{ClassificationID: "bk_container_manage", ClassificationName: "容器管理", ClassificationType: "inner", ClassificationIcon: "icon-cc-container"},
}

func addClassifications(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	tablename := common.BKTableNameObjClassification
	blog.Infof("add %s rows", tablename)
	for _, row := range classificationRows {
		if _, _, err = upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKClassificationIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	return nil
}

func getPropertyGroupData(ownerID string) []*metadata.Group {
	return []*metadata.Group{
		&metadata.Group{ObjectID: common.BKInnerObjIDPod, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
	}
}

func addPropertyGroupData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNamePropertyGroup
	blog.Errorf("add data for  %s table ", tablename)
	rows := getPropertyGroupData(conf.OwnerID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKObjIDField, "bk_group_id"}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	return nil
}

func getObjectDesData(ownerID string) []*metadata.Object {
	dataRows := []*metadata.Object{
		&metadata.Object{
			ObjCls:     "bk_container_manage",
			ObjectID:   common.BKInnerObjIDPod,
			ObjectName: "Pod",
			IsPre:      true,
			ObjIcon:    "",
			Position:   ``,
		},
	}
	t := metadata.Now()
	for _, r := range dataRows {
		r.CreateTime = &t
		r.LastTime = &t
		r.IsPaused = false
		r.Creator = common.CCSystemOperatorUserName
		r.OwnerID = ownerID
		r.Description = ""
		r.Modifier = ""
	}

	return dataRows
}

func addObjDesData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjDes
	blog.Errorf("add data for  %s table ", tablename)
	rows := getObjectDesData(conf.OwnerID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	return nil
}

func getObjAttDescData(ownerID string) []*Attribute {
	predataRows := PodRow()
	t := new(time.Time)
	*t = time.Now()
	for _, r := range predataRows {
		r.OwnerID = ownerID
		r.IsPre = true
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}
	return predataRows
}

func addObjAttDescData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAttDes
	blog.Infof("add data for %s table", tablename)
	rows := getObjAttDescData(conf.OwnerID)
	for _, row := range rows {
		_, _, err := upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
		if err != nil {
			blog.Errorf("add data for %s table error %#v", tablename, err)
			return err
		}
	}
	return nil
}

func addObjUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	oldAttributes := []Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(
		mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDPod}).All(ctx, oldAttributes)
	if err != nil {
		return err
	}

	propertyIDToProperty := make(map[string]Attribute{})
	var keyfunc = func(a, b string) string { return a + ":" + b }

	for _, oldAttr := range oldAttributes {
		ropertyIDToProperty[keyfunc(oldAttr.ObjectID, oldAttr.PropertyID)] = oldAttr
	}

	uniques := []metadata.ObjectUnique{
		// pod
		{
			ObjID:     common.BKInnerObjIDPod,
			MustCheck: true,
			Keys: []metadata.UniqueKey{
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPod, common.BKPodNameField)].ID),
				},
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPod, common.BKPodNamespaceField)].ID),
				},
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPod, common.BKPodClusterField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		{
			ObjID:     common.BKInnerObjIDPod,
			MustCheck: true,
			Keys: []metadata.UniqueKey{
				{
					Kind: metadata.UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPod, common.BKPodUUIDField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
	}

	for _, unique := range uniques {
		exists, err := isUniqueExists(ctx, db, conf, unique)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		uid, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
		if err != nil {
			return err
		}
		unique.ID = uid
		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	return nil
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
