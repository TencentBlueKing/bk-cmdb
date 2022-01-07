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

package y3_10_202112181521

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

var dataRows = map[string]*metadata.Attribute{
	common.BKBizSetNameField: {
		ObjectID:      common.BKInnerObjIDBizSet,
		PropertyID:    common.BKBizSetNameField,
		PropertyName:  "业务集名",
		IsRequired:    true,
		IsOnly:        true,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeSingleChar,
		Option:        `^[^\\\|\/:\*,<>"\?#\s]+$`,
	},
	common.BKBizSetIDField: {
		ObjectID:      common.BKInnerObjIDBizSet,
		PropertyID:    common.BKBizSetIDField,
		PropertyName:  "业务集ID",
		IsAPI:         true,
		IsRequired:    false,
		IsOnly:        true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeInt,
		Option:        metadata.IntOption{},
	},
	common.BKBizSetDescField: {
		ObjectID:      common.BKInnerObjIDBizSet,
		PropertyID:    common.BKBizSetDescField,
		PropertyName:  "业务集描述",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeSingleChar,
		Option:        "",
	},
	common.BKMaintainersField: {
		ObjectID:      common.BKInnerObjIDBizSet,
		PropertyID:    common.BKMaintainersField,
		PropertyName:  "运维人员",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.AppRole,
		PropertyType:  common.FieldTypeUser,
		Option:        "",
	},
	common.BKBizSetScopeField: {
		ObjectID:      common.BKInnerObjIDBizSet,
		PropertyID:    common.BKBizSetScopeField,
		PropertyName:  "业务范围",
		IsRequired:    true,
		IsOnly:        false,
		IsEditable:    true,
		IsAPI:         true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldObject,
		Option:        "",
		Placeholder:   "业务集所包含的业务的条件",
	},
}

func addBizSetObjectRow(ctx context.Context, db dal.RDB, ownerID string) error {

	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDBizSet}
	model := new(metadata.Object)

	// 判断是否有 biz_set 的对象表，如果没有需要初始化
	err := db.Table(common.BKTableNameObjDes).Find(filter).
		Fields(common.BKFieldID, common.BKObjNameField, common.CreatorField).One(ctx, model)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("count biz set object failed, err: %v", err)
		return err
	}
	if model.ID != 0 {
		if model.ObjectName == "业务集" && model.Creator == common.CCSystemOperatorUserName {
			return nil
		}
		blog.Errorf("the model %s already exists, and not system establishment, you must deal it.",
			common.BKInnerObjIDBizSet)
		return fmt.Errorf("model %s must be deal", common.BKInnerObjIDBizSet)
	}

	t := metadata.Now()
	dataRows := metadata.Object{
		ObjCls:      "bk_organization",
		ObjectID:    common.BKInnerObjIDBizSet,
		ObjectName:  "业务集",
		IsPre:       true,
		ObjIcon:     "icon-cc-business-set",
		CreateTime:  &t,
		LastTime:    &t,
		IsPaused:    false,
		Creator:     common.CCSystemOperatorUserName,
		OwnerID:     ownerID,
		Description: "",
		Modifier:    "",
	}
	uniqueKeys := []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}
	_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameObjDes, dataRows, "id", uniqueKeys, []string{"id"})
	if err != nil {
		blog.Errorf("add data for %s table error  %s", common.BKTableNameObjDes, err)
		return err
	}
	return nil
}

func addObjectUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	attrs := make([]metadata.Attribute, 0)

	cond := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDBizSet,
		"bk_property_id":    mapstr.MapStr{common.BKDBIN: []string{common.BKBizSetNameField, common.BKBizSetIDField}},
	}

	err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &attrs)
	if err != nil {
		return err
	}
	for _, attr := range attrs {
		keys := make([]metadata.UniqueKey, 0)

		keys = append(keys, metadata.UniqueKey{
			Kind: metadata.UniqueKeyKindProperty,
			ID:   uint64(attr.ID),
		})
		unique := metadata.ObjectUnique{
			ObjID:    common.BKInnerObjIDBizSet,
			Keys:     keys,
			Ispre:    true,
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
	}

	return nil

}

