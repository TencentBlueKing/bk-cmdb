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
	"fmt"
	"net/http"

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
	DeleteMainlineAssociaton(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error)

	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) error
	DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error
	UpdateAssociation(params types.ContextParams, data frtypes.MapStr, assoID int64) error
	SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error)

	DeleteAssociation(params types.ContextParams, cond condition.Condition) error
	SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error)
	CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error
	CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error
	DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error

	// 关联关系改造后的接口
	SearchType(ctx context.Context, h http.Header, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error)
	CreateType(ctx context.Context, h http.Header, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error)
	UpdateType(ctx context.Context, h http.Header, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error)
	DeleteType(ctx context.Context, h http.Header, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error)

	SearchObject(ctx context.Context, h http.Header, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error)
	CreateObject(ctx context.Context, h http.Header, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error)
	UpdateObject(ctx context.Context, h http.Header, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error)
	DeleteObject(ctx context.Context, h http.Header, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error)

	SearchInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	CreateInst(ctx context.Context, h http.Header, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error)
	DeleteInst(ctx context.Context, h http.Header, request *metadata.DeleteAssociationInstRequest) (resp *metadata.DeleteAssociationInstResult, err error)

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

func (a *association) CreateCommonAssociation(params types.ContextParams, data *metadata.Association) error {

	if len(data.AsstKindID) == 0 || len(data.AsstObjID) == 0 || len(data.ObjectID) == 0 {
		errmsg := fmt.Sprintf("[operation-asst] failed to create the association , association kind id associate/object id is required")
		blog.Error(errmsg)
		return params.Err.Error(common.CCErrorTopoAssociationMissingPrameters)
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
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data) > 0 {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , the associations %s->%s already exist ",
			cond.ToMapStr(), data.ObjectID, data.AsstObjID)
		return params.Err.Errorf(common.CCErrTopoAssociationAlreadyExist, data.ObjectID, data.AsstObjID)
	}

	// check object exists
	condObj := condition.CreateCondition()
	condObj.Field(common.BKObjIDField).Eq(data.ObjectID)

	rspObj, err := a.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, condObj.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !rspObj.Result {
		blog.Errorf("[operation-asst] create the association [%s->%s], but check object[%s] failed, err: %s", data.ObjectID,
			data.AsstObjID, data.ObjectID, rspObj.ErrMsg)
		return params.Err.New(rspObj.Code, rspObj.ErrMsg)
	}

	if len(rspObj.Data) != 1 {
		blog.Error("[operation-asst] create the association [%s->%s], but object[%s] do not exist.", data.ObjectID,
			data.AsstObjID, data.ObjectID, rspObj.ErrMsg)
		return params.Err.Error(common.CCErrTopoAssociationSourceObjectNotExist)
	}

	asstCond := condition.CreateCondition().Field(common.BKObjIDField).Eq(data.AsstObjID).ToMapStr()
	rspObj, err = a.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, asstCond)
	if nil != err {
		blog.Errorf("[operation-asst] create association [%s->%s], but get object[%s] failed, err: %s", data.AsstObjID,
			data.AsstObjID, data.AsstObjID, rspObj.ErrMsg)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !rspObj.Result {
		blog.Errorf("[operation-asst] create association [%s->%s], but check object[%s] failed, err: %s", data.AsstObjID,
			data.AsstObjID, data.AsstObjID, rspObj.ErrMsg)
		return params.Err.New(rspObj.Code, rspObj.ErrMsg)
	}

	if len(rspObj.Data) != 1 {
		blog.Error("[operation-asst] create association [%s->%s], but object[%s] do not exist.", data.AsstObjID,
			data.AsstObjID, data.AsstObjID)
		return params.Err.Error(common.CCErrTopoAssociationDestinationObjectNotExist)
	}

	// create a new
	rspAsst, err := a.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), params.Header, data)
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

	// find instance(s) belongs to this association
	cond = condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(result.Data[0].ObjectID)
	cond.Field(common.AssociatedObjectIDField).Eq(result.Data[0].AsstObjID)

	query := metadata.QueryInput{Condition: cond.ToMapStr()}
	resp, err := a.clientSet.ObjectController().Instance().SearchObjects(context.Background(), "module", params.Header, &query)
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but association instance(s) failed, err: %v", associationID, err)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association instances failed, err: %s", associationID, resp.ErrMsg)
		return params.Err.New(result.Code, result.ErrMsg)
	}

	if len(resp.Data.Info) != 0 {
		// object association has already been instantiated, associaton can not be deleted.
		blog.Errorf("[operation-asst] delete association with id[%d], but has multiple instances, can not be deleted.", associationID)
		return params.Err.Error(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// all the pre check has finished, delete the association now.
	cond = condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(associationID)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	return a.DeleteAssociation(params, cond)
}

func (a *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {

	// delete the object association
	rsp, err := a.clientSet.ObjectController().Meta().DeleteObjectAssociation(context.Background(), 0, params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

// TODO: need to confirm which fields can be updated.
func (a *association) UpdateAssociation(params types.ContextParams, data frtypes.MapStr, assoID int64) error {
	asst := &metadata.Association{}
	// TODO: can not update like this, this will cover fields that has not been updated.
	err := data.MarshalJSONInto(asst)
	if err != nil {
		errmsg := fmt.Sprintf("[operation-asst] update associaton with  %s", err.Error())
		blog.Error(errmsg)
		return params.Err.New(common.CCErrCommJSONUnmarshalFailed, errmsg)
	}

	asst.ID = assoID
	asst.OwnerID = params.SupplierAccount

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(assoID)
	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(asst.AsstObjID)
	cond.Field(metadata.AssociationFieldObjectID).Eq(asst.ObjectID)
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(params.SupplierAccount)

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to update the association (%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data) < 1 {
		errmsg := fmt.Sprintf("[operation-asst] failed to update the association , id %d not found", assoID)
		blog.Error(errmsg)
		return params.Err.New(common.CCErrAlreadyAssign, errmsg)
	}

	// check object exists

	condObj := condition.CreateCondition()
	condObj.Field(metadata.ModelFieldObjectID).Eq(asst.ObjectID)

	rspObj, err := a.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, condObj.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !rspObj.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspObj.ErrMsg)
		return params.Err.New(rspObj.Code, rspObj.ErrMsg)
	}

	condObjAsst := condition.CreateCondition()
	condObjAsst.Field(metadata.ModelFieldObjectID).Eq(asst.ObjectID)

	rspObjAsst, err := a.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, condObj.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !rspObjAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s", data, rspObjAsst.ErrMsg)
		return params.Err.New(rspObjAsst.Code, rspObjAsst.ErrMsg)
	}

	// if asst.AsstName == "" {
	// 	errmsg := fmt.Sprintf("[operation-asst] failed to create the association , asstname is required")
	// 	blog.Error(errmsg)
	// 	return params.Err.New(common.CCErrCommParamsInvalid, errmsg)
	//
	// }

	rspAsst, err := a.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), asst.ID, params.Header, asst.ToMapStr())
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

