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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.ModelOperation = (*modelManager)(nil)

type modelManager struct {
	*modelAttributeGroup
	*modelAttribute
	*modelClassification
	*modelAttrUnique
	dbProxy   dal.RDB
	dependent OperationDependences
}

// New create a new model manager instance
func New(dbProxy dal.RDB, dependent OperationDependences) core.ModelOperation {

	coreMgr := &modelManager{dbProxy: dbProxy, dependent: dependent}

	coreMgr.modelAttribute = &modelAttribute{dbProxy: dbProxy, model: coreMgr}
	coreMgr.modelClassification = &modelClassification{dbProxy: dbProxy, model: coreMgr}
	coreMgr.modelAttributeGroup = &modelAttributeGroup{dbProxy: dbProxy, model: coreMgr}
	coreMgr.modelAttrUnique = &modelAttrUnique{dbProxy: dbProxy}

	return coreMgr
}

func (m *modelManager) CreateModel(ctx core.ContextParams, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error) {

	dataResult := &metadata.CreateOneDataResult{}

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set", ctx.ReqID, inputParam.Spec.ObjectID)
		return dataResult, ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(ctx, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s", ctx.ReqID, inputParam.Spec.ObjCls, err.Error())
		return dataResult, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid", ctx.ReqID, inputParam.Spec.ObjCls)
		return dataResult, ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	// check the model if it is exists
	condCheckModel := mongo.NewCondition()
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	// ATTENTION: Currently only business dimension isolation is done,
	//           and there may be isolation requirements for other dimensions in the future.
	isExist, bizID := inputParam.Spec.Metadata.Label.Get(common.BKAppIDField)
	if isExist {
		_, metaCond := condCheckModel.Embed(metadata.BKMetadata)
		_, labelCond := metaCond.Embed(metadata.BKLabel)
		labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}

	_, exists, err := m.isExists(ctx, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the model (%s) is exists, error info is %s ", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
		return dataResult, err
	}
	if exists {
		blog.Warnf("request(%s): it is failed to  create a new model , because of the model (%s) is already exists ", ctx.ReqID, inputParam.Spec.ObjectID)
		return dataResult, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectID)
	}

	// 检查模型名称重复
	modelNameUniqueFilter := map[string]interface{}{
		common.BKObjNameField: inputParam.Spec.ObjectName,
	}
	bizFilter := metadata.PublicAndBizCondition(inputParam.Spec.Metadata)
	for key, value := range bizFilter {
		modelNameUniqueFilter[key] = value
	}
	sameNameCount, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(ctx)
	if err != nil {
		blog.Errorf("whether same name model exists, name: %s, err: %s, rid: %s", inputParam.Spec.ObjectName, err.Error(), ctx.ReqID)
		return dataResult, err
	}
	if sameNameCount > 0 {
		blog.Warnf("create model failed, field `%s` duplicated, rid: %s", inputParam.Spec.ObjectName, ctx.ReqID)
		return dataResult, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.ObjectName)
	}

	inputParam.Spec.OwnerID = ctx.SupplierAccount
	id, err := m.save(ctx, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", ctx.ReqID, inputParam.Spec, err.Error())
		return dataResult, err
	}

	_, err = m.modelAttribute.CreateModelAttributes(ctx, inputParam.Spec.ObjectID, metadata.CreateModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to create some attributes (%#v) for the model (%s), err: %v", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID, err)
		return dataResult, err
	}
	dataResult.Created.ID = id
	return dataResult, nil
}
func (m *modelManager) SetModel(ctx core.ContextParams, inputParam metadata.SetModel) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	// check the model attributes value
	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model, because of the modelID (%s) is not set", ctx.ReqID, inputParam.Spec.ObjectID)
		return dataResult, ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjectID)
	}

	// check the input classification ID
	isValid, err := m.modelClassification.isValid(ctx, inputParam.Spec.ObjCls)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classificationID(%s) is invalid, error info is %s", ctx.ReqID, inputParam.Spec.ObjCls, err.Error())
		return dataResult, err
	}

	if !isValid {
		blog.Warnf("request(%s): it is failed to create a new model, because of the classificationID (%s) is invalid", ctx.ReqID, inputParam.Spec.ObjCls)
		return dataResult, ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.ClassificationFieldID)
	}

	condCheckModel := mongo.NewCondition()
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	existsModel, exists, err := m.isExists(ctx, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check the model (%s) is exists, error info is %s ", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
		return &metadata.SetDataResult{}, err
	}

	inputParam.Spec.OwnerID = ctx.SupplierAccount
	// set model spec
	if exists {
		updateCond := mongo.NewCondition()
		updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
		updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

		_, err := m.update(ctx, mapstr.NewFromStruct(inputParam.Spec, "field"), updateCond)
		if nil != err {
			blog.Errorf("request(%s): it is failed to update some fields (%#v) for the model (%s), error info is %s", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}

		dataResult.UpdatedCount.Count++
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{OriginIndex: 0, ID: uint64(existsModel.ID)})
	} else {
		id, err := m.save(ctx, &inputParam.Spec)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the model (%#v), error info is %s", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}
		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{OriginIndex: 0, ID: id})
	}

	// set model attributes
	setAttrResult, err := m.modelAttribute.SetModelAttributes(ctx, inputParam.Spec.ObjectID, metadata.SetModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the attributes (%#v) for the model (%s), error info is %s", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
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

func (m *modelManager) UpdateModel(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	updateCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s ", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}
	updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	cnt, err := m.update(ctx, inputParam.Data, updateCond)
	return &metadata.UpdatedCount{Count: cnt}, err
}

func (m *modelManager) DeleteModel(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all models by the deletion condition
	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrCommParamsInvalid, err.Error())
	}

	modelItems, err := m.search(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to find the all models by the condition (%#v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	targetObjIDS := make([]string, 0)
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}

	// check if the model is used: firstly to check instance
	exists, err := m.dependent.HasInstance(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some instances, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some instances.", ctx.ReqID, targetObjIDS)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete)
	}

	// check if the model is used: secondly to check association
	exists, err = m.dependent.HasAssociation(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%#v) has some associations, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%#v), because of they have some associations.", ctx.ReqID, targetObjIDS)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	// delete model self
	cnt, err := m.deleteModelAndAttributes(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the models (%#v) and their attributes, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

// CascadeDeleteModel 将会删除模型/模型属性/属性分组/唯一校验
func (m *modelManager) CascadeDeleteModel(ctx core.ContextParams, modelID int64) (*metadata.DeletedCount, error) {

	deleteCond := mongo.NewCondition()
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldID, Val: modelID})

	// read all models by the deletion condition
	cnt, err := m.cascadeDelete(ctx, deleteCond)
	if nil != err {
		blog.ErrorJSON("CascadeDeleteModel failed, cascadeDelete failed, condition: %s, err: %s, rid: %s", deleteCond.ToMapStr(), err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, err
}

func (m *modelManager) SearchModel(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult, error) {

	dataResult := &metadata.QueryModelDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return dataResult, err
	}

	totalCount, err := m.count(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to get the count by the condition (%#v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	modelItems, err := m.search(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	dataResult.Info = modelItems
	return dataResult, nil
}

func (m *modelManager) SearchModelWithAttribute(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelWithAttributeDataResult, error) {

	dataResult := &metadata.QueryModelWithAttributeDataResult{}

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return dataResult, err
	}

	totalCount, err := m.count(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to get the count by the condition (%#v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	modelItems, err := m.search(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search models by the condition (%#v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	for _, modelItem := range modelItems {

		queryAttributeCond := mongo.NewCondition()
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelItem.ObjectID))
		queryAttributeCond.Element(mongo.Field(metadata.AttributeFieldSupplierAccount).Eq(modelItem.OwnerID))
		attributeItems, err := m.modelAttribute.search(ctx, queryAttributeCond)
		if nil != err {
			blog.Errorf("request(%s):it is failed to search the object(%s)'s attributes, error info is %s", modelItem.ObjectID, err.Error())
			return dataResult, err
		}
		dataResult.Info = append(dataResult.Info, metadata.SearchModelInfo{Spec: modelItem, Attributes: attributeItems})
	}

	return dataResult, nil
}
