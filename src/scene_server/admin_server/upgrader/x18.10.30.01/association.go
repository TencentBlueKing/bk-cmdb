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
package x18_10_30_01

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func createAssociationTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablenames := []string{common.BKTableNameAsstDes, common.BKTableNameObjAsst, common.BKTableNameInstAsst}
	for _, tablename := range tablenames {
		exists, err := db.HasTable(tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(tablename); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

func addPresetAssociationType(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameAsstDes

	asstTypes := []metadata.AssociationType{
		{
			AsstID:    "belong",
			AsstName:  "",
			OwnerID:   conf.OwnerID,
			SrcDes:    "属于",
			DestDes:   "包含",
			Direction: "src_to_dest",
		},
		{
			AsstID:    "group",
			AsstName:  "",
			OwnerID:   conf.OwnerID,
			SrcDes:    "组成",
			DestDes:   "组成于",
			Direction: "src_to_dest",
		},
		{
			AsstID:    "run",
			AsstName:  "",
			OwnerID:   conf.OwnerID,
			SrcDes:    "运行于",
			DestDes:   "运行",
			Direction: "src_to_dest",
		},
		{
			AsstID:    "belong",
			AsstName:  "",
			OwnerID:   conf.OwnerID,
			SrcDes:    "上联",
			DestDes:   "下联",
			Direction: "src_to_dest",
		},
	}

	for _, asstType := range asstTypes {
		_, _, err := upgrader.Upsert(ctx, db, tablename, asstType, "id", []string{"bk_asst_id"}, []string{"id"})
		if err != nil {
			return err
		}
	}
	return nil
}

func reconcilAsstData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAsst
	assts := []metadata.Association{}
	err := db.Table(tablename).Find(nil).All(ctx, &assts)
	if err != nil {
		return err
	}

	propertyCond := condition.CreateCondition()
	propertyCond.Field(common.BKPropertyTypeField).In([]string{common.FieldTypeSingleAsst, common.FieldTypeMultiAsst})
	propertys := []metadata.ObjAttDes{}
	err = db.Table(common.BKTableNameObjAttDes).Find(propertyCond.ToMapStr()).All(ctx, &propertys)
	if err != nil {
		return err
	}

	properyMap := map[string]metadata.ObjAttDes{}
	buildObjPropertyMapKey := func(objID string, propertyID string) string { return fmt.Sprintf("%s:%s") }
	for _, property := range propertys {
		properyMap[buildObjPropertyMapKey(property.ObjectID, property.PropertyID)] = property
	}

	for _, asst := range assts {
		if asst.ObjectAttID == common.BKChildStr {
			asst.AsstID = common.AssociationTypeGroup
			asst.ObjectAsstID = buildObjAsstID(asst.AsstObjID, asst.ObjectAttID)
			asst.Mapping = common.AssociationMappingOneToOne
			asst.OnDelete = common.AssociationOnDeleteNone
			switch asst.ObjectID {
			case common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDHost:
				asst.IsPre = true
			}
		} else {
			asst.AsstID = common.AssociationTypeDefault
			asst.ObjectAsstID = buildObjAsstID(asst.AsstObjID, asst.ObjectAttID)
			property := properyMap[buildObjPropertyMapKey(asst.ObjectID, asst.ObjectAttID)]
			switch property.PropertyType {
			case common.FieldTypeSingleAsst:
				asst.Mapping = common.AssociationMappingOneToOne
			case common.FieldTypeMultiAsst:
				asst.Mapping = common.AssociationMappingOneToMulti
			default:
				asst.Mapping = common.AssociationMappingOneToOne
			}
			asst.OnDelete = common.AssociationOnDeleteNone
		}
		_, _, err := upgrader.Upsert(ctx, db, tablename, asst, "id", []string{"bk_object_id", "bk_asst_object_id"}, []string{"id"})
		if err != nil {
			return err
		}

		updateInst := mapstr.New()
		updateInst.Set("bk_obj_asst_id", asst.ObjectAsstID)
		updateInst.Set("bk_asst_id", asst.AsstID)
		updateInst.Set("last_time", time.Now())

		updateInstCond := condition.CreateCondition()
		updateInstCond.Field("bk_obj_id").Eq(asst.ObjectID)
		updateInstCond.Field("bk_asst_id").Eq(asst.AsstObjID)
		err = db.Table(common.BKTableNameInstAsst).Update(ctx, updateInstCond.ToMapStr(), updateInst)
		if err != nil {
			return err
		}
	}

	err = db.Table(common.BKTableNameObjAttDes).Delete(ctx, propertyCond.ToMapStr())
	if err != nil {
		return err
	}
	return nil
}

func buildObjAsstID(srcObj string, propertyID string) string {
	return fmt.Sprintf("%s_%s", srcObj, propertyID)
}
