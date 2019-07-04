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
	"configcenter/src/common/blog"
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

	asstTypes := []metadata.AssociationKind{
		{
			AssociationKindID:       "belong",
			AssociationKindName:     "",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "属于",
			DestinationToSourceNote: "包含",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
		},
		{
			AssociationKindID:       "group",
			AssociationKindName:     "",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "组成",
			DestinationToSourceNote: "组成于",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
		},
		{
			AssociationKindID:       "bk_mainline",
			AssociationKindName:     "",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "组成",
			DestinationToSourceNote: "组成于",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
		},
		{
			AssociationKindID:       "run",
			AssociationKindName:     "",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "运行于",
			DestinationToSourceNote: "运行",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
		},
		{
			AssociationKindID:       "connect",
			AssociationKindName:     "",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "上联",
			DestinationToSourceNote: "下联",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
		},
		{
			AssociationKindID:       "default",
			AssociationKindName:     "默认关联",
			OwnerID:                 conf.OwnerID,
			SourceToDestinationNote: "关联",
			DestinationToSourceNote: "被关联",
			Direction:               metadata.DestinationToSource,
			IsPre:                   ptrue(),
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

type Association struct {
	metadata.Association `bson:",inline"`
	ObjectAttID          string `bson:"bk_object_att_id"`
}

func reconcilAsstData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	assts := []Association{}
	err := db.Table(common.BKTableNameObjAsst).Find(nil).All(ctx, &assts)
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
	buildObjPropertyMapKey := func(objID string, propertyID string) string { return fmt.Sprintf("%s:%s", objID, propertyID) }
	for _, property := range propertys {
		properyMap[buildObjPropertyMapKey(property.ObjectID, property.PropertyID)] = property
		blog.Infof("key %s: %+v", buildObjPropertyMapKey(property.ObjectID, property.PropertyID), property)
	}

	flag := "updateflag"
	for _, asst := range assts {
		if asst.ObjectAttID == common.BKChildStr {
			asst.AsstKindID = common.AssociationKindMainline
			asst.AssociationName = buildObjAsstID(asst)
			asst.Mapping = metadata.OneToOneMapping
			asst.OnDelete = metadata.NoAction
			if (asst.ObjectID == common.BKInnerObjIDModule && asst.AsstObjID == common.BKInnerObjIDSet) ||
				(asst.ObjectID == common.BKInnerObjIDHost && asst.AsstObjID == common.BKInnerObjIDModule) {
				asst.IsPre = ptrue()
			} else {
				asst.IsPre = pfalse()
			}

			// update ObjAsst
			updateCond := condition.CreateCondition()
			updateCond.Field("id").Eq(asst.ID)
			if err = db.Table(common.BKTableNameObjAsst).Update(ctx, updateCond.ToMapStr(), asst); err != nil {
				return err
			}

			// update InstAsst
			updateInst := mapstr.New()
			updateInst.Set("bk_obj_asst_id", asst.AssociationName)
			updateInst.Set("bk_asst_id", asst.AsstKindID)
			updateInst.Set("last_time", time.Now())
			err = db.Table(common.BKTableNameInstAsst).Update(ctx, updateCond.ToMapStr(), updateInst)
			if err != nil {
				return err
			}

		} else {
			asst.AsstKindID = common.AssociationTypeDefault
			property := properyMap[buildObjPropertyMapKey(asst.ObjectID, asst.ObjectAttID)]
			switch property.PropertyType {
			case common.FieldTypeSingleAsst:
				asst.Mapping = metadata.OneToManyMapping
			case common.FieldTypeMultiAsst:
				asst.Mapping = metadata.ManyToManyMapping
			default:
				blog.Warnf("property: %+v, asst: %+v, for key: %v", property, asst, buildObjPropertyMapKey(asst.ObjectID, asst.ObjectAttID))
				asst.Mapping = metadata.ManyToManyMapping
			}
			// 交换 源<->目标
			asst.AssociationAliasName = property.PropertyName
			asst.ObjectID, asst.AsstObjID = asst.AsstObjID, asst.ObjectID
			asst.OnDelete = metadata.NoAction
			asst.IsPre = pfalse()
			asst.AssociationName = buildObjAsstID(asst)

			blog.InfoJSON("obj: %s, att: %s to asst %s", asst.ObjectID, asst.ObjectAttID, asst)
			// update ObjAsst
			updateCond := condition.CreateCondition()
			updateCond.Field("id").Eq(asst.ID)
			if err = db.Table(common.BKTableNameObjAsst).Update(ctx, updateCond.ToMapStr(), asst); err != nil {
				return err
			}

			// update ObjAsst
			instAssts := []metadata.InstAsst{}

			instCond := condition.CreateCondition()
			instCond.Field("bk_obj_id").Eq(asst.AsstObjID)
			instCond.Field("bk_asst_obj_id").Eq(asst.ObjectID)
			instCond.Field(flag).NotEq(true)

			if err = db.Table(common.BKTableNameInstAsst).Find(instCond.ToMapStr()).All(ctx, &instAssts); err != nil {
				return err
			}
			for _, instAsst := range instAssts {
				updateInst := mapstr.New()
				updateInst.Set("bk_obj_asst_id", asst.AssociationName)
				updateInst.Set("bk_asst_id", asst.AsstKindID)

				// 交换 源<->目标
				updateInst.Set("bk_obj_id", instAsst.AsstObjectID)
				updateInst.Set("bk_asst_obj_id", instAsst.ObjectID)
				updateInst.Set("bk_inst_id", instAsst.AsstInstID)
				updateInst.Set("bk_asst_inst_id", instAsst.InstID)

				updateInst.Set(flag, true)

				updateInst.Set("last_time", time.Now())
				if err = db.Table(common.BKTableNameInstAsst).Update(ctx,
					mapstr.MapStr{
						"id": instAsst.ID,
					}, updateInst); err != nil {
					return err
				}
			}
		}
	}
	if err = db.Table(common.BKTableNameInstAsst).DropColumn(ctx, flag); err != nil {
		return err
	}

	// update bk_cloud_id to int
	cloudIDUpdateCond := condition.CreateCondition()
	cloudIDUpdateCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	cloudIDUpdateCond.Field(common.BKPropertyIDField).Eq(common.BKCloudIDField)
	cloudIDUpdateData := mapstr.New()
	cloudIDUpdateData.Set(common.BKPropertyTypeField, common.FieldTypeForeignKey)
	cloudIDUpdateData.Set(common.BKOptionField, nil)
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, cloudIDUpdateCond.ToMapStr(), cloudIDUpdateData)
	if err != nil {
		return err
	}
	deleteHostCloudAssociation := condition.CreateCondition()
	deleteHostCloudAssociation.Field("bk_obj_id").Eq(common.BKInnerObjIDHost)
	deleteHostCloudAssociation.Field("bk_asst_obj_id").Eq(common.BKInnerObjIDPlat)
	err = db.Table(common.BKTableNameObjAsst).Delete(ctx, deleteHostCloudAssociation.ToMapStr())
	if err != nil {
		return err
	}

	// drop outdate propertys
	err = db.Table(common.BKTableNameObjAttDes).Delete(ctx, propertyCond.ToMapStr())
	if err != nil {
		return err
	}

	// drop outdate column
	outdateColumns := []string{"bk_object_att_id", "bk_asst_forward", "bk_asst_name"}
	for _, column := range outdateColumns {
		if err = db.Table(common.BKTableNameObjAsst).DropColumn(ctx, column); err != nil {
			return err
		}
	}

	delCond := condition.CreateCondition()
	delCond.Field(common.AssociationKindIDField).Eq(nil)
	if err = db.Table(common.BKTableNameObjAsst).Delete(ctx, delCond.ToMapStr()); err != nil {
		return err
	}
	return nil
}

func buildObjAsstID(asst Association) string {
	return fmt.Sprintf("%s_%s_%s", asst.ObjectID, asst.AsstKindID, asst.AsstObjID)
}

func ptrue() *bool {
	tmp := true
	return &tmp
}
func pfalse() *bool {
	tmp := false
	return &tmp
}
