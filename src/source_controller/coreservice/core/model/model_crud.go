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

func (m *modelManager) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	data.Set(metadata.ModelFieldLastTime, time.Now())
	models := make([]metadata.Object, 0)
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(kit.Ctx, &models)
	if nil != err {
		blog.Errorf("find models failed, filter: %+v, err: %s, rid: %s", cond.ToMapStr(), err.Error(), kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	if objName, exist := data[common.BKObjNameField]; exist == true && len(util.GetStrByInterface(objName)) > 0 {
		for _, model := range models {
			modelName := data[common.BKObjNameField]

			// 检查模型名称重复
			modelNameUniqueFilter := map[string]interface{}{
				common.BKObjNameField: modelName,
				common.BKFieldID: map[string]interface{}{
					common.BKDBNE: model.ID,
				},
			}

			sameNameCount, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("check whether same name model exists failed, name: %s, filter: %+v, err: %s, rid: %s", modelName, modelNameUniqueFilter, err.Error(), kit.Rid)
				return 0, err
			}
			if sameNameCount > 0 {
				blog.Warnf("update model failed, field `%s` duplicated, rid: %s", modelName, kit.Rid)
				return 0, kit.CCError.Errorf(common.CCErrCommDuplicateItem, modelName)
			}

			// 一次更新多个模型的时候，唯一校验需要特别小心
			filter := map[string]interface{}{common.BKFieldID: model.ID}
			cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, filter, data)
			if nil != err {
				blog.Errorf("request(%s): it is failed to execute database update operation on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjDes, err.Error())
				return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
			}
		}
		return cnt, nil
	}

	cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database update operation on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjDes, err.Error())
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
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
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model property attribute
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model attribute error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model unique
	if err := mongodb.Client().Table(common.BKTableNameObjUnique).Delete(kit.Ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model
	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).DeleteMany(kit.Ctx, delCondMap)
	if err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	return cnt, nil
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
