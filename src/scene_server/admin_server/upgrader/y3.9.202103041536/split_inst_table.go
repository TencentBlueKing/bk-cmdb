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

package y3_9_202103041536

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	commIdx "configcenter/src/common/index"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func splitTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	objs := make([]object, 0)
	if err = db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre).All(ctx, &objs); err != nil {
		blog.Errorf("list all object id from db error. err: %s", err.Error())
		return
	}

	instIdxs := instanceDefaultIndexes
	instAsstIdxs := associationDefaultIndexes

	instTablePrefix := common.BKTableNameBaseInst + "_pub_"
	instAsstTablePrefix := common.BKTableNameInstAsst + "_pub_"
	for _, obj := range objs {
		objInstTable := instTablePrefix + obj.ObjectID
		objInstAsstTable := instAsstTablePrefix + obj.ObjectID

		if err = createTableFunc(ctx, objInstAsstTable, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if err = createTableIndex(ctx, objInstAsstTable, instAsstIdxs, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table index error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if obj.IsPre {
			// 内置模型只创建关联关系表
			continue
		}

		if err = createTableFunc(ctx, objInstTable, db); err != nil {
			blog.Errorf("create obj(%s) inst table error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if err = createTableIndex(ctx, objInstTable, instIdxs, db); err != nil {
			blog.Errorf("create obj(%s) inst table index error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if err = createTableLogicUniqueIndex(ctx, obj.ObjectID, objInstTable, db); err != nil {
			blog.Errorf("create obj(%s) inst unique table index error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if err = splitInstTable(ctx, obj.ObjectID, objInstTable, db); err != nil {
			return
		}

	}

	return nil
}

func splitInstTable(ctx context.Context, objID string, tableName string, db dal.RDB) error {

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}
	otps := types.FindOpts{
		WithObjectID: true,
	}
	const pageSize = uint64(1000)
	start := uint64(0)
	query := db.Table(common.BKTableNameBaseInst).Find(filter, otps).Limit(pageSize)
	for {
		insts := make([]map[string]interface{}, pageSize)
		if err := query.Start(start).All(ctx, &insts); err != nil {
			return fmt.Errorf("find obj(%s) inst list error. err: %s", objID, err.Error())
		}
		if len(insts) == 0 {
			break
		}
		start += pageSize
		// 为了可以重复执行，这里没有采用批量插入数据的方法
		for _, inst := range insts {
			filter := map[string]interface{}{
				"_id": inst["_id"],
			}
			cnt, err := db.Table(tableName).Find(filter).Count(ctx)
			if err != nil {
				return fmt.Errorf("check obj(%s) inst exists error. err: %s", objID, err.Error())
			}
			if cnt > 0 {
				continue
			}
			// TODO: test 1kw rows
			if err := db.Table(tableName).Insert(ctx, inst); err != nil {
				return fmt.Errorf("insert obj(%s) inst  error. err: %s", objID, err.Error())
			}
		}
		//time.Sleep(time.Millisecond * 10)
	}

	return nil

}

func createTableFunc(ctx context.Context, tableName string, db dal.DB) error {
	exist, err := db.HasTable(ctx, tableName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return db.CreateTable(ctx, tableName)
}

func createTableIndex(ctx context.Context, tableName string, idxs []types.Index, db dal.RDB) error {
	alreadyIdxNameMap, dbIndexList, err := getTableAlreadyIndexes(ctx, tableName, db)
	if err != nil {
		return err
	}

	for _, idx := range idxs {
		_, idExist := idx.Keys["_id"]
		if len(idx.Keys) == 1 && idExist {
			// create table index already exist
			continue
		}
		alreadyIdx, ok := alreadyIdxNameMap[idx.Name]
		if ok {
			if commIdx.IndexEqual(idx, alreadyIdx) {
				continue
			}
		}
		alreadyIdx, ok = commIdx.FindIndexByIndexFields(idx.Keys, dbIndexList)
		if ok {
			if commIdx.IndexEqual(idx, alreadyIdx) {
				continue
			}
		}
		if err := db.Table(tableName).CreateIndex(ctx, idx); err != nil {
			return fmt.Errorf("create index(%s) error. err: %s", idx.Name, err.Error())
		}
	}

	return nil
}

func getTableAlreadyIndexes(ctx context.Context, tableName string,
	db dal.RDB) (map[string]types.Index, []types.Index, error) {
	alreadyIdxs, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get collection(%s) index error. err: %s", tableName, err.Error())
	}
	alreadyIdxNameMap := make(map[string]types.Index, len(alreadyIdxs))
	for _, idx := range alreadyIdxs {
		alreadyIdxNameMap[idx.Name] = idx
	}
	return alreadyIdxNameMap, alreadyIdxs, nil
}

func createTableLogicUniqueIndex(ctx context.Context, objID string, tableName string, db dal.RDB) error {

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}
	uniqueIdxs := make([]objectUnique, 0)
	if err := db.Table(common.BKTableNameObjUnique).Find(filter).All(ctx, &uniqueIdxs); err != nil {
		return fmt.Errorf("get obj(%s) logic unique index error. err: %s", objID, err.Error())
	}

	alreadyIdxNameMap, _, err := getTableAlreadyIndexes(ctx, tableName, db)
	if err != nil {
		return err
	}

	// 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
	propertyIDTypeRelation, err := findObjAttrsIDRelation(ctx, objID, db)
	if err != nil {
		return err
	}

	for _, idx := range uniqueIdxs {
		newDBIndex, err := toDBUniqueIdx(idx, propertyIDTypeRelation)
		if err != nil {
			return fmt.Errorf("obj(%s). %s", objID, err.Error())
		}
		// 升级版本程序不考虑，数据不一致的情况
		if _, exists := alreadyIdxNameMap[newDBIndex.Name]; !exists {
			if err := db.Table(tableName).CreateIndex(ctx, newDBIndex); err != nil {
				return fmt.Errorf("create unique index(%s) error. err: %s", newDBIndex.Name, err.Error())
			}
		}
	}

	return nil
}

// 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
func findObjAttrsIDRelation(ctx context.Context, objID string, db dal.RDB) (map[int64]Attribute, error) {
	// 获取字段类型,只需要共有字段
	attrFilter := map[string]interface{}{
		common.BKObjIDField: objID,
		common.BKAppIDField: 0,
	}
	attrs := make([]Attribute, 0)
	fields := []string{common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID}
	if err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Fields(fields...).All(ctx, &attrs); err != nil {
		return nil, fmt.Errorf("get obj(%s) property error. err: %s", objID, err.Error())
	}

	attrIDMap := make(map[int64]Attribute, 0)
	for _, attr := range attrs {
		attrIDMap[attr.ID] = attr
	}

	return attrIDMap, nil
}

const CCLogicUniqueIdxNamePrefix = "bkcc_unique_"

func toDBUniqueIdx(idx objectUnique, attrIDMap map[int64]Attribute) (types.Index, error) {

	dbIdx := types.Index{
		Name:                    fmt.Sprintf("%s%d", CCLogicUniqueIdxNamePrefix, idx.ID),
		Unique:                  true,
		Background:              true,
		Keys:                    make(map[string]int32, len(idx.Keys)),
		PartialFilterExpression: make(map[string]interface{}, len(idx.Keys)),
	}

	// attrIDMap数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
	for _, key := range idx.Keys {
		attr := attrIDMap[int64(key.ID)]
		dbType := convFieldTypeToDBType(attr.PropertyType)
		if dbType == "" {
			blog.ErrorJSON("build unique index property id: %s type not support.", key.Kind)
			return dbIdx, fmt.Errorf("build unique index property(%s) type not support.", key.Kind)
		}
		dbIdx.Keys[attr.PropertyID] = 1
		dbIdx.PartialFilterExpression[attr.PropertyID] = map[string]interface{}{common.BKDBType: dbType}
	}

	return dbIdx, nil
}

type object struct {
	ID int64 `field:"id" json:"id" bson:"id"`

	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	ObjectName string `field:"bk_obj_name" json:"bk_obj_name" bson:"bk_obj_name"`

	// IsHidden front-end don't display the object if IsHidden is true
	IsHidden bool `field:"bk_ishidden" json:"bk_ishidden" bson:"bk_ishidden"`

	IsPre    bool `field:"ispre" json:"ispre" bson:"ispre"`
	IsPaused bool `field:"bk_ispaused" json:"bk_ispaused" bson:"bk_ispaused"`
}

type objectUnique struct {
	ID        uint64      `json:"id" bson:"id"`
	ObjID     string      `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool        `json:"must_check" bson:"must_check"`
	Keys      []uniqueKey `json:"keys" bson:"keys"`
	Ispre     bool        `json:"ispre" bson:"ispre"`
	OwnerID   string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  *time.Time  `json:"last_time" bson:"last_time"`
}

type uniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}

func convFieldTypeToDBType(typ string) string {
	switch typ {
	case FieldTypeSingleChar, FieldTypeLongChar:
		return "string"
	case FieldTypeInt, FieldTypeFloat, FieldTypeEnum, FieldTypeUser, FieldTypeTimeZone,
		FieldTypeList, FieldTypeOrganization:
		return "number"
	case FieldTypeDate, FieldTypeTime:
		return "date"
	case FieldTypeBool:
		return "bool"

	}

	// other type not support
	return ""
}

const (
	// FieldTypeSingleChar the single char filed type
	FieldTypeSingleChar string = "singlechar"

	// FieldTypeLongChar the long char field type
	FieldTypeLongChar string = "longchar"

	// FieldTypeInt the int field type
	FieldTypeInt string = "int"

	// FieldTypeFloat the float field type
	FieldTypeFloat string = "float"

	// FieldTypeEnum the enum field type
	FieldTypeEnum string = "enum"

	// FieldTypeDate the date field type
	FieldTypeDate string = "date"

	// FieldTypeTime the time field type
	FieldTypeTime string = "time"

	// FieldTypeUser the user field type
	FieldTypeUser string = "objuser"

	// FieldTypeTimeZone the timezone field type
	FieldTypeTimeZone string = "timezone"

	// FieldTypeBool the bool type
	FieldTypeBool string = "bool"

	// FieldTypeList the list type
	FieldTypeList string = "list"

	// FieldTypeOrganization the organization field type
	FieldTypeOrganization string = "organization"
)

var associationDefaultIndexes = []types.Index{
	{
		Name: common.CCLogicIndexNamePrefix + "bkObjId_bkInstID",
		Keys: map[string]int32{
			"bk_obj_id":  1,
			"bk_inst_id": 1,
		},
		Background: true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "id",
		Keys: map[string]int32{
			"id": 1,
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstId_bkObjId",
		Keys: map[string]int32{
			"bk_inst_id": 1,
			"bk_obj_id":  1,
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkAsstObjId_bkAsstInstId",
		Keys: map[string]int32{
			"bk_asst_obj_id":  1,
			"bk_asst_inst_id": 1,
		},
		Background: true,
	},
}

var instanceDefaultIndexes = []types.Index{
	{
		Name: common.CCLogicIndexNamePrefix + "bkObjId",
		Keys: map[string]int32{
			"bk_obj_id": 1,
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkSupplierAccount",
		Keys: map[string]int32{
			"bk_supplier_account": 1,
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstId",
		Keys: map[string]int32{
			"bk_inst_id": 1,
		},
		Background: true,
		// 新加 2021年03月11日
		Unique: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstName",
		Keys: map[string]int32{
			"bk_inst_name": 1,
		},
		Background: false,
	},
}

// Attribute attribute metadata definition
type Attribute struct {
	BizID             int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64       `field:"id" json:"id" bson:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option"`
	Description       string      `field:"description" json:"description" bson:"description"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
}
