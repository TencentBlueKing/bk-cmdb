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

package model

import (
	"strconv"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/logics/inst"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	DeleteAssociationType(kit *rest.Kit, asstTypeID int64) error
	// CreateOrUpdateAssociationType only allow import api to use
	CreateOrUpdateAssociationType(kit *rest.Kit, asst []metadata.AssociationKind) error
	CreateCommonAssociation(kit *rest.Kit, data *metadata.Association) (*metadata.Association, error)
	DeleteAssociationWithPreCheck(kit *rest.Kit, associationID int64) error
	UpdateObjectAssociation(kit *rest.Kit, data mapstr.MapStr, assoID int64) error
	SearchObjectAssocWithAssocKindList(kit *rest.Kit, asstKindIDs []string) (resp *metadata.AssociationList, err error)

	// CreateMainlineAssociation TODO
	// Mainline
	// CreateMainlineAssociation create mainline object association
	CreateMainlineAssociation(kit *rest.Kit, data *metadata.MainlineAssociation) (
		*metadata.Object, error)
	// DeleteMainlineAssociation delete mainline association by objID
	DeleteMainlineAssociation(kit *rest.Kit, objID string) error
	// SearchMainlineAssociationTopo get mainline topo of special model
	SearchMainlineAssociationTopo(kit *rest.Kit, targetObjID string) ([]*metadata.MainlineObjectTopo, error)
	// IsMainlineObject check whether objID is mainline object or not
	IsMainlineObject(kit *rest.Kit, objID string) (bool, error)

	// SetProxy proxy the interface
	SetProxy(object ObjectOperationInterface, inst inst.InstOperationInterface,
		instasst inst.AssociationOperationInterface)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) AssociationOperationInterface {
	return &association{
		clientSet:   client,
		authManager: authManager,
	}
}

type association struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	obj         ObjectOperationInterface
	inst        inst.InstOperationInterface
	instasst    inst.AssociationOperationInterface
}

// SetProxy proxy the interface
func (assoc *association) SetProxy(object ObjectOperationInterface, inst inst.InstOperationInterface,
	instasst inst.AssociationOperationInterface) {

	assoc.obj = object
	assoc.inst = inst
	assoc.instasst = instasst
}

// DeleteAssociationType delete association type except bk_mainline
func (assoc *association) DeleteAssociationType(kit *rest.Kit, asstTypeID int64) error {

	input := &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKFieldID: asstTypeID}}
	rsp, err := assoc.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search association kind by typeID[%d], but get detailed info failed, err: %v, rid: %s",
			asstTypeID, err, kit.Rid)
		return err
	}

	if len(rsp.Info) > 1 {
		blog.Errorf("search association kind by typeID[%d], but get multiple instance, rid: %s", asstTypeID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoGetMultipleAssocKindInstWithOneID)
	}

	if len(rsp.Info) == 0 {
		return nil
	}

	if rsp.Info[0].IsPre != nil && *rsp.Info[0].IsPre {
		blog.Errorf("typeID[%d] of association kind is a pre-defined association kind, rid: %s", asstTypeID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoDeletePredefinedAssociationKind)
	}

	// a already used association kind can not be deleted.
	assoFilter := []map[string]interface{}{mapstr.MapStr{common.AssociationKindIDField: rsp.Info[0].AssociationKindID}}
	asso, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjAsst,
		assoFilter)
	if err != nil {
		blog.Errorf("search association by association kind[%d] failed, err: %v, rid: %s",
			rsp.Info[0].AssociationKindID, err, kit.Rid)
		return err
	}

	if asso[0] != 0 {
		blog.Warnf("association kind[%d] has already been used, can not be deleted, rid: %s",
			rsp.Info[0].AssociationKindID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoAssociationKindHasBeenUsed)
	}

	deleteCond := &metadata.DeleteOption{Condition: mapstr.MapStr{common.BKFieldID: asstTypeID}}
	_, err = assoc.clientSet.CoreService().Association().DeleteAssociationType(kit.Ctx, kit.Header, deleteCond)
	if err != nil {
		blog.Errorf("delete association type failed, kind id: %d, err: %v, rid: %s", asstTypeID, err, kit.Rid)
		return err
	}

	return nil
}

