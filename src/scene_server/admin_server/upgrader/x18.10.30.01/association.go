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
			DestinationToSourceNote: "关联",
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
		} else {
			asst.AsstKindID = common.AssociationTypeDefault
			asst.AssociationName = buildObjAsstID(asst)
			asst.Mapping = metadata.OneToManyMapping
			asst.OnDelete = metadata.NoAction
			asst.IsPre = pfalse()
		}
		_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameObjAsst, asst, "id", []string{"bk_obj_id", "bk_asst_obj_id"}, []string{"id"})
		if err != nil {
			return err
		}

		updateInst := mapstr.New()
		updateInst.Set("bk_obj_asst_id", asst.AssociationName)
		updateInst.Set("bk_asst_id", asst.AsstKindID)
		updateInst.Set("last_time", time.Now())

		updateInstCond := condition.CreateCondition()
		updateInstCond.Field("bk_obj_id").Eq(asst.ObjectID)
		updateInstCond.Field("bk_asst_obj_id").Eq(asst.AsstObjID)
		err = db.Table(common.BKTableNameInstAsst).Update(ctx, updateInstCond.ToMapStr(), updateInst)
		if err != nil {
			return err
		}
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
	return nil
}

func buildObjAsstID(asst Association) string {
	return fmt.Sprintf("%s_%s_%s_%s", asst.ObjectID, asst.AsstKindID, asst.AsstObjID, asst.ObjectAttID)
}

func ptrue() *bool {
	tmp := true
	return &tmp
}
func pfalse() *bool {
	tmp := false
	return &tmp
}
