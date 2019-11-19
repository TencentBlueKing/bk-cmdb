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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type associationModel struct {
	dbProxy dal.RDB
}

func (m *associationModel) CreateModelAssociation(ctx core.ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {
	enableMainlineAssociationType := false
	return m.createModelAssociation(ctx, inputParam, enableMainlineAssociationType)
}

// CreateMainlineModelAssociation used for create association of type bk_mainline, as it can only create by special method,
// for example add a level to business modle
func (m *associationModel) CreateMainlineModelAssociation(ctx core.ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {
	enableMainlineAssociationType := true
	return m.createModelAssociation(ctx, inputParam, enableMainlineAssociationType)
}

func (m *associationModel) createModelAssociation(ctx core.ContextParams, inputParam metadata.CreateModelAssociation, enableMainlineAssociationType bool) (*metadata.CreateOneDataResult, error) {
	// enableMainlineAssociationType used for distinguish two creation mode
	// when enableMainlineAssociationType enabled, only bk_mainline type could be create
	// when enableMainlineAssociationType disabled, all type except bk_mainline could be create

	inputParam.Spec.OwnerID = ctx.SupplierAccount
	if err := m.isValid(ctx, inputParam); nil != err {
		return &metadata.CreateOneDataResult{}, err
	}

	inputParam.Spec.OwnerID = ctx.SupplierAccount

	exists, err := m.isExistsAssociationID(ctx, inputParam.Spec.AssociationName)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the association ID (%s) is exists, error info is %s", ctx.ReqID, inputParam.Spec.AssociationName, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}
	if exists {
		blog.Warnf("request(%s): it is failed create a new association, because of the association ID (%s) is exists", ctx.ReqID, inputParam.Spec.AsstKindID)
		return &metadata.CreateOneDataResult{}, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Spec.AssociationName)
	}

	exists, err = m.isExistsAssociationObjectWithAnotherObject(ctx, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID, inputParam.Spec.AsstKindID)
	if nil != err {
		blog.Errorf("request(%s): it is failed to create a new association, because of it is failed to check if the association (%s=>%s) is exists, error info is %s", ctx.ReqID, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}
	if exists {
		blog.Warnf("request(%s): it is failed to create a new association, because of it (%s=>%s) is exists", ctx.ReqID, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID)
		return &metadata.CreateOneDataResult{}, ctx.Error.Errorf(common.CCErrTopoAssociationAlreadyExist, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID)
	}

	asstKindID := inputParam.Spec.AsstKindID
	if enableMainlineAssociationType == false {
		// AsstKindID shouldn't be use bk_mainline
		if asstKindID == common.AssociationKindMainline {
			blog.Errorf("use inner association type: %v is forbidden, rid: %s", common.AssociationKindMainline, ctx.ReqID)
			return &metadata.CreateOneDataResult{}, ctx.Error.Errorf(common.CCErrorTopoAssociationKindMainlineUnavailable, asstKindID)
		}
	} else {
		// AsstKindID could only be bk_mainline
		if asstKindID != common.AssociationKindMainline {
			blog.Errorf("use CreateMainlineObjectAssociation method but bk_asst_id is: %s, rid: %s", asstKindID, ctx.ReqID)
			return &metadata.CreateOneDataResult{}, ctx.Error.Errorf(common.CCErrorTopoAssociationKindInconsistent, asstKindID)
		}
	}

	id, err := m.save(ctx, &inputParam.Spec)
	if nil != err {
		blog.Errorf("request(%s): it is failed to create a new association (%s=>%s), error info is %s", ctx.ReqID, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

func (m *associationModel) SetModelAssociation(ctx core.ContextParams, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error) {

	// TODO: need to care instance association, which used this model association

	return nil, nil
}

func (m *associationModel) UpdateModelAssociation(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	// ATTENTION: only to update the fields except bk_obj_asst_id, bk_obj_id, bk_asst_obj_id
	inputParam.Data.Remove(metadata.AssociationFieldObjectID)
	inputParam.Data.Remove(metadata.AssociationFieldAssociationObjectID)
	inputParam.Data.Remove(metadata.AssociationFieldSupplierAccount)
	inputParam.Data.Remove(metadata.AssociationFieldAsstID)

	updateCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is to failed to update the association by the condition (%v), error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, ctx.Error.New(common.CCErrCommPostInputParseError, err.Error())
	}

	// only field in white list could be update
	// bk_asst_obj_id is allowed for add business model level
	validFields := []string{"bk_obj_asst_name", "bk_asst_obj_id"}
	validData := map[string]interface{}{}
	filterOutFields := []string{}
	for key, val := range inputParam.Data {
		if isValidField := util.Contains(validFields, key); isValidField == false {
			filterOutFields = append(filterOutFields, key)
			continue
		}
		validData[key] = val
	}

	if len(filterOutFields) > 0 {
		blog.Warnf("update object association got invalid fields: %v, rid: %s", filterOutFields, ctx.ReqID)
	}

	cnt, err := m.update(ctx, validData, updateCond)
	if nil != err {
		blog.Errorf("request(%s): it is to update the association by the condition (%#v), error info is %s", ctx.ReqID, updateCond.ToMapStr(), err.Error())
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *associationModel) SearchModelAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {

	searchCond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is to convert the condition (%v) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.QueryResult{}, ctx.Error.New(common.CCErrCommPostInputParseError, err.Error())
	}

	resultItems, err := m.searchReturnMapStr(ctx, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is to search all associations by the condition (%#v), error info is %s", ctx.ReqID, searchCond.ToMapStr(), err.Error())
		return &metadata.QueryResult{}, err
	}

	return &metadata.QueryResult{Count: uint64(len(resultItems)), Info: resultItems}, nil
}

func (m *associationModel) DeleteModelAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all model associations
	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is to convert the condition (%s) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrCommPostInputParseError, err.Error())
	}

	needDeleteAssocaitionItems, err := m.search(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search all by the condition (%#v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	// if the model association was already used in instance association, then the deletion operation must be abandoned
	associationIDS := []string{}
	for _, assocaitionItem := range needDeleteAssocaitionItems {
		associationIDS = append(associationIDS, assocaitionItem.AssociationName)
	}

	exists, err := m.usedInSomeInstanceAssociation(ctx, associationIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check if the instances (%#v) is in used, error info is %s", ctx.ReqID, associationIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}
	if exists {
		blog.Warnf("request(%s): it is forbbiden to delete the model association by the instances (%#v)", ctx.ReqID, associationIDS)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// deletion operation
	cnt, err := m.delete(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is delete the instances by the condition (%#v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *associationModel) CascadeDeleteModelAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all model associations
	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is to convert the condition (%s) from mapstr into condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, ctx.Error.New(common.CCErrCommPostInputParseError, err.Error())
	}

	needDeleteAssocaitionItems, err := m.search(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is to search associations by the condition (%#v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	// if the model association was already used in instance association, then the deletion operation must be abandoned
	associationIDS := []string{}
	for _, assocaitionItem := range needDeleteAssocaitionItems {
		associationIDS = append(associationIDS, assocaitionItem.AssociationName)
	}

	// cascade deletion operation
	if err := m.cascadeInstanceAssociation(ctx, associationIDS); nil != err {
		blog.Errorf("request(%s): it is failed to cascade delete the assocaitions of the instances (%#v), error info is %s ", ctx.ReqID, associationIDS, err.Error())
		return &metadata.DeletedCount{}, err
	}

	// deletion operation
	cnt, err := m.delete(ctx, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is to delete some associations by the condition (%#v), error info is %s", ctx.ReqID, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}