// CheckBeAssociation and return error if the obj has been bind
func (a *association) CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	exists, err := a.SearchInstAssociation(params, &metatype.QueryInput{Condition: cond.ToMapStr()})
	if nil != err {
		return err
	}

	if len(exists) > 0 {
		beAsstObject := []string{}
		for _, asst := range exists {
			beAsstObject = append(beAsstObject, asst.ObjectID)
		}
		return params.Err.Errorf(common.CCErrTopoInstHasBeenAssociation, beAsstObject)
	}
	return nil
}

// 关联关系改造后的接口

func (a *association) SearchType(ctx context.Context, h http.Header, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Asst().SearchType(ctx, h, request)
}
func (a *association) CreateType(ctx context.Context, h http.Header, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Asst().CreateType(ctx, h, request)
}
func (a *association) UpdateType(ctx context.Context, h http.Header, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Asst().UpdateType(ctx, h, asstTypeID, request)
}
func (a *association) DeleteType(ctx context.Context, h http.Header, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error) {
	return a.clientSet.ObjectController().Asst().DeleteType(ctx, h, asstTypeID)
}
func (a *association) SearchObject(ctx context.Context, h http.Header, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Asst().SearchObject(ctx, h, request)
}
func (a *association) CreateObject(ctx context.Context, h http.Header, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Asst().CreateObject(ctx, h, request)
}
func (a *association) UpdateObject(ctx context.Context, h http.Header, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Asst().UpdateObject(ctx, h, asstID, request)
}
func (a *association) DeleteObject(ctx context.Context, h http.Header, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {
	return a.clientSet.ObjectController().Asst().DeleteObject(ctx, h, asstID)
}
func (a *association) SearchInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	return a.clientSet.ObjectController().Asst().SearchInst(ctx, h, request)
}
func (a *association) CreateInst(ctx context.Context, h http.Header, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	return a.clientSet.ObjectController().Asst().CreateInst(ctx, h, request)
}
func (a *association) DeleteInst(ctx context.Context, h http.Header, request *metadata.DeleteAssociationInstRequest) (resp *metadata.DeleteAssociationInstResult, err error) {
	return a.clientSet.ObjectController().Asst().DeleteInst(ctx, h, request)
}
