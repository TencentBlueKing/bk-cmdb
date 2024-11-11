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
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelManager) count(kit *rest.Kit, cond universalsql.Condition) (uint64, error) {

	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("it is failed to execute database count operation by the condition, err: %v, cond: %v, rid: %s",
			err, cond.ToMapStr(), kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) save(kit *rest.Kit, model *metadata.Object) (id uint64, err error) {
	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("it is failed to make sequence id on the table, err: %v, table: %s, rid: %s",
			err, common.BKTableNameObjDes, kit.Rid)
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	sortNum, err := m.GetSortNum(kit, *model, nil)
	if err != nil {
		blog.Errorf("set object sort number failed, err: %v, objectId: %s, rid: %s", err, model.ObjectID, kit.Rid)
		return id, err
	}

	model.ObjSortNumber = sortNum
	model.ID = int64(id)
	model.OwnerID = kit.SupplierAccount

	now := time.Now()
	if model.LastTime == nil {
		model.LastTime = &metadata.Time{}
		model.LastTime.Time = now
		model.Modifier = kit.User
	}
	if model.CreateTime == nil {
		model.CreateTime = &metadata.Time{}
		model.CreateTime.Time = now
		model.Creator = kit.User
	}

	if err = mongodb.Client().Table(common.BKTableNameObjDes).Insert(kit.Ctx, model); err != nil {
		return 0, err
	}

	assetIDGenerator := map[string]interface{}{
		common.BKFieldDBID:     metadata.GetIDRule(model.ObjectID),
		common.BKFieldSeqID:    0,
		common.CreateTimeField: now,
		common.LastTimeField:   now,
	}
	if err = mongodb.Client().Table(common.BKTableNameIDgenerator).Insert(kit.Ctx, assetIDGenerator); err != nil {
		blog.Errorf("add id generator data failed, data: %+v, err: %v, rid: %s", assetIDGenerator, err, kit.Rid)
		return 0, err
	}

	return id, err
}

// GetSortNum 获取模型排序字段
func (m *modelManager) GetSortNum(kit *rest.Kit, model metadata.Object, srcModel *metadata.Object) (int64, error) {
	// 查询当前分组下的模型信息
	modelInput := map[string]interface{}{metadata.ModelFieldObjCls: model.ObjCls}
	modelResult := make([]metadata.Object, 0)
	sortCond := "-obj_sort_number"
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelInput).Sort(sortCond).Fields(
		metadata.ModelFieldID, metadata.ModelFieldObjSortNumber).Limit(1).All(kit.Ctx, &modelResult); err != nil {
		blog.Error("get object sort number failed, database operation is failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}
	if len(modelResult) <= 0 {
		return 1, nil
	}

	if model.ObjSortNumber == 0 {
		return modelResult[0].ObjSortNumber + 1, nil
	}

	// 更新操作需要查使用模型原来所属的分组与目标分组做对比，查看是否是跨分组更新
	// 在当前分组移动
	if srcModel != nil && srcModel.ObjCls == model.ObjCls {
		sortNumCond := mapstr.MapStr{}
		step := int64(-1)
		if srcModel.ObjSortNumber > model.ObjSortNumber {
			// 前移
			sortNumCond.Set(common.BKDBGTE, model.ObjSortNumber)
			sortNumCond.Set(common.BKDBLT, srcModel.ObjSortNumber)
			step = 1
		} else {
			// 后移
			sortNumCond.Set(common.BKDBGT, srcModel.ObjSortNumber)
			sortNumCond.Set(common.BKDBLTE, model.ObjSortNumber)
		}

		if err := m.updateOtherObjSort(kit, model.ObjCls, sortNumCond, step); err != nil {
			blog.Errorf("update object sort number failed, err: %v, rid: %s", err, kit.Rid)
			return 0, err
		}
		return model.ObjSortNumber, nil
	}

	if srcModel != nil && srcModel.ObjCls != model.ObjCls {
		// 移动到其它分组：原分组中序号大于原序号的减一
		sortNumCond := mapstr.MapStr{common.BKDBGT: srcModel.ObjSortNumber}
		if err := m.updateOtherObjSort(kit, srcModel.ObjCls, sortNumCond, int64(-1)); err != nil {
			blog.Errorf("update object sort number failed, err: %v, rid: %s", err, kit.Rid)
			return 0, err
		}
	}

	if model.ObjSortNumber > modelResult[0].ObjSortNumber {
		return modelResult[0].ObjSortNumber + 1, nil
	}
	// 目标分组中序号大于等于参数序号的加一
	sortNumCond := mapstr.MapStr{common.BKDBGTE: model.ObjSortNumber}
	if err := m.updateOtherObjSort(kit, model.ObjCls, sortNumCond, int64(1)); err != nil {
		blog.Errorf("update object sort number failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}
	return model.ObjSortNumber, nil
}

// objSortNumberChange 根据传入条件更新其它模型的排序字段
func (m *modelManager) updateOtherObjSort(kit *rest.Kit, class string, sortNumCond mapstr.MapStr, step int64) error {
	incCond := mapstr.MapStr{
		metadata.ModelFieldObjCls:        class,
		metadata.ModelFieldObjSortNumber: sortNumCond,
	}
	incData := mapstr.MapStr{metadata.ModelFieldObjSortNumber: step}
	err := mongodb.Client().Table(common.BKTableNameObjDes).UpdateMultiModel(kit.Ctx, incCond, types.ModeUpdate{
		Op: "inc", Doc: incData})
	if err != nil {
		blog.Errorf("increase object sort number failed, err: %v, incCond: %v, rid: %s", err, incCond, kit.Rid)
		return err
	}
	return nil
}

func (m *modelManager) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {
	data.Set(metadata.ModelFieldModifier, kit.User)
	data.Set(metadata.ModelFieldLastTime, time.Now())
	models := make([]metadata.Object, 0)
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx, &models)
	if err != nil {
		blog.Errorf("find models failed, err: %v, filter: %v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

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
				blog.Errorf("failed to check the validity of the model name, filter: %v, err: %v, rid: %s",
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
				blog.Errorf("failed to update table (%s), err: %v, cond:  %v, data:  %v, err: %v,rid: %s",
					common.BKTableNameObjDes, filter, data, err, kit.Rid)
				return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
			}
			return cnt, nil
		}
	}

	// 如果有传递分组或排序字段条件，需要先设置好模型的排序位置再进行下一步更新
	if data.Exists(metadata.ModelFieldObjCls) || data.Exists(metadata.ModelFieldObjSortNumber) {
		for _, model := range models {
			if err := m.updateSortNum(kit, data, model); err != nil {
				blog.Errorf("set object sort number failed, err: %v, data: %s, rid: %s", err, data, kit.Rid)
				return 0, err
			}
		}
		data.Remove(metadata.ModelFieldObjSortNumber)
	}

	cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.Errorf("failed to update the table (%s), err: %v, rid: %s", common.BKTableNameObjDes, err, kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	return cnt, err
}

// updateSortNum 根据更新条件更新 obj_sort_number 值
func (m *modelManager) updateSortNum(kit *rest.Kit, data mapstr.MapStr,
	srcModel metadata.Object) error {

	object := metadata.Object{}
	if err := mapstruct.Decode2Struct(data, &object); err != nil {
		blog.Errorf("parsing data failed, err: %v, data: %v, rid: %s", err, data, kit.Rid)
		return err
	}

	if object.ObjSortNumber < 0 {
		blog.Errorf("obj sort number field invalid failed, err: obj sort number less than 0, obj_sort_number: %d, "+
			"rid: %s", object.ObjSortNumber, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	// 如果未传递了 bk_classification_id 字段，传递了 obj_sort_number 字段,则表示在当前分组下更新模型顺序
	// 更新当前模型 obj_sort_number 前先更新当前分组下其它模型 obj_sort_number
	if !data.Exists(metadata.ModelFieldObjCls) {
		object.ObjCls = srcModel.ObjCls
	}

	// 如果传递了 bk_classification_id 字段,按照更新模型所属分组处理
	sortNum, err := m.GetSortNum(kit, object, &srcModel)
	if err != nil {
		blog.Errorf("set object sort number failed, err: %v, object: %v, rid: %s", err, object, kit.Rid)
		return err
	}

	filter := map[string]interface{}{common.BKFieldID: srcModel.ID}
	doc := map[string]interface{}{metadata.ModelFieldObjSortNumber: sortNum}
	_, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, filter, doc)
	if err != nil {
		blog.Errorf("failed to update object sort number field, err: %v, cond: %v, data: %v, rid: %s",
			err, filter, doc, kit.Rid)
		return err
	}

	return nil
}

func (m *modelManager) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Object, error) {

	dataResult := make([]metadata.Object, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).
		All(kit.Ctx, &dataResult); err != nil {
		blog.Errorf("it is failed to find all models by the condition, err: %v, cond: %v, rid: %s",
			err, cond.ToMapStr(), kit.Rid)
		return dataResult, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	dataResult := make([]mapstr.MapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx,
		&dataResult); err != nil {
		blog.Errorf("it is failed to find all models by the condition, err: %v, cond: %v, rid: %s",
			err, cond.ToMapStr(), kit.Rid)
		return dataResult, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) delete(kit *rest.Kit, cond universalsql.Condition) (uint64, error) {

	if err := m.updateSortNumWhenDelete(kit, cond.ToMapStr()); err != nil {
		blog.Errorf("failed to update object sort number when delete object, err: %v, cond: %v, rid: %s",
			err, cond, kit.Rid)
		return 0, err
	}
	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, cond.ToMapStr())
	if err != nil {
		blog.Errorf("it is failed to execute a deletion operation on the table, err: %v, table: %s, rid: %s",
			err, common.BKTableNameObjDes, kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}

// updateSortNumWhenDelete 删除某个模型前，需要将模型分组下排在这个模型后面的前移
func (m *modelManager) updateSortNumWhenDelete(kit *rest.Kit, delCond map[string]interface{}) error {
	models := make([]metadata.Object, 0)
	err := mongodb.Client().Table(common.BKTableNameObjDes).Find(delCond).Fields(common.BKClassificationIDField,
		common.ObjSortNumberField).All(kit.Ctx, &models)
	if err != nil {
		blog.Errorf("find models failed, err: %v, filter: %v, rid: %s", err, delCond, kit.Rid)
		return err
	}

	if len(models) <= 0 {
		return nil
	}
	for _, model := range models {
		// 分组中序号大于参数序号的减一
		sortNumCond := mapstr.MapStr{common.BKDBGT: model.ObjSortNumber}
		if err := m.updateOtherObjSort(kit, model.ObjCls, sortNumCond, int64(-1)); err != nil {
			blog.Errorf("update object sort number failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}
	return nil
}

// cascadeDelete 删除模型的字段，分组，唯一校验。模型等。
func (m *modelManager) cascadeDelete(kit *rest.Kit, objIDs []string) (uint64, error) {
	delCond := mongo.NewCondition()
	delCond.Element(mongo.Field(common.BKObjIDField).In(objIDs))
	delCondMap := delCond.ToMapStr()

	// delete model property group
	if err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Delete(kit.Ctx, delCondMap); err != nil {
		blog.Errorf("delete model attribute group error. err: %v, cond: %s, rid: %s", err, delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model property attribute
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, delCondMap); err != nil {
		blog.Errorf("delete model attribute error. err: %v, cond: %s, rid: %s", err, delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	// delete model unique
	if err := mongodb.Client().Table(common.BKTableNameObjUnique).Delete(kit.Ctx, delCondMap); err != nil {
		blog.Errorf("delete model unique error. err: %v, cond: %s, rid: %s", err, delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	if err := m.updateSortNumWhenDelete(kit, delCondMap); err != nil {
		blog.Errorf("failed to update object sort number when delete object, err: %v, cond: %v, rid: %s", err,
			delCondMap, kit.Rid)
		return 0, err
	}

	// delete model
	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, delCondMap)
	if err != nil {
		blog.Errorf("delete model unique error. err: %v, cond: %s, rid: %s", err, delCondMap, kit.Rid)
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

	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, modelDelCond); err != nil {
		blog.Errorf("delete model attribute failed, err: %v, cond: %v, rid: %s", err, modelDelCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete table model quote relation.
	quoteCond := mapstr.MapStr{
		common.BKDestModelField:  obj,
		common.BKSrcModelField:   input.ObjID,
		common.BKPropertyIDField: input.PropertyID,
	}

	if err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Delete(kit.Ctx, quoteCond); err != nil {
		blog.Errorf("delete model quote relations failed, err: %v, filter: %v, rid: %v", err, quoteCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	cond := mapstr.MapStr{
		common.BKObjIDField: obj,
	}

	// delete model property group.
	if err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("delete model attribute group failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete table model.
	_, err = mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete model failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
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
		return fmt.Errorf("create object instance sharding table, err: %v, rid: %s", err, kit.Rid)
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
