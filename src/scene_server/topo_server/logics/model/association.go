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
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	SearchAssociationType(kit *rest.Kit,
		request *metadata.SearchAssociationTypeRequest) (*metadata.SearchAssociationTypeResult, error)
	CreateAssociationType(kit *rest.Kit,
		request *metadata.AssociationKind) (*metadata.CreateAssociationTypeResult, error)
	UpdateAssociationType(kit *rest.Kit, asstTypeID int64,
		request *metadata.UpdateAssociationTypeRequest) (*metadata.UpdateAssociationTypeResult, error)
	DeleteAssociationType(kit *rest.Kit, asstTypeID int64) (*metadata.DeleteAssociationTypeResult, error)
	SearchObjectAssociationByObjIDs(kit *rest.Kit,
		objIDs []interface{}) (*metadata.SearchAssociationObjectResult, error)
	SearchObjectAssociation(kit *rest.Kit,
		request mapstr.MapStr) (*metadata.SearchAssociationObjectResult, error)
	CreateCommonAssociation(kit *rest.Kit, data *metadata.Association) (*metadata.Association, error)
	DeleteAssociationWithPreCheck(kit *rest.Kit, associationID int64) error
	UpdateObjectAssociation(kit *rest.Kit, data mapstr.MapStr, assoID int64) error
	SearchObjectAssocWithAssocKindList(kit *rest.Kit, asstKindIDs []string) (resp *metadata.AssociationList, err error)
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
}

func (assoc *association) SearchAssociationType(kit *rest.Kit,
	request *metadata.SearchAssociationTypeRequest) (*metadata.SearchAssociationTypeResult, error) {

	input := metadata.QueryCondition{
		Condition: request.Condition,
		Page:      metadata.BasePage{Limit: request.Limit, Start: request.Start, Sort: request.Sort},
	}

	return assoc.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, &input)
}

