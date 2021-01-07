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

func (m *modelManager) CreateModel(kit *rest.Kit, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error) {

	locker := lock.NewLocker(redis.Client())
	// fmt.Sprintf("coreservice:create:model:%s", inputParam.Spec.ObjectID)
	redisKey := lock.GetLockKey(lock.CreateModelFormat, inputParam.Spec.ObjectID)

	looked, err := locker.Lock(redisKey, time.Second*35)
	defer locker.Unlock()
	if err != nil {
		blog.ErrorJSON("create model error. get create look error. err:%s, input:%s, rid:%s", err.Error(), inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommRedisOPErr)
	}
	if !looked {
		blog.ErrorJSON("create model have same task in progress. input:%s, rid:%s", inputParam, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommOPInProgressErr, fmt.Sprintf("create object(%s)", inputParam.Spec.ObjectID))
	}
	blog.V(5).Infof("create model redis look info. key:%s, bl:%v, err:%v, rid:%s", redisKey, looked, err, kit.Rid)

	dataResult := &metadata.CreateOneDataResult{}

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set", kit.Rid, inputParam.Spec.ObjectID)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	if !SatisfyMongoCollLimit(inputParam.Spec.ObjectID) {
		blog.Errorf("inputParam.Spec.ObjectID:%s not SatisfyMongoCollLimit", inputParam.Spec.ObjectID)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(kit, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s", kit.Rid, inputParam.Spec.ObjCls, err.Error())
		return dataResult, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid", kit.Rid, inputParam.Spec.ObjCls)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	// check the model if it is exists
	condCheckModelMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	condCheckModel, _ := mongo.NewConditionFromMapStr(condCheckModelMap)
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

	_, exists, err := m.isExists(kit, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the model (%s) is exists, error info is %s ", kit.Rid, inputParam.Spec.ObjectID, err.Error())
		return dataResult, err
	}
	if exists {
		blog.Warnf("request(%s): it is failed to  create a new model , because of the model (%s) is already exists ", kit.Rid, inputParam.Spec.ObjectID)
		return dataResult, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectID)
	}

	// 检查模型名称重复
	modelNameUniqueFilter := map[string]interface{}{
		common.BKObjNameField: inputParam.Spec.ObjectName,
	}

	sameNameCount, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("whether same name model exists, name: %s, err: %s, rid: %s", inputParam.Spec.ObjectName, err.Error(), kit.Rid)
		return dataResult, err
	}
	if sameNameCount > 0 {
		blog.Warnf("create model failed, field `%s` duplicated, rid: %s", inputParam.Spec.ObjectName, kit.Rid)
		return dataResult, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectName)
	}

	inputParam.Spec.OwnerID = kit.SupplierAccount
	id, err := m.save(kit, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", kit.Rid, inputParam.Spec, err.Error())
		return dataResult, err
	}

	if len(inputParam.Attributes) != 0 {
		_, err = m.modelAttribute.CreateModelAttributes(kit, inputParam.Spec.ObjectID, metadata.CreateModelAttributes{Attributes: inputParam.Attributes})
		if nil != err {
			blog.Errorf("request(%s): it is failed to create some attributes (%#v) for the model (%s), err: %v", kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err)
			return dataResult, err
		}
	}

	dataResult.Created.ID = id
	return dataResult, nil
}
func (m *modelManager) SetModel(kit *rest.Kit, inputParam metadata.SetModel) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set", kit.Rid, inputParam.Spec.ObjectID)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(kit, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s", kit.Rid, inputParam.Spec.ObjCls, err.Error())
		return dataResult, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid", kit.Rid, inputParam.Spec.ObjCls)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	condCheckModelMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	condCheckModel, _ := mongo.NewConditionFromMapStr(condCheckModelMap)
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

	existsModel, exists, err := m.isExists(kit, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check the model (%s) is exists, error info is %s ", kit.Rid, inputParam.Spec.ObjectID, err.Error())
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
			blog.Errorf("request(%s): it is failed to update some fields (%#v) for the model (%s), error info is %s", kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}

		dataResult.UpdatedCount.Count++
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{OriginIndex: 0, ID: uint64(existsModel.ID)})
	} else {
		id, err := m.save(kit, &inputParam.Spec)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", kit.Rid, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}
		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{OriginIndex: 0, ID: id})
	}

	// set model attributes
	setAttrResult, err := m.modelAttribute.SetModelAttributes(kit, inputParam.Spec.ObjectID, metadata.SetModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the attributes (%#v) for the model (%s), error info is %s", kit.Rid, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
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

func (m *modelManager) UpdateModel(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	updateCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s ", kit.Rid, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, updateCond)
	return &metadata.UpdatedCount{Count: cnt}, err
}

func (m *modelManager) DeleteModel(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all models by the deletion condition
	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrCommParamsInvalid, err.Error())
	}

	modelItems, err := m.search(kit, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to find the all models by the condition (%#v), error info is %s", kit.Rid, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	targetObjIDS := make([]string, 0)
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}

	// check if the model is used: firstly to check instance
	exists, err := m.dependent.HasInstance(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some instances, error info is %s", kit.Rid, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some instances.", kit.Rid, targetObjIDS)
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete)
	}

	// check if the model is used: secondly to check association
	exists, err = m.dependent.HasAssociation(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some associations, error info is %s", kit.Rid, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some associations.", kit.Rid, targetObjIDS)
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	// delete model self
	cnt, err := m.deleteModelAndAttributes(kit, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the models (%#v) and their attributes, error info is %s", kit.Rid, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

// CascadeDeleteModel 将会删除模型/模型属性/属性分组/唯一校验
func (m *modelManager) CascadeDeleteModel(kit *rest.Kit, modelID int64) (*metadata.DeletedCount, error) {

	deleteCondMap := util.SetQueryOwner(make(map[string]interface{}), kit.SupplierAccount)
	deleteCond, _ := mongo.NewConditionFromMapStr(deleteCondMap)
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldID, Val: modelID})

	// read all models by the deletion condition
	cnt, err := m.cascadeDelete(kit, deleteCond)
	if nil != err {
		blog.ErrorJSON("CascadeDeleteModel failed, cascadeDelete failed, condition: %s, err: %s, rid: %s", deleteCond.ToMapStr(), err.Error(), kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, err
}

func (m *modelManager) SearchModel(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult, error) {

	dataResult := &metadata.QueryModelDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return dataResult, err
	}

	totalCount, err := m.count(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to get the count by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	modelItems, err := m.search(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	dataResult.Info = modelItems
	return dataResult, nil
}

func (m *modelManager) SearchModelWithAttribute(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelWithAttributeDataResult, error) {

	dataResult := &metadata.QueryModelWithAttributeDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return dataResult, err
	}

	totalCount, err := m.count(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to get the count by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	modelItems, err := m.search(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	for _, modelItem := range modelItems {
		queryAttributeCondMap := util.SetQueryOwner(make(map[string]interface{}), modelItem.OwnerID)
		queryAttributeCond, _ := mongo.NewConditionFromMapStr(queryAttributeCondMap)
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelItem.ObjectID))
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldSupplierAccount).Eq(modelItem.OwnerID))
		attributeItems, err := m.modelAttribute.search(kit, queryAttributeCond)
		if nil != err {
			blog.Errorf("request(%s):it is failed to search the object(%s)'s attributes, error info is %s", modelItem.ObjectID, err.Error())
			return dataResult, err
		}
		dataResult.Info = append(dataResult.Info, metadata.SearchModelInfo{Spec: modelItem, Attributes: attributeItems})
	}

	return dataResult, nil
}
