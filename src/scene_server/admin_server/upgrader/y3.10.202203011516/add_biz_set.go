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

package y3_10_202203011516

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
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
		Creator:       common.CCSystemOperatorUserName,
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
		Creator:       common.CCSystemOperatorUserName,
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
		Creator:       common.CCSystemOperatorUserName,
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
		Creator:       common.CCSystemOperatorUserName,
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
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
		Placeholder:   "业务集所包含的业务的条件",
	},
}

const (
	// Business set initial ID
	bizSetInitialID = 10000000
)

func addBizSetObjectRow(ctx context.Context, db dal.RDB, ownerID string) error {

	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDBizSet}
	model := new(metadata.Object)

	// 判断是否有 BKInnerObjIDBizSet 的对象表，如果没有需要初始化
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
		blog.Errorf("the model %s already exists, but does not conform to the specification, object name: %s, "+
			"creator: %s, issue is #5923.", common.BKInnerObjIDBizSet, model.ObjectName, model.Creator)
		return fmt.Errorf("model %s failed to create", common.BKInnerObjIDBizSet)
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
	uniqueKeys := []string{common.BKObjIDField}
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
		common.BKPropertyIDField: mapstr.MapStr{
			common.BKDBIN: []string{common.BKBizSetNameField, common.BKBizSetIDField},
		},
	}
	if err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &attrs); err != nil {
		return err
	}

	// 需要判断 cc_ObjectUnique 中是否有关于业务集的值
	uniqueIdxs := make([]metadata.ObjectUnique, 0)
	condObjUnique := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDBizSet}

	if err := db.Table(common.BKTableNameObjUnique).Find(condObjUnique).Fields(common.BKObjectUniqueKeys).
		All(ctx, &uniqueIdxs); err != nil {
		return err
	}

	// if it exists, you need to determine whether the id is legal.
	if len(uniqueIdxs) > 0 {
		if len(uniqueIdxs) != len(attrs) {
			return errors.New("inconsistent number of unique indexes")
		}

		keysId := make(map[uint64]struct{})
		for _, index := range uniqueIdxs {
			for _, id := range index.Keys {
				keysId[id.ID] = struct{}{}
			}
		}
		// to prevent compliance with index scenarios, the number of indexes must be exactly the same.
		if len(keysId) != len(attrs) {
			return errors.New("invalid number of index ids")
		}
		for _, attr := range attrs {
			if _, ok := keysId[uint64(attr.ID)]; !ok {
				return errors.New("attr ID does not exist")
			}
		}
		return nil
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

// Idgen TODO
type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

func addBizSetIDToIdgenerator(ctx context.Context, db dal.RDB) error {

	// 1、find out whether there is data whose id is BKTableNameBaseBizSet in the cc_idgenerator table. If there is,
	// an error will be reported directly, which needs to be processed manually.
	filter := map[string]interface{}{
		"_id": common.BKTableNameBaseBizSet,
	}

	idGenerator := new(Idgen)
	err := db.Table(common.BKTableNameIDgenerator).Find(filter).Fields("SequenceID").One(ctx, idGenerator)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("count cc_BizSetBase id failed, err: %v", err)
		return err
	}
	if err != nil && db.IsNotFoundError(err) {
		// set the initialization value to 10000000
		data := map[string]interface{}{
			"_id":                  common.BKTableNameBaseBizSet,
			"SequenceID":           bizSetInitialID,
			common.CreateTimeField: time.Now(),
			common.LastTimeField:   time.Now(),
		}
		err = db.Table(common.BKTableNameIDgenerator).Insert(ctx, data)
		if nil != err {
			blog.Errorf("add data fail, error %s", err)
			return err
		}
		return nil
	}

	// Illegal if num is between 0 and bizSetInitialID
	if idGenerator.SequenceID < bizSetInitialID {
		blog.Errorf("cc_BizSetBase id should not exist, upgrade failed.")
		return errors.New("cc_BizSetBase id should not exist")
	}

	// If greater or equal to bizSetInitialID, it is considered that a legal business set has been created
	return nil
}

func addBizSetCollection(ctx context.Context, db dal.RDB) error {

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
	if err := db.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKPropertyIDField, common.BKPropertyNameField,
		common.CreatorField).All(ctx, &attrs); err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find object attribute describe fail, err: %v", err)
		return err
	}

	if len(attrs) > 0 {
		// 如果存在的话，数量必须一致。并且必须严格校验每个属性bk_property_name和creator必须完全一致，不一致直接报错需要先处理完毕后再升级
		if len(attrs) != len(dataRows) {
			blog.Errorf("Illegal number of business set model attributes, num is: %d", len(attrs))
			return errors.New("illegal number of business set model attributes")
		}

		for _, attr := range attrs {
			if data, ok := dataRows[attr.PropertyID]; ok {
				if attr.PropertyName != data.PropertyName || attr.Creator != data.Creator {
					blog.Errorf("the model biz set attribute %s already exists, but is illegal, name: %v, creator: %v",
						attr.PropertyID, attr.PropertyName, attr.Creator)
					return errors.New("model biz set attribute is invalid")
				}
			}
		}
		return nil
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}

	nowTime := metadata.Now()
	for _, row := range dataRows {
		row.OwnerID = ownerID
		row.IsPre = true
		row.IsReadOnly = false
		row.CreateTime = &nowTime
		row.LastTime = &nowTime
		row.Description = ""
		_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, row, "id", uniqueFields, []string{})
		if err != nil {
			blog.Errorf("add biz set attr failed, attribute: %v, err: %v", row, err)
			return err
		}
	}
	return nil
}

// addBizSetTableIndexes add indexes for common audit log query params
func addBizSetTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "biz_set_id",
			Keys:       bson.D{{common.BKBizSetIDField, 1}},
			Background: true,
			Unique:     true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "biz_set_name",
			Keys:       bson.D{{common.BKBizSetNameField, 1}},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_set_id_biz_set_name_owner_id",
			Keys: bson.D{
				{common.BKBizSetIDField, 1},
				{common.BKBizSetNameField, 1},
				{common.BKOwnerIDField, 1},
			},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(common.BKTableNameBaseBizSet).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for biz set table failed, err: %v", err)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		// skip the default "_id" index for the database
		if index.Name == "_id_" {
			continue
		}
		existIdxMap[index.Name] = true
	}

	// the number of indexes is not as expected.
	if len(existIdxMap) != 0 && (len(existIdxMap) < len(indexes)) {
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(common.BKTableNameBaseBizSet).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for biz set table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}

func addBizSetPropertyOption(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

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

	if err := addBizSetCollection(ctx, db); err != nil {
		return err
	}

	if err := addBizSetIDToIdgenerator(ctx, db); err != nil {
		return err
	}

	if err := addBizSetTableIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}
