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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	metatype "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Object, error)
	DeleteMainlineAssociation(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error)

	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (*metadata.Association, error)
	DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error
	UpdateAssociation(params types.ContextParams, data frtypes.MapStr, assoID int64) error
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
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	if 0 != len(objID) {
		cond.Field(common.BKObjIDField).Eq(objID)
	}
	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , err: %s", objID, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}

func (a *association) SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error) {

	rsp, err := a.clientSet.ObjectController().Instance().SearchObjects(context.Background(), common.BKTableNameInstAsst, params.Header, query)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the association info, query: %#v, err: %s", query, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	var instAsst []metadata.InstAsst
	for _, info := range rsp.Data.Info {
		asst := metadata.InstAsst{}
		if err := info.MarshalJSONInto(&asst); nil != err {
			return nil, err
		}
		instAsst = append(instAsst, asst)
	}
	blog.V(4).Infof("[SearchInstAssociation] search association, condition: %#v, results: %#v, unmarshal to: %#v", query, rsp.Data.Info, instAsst)
	return instAsst, nil
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
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.AssociationKindIDField).Eq(data.AsstKindID)

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data) > 0 {
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
	rspAsst, err := a.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), params.Header, data)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspAsst.ErrMsg)
		return nil, params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	if len(rspAsst.Data) == 0 {
		return nil, params.Err.Error(common.CCErrCommNotFound)
	}

	return &rspAsst.Data[0], nil
}

