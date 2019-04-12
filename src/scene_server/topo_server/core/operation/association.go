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

package operation

import (
	"context"
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Object, error)
	DeleteMainlineAssociaton(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error)

	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (*metadata.Association, error)
	DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error
	UpdateAssociation(params types.ContextParams, data mapstr.MapStr, assoID int64) error
	SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error)

	DeleteAssociation(params types.ContextParams, cond condition.Condition) error
	SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error)
	CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error
	CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error
	DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error

	// 关联关系改造后的接口
	SearchObjectAssoWithAssoKindList(params types.ContextParams, asstKindIDs []string) (resp *metadata.AssociationList, err error)
	SearchType(params types.ContextParams, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error)
	CreateType(cparams types.ContextParams, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error)
	UpdateType(params types.ContextParams, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error)
	DeleteType(params types.ContextParams, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error)

	SearchObject(params types.ContextParams, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error)
	CreateObject(params types.ContextParams, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error)
	UpdateObject(params types.ContextParams, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error)
	DeleteObject(params types.ContextParams, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error)

	SearchInst(params types.ContextParams, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	CreateInst(params types.ContextParams, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error)
	DeleteInst(params types.ContextParams, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error)

	ImportInstAssociation(ctx context.Context, params types.ContextParams, objID string, importData map[int]metadata.ExcelAssocation) (resp metadata.ResponeImportAssociationData, err error)

	SetProxy(cls ClassificationOperationInterface, obj ObjectOperationInterface, grp GroupOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface) AssociationOperationInterface {
	return &association{
		clientSet: client,
	}
}

type association struct {
	clientSet    apimachinery.ClientSetInterface
	cls          ClassificationOperationInterface
	obj          ObjectOperationInterface
	grp          GroupOperationInterface
	attr         AttributeOperationInterface
	inst         InstOperationInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

func (a *association) SetProxy(cls ClassificationOperationInterface, obj ObjectOperationInterface, grp GroupOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory) {
	a.cls = cls
	a.obj = obj
	a.attr = attr
	a.inst = inst
	a.grp = grp
	a.modelFactory = targetModel
	a.instFactory = targetInst
}

func (a *association) SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error) {

	cond := condition.CreateCondition()
	if 0 != len(objID) {
		cond.Field(common.BKObjIDField).Eq(objID)
	}

	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , err: %s", objID, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

func (a *association) SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error) {
	intput, err := mapstr.NewFromInterface(query.Condition)
	rsp, err := a.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: intput})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the association info, query: %#v, err: %s", query, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

// CreateCommonAssociation create a common association, in topo model scene, which doesn't include bk_mainline association type
func (a *association) CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (*metadata.Association, error) {
	if data.AsstKindID == common.AssociationKindMainline {
		return nil, params.Err.Error(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}
	if len(data.AsstKindID) == 0 || len(data.AsstObjID) == 0 || len(data.ObjectID) == 0 {
		blog.Errorf("[operation-asst] failed to create the association , association kind id associate/object id is required")
		return nil, params.Err.Error(common.CCErrorTopoAssociationMissingParameters)
	}

	// if the on delete action is empty, set none as default.
	if len(data.OnDelete) == 0 {
		data.OnDelete = metadata.NoAction
	}

	//  check the association
	cond := condition.CreateCondition()
	cond.Field(common.AssociatedObjectIDField).Eq(data.AsstObjID)
	cond.Field(common.BKObjIDField).Eq(data.ObjectID)
	cond.Field(common.AssociationKindIDField).Eq(data.AsstKindID)

	rsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data.Info) > 0 {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , the associations %s->%s already exist ",
			cond.ToMapStr(), data.ObjectID, data.AsstObjID)
		return nil, params.Err.Errorf(common.CCErrTopoAssociationAlreadyExist, data.ObjectID, data.AsstObjID)
	}

	// check source object exists
	if err := a.obj.IsValidObject(params, data.ObjectID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, err: %s", data.ObjectID, err.Error())
		return nil, err
	}

	if err := a.obj.IsValidObject(params, data.AsstObjID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, err: %s", data.AsstObjID, err.Error())
		return nil, err
	}

	// create a new
	rspAsst, err := a.clientSet.CoreService().Association().CreateModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: *data})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspAsst.ErrMsg)
		return nil, params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	return data, nil
}

