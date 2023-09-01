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

// Package model TODO
package model

import (
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

var _ core.ModelOperation = (*modelManager)(nil)

type modelManager struct {
	*modelAttributeGroup
	*modelAttribute
	*modelClassification
	*modelAttrUnique
	language  language.CCLanguageIf
	dependent OperationDependences
}

// New create a new model manager instance
func New(dependent OperationDependences, language language.CCLanguageIf) core.ModelOperation {

	coreMgr := &modelManager{dependent: dependent, language: language}
	coreMgr.modelAttribute = &modelAttribute{model: coreMgr, language: language}
	coreMgr.modelClassification = &modelClassification{model: coreMgr}
	coreMgr.modelAttributeGroup = &modelAttributeGroup{model: coreMgr}
	coreMgr.modelAttrUnique = &modelAttrUnique{}

	return coreMgr
}

// CreateTableModel create a new table model
func (m *modelManager) CreateTableModel(kit *rest.Kit, inputParam metadata.CreateModel) (
	*metadata.CreateOneDataResult, error) {

	// check the model attributes value
	if len(inputParam.Spec.ObjectID) == 0 {
		blog.Errorf("table model object %s is not set, rid: %s", inputParam.Spec.ObjectID, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	if len(inputParam.Attributes) == 0 {
		blog.Errorf("table model attr is not set, rid: %s", kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "attribute")
	}

	locker := lock.NewLocker(redis.Client())
	redisKey := lock.GetLockKey(lock.CreateModelFormat, inputParam.Spec.ObjectID)

	locked, err := locker.Lock(redisKey, time.Second*35)
	defer locker.Unlock()
	if err != nil {
		blog.Errorf("get create table model lock failed, err: %v, input: %v, rid: %s", err, inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommRedisOPErr)
	}

	if !locked {
		blog.Errorf("create table model have same task in progress, input: %v, rid:%s", inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommOPInProgressErr,
			fmt.Sprintf("create table object(%s)", inputParam.Spec.ObjectID))
	}

	blog.V(5).Infof("create table model redis lock info, key: %s, bl: %v, err: %v, rid: %s",
		redisKey, locked, err, kit.Rid)

	originObjID := inputParam.Spec.ObjectID
	inputParam.Spec.ObjectID = metadata.GenerateModelQuoteObjID(inputParam.Spec.ObjectID,
		inputParam.Attributes[0].PropertyID)
	inputParam.Spec.ObjectName = metadata.GenerateModelQuoteObjID(inputParam.Spec.ObjectID,
		inputParam.Attributes[0].PropertyName)

	// check the model if it is exists
	condCheckModelMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	condCheckModel, _ := mongo.NewConditionFromMapStr(condCheckModelMap)
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	_, exists, err := m.isExists(kit, condCheckModel)
	if nil != err {
		blog.Errorf("failed to check whether the table model (%s) is exists, err: %v, rid: %s",
			inputParam.Spec.ObjectID, err, kit.Rid)
		return nil, err
	}
	if exists {
		blog.Errorf("failed to create a table model, model (%s) is already exists, rid: %s ",
			inputParam.Spec.ObjectID, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectID)
	}

	// check for duplicate model names
	modelNameUniqueFilter := map[string]interface{}{
		common.BKObjNameField: inputParam.Spec.ObjectName,
	}
	sameNameCount, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get model count failed, name: %s, err: %v, rid: %s", inputParam.Spec.ObjectName, err, kit.Rid)
		return nil, err
	}
	if sameNameCount > 0 {
		blog.Warnf("create model failed, field `%s` duplicated, rid: %s", inputParam.Spec.ObjectName, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectName)
	}

	// create new table model after checking base information and sharding table operation.
	inputParam.Spec.OwnerID = kit.SupplierAccount
	id, err := m.save(kit, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the model (%#v), err: %v", kit.Rid, inputParam.Spec, err)
		return nil, err
	}

	// 创建源模型的模型属性
	_, err = m.modelAttribute.CreateTableModelAttributes(kit, originObjID,
		metadata.CreateModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("it is failed to create some attributes (%#v) for the model (%s), err: %v, rid: %s",
			inputParam.Attributes, inputParam.Spec.ObjectID, err, kit.Rid)
		return nil, err
	}

	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

// CreateModel create a new model
func (m *modelManager) CreateModel(kit *rest.Kit, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error) {

	locker := lock.NewLocker(redis.Client())
	redisKey := lock.GetLockKey(lock.CreateModelFormat, inputParam.Spec.ObjectID)

	locked, err := locker.Lock(redisKey, time.Second*35)
	defer locker.Unlock()
	if err != nil {
		blog.ErrorJSON("create model error. get create look error. err:%s, input:%s, rid:%s", err, inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommRedisOPErr)
	}
	if !locked {
		blog.ErrorJSON("create model have same task in progress. input:%s, rid:%s", inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommOPInProgressErr, fmt.Sprintf("create object(%s)",
			inputParam.Spec.ObjectID))
	}
	blog.V(5).Infof("create model redis look info. key:%s, bl:%v, err:%v, rid:%s", redisKey, locked, err, kit.Rid)

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set", kit.Rid, inputParam.Spec.ObjectID)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	// 禁止创建bk或者BK打头的模型，用于后续创建内置模型，避免冲突
	if err := validIDStartWithBK(kit, inputParam.Spec.ObjectID); err != nil {
		return nil, err
	}

	// 因为模型名称会用于生成实例和实例关联的mongodb表名，所以需要校验模型对应的实例表和实例关联表名均不超过mongodb的长度限制
	if !SatisfyMongoCollLimit(common.GetObjectInstTableName(inputParam.Spec.ObjectID, kit.SupplierAccount)) ||
		!SatisfyMongoCollLimit(common.GetObjectInstAsstTableName(inputParam.Spec.ObjectID, kit.SupplierAccount)) {
		blog.Errorf("inputParam.Spec.ObjectID:%s not SatisfyMongoCollLimit", inputParam.Spec.ObjectID)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(kit, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s",
			kit.Rid, inputParam.Spec.ObjCls, err)
		return nil, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid",
			kit.Rid, inputParam.Spec.ObjCls)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	// check the model if it is exists
	condCheckModelMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	condCheckModel, _ := mongo.NewConditionFromMapStr(condCheckModelMap)
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	_, exists, err := m.isExists(kit, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the model (%s) is exists, error info is %s ", kit.Rid, inputParam.Spec.ObjectID, err)
		return nil, err
	}
	if exists {
		blog.Warnf("request(%s): it is failed to  create a new model , because of the model (%s) is already exists ",
			kit.Rid, inputParam.Spec.ObjectID)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectID)
	}

	// 检查模型名称重复
	modelNameUniqueFilter := map[string]interface{}{
		common.BKObjNameField: inputParam.Spec.ObjectName,
	}
	sameNameCount, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("whether same name model exists, name: %s, err: %v, rid: %s", inputParam.Spec.ObjectName,
			err, kit.Rid)
		return nil, err
	}
	if sameNameCount > 0 {
		blog.Warnf("create model failed, field `%s` duplicated, rid: %s", inputParam.Spec.ObjectName, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectName)
	}

	// check object instance and instance association table.
	/* 	if err := m.createObjectShardingTables(kit, inputParam.Spec.ObjectID); err != nil {
		blog.Errorf("handle object sharding tables failed, object: %s name: %s, err: %v, rid: %s",
			inputParam.Spec.ObjectID, inputParam.Spec.ObjectName, err, kit.Rid)
		return nil, err
	} */

	// create new model after checking base informations and sharding table operation.
	inputParam.Spec.OwnerID = kit.SupplierAccount
	id, err := m.save(kit, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", kit.Rid, inputParam.Spec,
			err)
		return nil, err
	}

	// create initial phase model attributes.
	if len(inputParam.Attributes) != 0 {
		_, err = m.modelAttribute.CreateModelAttributes(kit, inputParam.Spec.ObjectID, metadata.CreateModelAttributes{
			Attributes: inputParam.Attributes})
		if nil != err {
			blog.Errorf("request(%s): it is failed to create some attributes (%#v) for the model (%s), err: %v",
				kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err)
			return nil, err
		}
	}

	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

// SetModel TODO
func (m *modelManager) SetModel(kit *rest.Kit, inputParam metadata.SetModel) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set",
			kit.Rid, inputParam.Spec.ObjectID)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(kit, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s",
			kit.Rid, inputParam.Spec.ObjCls, err)
		return dataResult, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid",
			kit.Rid, inputParam.Spec.ObjCls)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	condCheckModelMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	condCheckModel, _ := mongo.NewConditionFromMapStr(condCheckModelMap)
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

	existsModel, exists, err := m.isExists(kit, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check the model (%s) is exists, error info is %s ", kit.Rid,
			inputParam.Spec.ObjectID, err)
		return &metadata.SetDataResult{}, err
	}

	inputParam.Spec.OwnerID = kit.SupplierAccount
	// set model spec
	if exists {
		updateCondMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
		updateCond, _ := mongo.NewConditionFromMapStr(updateCondMap)
		updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

		_, err := m.update(kit, mapstr.NewFromStruct(inputParam.Spec, "field"), updateCond)
		if nil != err {
			blog.Errorf("request(%s): it is failed to update some fields (%#v) for the model (%s), error info is %s",
				kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err)
			return dataResult, err
		}

		dataResult.UpdatedCount.Count++
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{OriginIndex: 0,
			ID: uint64(existsModel.ID)})
	} else {
		id, err := m.save(kit, &inputParam.Spec)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", kit.Rid,
				inputParam.Spec.ObjectID, err)
			return dataResult, err
		}
		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{OriginIndex: 0, ID: id})
	}

	// set model attributes
	setAttrResult, err := m.modelAttribute.SetModelAttributes(kit, inputParam.Spec.ObjectID, metadata.SetModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the attributes (%#v) for the model (%s), error info is %s",
			kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err)
		return dataResult, err
	}
	_ = setAttrResult // TODO: how to return this result ? let me think about it;
	/*
		// set attribute result, ignore model operation result
		dataResult.CreatedCount = setAttrResult.CreatedCount
		dataResult.UpdatedCount = setAttrResult.UpdatedCount

		dataResult.Created = append(dataResult.Created, setAttrResult.Created...)
		dataResult.Updated = append(dataResult.Updated, setAttrResult.Updated...)
		dataResult.Exceptions = append(dataResult.Exceptions, setAttrResult.Exceptions...)
	*/
	return dataResult, err
}

