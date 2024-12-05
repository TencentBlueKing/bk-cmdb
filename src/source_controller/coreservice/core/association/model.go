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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

type associationModel struct {
}

// CreateModelAssociation TODO
func (m *associationModel) CreateModelAssociation(kit *rest.Kit,
	inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {
	enableMainlineAssociationType := false
	return m.createModelAssociation(kit, inputParam, enableMainlineAssociationType)
}

// CreateMainlineModelAssociation used for create association of type bk_mainline, as it can only create by special method,
// for example add a level to business modle
func (m *associationModel) CreateMainlineModelAssociation(kit *rest.Kit,
	inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {
	enableMainlineAssociationType := true
	return m.createModelAssociation(kit, inputParam, enableMainlineAssociationType)
}

var forbiddenCreateAssociationObjList = []string{
	common.BKInnerObjIDProject,
}

func (m *associationModel) createModelAssociation(kit *rest.Kit, inputParam metadata.CreateModelAssociation,
	enableMainlineAssociationType bool) (*metadata.CreateOneDataResult, error) {
	// enableMainlineAssociationType used for distinguish two creation mode
	// when enableMainlineAssociationType enabled, only bk_mainline type could be create
	// when enableMainlineAssociationType disabled, all type except bk_mainline could be create

	inputParam.Spec.TenantID = kit.TenantID
	if err := m.isValid(kit, inputParam); err != nil {
		return &metadata.CreateOneDataResult{}, err
	}

	exists, err := m.isExistsAssociationID(kit, inputParam.Spec.AssociationName)
	if err != nil {
		blog.Errorf("failed to check whether the association ID (%s) is exists, error: %v, rid: %s",
			inputParam.Spec.AssociationName, err, kit.Rid)
		return &metadata.CreateOneDataResult{}, err
	}
	if exists {
		blog.Errorf("the association ID (%s) is exists, rid: %s", inputParam.Spec.AsstKindID, kit.Rid)
		return &metadata.CreateOneDataResult{}, kit.CCError.Errorf(common.CCErrCommDuplicateItem,
			inputParam.Spec.AssociationName)
	}

	exists, err = m.isExistsAssociationObjectWithAnotherObject(kit, inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID,
		inputParam.Spec.AsstKindID)
	if err != nil {
		blog.Errorf("failed to check if the association (%s=>%s) is exists, error: %v, rid: %s",
			inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID, err, kit.Rid)
		return &metadata.CreateOneDataResult{}, err
	}
	if exists {
		blog.Errorf("(%s=>%s) is exists, rid: %s", inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID, kit.Rid)
		return &metadata.CreateOneDataResult{}, kit.CCError.Errorf(common.CCErrTopoAssociationAlreadyExist,
			inputParam.Spec.ObjectID, inputParam.Spec.AsstObjID)
	}

	asstKindID := inputParam.Spec.AsstKindID
	if enableMainlineAssociationType == false {
		// AsstKindID shouldn't be use bk_mainline
		if asstKindID == common.AssociationKindMainline {
			blog.Errorf("use inner association type: %v is forbidden, rid: %s", common.AssociationKindMainline,
				kit.Rid)
			return &metadata.CreateOneDataResult{}, kit.CCError.Errorf(
				common.CCErrorTopoAssociationKindMainlineUnavailable, asstKindID)
		}
	} else {
		// AsstKindID could only be bk_mainline
		if asstKindID != common.AssociationKindMainline {
			blog.Errorf("use CreateMainlineObjectAssociation method but bk_asst_id is: %s, rid: %s", asstKindID,
				kit.Rid)
			return &metadata.CreateOneDataResult{}, kit.CCError.Errorf(common.CCErrorTopoAssociationKindInconsistent,
				asstKindID)
		}
	}

	id, err := m.save(kit, &inputParam.Spec)
	if err != nil {
		blog.Errorf("failed to create a new association (%s=>%s), error: %v", inputParam.Spec.ObjectID,
			inputParam.Spec.AsstObjID, err, kit.Rid)
		return &metadata.CreateOneDataResult{}, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

// SetModelAssociation TODO
func (m *associationModel) SetModelAssociation(kit *rest.Kit,
	inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error) {

	// TODO: need to care instance association, which used this model association

	return nil, nil
}

// UpdateModelAssociation TODO
func (m *associationModel) UpdateModelAssociation(kit *rest.Kit,
	inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	// ATTENTION: only to update the fields except bk_obj_asst_id, bk_obj_id, bk_asst_obj_id
	inputParam.Data.Remove(metadata.AssociationFieldObjectID)
	inputParam.Data.Remove(metadata.AssociationFieldAssociationObjectID)
	inputParam.Data.Remove(common.TenantID)
	inputParam.Data.Remove(metadata.AssociationFieldAsstID)

	updateCond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("update association failed, err: %v, condition: %v, rid: %s", err,
			inputParam.Condition, kit.Rid)
		return &metadata.UpdatedCount{}, kit.CCError.New(common.CCErrCommPostInputParseError, err.Error())
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
		blog.Warnf("update object association got invalid fields: %v, rid: %s", filterOutFields, kit.Rid)
	}

	cnt, err := m.update(kit, validData, updateCond)
	if err != nil {
		blog.Errorf("update association failed, err: %v, condition: %v, rid: %s", err, updateCond, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

// SearchModelAssociation search model associations
func (m *associationModel) SearchModelAssociation(kit *rest.Kit,
	inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {

	searchCond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("convert the condition from mapstr failed, err: %v, condition: %v, rid: %s", err,
			inputParam.Condition, kit.Rid)
		return &metadata.QueryResult{}, kit.CCError.New(common.CCErrCommPostInputParseError, err.Error())
	}

	resultItems, err := m.searchReturnMapStr(kit, searchCond)
	if err != nil {
		blog.Errorf("search all associations failed, err: %v, condition: %v, rid: %s", err, searchCond, kit.Rid)
		return &metadata.QueryResult{}, err
	}

	return &metadata.QueryResult{Count: uint64(len(resultItems)), Info: resultItems}, nil
}

// CountModelAssociations counts target model associations num
func (m *associationModel) CountModelAssociations(kit *rest.Kit, input *metadata.Condition) (
	*metadata.CommonCountResult, error) {

	cond, err := mongo.NewConditionFromMapStr(input.Condition)
	if err != nil {
		blog.Errorf("convert the condition from mapstr failed, err: %v, cond: %v, rid: %s", err, input.Condition,
			kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommPostInputParseError, err.Error())
	}

	count, err := m.count(kit, cond)
	if err != nil {
		blog.Errorf("count model associations failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	return &metadata.CommonCountResult{Count: count}, nil
}

// DeleteModelAssociation TODO
func (m *associationModel) DeleteModelAssociation(kit *rest.Kit,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all model associations
	deleteCond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("convert the condition from mapstr failed, err: %v, condition: %v, rid: %s", err,
			inputParam.Condition, kit.Rid)
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrCommPostInputParseError, err.Error())
	}

	needDeleteAssocaitionItems, err := m.search(kit, deleteCond)
	if err != nil {
		blog.Errorf("search all failed, err: %v, condition: %v, rid: %s", err, deleteCond, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	// if the model association was already used in instance association, then the deletion operation must be abandoned
	associationIDS := []string{}
	for _, assocaitionItem := range needDeleteAssocaitionItems {
		associationIDS = append(associationIDS, assocaitionItem.AssociationName)
	}

	exists, err := m.usedInSomeInstanceAssociation(kit, associationIDS)
	if err != nil {
		blog.Errorf("check if the instances is in used failed, err: %v, IDs: %v, rid: %s", err, associationIDS,
			kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	if exists {
		blog.Warnf("it is forbbiden to delete the model association by the instances, IDs: %v, rid: %s",
			associationIDS, kit.Rid)
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// deletion operation
	cnt, err := m.delete(kit, deleteCond)
	if err != nil {
		blog.Errorf("delete the instances failed, err: %v, cond: %v, rid: %s", err, deleteCond, kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}

// CascadeDeleteModelAssociation TODO
func (m *associationModel) CascadeDeleteModelAssociation(kit *rest.Kit,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	// read all model associations
	deleteCond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("convert the condition from mapstr failed, err: %v, condition: %v, rid: %s", err,
			inputParam.Condition, kit.Rid)
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrCommPostInputParseError, err.Error())
	}

	needDeleteAssocaitionItems, err := m.search(kit, deleteCond)
	if err != nil {
		blog.Errorf("search associations failed, err: %v, condition: %v, rid: %s", err, deleteCond, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	// if the model association was already used in instance association, then the deletion operation must be abandoned
	associationIDS := []string{}
	for _, assocaitionItem := range needDeleteAssocaitionItems {
		associationIDS = append(associationIDS, assocaitionItem.AssociationName)
	}

	// cascade deletion operation
	if err := m.cascadeInstanceAssociation(kit, associationIDS); err != nil {
		blog.Errorf("cascade delete the associations failed, err: %v, condition: %v, rid: %s", err, deleteCond,
			kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	// deletion operation
	cnt, err := m.delete(kit, deleteCond)
	if err != nil {
		blog.Errorf("delete associations failed, err: %v, condition: %v, rid: %s", err, deleteCond, kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}
