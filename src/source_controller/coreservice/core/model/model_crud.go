/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"fmt"
	"reflect"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	dbindex "configcenter/src/common/index"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelManager) count(kit *rest.Kit, cond universalsql.Condition) (uint64, error) {

	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation by the condition (%#v), error info is %s", kit.Rid, cond.ToMapStr(), err.Error())
		return 0, kit.CCError.Errorf(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) save(kit *rest.Kit, model *metadata.Object) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("request(%s): it is failed to make sequence id on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjDes, err.Error())
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	sortNum, err := m.GetModelLastNum(kit, *model)
	if err != nil {
		blog.Errorf("set object sort number failed, err: %v, objectId: %s, rid: %s", err, model.ObjectID, kit.Rid)
		return id, err
	}

	model.ObjSortNumber = sortNum
	model.ID = int64(id)
	model.OwnerID = kit.SupplierAccount

	if nil == model.LastTime {
		model.LastTime = &metadata.Time{}
		model.LastTime.Time = time.Now()
	}
	if nil == model.CreateTime {
		model.CreateTime = &metadata.Time{}
		model.CreateTime.Time = time.Now()
	}

	err = mongodb.Client().Table(common.BKTableNameObjDes).Insert(kit.Ctx, model)
	return id, err
}

// GetModelLastNum 获取模型排序字段
func (m *modelManager) GetModelLastNum(kit *rest.Kit, model metadata.Object) (int64, error) {
	//查询当前分组下的模型信息
	modelInput := map[string]interface{}{metadata.ModelFieldObjCls: model.ObjCls}
	modelResult := make([]metadata.Object, 10)
	sortCond := "-obj_sort_number"
	if findErr := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelInput).Sort(sortCond).
		Fields(metadata.ModelFieldID, metadata.ModelFieldObjSortNumber).All(kit.Ctx, &modelResult); findErr != nil {
		blog.Error("get object sort number failed, database operation is failed, err: %v, rid: %s", findErr, kit.Rid)
		return 0, findErr
	}

	switch {
	case model.ObjSortNumber == 0 && len(modelResult) <= 0:
		return 0, nil

	case model.ObjSortNumber == 0 && len(modelResult) > 0:
		return modelResult[0].ObjSortNumber + 1, nil

	case model.ObjSortNumber > 0 && len(modelResult) <= 0:
		return model.ObjSortNumber, nil

	case model.ObjSortNumber > 0 && len(modelResult) > 0:
		if updateErr := m.valueAdd(kit, model, modelResult); updateErr != nil {
			blog.Errorf("update object sort number failed, err: %v, ctx:%v, rid: %s", updateErr, kit.Rid)
			return 0, updateErr
		}
		return model.ObjSortNumber, nil

	default:
		// 按逻辑不应触达此处
		return 0, kit.CCError.Error(common.CCErrCommParamsIsInvalid)
	}
}

// valueAdd 大于等于model值的obj_sort_number字段值加一
func (m *modelManager) valueAdd(kit *rest.Kit, model metadata.Object, modelArr []metadata.Object) error {
	for _, mod := range modelArr {
		if mod.ObjSortNumber < model.ObjSortNumber {
			continue
		}
		updateFilter := map[string]interface{}{metadata.ModelFieldID: mod.ID}
		updateData := map[string]interface{}{metadata.ModelFieldObjSortNumber: mod.ObjSortNumber + 1}
		updateErr := mongodb.Client().Table(common.BKTableNameObjDes).Update(kit.Ctx, updateFilter, updateData)
		if updateErr != nil {
			blog.Errorf("update object sort number failed, err: %v, ctx:%v, rid: %s",
				updateErr, updateFilter, kit.Rid)
			return updateErr
		}
	}
	return nil
}

func (m *modelManager) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	data.Set(metadata.ModelFieldLastTime, time.Now())
	models := make([]metadata.Object, 0)
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx, &models)
	if err != nil {
		blog.Errorf("find models failed, err: %s, filter: %+v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	if err := m.setUpdateObjectSortNumber(kit, &data); err != nil {
		blog.Errorf("set object sort number failed, err: %v, data: %s, rid: %s", err, data, kit.Rid)
		return 0, err
	}

	// remove unchangeable fields.
	data.Remove(metadata.ModelFieldObjectID)
	data.Remove(metadata.ModelFieldID)

	// 停用模型 pausedFlag 为 true
	pausedFlag := false

	paused, exist := data[metadata.ModelFieldIsPaused]
	if exist {
		flag, ok := paused.(bool)
		if exist && !ok {
			blog.Errorf("attr(%v) type error, type: %v, rid: %s", metadata.ModelFieldIsPaused,
				reflect.TypeOf(paused), kit.Rid)
			return 0, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ModelFieldIsPaused)
		}
		pausedFlag = flag
	}

	objName, objNameExist := data[common.BKObjNameField]

	if (objNameExist && len(util.GetStrByInterface(objName)) > 0) || pausedFlag {
		for _, model := range models {
			if err := m.isExistProcessingTask(kit, model.ID, pausedFlag); err != nil {
				return 0, err
			}

			// 检查模型名称重复
			nameCond := map[string]interface{}{
				common.BKObjNameField: objName,
				common.BKFieldID: map[string]interface{}{
					common.BKDBNE: model.ID,
				},
			}

			count, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(nameCond).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("failed to check the validity of the model name, filter: %+v, err: %v, rid: %s",
					nameCond, err, kit.Rid)
				return 0, err
			}
			if count > 0 {
				blog.Warnf("update model failed, field `%s` duplicated, rid: %s", objName, kit.Rid)
				return 0, kit.CCError.Errorf(common.CCErrCommDuplicateItem, objName)
			}
			// 一次更新多个模型的时候，唯一校验需要特别小心
			filter := map[string]interface{}{common.BKFieldID: model.ID}
			cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, filter, data)
			if err != nil {
				blog.Errorf("failed to update table (%s), err: %v, cond: %+v, data: %+v, err: %v,rid: %s",
					common.BKTableNameObjDes, filter, data, err, kit.Rid)
				return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
			}
			return cnt, nil
		}
	}

	cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.Errorf("failed to update the table (%s), err: %s, rid: %s", common.BKTableNameObjDes, err, kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

// setUpdateObjectSortNumber 根据更新条件设置 obj_sort_number 值
func (m *modelManager) setUpdateObjectSortNumber(kit *rest.Kit, data *mapstr.MapStr) error {
	if !data.Exists(metadata.ModelFieldObjSortNumber) && !data.Exists(metadata.ModelFieldObjCls) {
		return nil
	}

	object := metadata.Object{}
	if err := data.ToStructByTag(&object, "field"); err != nil {
		blog.Errorf("parsing data failed, err: %v, data: %v, rid: %s", err, data, kit.Rid)
		return err
	}
	//如果传递了 bk_classification_id 字段,按照更新模型所属分组处理
	if data.Exists(metadata.ModelFieldObjCls) {
		sortNum, err := m.GetModelLastNum(kit, object)
		if err != nil {
			blog.Errorf("set object sort number failed, err: %v, objectId: %s, rid: %s", err, model.ObjectID, kit.Rid)
			return err
		}
		data.Set(metadata.ModelFieldObjSortNumber, sortNum)
		return nil
	}

	//如果未传递了 bk_classification_id 字段，传递了 obj_sort_number 字段,则表示在当前分组下更新模型顺序
	//更新当前模型 obj_sort_number 前先更新当前分组下其它模型 obj_sort_number
	if object.ObjSortNumber < 0 {
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	//查询当前模型的分组信息
	clsInput := map[string]interface{}{metadata.ModelFieldID: object.ID}
	clsResult := make([]metadata.Object, 0)
	if findErr := mongodb.Client().Table(common.BKTableNameObjDes).Find(clsInput).Fields(metadata.ModelFieldObjCls).
		All(kit.Ctx, &clsResult); findErr != nil {
		blog.Error("get object classification failed, err: %v, objID: %d, rid: %s", findErr, object.ID, kit.Rid)
		return findErr
	}
	if len(clsResult) <= 0 {
		blog.Errorf("no model classification id founded, err: model classification no founded, objID: %d, rid: %s",
			object.ID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorModelClassificationNotFound)
	}

	//查询当前分组下的所有模型信息
	modelInput := map[string]interface{}{metadata.ModelFieldObjCls: clsResult[0].ObjCls}
	modelResult := make([]metadata.Object, 10)
	if findErr := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelInput).
		Fields(metadata.ModelFieldID, metadata.ModelFieldObjSortNumber).All(kit.Ctx, &modelResult); findErr != nil {
		blog.Error("get object sort number failed, database operation is failed, err: %v, rid: %s", findErr, kit.Rid)
		return findErr
	}

	if updateErr := m.valueAdd(kit, object, modelResult); updateErr != nil {
		blog.Errorf("update object sort number failed, err: %v, ctx:%v, rid: %s", updateErr, kit.Rid)
		return updateErr
	}
	return nil
}

func (m *modelManager) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Object, error) {

	dataResult := make([]metadata.Object, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", kit.Rid, cond.ToMapStr(), err.Error())
		return dataResult, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	dataResult := make([]mapstr.MapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", kit.Rid, cond.ToMapStr(), err.Error())
		return dataResult, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) delete(kit *rest.Kit, cond universalsql.Condition) (uint64, error) {

	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a deletion operation on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjDes, err.Error())
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}

// cascadeDelete 删除模型的字段，分组，唯一校验。模型等。
func (m *modelManager) cascadeDelete(kit *rest.Kit, objIDs []string) (uint64, error) {
	delCond := mongo.NewCondition()
	delCond.Element(mongo.Field(common.BKObjIDField).In(objIDs))
	delCondMap := util.SetQueryOwner(delCond.ToMapStr(), kit.SupplierAccount)

	// delete model property group
	if err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Delete(kit.Ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model attribute group error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model property attribute
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model attribute error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model unique
	if err := mongodb.Client().Table(common.BKTableNameObjUnique).Delete(kit.Ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model
	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, delCondMap)
	if err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	return cnt, nil
}

// cascadeDeleteTable delete the fields of the tabular model, grouping. model etc.
func (m *modelManager) cascadeDeleteTable(kit *rest.Kit, input metadata.DeleteTableOption) error {

	obj := metadata.GenerateModelQuoteObjID(input.ObjID, input.PropertyID)

	// delete quoted instance table
	instTable := common.GetInstTableName(obj, kit.SupplierAccount)
	err := mongodb.Client().DropTable(kit.Ctx, instTable)
	if err != nil {
		blog.Errorf("drop instance table failed, err: %v, table: %s, rid: %s", err, instTable, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model property attribute.
	modelDelCond := mapstr.MapStr{
		common.BKFieldID: input.ID,
	}
	modelDelCond = util.SetQueryOwner(modelDelCond, kit.SupplierAccount)

	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, modelDelCond); err != nil {
		blog.Errorf("delete model attribute failed, err: %v, cond: %+v, rid: %s", err, modelDelCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete table model quote relation.
	quoteCond := mapstr.MapStr{
		common.BKDestModelField:  obj,
		common.BKSrcModelField:   input.ObjID,
		common.BKPropertyIDField: input.PropertyID,
	}
	quoteCond = util.SetQueryOwner(quoteCond, kit.SupplierAccount)

	if err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Delete(kit.Ctx, quoteCond); err != nil {
		blog.Errorf("delete model quote relations failed, err: %v, filter: %+v, rid: %v", err, quoteCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	cond := mapstr.MapStr{
		common.BKObjIDField: obj,
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	// delete model property group.
	if err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("delete model attribute group failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete table model.
	_, err = mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete model failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	return nil
}

// createTableObjectShardingTables creates new collections for new table model,
// which create new object instance and association collections, and fix missing indexes.
func (m *modelManager) createTableObjectShardingTables(kit *rest.Kit, objID string) error {
	// table collection names.
	instTableName := common.GetObjectInstTableName(objID, kit.SupplierAccount)
	// table collections indexes.
	instTableIndexes := dbindex.TableInstanceIndexes()
	// create table object instance collection.
	err := m.createShardingTable(kit, instTableName, instTableIndexes)
	if err != nil {
		return fmt.Errorf("create object instance sharding table, %+v", err)
	}
	return nil
}

// createObjectShardingTables creates new collections for new model,
// which create new object instance and association collections, and fix missing indexes.
func (m *modelManager) createObjectShardingTables(kit *rest.Kit, objID string, isMainLine bool) error {
	// collection names.
	instTableName := common.GetObjectInstTableName(objID, kit.SupplierAccount)
	instAsstTableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)

	// collections indexes.
	instTableIndexes := dbindex.InstanceIndexes()
	instAsstTableIndexes := dbindex.InstanceAssociationIndexes()
	// 主线模型和非主线模型的唯一索引不一样
	if isMainLine {
		instTableIndexes = append(instTableIndexes, index.MainLineInstanceUniqueIndex()...)
	} else {
		instTableIndexes = append(instTableIndexes, index.InstanceUniqueIndex()...)
	}

	// create object instance table.
	err := m.createShardingTable(kit, instTableName, instTableIndexes)
	if err != nil {
		return fmt.Errorf("create object instance sharding table, %+v", err)
	}

	// create object instance association table.
	err = m.createShardingTable(kit, instAsstTableName, instAsstTableIndexes)
	if err != nil {
		return fmt.Errorf("create object instance association sharding table, %+v", err)
	}

	return nil
}

// dropObjectShardingTables drops the collections of target model.
func (m *modelManager) dropObjectShardingTables(kit *rest.Kit, objID string) error {
	// collection names.
	instTableName := common.GetObjectInstTableName(objID, kit.SupplierAccount)
	instAsstTableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)

	// drop object instance table.
	err := m.dropShardingTable(kit, instTableName)
	if err != nil {
		return fmt.Errorf("drop object instance sharding table, %+v", err)
	}

	// drop object instance association table.
	err = m.dropShardingTable(kit, instAsstTableName)
	if err != nil {
		return fmt.Errorf("drop object instance association sharding table, %+v", err)
	}

	return nil
}

// createShardingTable creates a new collection with target name, and fix missing indexes base on given index list.
func (m *modelManager) createShardingTable(kit *rest.Kit, tableName string, indexes []types.Index) error {
	// check table existence.
	tableExists, err := mongodb.Client().HasTable(kit.Ctx, tableName)
	if err != nil {
		return fmt.Errorf("check sharding table existence failed, %+v", err)
	}
	if !tableExists {
		err = mongodb.Client().CreateTable(kit.Ctx, tableName)
		if err != nil && !mongodb.Client().IsDuplicatedError(err) {
			return fmt.Errorf("create sharding table failed, %+v", err)
		}
	}

	// target collection is exist, try to check and fix the missing indexes now.
	missingIndexes := []types.Index{}

	// get all created table indexes.
	createdIndexes, err := mongodb.Client().Table(tableName).Indexes(kit.Ctx)
	if err != nil {
		return fmt.Errorf("get created sharding table[%s] indexes failed, %+v", tableName, err)
	}

	// find missing indexes.
	for _, index := range indexes {
		createdIndex, indexExists := dbindex.FindIndexByIndexFields(index.Keys, createdIndexes)
		if !indexExists || !dbindex.IndexEqual(index, createdIndex) {
			missingIndexes = append(missingIndexes, index)
		}
		// NOTE: DO NOT delete index, maybe it's created by other way.
	}

	// create missing indexes.
	for _, index := range missingIndexes {
		err = mongodb.Client().Table(tableName).CreateIndex(kit.Ctx, index)
		if err != nil {
			return fmt.Errorf("create sharding table[%s] index failed, index: %+v, %+v", tableName, index, err)
		}
	}

	return nil
}

// dropShardingTable drops the sharding table with target name.
func (m *modelManager) dropShardingTable(kit *rest.Kit, tableName string) error {
	if !common.IsObjectShardingTable(tableName) {
		return fmt.Errorf("not sharding table, can't drop it")
	}

	// check remain data.
	err := mongodb.Client().Table(tableName).Find(common.KvMap{}).One(kit.Ctx, &common.KvMap{})
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		return fmt.Errorf("check data failed, can't drop the sharding table[%s], %+v", tableName, err)
	}
	if err == nil {
		return fmt.Errorf("can't drop the non-empty sharding table[%s]", tableName)
	}

	// drop the empty table.
	if err := mongodb.Client().DropTable(kit.Ctx, tableName); err != nil {
		return fmt.Errorf("drop sharding table[%s] failed, %+v", tableName, err)
	}
	return nil
}
