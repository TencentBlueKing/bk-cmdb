/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.ModelOperation = (*modelManager)(nil)

type modelManager struct {
	*modelAttribute
	*modelClassification
	dbProxy   dal.RDB
	dependent OperationDependences
}

// New create a new model manager instance
func New(dbProxy dal.RDB, dependent OperationDependences) core.ModelOperation {

	coreMgr := &modelManager{dbProxy: dbProxy, dependent: dependent}

	coreMgr.modelAttribute = &modelAttribute{dbProxy: dbProxy, model: coreMgr}
	coreMgr.modelClassification = &modelClassification{dbProxy: dbProxy, model: coreMgr}

	return coreMgr
}

func (m *modelManager) CreateModel(ctx core.ContextParams, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error) {

	condCheckModel := mongo.NewCondition()
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	_, exists, err := m.isExists(ctx, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the model (%s) is exists, error info is %s ", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}

	if exists {
		return &metadata.CreateOneDataResult{}, ctx.Error.Error(common.CCErrCommDuplicateItem)
	}

	id, err := m.save(ctx, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the model (%v), error info is %s", ctx.ReqID, inputParam.Spec, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}

	_, err = m.modelAttribute.CreateModelAttributes(ctx, inputParam.Spec.ObjectID, metadata.CreateModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to create some attributes (%v) for the model (%s), error info is %s", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID)
		return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
	}

	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}
func (m *modelManager) SetModel(ctx core.ContextParams, inputParam metadata.SetModel) (*metadata.SetDataResult, error) {

	condCheckModel := mongo.NewCondition()
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})
	condCheckModel.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	existsModel, exists, err := m.isExists(ctx, condCheckModel)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check the model (%s) is exists, error info is %s ", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
		return &metadata.SetDataResult{}, err
	}

	dataResult := &metadata.SetDataResult{}

	// set model spec
	if !exists {
		updateCond := mongo.NewCondition()
		updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
		updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: inputParam.Spec.ObjectID})

		_, err := m.update(ctx, mapstr.NewFromStruct(inputParam.Spec, "field"), updateCond)
		if nil != err {
			blog.Errorf("request(%s): it is failed to update some fields (%v) for the model (%s), error info is %s", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}

		dataResult.UpdatedCount.Count++
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{OriginIndex: 0, ID: uint64(existsModel.ID)})
	} else {
		id, err := m.save(ctx, &inputParam.Spec)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the model (%v), error info is %s", ctx.ReqID, inputParam.Spec.ObjectID, err.Error())
			return dataResult, err
		}
		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{OriginIndex: 0, ID: id})
	}

	// set model attributes
	setAttrResult, err := m.modelAttribute.SetModelAttributes(ctx, inputParam.Spec.ObjectID, metadata.SetModelAttributes{Attributes: inputParam.Attributes})
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the attributes (%v) for the model (%s), error info is %s", ctx.ReqID, inputParam.Attributes, inputParam.Spec.ObjectID, err.Error())
		return dataResult, err
	}

	// set attribute result, ignore model operation result
	dataResult.CreatedCount = setAttrResult.CreatedCount
	dataResult.UpdatedCount = setAttrResult.UpdatedCount

	dataResult.Created = append(dataResult.Created, setAttrResult.Created...)
	dataResult.Updated = append(dataResult.Updated, setAttrResult.Updated...)
	dataResult.Exceptions = append(dataResult.Exceptions, setAttrResult.Exceptions...)

	return dataResult, err
}

func (m *modelManager) UpdateModel(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	updateCond, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%v) from mapstr into condition object, error info is %s ", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}
	updateCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})

	cnt, err := m.update(ctx, inputParam.Data, updateCond)
	return &metadata.UpdatedCount{Count: cnt}, err
}

func (m *modelManager) DeleteModel(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all models by the deletion conditon
	deleteCond, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrCommParamsInvalid, err.Error())
	}

	modelItems := []metadata.ObjectDes{}
	err = m.dbProxy.Table(common.BKTableNameObjDes).Find(deleteCond.ToMapStr()).All(ctx, &modelItems)
	if nil != err {
		blog.Errorf("request(%s): it is failed to find the all models by the condition (%v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	targetObjIDS := []string{}
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}

	// check if the model is used: firstly to check instance
	exists, err := m.dependent.HasInstance(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%v) has some instances, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%v), because of they have some instances.", ctx.ReqID, targetObjIDS)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete)
	}

	// check if the model is used: secondly to check association
	exists, err = m.dependent.HasAssociation(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether some models (%v) has some associations, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if exists {
		blog.Warnf("request(%s): it is forbidden to delete the models (%v), because of they have some associations.", ctx.ReqID, targetObjIDS)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	// delete model self
	cnt, err := m.deleteModelAndAttributes(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the models (%v) and their attributes, error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *modelManager) CascadeDeleteModel(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all models by the deletion condition
	deleteCond, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrCommParamsInvalid, err.Error())
	}
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: ctx.SupplierAccount})

	cnt, err := m.cascadeDelete(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to cascade delete some models by the condition (%v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, err
}

func (m *modelManager) SearchModel(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {

	searchCond, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.QueryResult{}, err
	}

	dataResult, err := m.searchReturnMapStr(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is faield to search models by the condition (%v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return &metadata.QueryResult{}, err
	}
	return &metadata.QueryResult{Count: uint64(len(dataResult)), Info: dataResult}, nil
}