// CreateOrUpdateAssociationType only allow import api to use
// CreateOrUpdateAssociationType create association type, if association type exist, update it
func (assoc *association) CreateOrUpdateAssociationType(kit *rest.Kit, asst []metadata.AssociationKind) error {

	if len(asst) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "asst")
	}

	asstKindMap := make(map[string]metadata.AssociationKind)
	asstKindIDMap := make(map[string]struct{}, 0)
	asstKindID := make([]string, 0)
	for index, item := range asst {
		asstKindMap[item.AssociationKindID] = asst[index]
		asstKindIDMap[item.AssociationKindID] = struct{}{}
		asstKindID = append(asstKindID, item.AssociationKindID)
	}

	asstQuery := &metadata.QueryCondition{
		Condition:      map[string]interface{}{common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKindID}},
		DisableCounter: true,
	}
	rsp, err := assoc.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, asstQuery)
	if err != nil {
		blog.Errorf("search asstkind failed, cond: %v, err: %v, rid: %s", asstQuery, err, kit.Rid)
		return err
	}

	updateCond := &metadata.UpdateOption{}
	for _, item := range rsp.Info {
		delete(asstKindIDMap, item.AssociationKindID)

		data := asstKindMap[item.AssociationKindID]
		dataMapstr := data.ToMapStr()
		delete(dataMapstr, common.BKFieldID)
		updateCond.Condition = mapstr.MapStr{common.AssociationKindIDField: item}
		updateCond.Data = dataMapstr

		_, err := assoc.clientSet.CoreService().Association().UpdateAssociationType(kit.Ctx, kit.Header, updateCond)
		if err != nil {
			blog.Errorf("update asstkind failed, cond: %+v, err: %v, rid: %s", updateCond, err, kit.Rid)
			return err
		}
	}

	if len(asstKindIDMap) == 0 {
		return nil
	}

	createCond := &metadata.CreateManyAssociationKind{Datas: make([]metadata.AssociationKind, 0)}
	for key := range asstKindIDMap {
		createCond.Datas = append(createCond.Datas, asstKindMap[key])
	}

	createRsp, err := assoc.clientSet.CoreService().Association().CreateManyAssociation(kit.Ctx, kit.Header, createCond)
	if err != nil {
		blog.Errorf("create asstkind failed, cond: %+v, err: %v, rid: %s", createCond, err, kit.Rid)
		return err
	}

	if len(createRsp.Repeated) > 0 {
		blog.Errorf("asst kind repeated, data: %+v, rid: %s", createRsp.Repeated, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDuplicateItem)
	}

	if len(createRsp.Exceptions) > 0 {
		blog.Errorf("asst kind failed, data: %+v, rid: %s", createRsp.Exceptions, kit.Rid)
		return kit.CCError.CCErrorf(int(createRsp.Exceptions[0].Code), createRsp.Exceptions[0].Message)
	}

	// register association type resource creator action to iam
	if auth.EnableAuthorize() {
		indexID := make(map[int64]int64)
		for _, item := range createRsp.Created {
			indexID[item.OriginIndex] = int64(item.ID)
		}

		for index, item := range createCond.Datas {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.SysAssociationType),
				ID:      strconv.FormatInt(indexID[int64(index)], 10),
				Name:    item.AssociationKindName,
				Creator: kit.User,
			}
			if _, err = assoc.authManager.Authorizer.RegisterResourceCreatorAction(kit.Ctx, kit.Header,
				iamInstance); err != nil {
				blog.Errorf("register created association type to iam failed, err: %v, rid: %s", err, kit.Rid)
				return err
			}
		}
	}

	return nil
}

