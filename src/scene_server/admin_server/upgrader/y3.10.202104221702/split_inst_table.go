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

package y3_10_202104221702

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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	oldInstTable        = "cc_ObjectBase"
	oldInstAsstTable    = "cc_InstAsst"
	instAsstTableFormat = "cc_InstAsst_%v_pub_%v"
	objectBaseMapping   = "cc_ObjectBaseMapping"
	maxWorkNumber       = 200
	pageSize            = uint64(5000)
)

func splitTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	objs := make([]object, 0)
	if err = db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre, common.BKOwnerIDField).All(ctx, &objs); err != nil {
		blog.Errorf("list all object id from db error. err: %s", err.Error())
		return
	}

	instIdxs := instanceDefaultIndexes
	instAsstIdxs := associationDefaultIndexes

	var objectIDs []string
	for _, obj := range objs {
		objInstTable := buildInstTableName(obj.ObjectID, obj.OwnerID)         // instTablePrefix + obj.ObjectID
		objInstAsstTable := buildInstAsstTableName(obj.ObjectID, obj.OwnerID) // instAsstTablePrefix + obj.ObjectID

		objectIDs = append(objectIDs, obj.ObjectID)
		if err = createTableFunc(ctx, objInstAsstTable, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if err = createTableIndex(ctx, objInstAsstTable, instAsstIdxs, db); err != nil {
			blog.Errorf("create obj(%s) inst asst table index error. err: %s", obj.ObjectID, err.Error())
			return
		}

		if obj.IsPre {
			objInstTable := innerObjIDTableNameRelation[obj.ObjectID]
			if objInstTable != "" {
				if err = createTableLogicUniqueIndex(ctx, obj.ObjectID, objInstTable, db); err != nil {
					blog.Errorf("create obj(%s) inst unique table index error. err: %s", obj.ObjectID, err.Error())
					return
				}
			}

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

	}

	blog.Info("start copy instance to sharding table")
	if err = splitInstTable(ctx, db); err != nil {
		return err
	}

	blog.Info("start copy instance association to sharding table")
	if err = splitInstAsstTable(ctx, db); err != nil {
		return err
	}

	return nil
}

func initWorkerChn(workerNum int) chan struct{} {
	workerNumCtlChn := make(chan struct{}, workerNum)
	for idx := 0; idx < workerNum; idx++ {
		workerNumCtlChn <- struct{}{}
	}
	return workerNumCtlChn
}

func splitInstAsstTable(ctx context.Context, db dal.RDB) error {
	filter := map[string]interface{}{}
	opts := types.NewFindOpts().SetWithObjectID(true)
	start := uint64(0)
	query := db.Table(oldInstAsstTable).Find(filter, opts).Limit(pageSize).Sort("id")
	workerChn := initWorkerChn(maxWorkNumber)
	errChn := make(chan error, maxWorkNumber)
	for {
		blog.Infof("find instance association detail info list. start: %d", start)
		asstList := make([]map[string]interface{}, pageSize)
		if err := query.Start(start).All(ctx, &asstList); err != nil {
			return fmt.Errorf("find inst association list error. err: %s", err.Error())
		}

		if len(asstList) == 0 {
			// 没有数据了，
			break
		}

		for _, asst := range asstList {
			if len(errChn) > 0 {
				return <-errChn
			}
			<-workerChn
			go func(copyAsst map[string]interface{}) {
				defer func() {
					workerChn <- struct{}{}
				}()
				if err := copyInstanceAssociationToShardingTable(ctx, copyAsst, db); err != nil {
					errChn <- err
				}
			}(asst)
		}
		start += pageSize

	}

	// 等待所有任务的结束
	for idx := 0; idx < maxWorkNumber; idx++ {
		<-workerChn
	}

	// 避免最后一个任务的时候出现错误
	if len(errChn) > 0 {
		return <-errChn
	}

	return nil

}

func buildInstTableName(objID, supplierAccount interface{}) string {
	return fmt.Sprintf("cc_ObjectBase_%v_pub_%v", supplierAccount, objID)
}

func buildInstAsstTableName(objID, supplierAccount interface{}) string {
	return fmt.Sprintf("cc_InstAsst_%v_pub_%v", supplierAccount, objID)
}

func copyInstanceAssociationToShardingTable(ctx context.Context, association map[string]interface{}, db dal.RDB) error {

	objTableName := buildInstAsstTableName(association[common.BKObjIDField], association[common.BKOwnerIDField])
	asstObjTableName := buildInstAsstTableName(association[common.BKAsstObjIDField], association[common.BKOwnerIDField])

	filter := map[string]interface{}{
		"_id": association["_id"],
	}

	if err := db.Table(objTableName).Upsert(ctx, filter, association); err != nil {
		blog.ErrorJSON("copy instance association error. association: %s, err: %s", association, err.Error())
		return fmt.Errorf("copy instance association error . err: %s", err.Error())
	}
	if err := db.Table(asstObjTableName).Upsert(ctx, filter, association); err != nil {
		blog.ErrorJSON("copy instance association error. association: %s, err: %s", association, err.Error())
		return fmt.Errorf("copy instance association error . err: %s", err.Error())
	}

	return nil

}

func splitInstTable(ctx context.Context, db dal.RDB) error {
	filter := map[string]interface{}{}
	opts := types.NewFindOpts().SetWithObjectID(true)
	start := uint64(0)
	query := db.Table(oldInstTable).Find(filter, opts).Limit(pageSize).Sort("bk_inst_id")

	workerChn := initWorkerChn(maxWorkNumber)
	errChn := make(chan error, maxWorkNumber)

	for {
		blog.Infof("find instance detail info list. start: %d", start)
		insts := make([]map[string]interface{}, pageSize)
		if err := query.Start(start).All(ctx, &insts); err != nil {
			return fmt.Errorf("find inst list error. err: %s", err.Error())
		}
		if len(insts) == 0 {
			// 没有数据了
			break
		}
		for _, inst := range insts {
			if len(errChn) > 0 {
				return <-errChn
			}
			<-workerChn
			go func(copyInst map[string]interface{}) {
				defer func() {
					workerChn <- struct{}{}
				}()
				if err := copyInstanceToShardingTable(ctx, copyInst, db); err != nil {
					errChn <- err
				}
			}(inst)
		}
		start += pageSize

	}

	// 等待所有任务的结束
	for idx := 0; idx < maxWorkNumber; idx++ {
		<-workerChn
	}

	// 避免最后一个任务的时候出现错误
	if len(errChn) > 0 {
		return <-errChn
	}
	return nil
}

func copyInstanceToShardingTable(ctx context.Context, inst map[string]interface{}, db dal.RDB) error {

	objID := fmt.Sprintf("%v", inst[common.BKObjIDField])
	tableName := buildInstTableName(objID, inst[common.BKOwnerIDField]) // instTablePrefix + objID

	mappingFilter := map[string]interface{}{
		common.BKInstIDField: inst[common.BKInstIDField],
	}
	doc := map[string]interface{}{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: inst[common.BKOwnerIDField],
	}

	if err := db.Table(objectBaseMapping).Upsert(ctx, mappingFilter, doc); err != nil {
		return fmt.Errorf("upsert instance id and object id mapping error. row object id: %v, err: %s",
			inst["_id"], err.Error())
	}

	filter := map[string]interface{}{
		"_id": inst["_id"],
	}
	if err := db.Table(tableName).Upsert(ctx, filter, inst); err != nil {
		return fmt.Errorf("insert obj(%s) inst  error. err: %s", objID, err.Error())
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

	isIDUniqueFunc := func(idx types.Index) bool {
		if len(idx.Keys) != 1 {
			return false
		}

		idxKeyMap := idx.Keys.Map()
		if _, exist := idxKeyMap["id"]; exist {
			return true
		}
		return false
	}

	for _, idx := range idxs {
		idxKeyMap := idx.Keys.Map()
		_, idExist := idxKeyMap["_id"]
		if len(idx.Keys) == 1 && idExist {
			// create table index already exist
			continue
		}
		alreadyIdx, ok := alreadyIdxNameMap[idx.Name]
		if ok {
			if commIdx.IndexEqual(idx, alreadyIdx) {
				continue
			} else {
				if err := db.Table(tableName).DropIndex(ctx, idx.Name); err != nil {
					return fmt.Errorf("delete index(%s) error. err: %s", idx.Name, err.Error())
				}
			}
		}
		alreadyIdx, ok = FindIndexByIndexFields(idx.Keys, dbIndexList)
		if ok {

			if isIDUniqueFunc(idx) {
				continue
			} else {
				if err := db.Table(tableName).DropIndex(ctx, alreadyIdx.Name); err != nil {
					return fmt.Errorf("delete index(%s) error. err: %s", alreadyIdx.Name, err.Error())
				}
			}
		}

		if err := db.Table(tableName).CreateIndex(ctx, idx); err != nil && !db.IsDuplicatedError(err) {
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

	alreadyIdxNameMap, dbTableIndexes, err := getTableAlreadyIndexes(ctx, tableName, db)
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
		blog.InfoJSON("create unique index table, table: %s, index: %s", tableName, newDBIndex)

		dbTableIndex, exist := FindIndexByIndexFields(newDBIndex.Keys, dbTableIndexes)
		if exist {
			if dbTableIndex.Name == newDBIndex.Name {
				continue
			}
			if err := db.Table(tableName).DropIndex(ctx, dbTableIndex.Name); err != nil {
				return fmt.Errorf("delete unique index(%s) error. err: %s", dbTableIndex.Name, err.Error())
			}
		}

		// 升级版本程序不考虑，数据不一致的情况
		if _, exists := alreadyIdxNameMap[newDBIndex.Name]; !exists {
			if err := db.Table(tableName).CreateIndex(ctx, newDBIndex); err != nil && !db.IsDuplicatedError(err) {
				return fmt.Errorf("create unique index(%s) error. err: %s", newDBIndex.Name, err.Error())
			}
		}
	}

	return nil
}

// findObjAttrsIDRelation 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
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

// CCLogicUniqueIdxNamePrefix TODO
const CCLogicUniqueIdxNamePrefix = "bkcc_unique_"

func toDBUniqueIdx(idx objectUnique, attrIDMap map[int64]Attribute) (types.Index, error) {

	dbIdx := types.Index{
		Name:                    fmt.Sprintf("%s%d", CCLogicUniqueIdxNamePrefix, idx.ID),
		Unique:                  true,
		Background:              true,
		Keys:                    make(bson.D, 0),
		PartialFilterExpression: make(map[string]interface{}, len(idx.Keys)),
	}

	// attrIDMap数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
	for _, key := range idx.Keys {
		attr := attrIDMap[int64(key.ID)]
		if idx.ObjID == common.BKInnerObjIDHost && attr.PropertyID == common.BKCloudIDField {
			// NOTEICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未foreignkey 特殊的类型
			attr.PropertyType = common.FieldTypeInt
		}
		if idx.ObjID == common.BKInnerObjIDHost &&
			(attr.PropertyID == common.BKHostInnerIPField || attr.PropertyID == common.BKHostOuterIPField ||
				attr.PropertyID == common.BKOperatorField || attr.PropertyID == common.BKBakOperatorField) {
			// NOTEICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未innerIP,OuterIP 特殊的类型
			attr.PropertyType = common.FieldTypeList
		}
		dbType := convFieldTypeToDBType(attr.PropertyType)
		if dbType == "" {
			blog.ErrorJSON("build unique index property id: %s type: %s not support.", key.Kind, attr.PropertyType)
			return dbIdx, fmt.Errorf("build unique index property(%s) type(%s) not support.",
				key.Kind, attr.PropertyType)
		}
		dbIdx.Keys = append(dbIdx.Keys, primitive.E{
			Key:   attr.PropertyID,
			Value: 1,
		})
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

	IsPre    bool   `field:"ispre" json:"ispre" bson:"ispre"`
	IsPaused bool   `field:"bk_ispaused" json:"bk_ispaused" bson:"bk_ispaused"`
	OwnerID  string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
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
	case FieldTypeSingleChar, FieldTypeEnum, FieldTypeDate, FieldTypeList:
		return "string"
	case FieldTypeInt, FieldTypeFloat:
		return "number"
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
		Keys: bson.D{
			{"bk_obj_id", 1},
			{"bk_inst_id", 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "id",
		Keys: bson.D{
			{"id", 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstId_bkObjId",
		Keys: bson.D{
			{"bk_inst_id", 1},
			{"bk_obj_id", 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkAsstObjId_bkAsstInstId",
		Keys: bson.D{
			{"bk_asst_obj_id", 1},
			{"bk_asst_inst_id", 1},
		},
		Background: true,
	},
}

var instanceDefaultIndexes = []types.Index{
	{
		Name: common.CCLogicIndexNamePrefix + "bkObjId",
		Keys: bson.D{
			{"bk_obj_id", 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkSupplierAccount",
		Keys: bson.D{
			{"bk_supplier_account", 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstId",
		Keys: bson.D{
			{"bk_inst_id", 1},
		},
		Background: true,
		// 新加 2021年03月11日
		Unique: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkInstName",
		Keys: bson.D{
			{"bk_inst_name", 1},
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

// FindIndexByIndexFields find index by index fields
func FindIndexByIndexFields(keys bson.D, indexList []types.Index) (dbIndex types.Index, exists bool) {
	targetIdxMap := keys.Map()
	for _, idx := range indexList {
		idxMap := idx.Keys.Map()
		if len(targetIdxMap) != len(idxMap) {
			continue
		}
		exists = true
		for key := range idxMap {
			if _, keyExists := targetIdxMap[key]; !keyExists {
				exists = false
				break
			}
		}
		if exists {
			return idx, exists
		}

	}

	return types.Index{}, false
}