func (a *association) DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error {

	rsp, err := a.clientSet.CoreService().Association().DeleteInstAssociation(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to delete the inst association info , err: %s", rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (a *association) CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error {
	// create a new
	rspAsst, err := a.clientSet.CoreService().Association().CreateInstAssociation(context.Background(), params.Header, &metadata.CreateOneInstanceAssociation{Data: *data})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspAsst.ErrMsg)
		return params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	return nil
}

func (a *association) DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error {
	// if this association has already been instantiated, then this association should not be deleted.
	// get the association with id at first.
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(associationID)
	result, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %v", associationID, err)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !result.Result {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %s", associationID, result.ErrMsg)
		return params.Err.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("[operation-asst] delete association with id[%d], but can not find this association, return now.", associationID)
		return nil
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("[operation-asst] delete association with id[%d], but got multiple association", associationID)
		return params.Err.Error(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if result.Data.Info[0].AsstKindID == common.AssociationKindMainline {
		return params.Err.Error(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}

	// find instance(s) belongs to this association
	cond = condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(result.Data.Info[0].AssociationName)
	query := metadata.QueryInput{Condition: cond.ToMapStr()}
	insts, err := a.SearchInstAssociation(params, &query)
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but association instance(s) failed, err: %v", associationID, err)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if len(insts) != 0 {
		// object association has already been instantiated, association can not be deleted.
		blog.Errorf("[operation-asst] delete association with id[%d], but has multiple instances, can not be deleted.", associationID)
		return params.Err.Error(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// TODO: check association on_delete action before really delete this association.
	// all the pre check has finished, delete the association now.
	cond = condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(associationID)
	return a.DeleteAssociation(params, cond)
}

func (a *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {
	rsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %v", cond.ToMapStr(), err)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		// we assume this association has already been deleted.
		blog.Warnf("delete object association, but can not get association with cond[%v] ", cond.ToMapStr())
		return nil
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data.Info[0].IsPre && *rsp.Data.Info[0].IsPre {
		blog.Errorf("delete object association with cond[%v], but it's a pre-defined association, can not be deleted.", cond.ToMapStr())
		return params.Err.Error(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// delete the object association
	result, err := a.clientSet.CoreService().Association().DeleteModelAssociation(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", cond.ToMapStr(), result.ErrMsg)
		return params.Err.Error(result.Code)
	}

	return nil
}

func (a *association) UpdateAssociation(params types.ContextParams, data mapstr.MapStr, assoID int64) error {
	asst := &metadata.Association{}
	err := data.MarshalJSONInto(asst)
	if err != nil {
		blog.Errorf("[operation-asst] update association with  %s", err.Error())
		return params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if field, can := asst.CanUpdate(); !can {
		blog.Warnf("update association[%d], but request to update a forbidden update field[%s].", assoID, field)
		return params.Err.Error(common.CCErrorTopoObjectAssociationUpdateForbiddenFields)
	}

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(assoID)
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(params.SupplierAccount)

	rsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to update the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		blog.Errorf("[operation-asst] failed to update the object association , id %d not found", assoID)
		return params.Err.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data.Info[0].IsPre && *rsp.Data.Info[0].IsPre {
		blog.Errorf("update object association[%d], but it's a pre-defined association, can not be updated.", assoID)
		return params.Err.Error(common.CCErrorTopoUpdatePredefinedAssociation)
	}

	// check object exists
	if err := a.obj.IsValidObject(params, rsp.Data.Info[0].ObjectID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s", rsp.Data.Info[0].ObjectID, err.Error())
		return err
	}

	if err := a.obj.IsValidObject(params, rsp.Data.Info[0].AsstObjID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s", rsp.Data.Info[0].AsstObjID, err.Error())
		return err
	}

	updateopt := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(assoID).ToMapStr(),
		Data:      data,
	}
	rspAsst, err := a.clientSet.CoreService().Association().UpdateModelAssociation(context.Background(), params.Header, &updateopt)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspAsst.ErrMsg)
		return params.Err.Error(rspAsst.Code)
	}

	return nil
}

// CheckBeAssociation and return error if the obj has been bind
func (a *association) CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	exists, err := a.SearchInstAssociation(params, &metadata.QueryInput{Condition: cond.ToMapStr()})
	if nil != err {
		return err
	}

	if len(exists) > 0 {
		beAsstObject := []string{}
		for _, asst := range exists {
			instRsp, err := a.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, asst.ObjectID,
				&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKInstIDField: asst.InstID}})
			if err != nil {
				return params.Err.Error(common.CCErrObjectSelectInstFailed)
			}
			if !instRsp.Result {
				return params.Err.New(instRsp.Code, instRsp.ErrMsg)
			}
			if len(instRsp.Data.Info) <= 0 {
				// 作为补充而存在，删除实例主机已经不存在的脏实例关联
				if delErr := a.DeleteInstAssociation(params, condition.CreateCondition().
					Field(common.BKObjIDField).Eq(asst.ObjectID).Field(common.BKAsstInstIDField).Eq(asst.InstID)); delErr != nil {
					return delErr
				}
				continue
			}
			beAsstObject = append(beAsstObject, asst.ObjectID)
		}
		if len(beAsstObject) > 0 {
			return params.Err.Errorf(common.CCErrTopoInstHasBeenAssociation, beAsstObject)
		}
	}
	return nil
}