// UpdateModel TODO
func (m *modelManager) UpdateModel(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	updateCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(),
		kit.SupplierAccount))
	if err != nil {
		blog.Errorf("convert the condition from mapstr into condition object failed, err: %v, condition: %v, rid: %s",
			err, inputParam.Condition, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, updateCond)
	return &metadata.UpdatedCount{Count: cnt}, err
}

// DeleteModel TODO
func (m *modelManager) DeleteModel(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all models by the deletion condition
	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(),
		kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, "+
			"error info is %s", kit.Rid, inputParam.Condition, err)
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrCommParamsInvalid, err.Error())
	}

	modelItems, err := m.search(kit, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to find the all models by the condition (%#v), error info is %s",
			kit.Rid, deleteCond.ToMapStr(), err)
		return &metadata.DeletedCount{}, err
	}

	targetObjIDS := make([]string, 0)
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}

	// check if the model is used: firstly to check instance
	exists, err := m.dependent.HasInstance(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some instances, "+
			"error info is %s", kit.Rid, targetObjIDS, err)
		return &metadata.DeletedCount{}, err
	}
	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some instances.",
			kit.Rid, targetObjIDS)
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrCoreServiceModelHasInstanceErr)
	}

	// check if the model is used: secondly to check association
	exists, err = m.dependent.HasAssociation(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some associations, "+
			"error info is %s", kit.Rid, targetObjIDS, err)
		return &metadata.DeletedCount{}, err
	}
	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some associations.",
			kit.Rid, targetObjIDS)
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrCoreServiceModelHasAssociationErr)
	}

	// delete model self
	cnt, err := m.deleteModelAndAttributes(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the models (%#v) and their attributes, error info is %s",
			kit.Rid, targetObjIDS, err)
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

// CascadeDeleteModel 将会删除模型/模型属性/属性分组/唯一校验/模型继承的模版相关数据
func (m *modelManager) CascadeDeleteModel(kit *rest.Kit, modelID int64) (*metadata.DeletedCount, error) {
	// NOTE: just single model cascade delete action now.
	deleteCondMap := util.SetQueryOwner(make(map[string]interface{}), kit.SupplierAccount)
	deleteCond, _ := mongo.NewConditionFromMapStr(deleteCondMap)
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldID, Val: modelID})

	// NOTE: the func logics supports cascade delete models in batch mode.
	// You can change the condition to index many models.
	models, err := m.search(kit, deleteCond)
	if err != nil {
		blog.Errorf("cascade delete model, search target models failed, condition: %s, error: %s, rid: %s",
			deleteCond.ToMapStr(), err, kit.Rid)
		return nil, err
	}

	objIDs := make([]string, 0)
	for _, model := range models {
		objIDs = append(objIDs, model.ObjectID)
	}
	if len(objIDs) == 0 {
		return &metadata.DeletedCount{}, nil
	}

	// check object instances.
	exists, err := m.dependent.HasInstance(kit, objIDs)
	if err != nil {
		blog.Errorf("cascade delete model, check object instance failed, objects: %+v, error: %s, rid: %s",
			objIDs, err, kit.Rid)
		return nil, err
	}
	if exists {
		blog.Errorf("cascade delete model failed, there are vestigital object instances, objects: %+v, rid: %s",
			objIDs, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCoreServiceModelHasInstanceErr)
	}

	// check model associations.
	exists, err = m.dependent.HasAssociation(kit, objIDs)
	if err != nil {
		blog.Errorf("cascade delete model, check model associations failed, objects: %+v, error: %s, rid: %s",
			objIDs, err, kit.Rid)
		return nil, err
	}
	if exists {
		blog.Errorf("cascade delete model failed, there are vestigital model associations, objects: %+v, rid: %s",
			objIDs, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCoreServiceModelHasAssociationErr)
	}

	// cascade delete.
	cnt, err := m.cascadeDelete(kit, objIDs)
	if err != nil {
		blog.Errorf("cascade delete model failed, objects: %+v, err: %v, rid: %s", objIDs, err, kit.Rid)
		return nil, err
	}

	if err := m.deleteModelAndFieldTemplateRelation(kit, modelID); err != nil {
		return nil, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *modelManager) getObjFieldTemplateRelation(kit *rest.Kit, modelID int64) ([]int64, error) {

	result := make([]metadata.ObjFieldTemplateRelation, 0)
	filter := mapstr.MapStr{
		common.ObjectIDField: modelID,
	}

	if err := mongodb.Client().Table(common.BKTableNameObjFieldTemplateRelation).Find(filter).
		Fields(common.BKTemplateID).All(kit.Ctx, &result); err != nil {
		blog.Errorf("failed to get object and relation, filter: (%#v), err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	templateIDs := make([]int64, 0)

	if len(result) == 0 {
		return templateIDs, nil
	}

	for _, relation := range result {
		templateIDs = append(templateIDs, relation.TemplateID)
	}
	return templateIDs, nil
}

func (m *modelManager) isExistProcessingTask(kit *rest.Kit, modelID int64, isStop bool) error {
	if !isStop {
		return nil
	}

	templateIDs, err := m.getObjFieldTemplateRelation(kit, modelID)
	if err != nil {
		return err
	}

	if err := dealProcessRunningTasks(kit, templateIDs, modelID, isStop); err != nil {
		return err
	}
	return nil
}

func (m *modelManager) deleteModelAndFieldTemplateRelation(kit *rest.Kit, modelID int64) error {

	templateIDs, err := m.getObjFieldTemplateRelation(kit, modelID)
	if err != nil {
		return err
	}

	if err := dealProcessRunningTasks(kit, templateIDs, modelID, false); err != nil {
		return err
	}

	filter := mapstr.MapStr{
		common.ObjectIDField: modelID,
	}

	// delete object field template relation
	if err := mongodb.Client().Table(common.BKTableNameObjFieldTemplateRelation).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("delete model field template relation failed, cond: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// CascadeDeleteTableModel delete table models in a cascading manner
func (m *modelManager) CascadeDeleteTableModel(kit *rest.Kit, intput metadata.DeleteTableOption) error {
	// NOTE: just single model cascade delete action now.
	condMap := util.SetQueryOwner(make(map[string]interface{}), kit.SupplierAccount)
	cond, _ := mongo.NewConditionFromMapStr(condMap)
	cond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: intput.ObjID})

	// NOTE: the func logics supports cascade delete models in batch mode.
	// You can change the condition to index many models.
	models, err := m.search(kit, cond)
	if err != nil {
		blog.Errorf("search target models failed, cond: %s, err: %v, rid: %s", cond.ToMapStr(), err, kit.Rid)
		return err
	}

	if len(models) == 0 {
		return nil
	}

	// cascade delete table related resources.
	if err := m.cascadeDeleteTable(kit, intput); err != nil {
		blog.Errorf("cascade delete model failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// SearchModel TODO
func (m *modelManager) SearchModel(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult,
	error) {

	dataResult := &metadata.QueryModelDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(),
		kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object,"+
			" error info is %s", kit.Rid, inputParam.Condition, err)
		return dataResult, err
	}

	totalCount, err := m.count(kit, searchCond)
	if nil != err {
		blog.Errorf("failed to get count by the cond (%#v), err: %v, rid: %s", searchCond.ToMapStr(), err, kit.Rid)
		return dataResult, err
	}

	modelItems, err := m.search(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", kit.Rid,
			searchCond.ToMapStr(), err)
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	dataResult.Info = modelItems
	return dataResult, nil
}

// SearchModelWithAttribute TODO
func (m *modelManager) SearchModelWithAttribute(kit *rest.Kit, inputParam metadata.QueryCondition) (
	*metadata.QueryModelWithAttributeDataResult, error) {

	dataResult := &metadata.QueryModelWithAttributeDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(),
		kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, "+
			"error info is %s", kit.Rid, inputParam.Condition, err)
		return dataResult, err
	}

	totalCount, err := m.count(kit, searchCond)
	if nil != err {
		blog.Errorf("failed to get count by the cond (%#v), err: %v, rid: %s", searchCond.ToMapStr(), err, kit.Rid)
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	modelItems, err := m.search(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", kit.Rid,
			searchCond.ToMapStr(), err)
		return dataResult, err
	}

	for _, modelItem := range modelItems {
		queryAttributeCondMap := util.SetQueryOwner(make(map[string]interface{}), modelItem.OwnerID)
		queryAttributeCond, _ := mongo.NewConditionFromMapStr(queryAttributeCondMap)
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelItem.ObjectID))
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldSupplierAccount).Eq(modelItem.OwnerID))
		attributeItems, err := m.modelAttribute.search(kit, queryAttributeCond)
		if nil != err {
			blog.Errorf("request(%s):it is failed to search the object(%s)'s attributes, error info is %s",
				modelItem.ObjectID, err)
			return dataResult, err
		}
		dataResult.Info = append(dataResult.Info, metadata.SearchModelInfo{Spec: modelItem, Attributes: attributeItems})
	}

	return dataResult, nil
}