func (assoc *association) CreateAssociationType(kit *rest.Kit,
	request *metadata.AssociationKind) (*metadata.CreateAssociationTypeResult, error) {

	rsp, err := assoc.clientSet.CoreService().Association().CreateAssociationType(kit.Ctx, kit.Header,
		&metadata.CreateAssociationKind{Data: *request})
	if err != nil {
		blog.Errorf("create association type failed, kind id: %s, err: %v, rid: %s",
			request.AssociationKindID, err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}
	if rsp.Result == false || rsp.Code != 0 {
		blog.Errorf("create association type failed, request: %s, response: %s, rid: %s", request, rsp, kit.Rid)
		return nil, errors.NewCCError(rsp.Code, rsp.ErrMsg)
	}
	resp := &metadata.CreateAssociationTypeResult{BaseResp: rsp.BaseResp}
	resp.Data.ID = int64(rsp.Data.Created.ID)

	return resp, nil

}

func (assoc *association) UpdateAssociationType(kit *rest.Kit, asstTypeID int64,
	request *metadata.UpdateAssociationTypeRequest) (*metadata.UpdateAssociationTypeResult, error) {

	input := metadata.UpdateOption{
		Condition: mapstr.MapStr{common.BKFieldID: asstTypeID},
		Data:      mapstr.NewFromStruct(request, "json"),
	}

	rsp, err := assoc.clientSet.CoreService().Association().UpdateAssociationType(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("update association type failed, kind id: %d, err: %v, rid: %s", asstTypeID, err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}
	resp := &metadata.UpdateAssociationTypeResult{BaseResp: rsp.BaseResp}
	return resp, nil
}

func (assoc *association) DeleteAssociationType(kit *rest.Kit,
	asstTypeID int64) (*metadata.DeleteAssociationTypeResult, error) {

	result, err := assoc.SearchAssociationType(kit, &metadata.SearchAssociationTypeRequest{
		Condition: mapstr.MapStr{common.BKFieldID: asstTypeID},
	})
	if err != nil {
		blog.Errorf("search association kind by typeID[%d], but get detailed info failed, err: %v, rid: %s",
			asstTypeID, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("search association kind by typeID[%d], but get detailed info failed, err: %s, rid: %s",
			asstTypeID, result.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("search association kind by typeID[%d], but get multiple instance, rid: %s", asstTypeID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoGetMultipleAssocKindInstWithOneID)
	}

	if len(result.Data.Info) == 0 {
		return &metadata.DeleteAssociationTypeResult{BaseResp: metadata.SuccessBaseResp, Data: common.CCSuccessStr}, nil
	}

	if result.Data.Info[0].IsPre != nil && *result.Data.Info[0].IsPre {
		blog.Errorf("typeID[%d] of association kind is a pre-defined association kind, rid: %s", asstTypeID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoDeletePredefinedAssociationKind)
	}

	// a already used association kind can not be deleted.
	asso, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{common.AssociationKindIDField: result.Data.Info[0].AssociationKindID},
		})
	if err != nil {
		blog.Errorf("get objects that used this association kind[%d] failed, err: %v, rid: %s",
			asstTypeID, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("get objects that used this association kind[%d] failed, err: %v, rid: %s",
			asstTypeID, result.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(asso.Data.Info) != 0 {
		blog.Warnf("association kind[%d] has already been used, can not be deleted, rid: %s", asstTypeID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoAssociationKindHasBeenUsed)
	}

	rsp, err := assoc.clientSet.CoreService().Association().DeleteAssociationType(kit.Ctx, kit.Header,
		&metadata.DeleteOption{
			Condition: mapstr.MapStr{common.BKFieldID: asstTypeID},
		},
	)
	if err != nil {
		blog.Errorf("delete association type failed, kind id: %d, err: %v, rid: %s", asstTypeID, err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}

	return &metadata.DeleteAssociationTypeResult{BaseResp: rsp.BaseResp}, nil
}

func (assoc *association) SearchObjectAssociationByObjIDs(kit *rest.Kit,
	objIDs []interface{}) (*metadata.SearchAssociationObjectResult, error) {

	cond := mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{
			common.BKDBIN: objIDs,
		},
	}

	return assoc.SearchObjectAssociation(kit, cond)
}

func (assoc *association) SearchObjectAssociation(kit *rest.Kit,
	request mapstr.MapStr) (*metadata.SearchAssociationObjectResult, error) {

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: request})
	if err != nil {
		blog.Errorf("search object association failed, cond: %v, err: %s, rid: %s", request, err.Error(), kit.Rid)
		return nil, err
	}

	resp := &metadata.SearchAssociationObjectResult{BaseResp: rsp.BaseResp, Data: []*metadata.Association{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, nil
}

func (assoc *association) CreateCommonAssociation(kit *rest.Kit,
	data *metadata.Association) (*metadata.Association, error) {

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
	cond := mapstr.MapStr{
		common.AssociatedObjectIDField: data.AsstObjID,
		common.BKObjIDField:            data.ObjectID,
		common.AssociationKindIDField:  data.AsstKindID,
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("read object association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("failed to create the association (%#v) , err: %s, rid: %s", cond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data.Info) > 0 {
		blog.Errorf("failed to create the association (%#v), the associations %s->%s already exist, rid: %s",
			kit.Rid, cond, data.ObjectID, data.AsstObjID)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoAssociationAlreadyExist, data.ObjectID, data.AsstObjID)
	}

	if err := assoc.isObjectInAssocValid(kit, data.ObjectID, data.AsstObjID); err != nil {
		blog.Errorf("objectID or asstObjectID is invalid, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	// create a new
	rspAsst, err := assoc.clientSet.CoreService().Association().CreateModelAssociation(kit.Ctx, kit.Header,
		&metadata.CreateModelAssociation{Spec: *data})
	if err != nil {
		blog.Errorf("create object association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("create object association failed, param: %#v , err: %s, rid: %s", data, rspAsst.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	data.ID = int64(rspAsst.Data.Created.ID)
	return data, nil
}

func (assoc *association) DeleteAssociation(kit *rest.Kit, cond mapstr.MapStr) error {

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("get association with cond[%v] failed, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("get association with cond[%v] failed, err: %s, rid: %s", cond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		// we assume this association has already been deleted.
		blog.Warnf("get association with cond[%v] failed, return is empty, rid: %s", cond, kit.Rid)
		return nil
	}

	// a pre-defined association can not be updated.
	if rsp.Data.Info[0].IsPre != nil && *rsp.Data.Info[0].IsPre {
		blog.Errorf("object association with cond[%v] is a pre-defined association, rid: %s", cond, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// delete the object association
	result, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header,
		&metadata.DeleteOption{Condition: cond})
	if err != nil {
		blog.Errorf("delete object association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete object association failed, err: %s, rid: %s", cond, result.ErrMsg, kit.Rid)
		return kit.CCError.Error(result.Code)
	}

	return nil
}

func (assoc *association) DeleteAssociationWithPreCheck(kit *rest.Kit, associationID int64) error {
	// if this association has already been instantiated, then this association should not be deleted.
	// get the association with id at first.
	result, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{metadata.AssociationFieldAssociationId: associationID}})
	if err != nil {
		blog.Errorf("get this association for pre check failed, err: %v, rid: %s", associationID, err, kit.Rid)
		return kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !result.Result {
		blog.Errorf("get this association for pre check failed, err: %s, rid: %s",
			associationID, result.ErrMsg, kit.Rid)
		return kit.CCError.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("can not find this id[%d] of association, rid: %s", associationID, kit.Rid)
		return nil
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("search inst by associationID[%d] got multiple association, rid: %s", associationID, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if result.Data.Info[0].AsstKindID == common.AssociationKindMainline {
		return kit.CCError.CCError(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}

	// find instance(s) belongs to this association
	// TODO after merge change this interface call to SearchInstanceAssociation of inst/association.go
	queryCond := &metadata.InstAsstQueryCondition{
		Cond: metadata.QueryCondition{Condition: mapstr.MapStr{common.AssociationObjAsstIDField: result.Data.
			Info[0].AssociationName}},
		ObjID: result.Data.Info[0].ObjectID,
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search inst association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("search association info failed, query: %#v, err: %s, rid: %s", queryCond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) != 0 {
		// object association has already been instantiated, association can not be deleted.
		blog.Errorf("search association by associationID[%d], got instances, rid: %s", associationID, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// TODO: check association on_delete action before really delete this association.
	// all the pre check has finished, delete the association now.
	return assoc.DeleteAssociation(kit, mapstr.MapStr{metadata.AssociationFieldAssociationId: associationID})
}

func (assoc *association) UpdateObjectAssociation(kit *rest.Kit, data mapstr.MapStr, assoID int64) error {

	asst := new(metadata.Association)
	err := data.MarshalJSONInto(asst)
	if err != nil {
		blog.Errorf("marshal data into asst failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if field, can := asst.CanUpdate(); !can {
		blog.Warnf("request to update a forbidden update field[%s], rid: %s", assoID, field, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoObjectAssociationUpdateForbiddenFields)
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{metadata.AssociationFieldAssociationId: assoID}})
	if err != nil {
		blog.Errorf("request to search object association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("update the association by association ID (%s) failed, err: %s, rid: %s",
			assoID, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		blog.Errorf("update the object association failed, id %d not found, rid: %s", assoID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoObjectAssociationNotExist)
	}

	// a pre-defined association can not be updated.
	if rsp.Data.Info[0].IsPre != nil && *rsp.Data.Info[0].IsPre {
		blog.Errorf("object association[%d] is a pre-defined association, rid: %s", assoID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoUpdatePredefinedAssociation)
	}

	// check object exists
	if err := assoc.isObjectInAssocValid(kit, rsp.Data.Info[0].ObjectID, rsp.Data.Info[0].AsstObjID); err != nil {
		blog.Errorf("objectID or asstObjectID is invalid, err: %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	rspAsst, err := assoc.clientSet.CoreService().Association().UpdateModelAssociation(kit.Ctx, kit.Header,
		&metadata.UpdateOption{
			Condition: mapstr.MapStr{common.BKFieldID: assoID},
			Data:      data,
		})
	if err != nil {
		blog.Errorf("request to update object association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspAsst.Result {
		blog.Errorf("update the association (%#v) failed, err: %s, rid: %s", data, rspAsst.ErrMsg, kit.Rid)
		return kit.CCError.CCError(rspAsst.Code)
	}

	return nil
}

func (assoc *association) SearchObjectAssocWithAssocKindList(kit *rest.Kit,
	asstKindIDs []string) (resp *metadata.AssociationList, err error) {

	if len(asstKindIDs) == 0 {
		return &metadata.AssociationList{Associations: make([]metadata.AssociationDetail, 0)}, nil
	}

	cond := mapstr.MapStr{
		common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKindIDs},
	}
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})

	if err != nil {
		blog.Errorf("get object association list failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("get object association list failed, err: %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCErrorf(rsp.Code, rsp.ErrMsg)
	}

	asso := make([]metadata.AssociationDetail, 0)
	for _, association := range rsp.Data.Info {
		asso = append(asso, metadata.AssociationDetail{
			AssociationKindID: association.AsstKindID,
			Associations:      []metadata.Association{association},
		})
	}

	return &metadata.AssociationList{Associations: asso}, nil
}

func (assoc *association) isObjectInAssocValid(kit *rest.Kit, objectID, asstObjectID string) error {
	// check source object exists
	objRsp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: objectID},
	})
	if err != nil {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", objectID, err.Error(), kit.Rid)
		return err
	}

	if !objRsp.Result {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", objectID, objRsp.ErrMsg, kit.Rid)
		return err
	}

	if len(objRsp.Data.Info) == 0 {
		blog.Errorf("the object(%s) is invalid, return is empty, rid: %s", objectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	// check target object exists
	asstObjRsp, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: asstObjectID},
	})
	if err != nil {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", asstObjectID, err.Error(), kit.Rid)
		return err
	}

	if !asstObjRsp.Result {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", asstObjectID, asstObjRsp.ErrMsg, kit.Rid)
		return err
	}

	if len(asstObjRsp.Data.Info) == 0 {
		blog.Errorf("the object(%s) is invalid, return is empty, rid: %s", asstObjectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAsstObjIDField)
	}

	return nil
}