func (a *association) DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error {

	rsp, err := a.clientSet.ObjectController().Instance().DelObject(context.Background(), common.BKTableNameInstAsst, params.Header, cond.ToMapStr())
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

	rspAsst, err := a.clientSet.ObjectController().Instance().CreateObject(context.Background(), common.BKTableNameInstAsst, params.Header, data)
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
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	result, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %v", associationID, err)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !result.Result {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %s", associationID, result.ErrMsg)
		return params.Err.New(result.Code, result.ErrMsg)
	}

	if len(result.Data) == 0 {
		blog.Errorf("[operation-asst] delete association with id[%d], but can not find this association, return now.", associationID)
		return nil
	}

	if len(result.Data) > 1 {
		blog.Errorf("[operation-asst] delete association with id[%d], but got multiple association", associationID)
		return params.Err.Error(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if result.Data[0].AsstKindID == common.AssociationKindMainline {
		return params.Err.Error(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}

	// find instance(s) belongs to this association
	cond = condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(result.Data[0].AssociationName)
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
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	return a.DeleteAssociation(params, cond)
}

func (a *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %v", cond.ToMapStr(), err)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	if len(rsp.Data) < 1 {
		// we assume this association has already been deleted.
		blog.Warnf("delete object association, but can not get association with cond[%v] ", cond.ToMapStr())
		return params.Err.Error(common.CCErrorTopoAssociationDoNotExist)
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data[0].IsPre && *rsp.Data[0].IsPre {
		blog.Errorf("delete object association with cond[%v], but it's a pre-defined association, can not be deleted.", cond.ToMapStr())
		return params.Err.Error(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// delete the object association
	result, err := a.clientSet.ObjectController().Meta().DeleteObjectAssociation(context.Background(), 0, params.Header, cond.ToMapStr())
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

func (a *association) UpdateAssociation(params types.ContextParams, data frtypes.MapStr, assoID int64) error {
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

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to update the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	if len(rsp.Data) < 1 {
		blog.Errorf("[operation-asst] failed to update the object association , id %d not found", assoID)
		return params.Err.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data[0].IsPre && *rsp.Data[0].IsPre {
		blog.Errorf("update object association[%d], but it's a pre-defined association, can not be updated.", assoID)
		return params.Err.Error(common.CCErrorTopoUpdatePredefinedAssociation)
	}

	// check object exists
	if err := a.obj.IsValidObject(params, rsp.Data[0].ObjectID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s", rsp.Data[0].ObjectID, err.Error())
		return err
	}

	if err := a.obj.IsValidObject(params, rsp.Data[0].AsstObjID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s", rsp.Data[0].AsstObjID, err.Error())
		return err
	}

	rspAsst, err := a.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), assoID, params.Header, data)
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
	exists, err := a.SearchInstAssociation(params, &metatype.QueryInput{Condition: cond.ToMapStr()})
	if nil != err {
		return err
	}

	if len(exists) > 0 {
		beAsstObject := []string{}
		for _, asst := range exists {
			instRsp, err := a.clientSet.ObjectController().Instance().SearchObjects(context.Background(), asst.ObjectID, params.Header,
				&metadata.QueryInput{Condition: frtypes.MapStr{common.BKInstIDField: asst.InstID}})
			if err != nil {
				return params.Err.Error(common.CCErrObjectSelectInstFailed)
			}
			if !instRsp.Result {
				return params.Err.New(instRsp.Code, instRsp.ErrMsg)
			}
			if len(instRsp.Data.Info) <= 0 {
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
		cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		cond.Field(common.AssociationKindIDField).Eq(id)

		r, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
		if err != nil {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v", id, err)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !r.Result {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v", id, r.ErrMsg)
			return nil, params.Err.Errorf(r.Code, r.ErrMsg)
		}

		asso = append(asso, metadata.AssociationDetail{AssociationKindID: id, Associations: r.Data})
	}

	return &metadata.AssociationList{Associations: asso}, nil
}

func (a *association) SearchType(params types.ContextParams, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Association().SearchType(context.TODO(), params.Header, request)
}
func (a *association) CreateType(params types.ContextParams, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Association().CreateType(context.TODO(), params.Header, request)
}
func (a *association) UpdateType(params types.ContextParams, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Association().UpdateType(context.TODO(), params.Header, asstTypeID, request)
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
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.AssociationKindIDField).Eq(result.Data.Info[0].AssociationKindID)
	filter := metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()}
	asso, err := a.clientSet.ObjectController().Association().SearchObject(context.TODO(), params.Header, &filter)
	if err != nil {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %v", asstTypeID, err)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %s", asstTypeID, result.ErrMsg)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(asso.Data) != 0 {
		blog.Warnf("delete association kind[%d], but it has already been used, can not be deleted.", asstTypeID)
		return nil, params.Err.Error(common.CCErrorTopoAssociationKindHasBeenUsed)
	}

	return a.clientSet.ObjectController().Association().DeleteType(context.TODO(), params.Header, asstTypeID)
}

func (a *association) SearchObject(params types.ContextParams, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Association().SearchObject(context.TODO(), params.Header, request)
}
func (a *association) CreateObject(params types.ContextParams, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Association().CreateObject(context.TODO(), params.Header, request)
}
func (a *association) UpdateObject(params types.ContextParams, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Association().UpdateObject(context.TODO(), params.Header, asstID, request)
}
func (a *association) DeleteObject(params types.ContextParams, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Association().DeleteObject(context.TODO(), params.Header, asstID)
}

func (a *association) SearchInst(params types.ContextParams, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	return a.clientSet.ObjectController().Association().SearchInst(context.TODO(), params.Header, request)
}

func (a *association) checkObjectIsPause(params types.ContextParams, cond condition.Condition) (err error) {
	model, err := a.obj.FindObject(params, cond)
	if err != nil {
		return err
	}
	if len(model) == 0 {
		return params.Err.Error(common.CCErrCommNotFound)
	}
	if model[0].GetIsPaused() {
		return params.Err.Error(common.CCErrorTopoModleStopped)
	}
	return nil
}

func (a *association) CreateInst(params types.ContextParams, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	cond := condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstId)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	result, err := a.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %v", cond, err)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %s", cond, resp.ErrMsg)
		return nil, params.Err.Error(resp.Code)
	}

	if len(result.Data) == 0 {
		blog.Errorf("create instance association, but can not find object association[%s]. ", request.ObjectAsstId)
		return nil, params.Err.Error(common.CCErrorTopoObjectAssociationNotExist)
	}
	modelAsst := result.Data[0]
	if err := a.checkObjectIsPause(params, condition.CreateCondition().Field(common.BKObjIDField).Eq(modelAsst.ObjectID)); err != nil {
		blog.Errorf("create instance association, but model check for %s Failed: %v", modelAsst.ObjectID, err)
		return nil, err
	}

	if err := a.checkObjectIsPause(params, condition.CreateCondition().Field(common.BKObjIDField).Eq(modelAsst.AsstObjID)); err != nil {
		blog.Errorf("create instance association, but model check for %s Failed: %v", modelAsst.AsstObjID, err)
		return nil, err
	}

	switch modelAsst.Mapping {

	case metatype.OneToOneMapping:
		cond := condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstId)
		cond.Field(common.BKInstIDField).Eq(request.InstId)
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
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstId)
		cond.Field(common.BKAsstInstIDField).Eq(request.AsstInstId)

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

	case metadata.OneToManyMapping:
		cond = condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstId)
		cond.Field(common.BKAsstInstIDField).Eq(request.AsstInstId)

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
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToManyAssociation)
		}

	default:
		// after all the check, new association instance can be created.
	}

	return a.clientSet.ObjectController().Association().CreateInst(context.TODO(), params.Header, request)
}

func (a *association) DeleteInst(params types.ContextParams, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error) {
	return a.clientSet.ObjectController().Association().DeleteInst(context.TODO(), params.Header, assoID)
}