// validIDStartWithBK validate the id start with bk or BK
func validIDStartWithBK(kit *rest.Kit, modelID string) error {
	if strings.HasPrefix(strings.ToLower(modelID), "bk") {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_obj_id value can not start with bk or BK")
	}
	return nil
}

func dealProcessRunningTasks(kit *rest.Kit, ids []int64, objectID int64, isStop bool) error {

	if len(ids) == 0 {
		return nil
	}

	// 1、get the status of the task
	cond := mapstr.MapStr{
		common.BKInstIDField: mapstr.MapStr{
			common.BKDBIN: ids,
		},
		metadata.APITaskExtraField: objectID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{
				metadata.APITaskStatusExecute, metadata.APITaskStatusWaitExecute,
				metadata.APITaskStatusNew},
		},
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	result := make([]metadata.APITaskSyncStatus, 0)
	if err := mongodb.Client().Table(common.BKTableNameAPITaskSyncHistory).Find(cond).
		Fields(common.BKStatusField, common.BKTaskIDField).
		All(kit.Ctx, &result); err != nil {
		blog.Errorf("search task failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	// 2、the possible task status scenarios are: one is executing,
	// one is waiting or new, but there will be no more than two tasks.
	if len(result) > metadata.MaxFieldTemplateTaskNum {
		blog.Errorf("task num incorrect, template IDs: %v, objID: %d, rid: %s", ids, objectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommGetMultipleObject,
			fmt.Sprintf("template IDs: %v, objID: %v", ids, objectID))
	}

	if len(result) == 0 {
		return nil
	}

	// 3. if it is a deactivated model scenario, deactivation is not allowed as long as
	// there is a status of running and waiting for execution.
	if isStop {
		tasksIDs := make([]string, 0)
		for _, task := range result {
			tasksIDs = append(tasksIDs, task.TaskID)
		}
		blog.Errorf("there are tasks running, tasksIDs: %v, template IDs: %v, objID: %d, rid: %s", tasksIDs,
			ids, objectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrorTopoFieldTemplateForbiddenPauseModel,
			fmt.Sprintf("task IDs: %v, objID: %v", tasksIDs, objectID))
	}

	// 4. if there is a running task in the deleted model scene,
	// an error is returned. If it is a task status waiting to be
	// executed, the task needs to be deleted	var taskID string.
	var taskID string
	for _, info := range result {
		if info.Status == metadata.APITaskStatusExecute {
			blog.Errorf("unbinding failed, sync task(%s) is running, template IDs: %v, objID: %d, rid: %d",
				info.TaskID, ids, objectID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrTaskDeleteConflict, info.TaskID)
		}
		taskID = info.TaskID
	}

	// 4、if there is a queued task that needs to be cleared.
	delCond := mapstr.MapStr{
		common.BKTaskIDField: taskID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusWaitExecute, metadata.APITaskStatusNew},
		},
	}
	err := mongodb.Client().Table(common.BKTableNameAPITask).Delete(kit.Ctx, delCond)
	if err != nil {
		blog.Errorf("delete task failed, cond: %#v, err: %v, rid: %s", delCond, err, kit.Rid)
		return err
	}
	return nil
}
