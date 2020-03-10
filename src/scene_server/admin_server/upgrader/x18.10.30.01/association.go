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
	"configcenter/src/storage/dal/types"
)

func createAssociationTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablenames := []string{common.BKTableNameAsstDes, common.BKTableNameObjAsst, common.BKTableNameInstAsst}
	for _, tablename := range tablenames {
		exists, err := db.HasTable(ctx, tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tablename); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

func createInstanceAssociationIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	idxArr, err := db.Table(common.BKTableNameInstAsst).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error. err:%s", common.BKTableNameInstAsst, err.Error())
		return err
	}

	createIdxArr := []types.Index{
		types.Index{Name: "idx_id", Keys: map[string]int32{"id": -1}, Background: true, Unique: true},
		types.Index{Name: "idx_objID_asstObjID_asstID", Keys: map[string]int32{"bk_obj_id": -1, "bk_asst_obj_id": -1, "bk_asst_id": -1}},
	}
	for _, idx := range createIdxArr {
		exist := false
		for _, existIdx := range idxArr {
			if existIdx.Name == idx.Name {
				exist = true
				break
			}
		}
		// index already exist, skip create
		if exist {
			continue
		}
		if err := db.Table(common.BKTableNameInstAsst).CreateIndex(ctx, idx); err != nil {
			blog.ErrorJSON("create index to cc_InstAsst error, err:%s, current index:%s, all create index:%s", err.Error(), idx, createIdxArr)
			return err
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
	propertyCond.Field(common.BKPropertyTypeField).In([]string{"multiasst", "singleasst"})
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
			case "singleasst":
				asst.Mapping = metadata.OneToManyMapping
			case "multiasst":
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

			instCond := condition.CreateCondition()
			instCond.Field("bk_obj_id").Eq(asst.AsstObjID)
			instCond.Field("bk_asst_obj_id").Eq(asst.ObjectID)
			instCond.Field(flag).NotEq(true)

			pageSize := uint64(2000)
			page := 0
			for {
				page += 1
				// update ObjAsst
				instAssts := []metadata.InstAsst{}
				blog.InfoJSON("find  data from table:%s, page:%s, cond:%s", common.BKTableNameInstAsst, page, instCond.ToMapStr())
				if err = db.Table(common.BKTableNameInstAsst).Find(instCond.ToMapStr()).Limit(pageSize).All(ctx, &instAssts); err != nil {
					return err
				}

				blog.InfoJSON("find  data from table:%s, cond:%s, result count:%s", common.BKTableNameInstAsst, instCond.ToMapStr(), len(instAssts))
				if len(instAssts) == 0 {
					break
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
					blog.InfoJSON("update instasst, id:%s, updateInst:%s", instAsst.ID, updateInst)
					if err = db.Table(common.BKTableNameInstAsst).Update(ctx,
						mapstr.MapStr{
							"id": instAsst.ID,
						}, updateInst); err != nil {
						return err
					}
				}
			}

		}
	}
	blog.InfoJSON("start drop column cond:%s", flag)
	if err = db.Table(common.BKTableNameInstAsst).DropColumn(ctx, flag); err != nil {
		return err
	}

	// update bk_cloud_id to int
	cloudIDUpdateCond := condition.CreateCondition()
	cloudIDUpdateCond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	cloudIDUpdateCond.Field(common.BKPropertyIDField).Eq(common.BKCloudIDField)
	cloudIDUpdateData := mapstr.New()
	cloudIDUpdateData.Set(common.BKPropertyTypeField, "foreignkey")
	cloudIDUpdateData.Set(common.BKOptionField, nil)
	blog.InfoJSON("update host cloud association cond:%s, data:%s", cloudIDUpdateCond.ToMapStr(), cloudIDUpdateData)
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, cloudIDUpdateCond.ToMapStr(), cloudIDUpdateData)
	if err != nil {
		return err
	}
	deleteHostCloudAssociation := condition.CreateCondition()
	deleteHostCloudAssociation.Field("bk_obj_id").Eq(common.BKInnerObjIDHost)
	deleteHostCloudAssociation.Field("bk_asst_obj_id").Eq(common.BKInnerObjIDPlat)
	blog.InfoJSON("delete host cloud association table:%s, cond:%s", common.BKTableNameObjAsst, deleteHostCloudAssociation.ToMapStr())
	err = db.Table(common.BKTableNameObjAsst).Delete(ctx, deleteHostCloudAssociation.ToMapStr())
	if err != nil {
		return err
	}

	blog.InfoJSON("delete host cloud association table:%s, cond:%s", common.BKTableNameObjAttDes, propertyCond.ToMapStr())
	// drop outdate propertys
	err = db.Table(common.BKTableNameObjAttDes).Delete(ctx, propertyCond.ToMapStr())
	if err != nil {
		return err
	}

	// drop outdate column
	outdateColumns := []string{"bk_object_att_id", "bk_asst_forward", "bk_asst_name"}
	for _, column := range outdateColumns {
		blog.InfoJSON("delete field from table:%s, cond:%s", common.BKTableNameObjAsst, column)
		if err = db.Table(common.BKTableNameObjAsst).DropColumn(ctx, column); err != nil {
			return err
		}
	}

	delCond := condition.CreateCondition()
	delCond.Field(common.AssociationKindIDField).Eq(nil)
	blog.InfoJSON("delete host cloud association table:%s, cond:%s", common.BKTableNameObjAsst, delCond.ToMapStr())
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