func addBizSetCollection(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	exists, err := db.HasTable(ctx, common.BKTableNameBaseBizSet)
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v, rid: %s", common.BKTableNameWatchToken, err)
		return err
	}

	if !exists {
		if err := db.CreateTable(ctx, common.BKTableNameBaseBizSet); err != nil {
			return err
		}
		return nil
	}
	blog.Infof("biz set collection has been created")
	return nil
}

// addBizSetPropertyGroup add biz set property group.
func addBizSetPropertyGroup(ctx context.Context, db dal.RDB, ownerID string) error {

	rows := []*metadata.Group{
		{
			ObjectID:   common.BKInnerObjIDBizSet,
			GroupID:    mCommon.BaseInfo,
			GroupName:  mCommon.BaseInfoName,
			GroupIndex: 1,
			OwnerID:    ownerID,
			IsDefault:  true,
		}, {
			ObjectID:   common.BKInnerObjIDBizSet,
			GroupID:    mCommon.AppRole,
			GroupName:  mCommon.AppRoleName,
			GroupIndex: 2,
			OwnerID:    ownerID,
			IsDefault:  true,
		},
	}

	for _, row := range rows {
		if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNamePropertyGroup, row, "id",
			[]string{common.BKObjIDField, common.BKPropertyGroupIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", common.BKTableNamePropertyGroup, err)
			return err
		}
	}

	return nil
}

// addBizSetObjectAttrRow update process bind info attribute
func addBizSetObjectAttrRow(ctx context.Context, db dal.RDB, ownerID string) error {

	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDBizSet}
	attrs := make([]metadata.Attribute, 0)
	// 判断是否有bizSet的对象属性表，如果没有需要初始化
	if err := db.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKPropertyIDField, common.BKObjNameField,
		common.CreatorField).All(ctx, &attrs); err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find object attribute describe fail, err: %v", err)
		return err
	}

	if len(attrs) > 0 {
		// 如果存在的话，数量必须一致。并且必须严格校验每个属性bk_property_name和creator必须完全一致，不一致直接报错需要先处理完毕后再升级
		if len(attrs) != len(dataRows) {
			blog.Errorf("the biz set model attrs num is incorrect")
			return fmt.Errorf("the biz set model attrs num is incorrect")
		}

		for _, attr := range attrs {
			if data, ok := dataRows[attr.PropertyID]; ok {
				if attr.PropertyName != data.PropertyName || attr.Creator != attr.Creator {
					blog.Errorf("the model biz set attr %s already exists, and not system establishment, you must"+
						" deal it.", attr.PropertyID)
					return fmt.Errorf("the model biz set attr  already exists")
				}
			}
		}
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	nowTime := metadata.Now()
	for _, row := range dataRows {
		row.OwnerID = ownerID
		row.IsPre = true
		row.IsReadOnly = false
		row.CreateTime = &nowTime
		row.LastTime = &nowTime
		row.Creator = common.CCSystemOperatorUserName
		row.Description = ""

		_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, row, "id", uniqueFields, []string{})
		if err != nil {
			blog.Errorf("add biz set attr failed, attribute: %v, err: %v", row, err)
			return err
		}
	}
	return nil
}

func addBizSePropertytOption(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	if err := addBizSetObjectRow(ctx, db, conf.OwnerID); err != nil {
		return err
	}

	if err := addBizSetPropertyGroup(ctx, db, conf.OwnerID); err != nil {
		return err
	}

	if err := addBizSetObjectAttrRow(ctx, db, conf.OwnerID); err != nil {
		return err
	}

	if err := addObjectUnique(ctx, db, conf); err != nil {
		return err
	}

	if err := addBizSetCollection(ctx, db, conf); err != nil {
		return err
	}
	return nil
}