// 关联关系改造后的接口
func (a *association) SearchObjectAssoWithAssoKindList(params types.ContextParams, asstKindIDs []string) (resp *metadata.AssociationList, err error) {
	if len(asstKindIDs) == 0 {
		return &metadata.AssociationList{Associations: make([]metadata.AssociationDetail, 0)}, nil
	}

	asso := make([]metadata.AssociationDetail, 0)
	for _, id := range asstKindIDs {
		cond := condition.CreateCondition()
		cond.Field(common.AssociationKindIDField).Eq(id)

		r, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v", id, err)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !r.Result {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v", id, r.ErrMsg)
			return nil, params.Err.Errorf(r.Code, r.ErrMsg)
		}

		asso = append(asso, metadata.AssociationDetail{AssociationKindID: id, Associations: r.Data.Info})
	}

	return &metadata.AssociationList{Associations: asso}, nil
}

func (a *association) SearchType(params types.ContextParams, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
	input := metadata.QueryCondition{
		Condition: request.Condition,
		Limit:     metadata.SearchLimit{Limit: int64(request.Limit), Offset: int64(request.Start)},
	}

	for _, key := range strings.Split(request.Sort, ",") {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		var isDesc bool
		switch key[0] {
		case '-':
			key = strings.TrimLeft(key, "-")
			isDesc = true
		case '+':
			key = strings.TrimLeft(key, "+")
		}
		input.SortArr = append(input.SortArr, metadata.SearchSort{IsDsc: isDesc, Field: key})
	}

	return a.clientSet.CoreService().Association().ReadAssociation(context.Background(), params.Header, &input)

}

func (a *association) CreateType(params types.ContextParams, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {
	rsp, err := a.clientSet.CoreService().Association().CreateAssociation(context.Background(), params.Header, &metadata.CreateAssociationKind{Data: *request})
	resp = &metadata.CreateAssociationTypeResult{BaseResp: rsp.BaseResp}
	resp.Data.ID = int64(rsp.Data.Created.ID)
	return resp, err

}

func (a *association) UpdateType(params types.ContextParams, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstTypeID).ToMapStr(),
		Data:      mapstr.NewFromStruct(request, "json"),
	}

	rsp, err := a.clientSet.CoreService().Association().UpdateAssociation(context.Background(), params.Header, &input)
	resp = &metadata.UpdateAssociationTypeResult{BaseResp: rsp.BaseResp}
	return resp, err
}

func (a *association) DeleteType(params types.ContextParams, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error) {
	cond := condition.CreateCondition()
	cond.Field("id").Eq(asstTypeID)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	query := &metadata.SearchAssociationTypeRequest{
		Condition: cond.ToMapStr(),
	}

	result, err := a.SearchType(params, query)
	if err != nil {
		blog.Errorf("delete association kind[%d], but get detailed info failed, err: %v", asstTypeID, err)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete association kind[%d], but get detailed info failed, err: %s", asstTypeID, result.ErrMsg)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("delete association kind[%d], but get multiple instance", asstTypeID)
		return nil, params.Err.Error(common.CCErrorTopoGetMultipleAssoKindInstWithOneID)
	}

	if len(result.Data.Info) == 0 {
		return &metadata.DeleteAssociationTypeResult{BaseResp: metadata.SuccessBaseResp, Data: common.CCSuccessStr}, nil
	}

	if result.Data.Info[0].IsPre != nil && *result.Data.Info[0].IsPre {
		blog.Errorf("delete association kind[%d], but this is a pre-defined association kind, can not be deleted.", asstTypeID)
		return nil, params.Err.Error(common.CCErrorTopoDeletePredefinedAssociationKind)
	}

	// a already used association kind can not be deleted.
	cond = condition.CreateCondition()
	cond.Field(common.AssociationKindIDField).Eq(result.Data.Info[0].AssociationKindID)
	asso, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %v", asstTypeID, err)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %s", asstTypeID, result.ErrMsg)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(asso.Data.Info) != 0 {
		blog.Warnf("delete association kind[%d], but it has already been used, can not be deleted.", asstTypeID)
		return nil, params.Err.Error(common.CCErrorTopoAssociationKindHasBeenUsed)
	}

	rsp, err := a.clientSet.CoreService().Association().DeleteAssociation(
		context.Background(), params.Header, &metadata.DeleteOption{
			Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstTypeID).ToMapStr(),
		},
	)

	return &metadata.DeleteAssociationTypeResult{BaseResp: rsp.BaseResp}, err
}

