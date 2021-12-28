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
	"configcenter/src/common/condition"
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// default group
var (
	groupBaseInfo = mCommon.BaseInfo
)

func addBizSetObjectRow(ctx context.Context, db dal.RDB, ownerID string) error {

	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDBizSet}

	// 判断是否有 biz_set 的对象表，如果没有需要初始化
	if count, err := db.Table(common.BKTableNameObjDes).Find(filter).Count(ctx); err != nil {
		blog.Errorf("find object fail,err: %v", err)
		return err
	} else if count >= 1 {
		return nil
	}

	t := metadata.Now()
	dataRows := metadata.Object{
		ObjCls:      "bk_organization",
		ObjectID:    common.BKInnerObjIDBizSet,
		ObjectName:  "业务集",
		IsPre:       true,
		ObjIcon:     "icon-cc-business-set",
		Position:    `{"bk_organization":{"x":-100,"y":-100}}`,
		CreateTime:  &t,
		LastTime:    &t,
		IsPaused:    false,
		Creator:     common.CCSystemOperatorUserName,
		OwnerID:     ownerID,
		Description: "",
		Modifier:    "",
	}
	uniqueKeys := []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjDes, dataRows, "id", uniqueKeys, []string{"id"})
	if err != nil {
		blog.Errorf("add data for %s table error  %s", common.BKTableNameObjDes, err)
		return err
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

func addObjectUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	oldAttributes := make([]metadata.Attribute, 0)

	cond := mapstr.New()
	cond.Set(common.BKObjIDField, common.BKInnerObjIDBizSet)
	err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &oldAttributes)
	if err != nil {
		return err
	}

	keys := make([]metadata.UniqueKey, 0)

	for _, oldAttr := range oldAttributes {

		keys = append(keys, metadata.UniqueKey{
			Kind: metadata.UniqueKeyKindProperty,
			ID:   uint64(oldAttr.ID),
		})
	}
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

	return nil

}

func addDefaultBiz(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	if count, err := db.Table(common.BKTableNameBaseBizSet).Find(nil).Count(ctx); err != nil {
		return err
	} else if count >= 1 {
		return nil
	}

	// add default biz set
	defaultBizSet := map[string]interface{}{}

	if err := db.Table(common.BKTableNameBaseBizSet).InsertOne(ctx, defaultBizSet); err != nil {
		return err
	}
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

	// 判断是否有bizSet的对象属性表，如果没有需要初始化
	if count, err := db.Table(common.BKTableNameObjAttDes).Find(filter).Count(ctx); err != nil {
		blog.Errorf("find object attribute describe fail,err: %v", err)
		return err
	} else if count >= 4 {
		return nil
	}

	objID := common.BKInnerObjIDBizSet

	dataRows := []*metadata.Attribute{
		{
			ObjectID:      objID,
			PropertyID:    common.BKAppSetNameField,
			PropertyName:  "业务集名",
			IsRequired:    true,
			IsOnly:        true,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		{
			ObjectID:      objID,
			PropertyID:    common.BKBizSetDescField,
			PropertyName:  "业务集描述",
			IsRequired:    false,
			IsOnly:        false,
			IsPre:         true,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},

		{
			ObjectID:      objID,
			PropertyID:    common.BKMaintainersField,
			PropertyName:  "运维人员",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.AppRole,
			PropertyType:  common.FieldTypeUser,
			Option:        "",
		},
		{
			ObjectID:      objID,
			PropertyID:    common.BKScopeField,
			PropertyName:  "条件范围",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    true,
			IsAPI:         true,
			IsPre:         true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldCondition,
			Option:        "",
		},
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKPropertyNameField,
		common.BKOwnerIDField}
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

	if err := addDefaultBiz(ctx, db, conf); err != nil {
		return err
	}
	return nil
}