// CreateCommonAssociation create common object association
func (assoc *association) CreateCommonAssociation(kit *rest.Kit, data *metadata.Association) (*metadata.Association,
	error) {

	if data.AsstKindID == common.AssociationKindMainline {
		return nil, kit.CCError.CCError(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}
	if len(data.AsstKindID) == 0 || len(data.AsstObjID) == 0 || len(data.ObjectID) == 0 {
		blog.Errorf("association kind id、 associate/object id is required, rid: %s", kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoAssociationMissingParameters)
	}

	// if the on delete action is empty, set none as default.
	if len(data.OnDelete) == 0 {
		data.OnDelete = metadata.NoAction
	}

	// check if this association has already exist,
	// if yes, it's not allowed to create this association

	//  check the association
	filter := []map[string]interface{}{mapstr.MapStr{
		common.AssociatedObjectIDField: data.AsstObjID,
		common.BKObjIDField:            data.ObjectID,
		common.AssociationKindIDField:  data.AsstKindID,
	}}
	rsp, ccErr := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjAsst,
		filter)
	if ccErr != nil {
		blog.Errorf("failed to create the association (%#v) , err: %v, rid: %s", filter, ccErr, kit.Rid)
		return nil, ccErr
	}

	if rsp[0] != 0 {
		blog.Errorf("failed to create the association (%#v), the associations %s->%s already exist, rid: %s",
			filter, data.ObjectID, data.AsstObjID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoAssociationAlreadyExist, data.ObjectID, data.AsstObjID)
	}

	if err := assoc.isObjectInAssocValid(kit, data.ObjectID, data.AsstObjID); err != nil {
		blog.Errorf("objectID or asstObjectID is invalid, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// create a new
	cond := &metadata.CreateModelAssociation{Spec: *data}
	rspAsst, err := assoc.clientSet.CoreService().Association().CreateModelAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("create object association failed, param: %#v , err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	data.ID = int64(rspAsst.Created.ID)
	return data, nil
}

// DeleteAssociationWithPreCheck delete association after pre-check
func (assoc *association) DeleteAssociationWithPreCheck(kit *rest.Kit, associationID int64) error {

	cond := &metadata.QueryCondition{Condition: mapstr.MapStr{metadata.AssociationFieldAssociationId: associationID}}
	// if this association has already been instantiated, then this association should not be deleted.
	// get the association with id at first.
	result, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("get this association for pre check failed, err: %v, rid: %s", associationID, err, kit.Rid)
		return err
	}

	if len(result.Info) == 0 {
		blog.Errorf("can not find this id[%d] of association, rid: %s", associationID, kit.Rid)
		return nil
	}

	if len(result.Info) > 1 {
		blog.Errorf("search inst by associationID[%d] got multiple association, rid: %s", associationID, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if result.Info[0].AsstKindID == common.AssociationKindMainline {
		return kit.CCError.CCError(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}

	// a pre-defined association can not be updated.
	if result.Info[0].IsPre != nil && *result.Info[0].IsPre {
		blog.Errorf("object association id[%d] is a pre-defined association, rid: %s", associationID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// find instance(s) belongs to this association
	params := &metadata.Condition{
		Condition: mapstr.MapStr{common.AssociationObjAsstIDField: result.Info[0].AssociationName},
	}
	rsp, err := assoc.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header,
		result.Info[0].ObjectID, params)
	if err != nil {
		blog.Errorf("count inst association failed, objID: %s, params: %#v err: %v, rid: %s", result.Info[0].ObjectID,
			params, err, kit.Rid)
		return err
	}

	if rsp.Count != 0 {
		// object association has already been instantiated, association can not be deleted.
		blog.Errorf("search association by associationID[%d], got instances, rid: %s", associationID, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// TODO: check association on_delete action before really delete this association.
	// all the pre check has finished, delete the association now.
	deleteCond := &metadata.DeleteOption{Condition: cond.Condition}
	_, err = assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header, deleteCond)
	if err != nil {
		blog.Errorf("delete object association failed, err: %v, rid: %s", deleteCond, err, kit.Rid)
		return err
	}

	return nil
}

// UpdateObjectAssociation update object association by assoID
func (assoc *association) UpdateObjectAssociation(kit *rest.Kit, data mapstr.MapStr, assoID int64) error {

	if field, can := canUpdate(data); !can {
		blog.Errorf("request to update a forbidden update, id:[%d], field[%s], rid: %s", assoID, field, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoObjectAssociationUpdateForbiddenFields)
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{metadata.AssociationFieldAssociationId: assoID}})
	if err != nil {
		blog.Errorf("search object association by id(%d) failed, err: %v, rid: %s", assoID, err, kit.Rid)
		return err
	}

	if len(rsp.Info) < 1 {
		blog.Errorf("update the object association failed, id %d not found, rid: %s", assoID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoObjectAssociationNotExist)
	}

	// a pre-defined association can not be updated.
	if rsp.Info[0].IsPre != nil && *rsp.Info[0].IsPre {
		blog.Errorf("object association[%d] is a pre-defined association, rid: %s", assoID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoUpdatePredefinedAssociation)
	}

	// check object exists
	if err = assoc.isObjectInAssocValid(kit, rsp.Info[0].ObjectID, rsp.Info[0].AsstObjID); err != nil {
		blog.Errorf("objectID or asstObjectID is invalid, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	updateCond := &metadata.UpdateOption{
		Condition: mapstr.MapStr{common.BKFieldID: assoID},
		Data:      data,
	}
	_, err = assoc.clientSet.CoreService().Association().UpdateModelAssociation(kit.Ctx, kit.Header, updateCond)
	if err != nil {
		blog.Errorf("update the association (%#v) failed, err: %v, rid: %s", updateCond, err, kit.Rid)
		return err
	}

	return nil
}

// SearchObjectAssocWithAssocKindList search object associtaion by asstkind ids
func (assoc *association) SearchObjectAssocWithAssocKindList(kit *rest.Kit, asstKindIDs []string) (
	*metadata.AssociationList, error) {

	if len(asstKindIDs) == 0 {
		return &metadata.AssociationList{Associations: make([]metadata.AssociationDetail, 0)}, nil
	}

	queryCond := &metadata.QueryCondition{Condition: mapstr.MapStr{
		common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKindIDs},
	}}
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("get object association list failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	asso := make([]metadata.AssociationDetail, 0)
	for _, association := range rsp.Info {
		asso = append(asso, metadata.AssociationDetail{
			AssociationKindID: association.AsstKindID,
			Associations:      []metadata.Association{association},
		})
	}

	return &metadata.AssociationList{Associations: asso}, nil
}

func (assoc *association) isObjectInAssocValid(kit *rest.Kit, objectID, asstObjectID string) error {
	// check source object exists
	queryCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: []string{objectID, asstObjectID}}},
	}
	objRsp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("read the object(%s) failed, err: %v, rid: %s", objectID, err, kit.Rid)
		return err
	}

	if len(objRsp.Info) == 0 {
		blog.Errorf("object(%s) and asstObject(%s) is invalid, rid: %s", objectID, asstObjectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_obj_id&bk_asst_obj_id")
	}

	checkMap := make(map[string]struct{})
	for _, item := range objRsp.Info {
		checkMap[item.ObjectID] = struct{}{}
	}

	if _, exist := checkMap[objectID]; !exist {
		blog.Errorf("object(%s) is invalid, rid: %s", objectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	if _, exist := checkMap[asstObjectID]; !exist {
		blog.Errorf("object(%s) is invalid, rid: %s", asstObjectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAsstObjIDField)
	}

	return nil
}

func canUpdate(data mapstr.MapStr) (field string, can bool) {
	id, exist := data.Get(common.BKFieldID)
	if exist {
		if idInt, err := util.GetInt64ByInterface(id); err != nil || idInt != 0 {
			return common.BKFieldID, false
		}
	}

	_, exist = data.Get(common.BkSupplierAccount)
	if exist {
		return common.BkSupplierAccount, false
	}

	_, exist = data.Get(common.AssociationObjAsstIDField)
	if exist {
		return common.AssociationObjAsstIDField, false
	}

	_, exist = data.Get(common.BKObjIDField)
	if exist {
		return common.BKObjIDField, false
	}

	_, exist = data.Get(common.BKAsstObjIDField)
	if exist {
		return common.BKAsstObjIDField, false
	}

	_, exist = data.Get("mapping")
	if exist {
		return "mapping", false
	}

	_, exist = data.Get(common.BKIsPre)
	if exist {
		return common.BKIsPre, false
	}

	// only on delete, association kind id, alias name can be update.
	return "", true
}