func (a *association) SearchObject(params types.ContextParams, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	rsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: request.Condition})

	resp = &metadata.SearchAssociationObjectResult{BaseResp: rsp.BaseResp, Data: []*metadata.Association{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, err
}

func (a *association) CreateObject(params types.ContextParams, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	rsp, err := a.clientSet.CoreService().Association().CreateModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: *request})

	resp = &metadata.CreateAssociationObjectResult{
		BaseResp: rsp.BaseResp,
	}
	resp.Data.ID = int64(rsp.Data.Created.ID)
	return resp, err
}

func (a *association) UpdateObject(params types.ContextParams, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstID).ToMapStr(),
		Data:      mapstr.NewFromStruct(request, "json"),
	}

	rsp, err := a.clientSet.CoreService().Association().UpdateModelAssociation(context.Background(), params.Header, &input)
	resp = &metadata.UpdateAssociationObjectResult{
		BaseResp: rsp.BaseResp,
	}
	return resp, err
}

func (a *association) DeleteObject(params types.ContextParams, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {

	input := metadata.DeleteOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstID).ToMapStr(),
	}
	rsp, err := a.clientSet.CoreService().Association().DeleteModelAssociation(context.Background(), params.Header, &input)
	return &metadata.DeleteAssociationObjectResult{BaseResp: rsp.BaseResp}, err

}

func (a *association) SearchInst(params types.ContextParams, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	rsp, err := a.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: request.Condition})

	resp = &metadata.SearchAssociationInstResult{BaseResp: rsp.BaseResp, Data: []*metadata.InstAsst{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, err
}

func (a *association) CreateInst(params types.ContextParams, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	cond := condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
	result, err := a.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %v", cond, err)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %s", cond, resp.ErrMsg)
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}

	if len(result.Data) == 0 {
		blog.Errorf("create instance association, but can not find object association[%s]. ", request.ObjectAsstID)
		return nil, params.Err.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	objectAsst := result.Data[0]

	objID := objectAsst.ObjectID
	asstObjID := objectAsst.AsstObjID

	switch result.Data[0].Mapping {
	case metadata.OneToOneMapping:
		// search instances belongs to this association.
		cond := condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
		cond.Field(common.BKInstIDField).Eq(request.InstID)
		inst, err := a.SearchInst(params, &metadata.SearchAssociationInstRequest{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v", cond, err)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !inst.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s", cond, resp.ErrMsg)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}
		if len(inst.Data) >= 1 {
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}

		cond = condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
		cond.Field(common.BKAsstInstIDField).Eq(request.AsstInstID)

		inst, err = a.SearchInst(params, &metadata.SearchAssociationInstRequest{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v", cond, err)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !inst.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s", cond, resp.ErrMsg)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}
		if len(inst.Data) >= 1 {
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}

	default:
		// after all the check, new association instance can be created.
	}

	input := metadata.CreateOneInstanceAssociation{
		Data: metadata.InstAsst{
			ObjectAsstID:      request.ObjectAsstID,
			InstID:            request.InstID,
			AsstInstID:        request.AsstInstID,
			ObjectID:          objID,
			AsstObjectID:      asstObjID,
			AssociationKindID: objectAsst.AsstKindID,
		},
	}
	rsp, err := a.clientSet.CoreService().Association().CreateInstAssociation(context.Background(), params.Header, &input)

	resp = &metadata.CreateAssociationInstResult{BaseResp: rsp.BaseResp}
	resp.Data.ID = int64(rsp.Data.Created.ID)
	return resp, err
}

func (a *association) DeleteInst(params types.ContextParams, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error) {
	input := metadata.DeleteOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(assoID).ToMapStr(),
	}
	rsp, err := a.clientSet.CoreService().Association().DeleteInstAssociation(context.Background(), params.Header, &input)
	resp = &metadata.DeleteAssociationInstResult{
		BaseResp: rsp.BaseResp,
	}

	return resp, err
}
