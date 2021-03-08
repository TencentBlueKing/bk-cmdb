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
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func splitTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	objs := make([]object, 0)
	if err = db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre).All(ctx, &objs); err != nil {
		blog.Errorf("list all obect id from db error. err: %s", err.Error())
		return
	}

	instIdxs := instanceDefaultIndex
	instAsstIdxs := assoicationDefaultIndex

	instTablePrefix := common.BKTableNameBaseInst + "_pub_"
	instAsstTablePrefix := common.BKTableNameInstAsst + "_pub_"
	for _, obj := range objs {
		objInstTable := instTablePrefix + obj.ObjectID
		objInstAsstTable := instAsstTablePrefix + obj.ObjectID

		if err = createTableFunc(ctx, objInstAsstTable, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table error. err: %s", obj.ObjectID, err.Error())
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
		if err = createTableIndex(ctx, objInstAsstTable, instAsstIdxs, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table index error. err: %s", obj.ObjectID, err.Error())
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
			if err := db.Table(tableName).Insert(ctx, inst); err != nil {
				return fmt.Errorf("insert obj(%s) inst  error. err: %s", objID, err.Error())
			}
		}
		time.Sleep(time.Millisecond * 10)
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
	alreadyIdxNameMap, dbIndexList, err := getTableAlreadyIndexs(ctx, tableName, db)
	if err != nil {
		return err
	}

	for _, idx := range idxs {
		_, idExist := idx.Keys["_id"]
		if len(idx.Keys) == 1 && idExist {
			// create table already exist
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

func getTableAlreadyIndexs(ctx context.Context, tableName string,
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

	alreadyIdxNameMap, _, err := getTableAlreadyIndexs(ctx, tableName, db)
	if err != nil {
		return err
	}

	// 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
	propertyIDTypeRelation, err := findObjAttrsIDRelation(ctx, objID, db)
	if err != nil {
		return err
	}

	for _, idx := range uniqueIdxs {
		newDBIndx, err := toDBUniqueIdx(idx, propertyIDTypeRelation)
		if err != nil {
			return fmt.Errorf("obj(%s). %s", objID, err.Error())
		}
		// 升级版本程序不考虑，数据不一致的情况
		if _, exists := alreadyIdxNameMap[newDBIndx.Name]; !exists {
			if err := db.Table(tableName).CreateIndex(ctx, newDBIndx); err != nil {
				return fmt.Errorf("create unique index(%s) error. err: %s", newDBIndx.Name, err.Error())
			}
		}
	}

	return nil
}

// 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
func findObjAttrsIDRelation(ctx context.Context, objID string, db dal.RDB) (map[int64]metadata.Attribute, error) {
	// 获取字段类型,只需要共有字段
	attrFilter := map[string]interface{}{
		common.BKObjIDField: objID,
		common.BKAppIDField: 0,
	}
	attrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID}
	if err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Fields(fields...).All(ctx, &attrs); err != nil {
		return nil, fmt.Errorf("get obj(%s) property error. err: %s", objID, err.Error())
	}

	attrIDMap := make(map[int64]metadata.Attribute, 0)
	for _, attr := range attrs {
		attrIDMap[attr.ID] = attr
	}

	return attrIDMap, nil
}

const CCLogicUniqueIdxNamePrefix = "bkcc_unique_"

func toDBUniqueIdx(idx objectUnique, attrIDMap map[int64]metadata.Attribute) (types.Index, error) {

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

var assoicationDefaultIndex = []types.Index{
	{
		Keys: map[string]int32{
			common.BKOwnerIDField: 1,
			common.BKInstIDField:  1,
		},
		Name:       "bkcc_idx_ObjID_InstID",
		Background: true,
	},
	{
		Keys: map[string]int32{
			common.BKFieldID: 1,
		},
		Name:       "bkcc_unique_ID",
		Unique:     true,
		Background: true,
	},
	{
		Keys: map[string]int32{
			common.BKAsstObjIDField:  1,
			common.BKAsstInstIDField: 1,
		},
		Name:       "bkcc_idx_AsstObjID_AsstInstID",
		Unique:     true,
		Background: true,
	},
}

var instanceDefaultIndex = []types.Index{
	{
		Keys: map[string]int32{
			common.BKObjIDField: 1,
		},
		Name:       "bkcc_idx_ObjID",
		Background: true,
	},
	{
		Keys: map[string]int32{
			common.BKOwnerIDField: 1,
		},
		Name:       "bkcc_idx_supplierAccount",
		Background: true,
	},
	{
		Keys: map[string]int32{
			common.BKInstIDField: 1,
		},
		Name:       "bkcc_idx_InstId",
		Background: true,
		Unique:     true,
	},
	{
		Keys: map[string]int32{
			common.BKInstNameField: 1,
		},
		Name:       "bkcc_idx_InstName",
		Background: true,
	},
}
