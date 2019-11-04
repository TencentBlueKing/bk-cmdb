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
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

var (
	InstanceAssociationAuditHeaders = []metadata.Header{
		{
			PropertyName: "association kind",
			PropertyID:   common.AssociationKindIDField,
		},
		{
			PropertyName: "association instance id",
			PropertyID:   common.BKAsstInstIDField,
		},
		{
			PropertyName: "association model id",
			PropertyID:   common.BKAsstObjIDField,
		},
		{
			PropertyName: "instance id",
			PropertyID:   common.BKInstIDField,
		},
		{
			PropertyName: "association id",
			PropertyID:   common.AssociationObjAsstIDField,
		},
		{
			PropertyName: common.AssociationKindIDField,
			PropertyID:   "name",
		},
	}
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Object, error)
	DeleteMainlineAssociation(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64, withStatistics bool) ([]*metadata.TopoInstRst, error)
	IsMainlineObject(params types.ContextParams, objID string) (bool, error)

	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (*metadata.Association, error)
	DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error
	UpdateAssociation(params types.ContextParams, data mapstr.MapStr, assoID int64) error
	SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error)
	SearchObjectsAssociation(params types.ContextParams, objIDs []string) ([]metadata.Association, error)

	DeleteAssociation(params types.ContextParams, cond condition.Condition) error
	SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error)
	SearchInstAssociationList(params types.ContextParams, query *metadata.QueryCondition) ([]metadata.InstAsst, uint64, error)
	SearchInstAssociationUIList(params types.ContextParams, objID string, query *metadata.QueryCondition) (result interface{}, asstCnt uint64, err error)
	SearchInstAssociationSingleObjectInstInfo(params types.ContextParams, returnInstInfoObjID string, query *metadata.QueryCondition) (result []metadata.InstBaseInfo, cnt uint64, err error)
	CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error
	CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error
	DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error

	// 关联关系改造后的接口
	SearchObjectAssocWithAssocKindList(params types.ContextParams, asstKindIDs []string) (resp *metadata.AssociationList, err error)
	SearchType(params types.ContextParams, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error)
	CreateType(params types.ContextParams, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error)
	UpdateType(params types.ContextParams, asstTypeID int64, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error)
	DeleteType(params types.ContextParams, asstTypeID int64) (resp *metadata.DeleteAssociationTypeResult, err error)

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
func NewAssociationOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) AssociationOperationInterface {
	return &association{
		clientSet:   client,
		authManager: authManager,
	}
}

type association struct {
	clientSet    apimachinery.ClientSetInterface
	authManager  *extensions.AuthManager
	cls          ClassificationOperationInterface
	obj          ObjectOperationInterface
	grp          GroupOperationInterface
	attr         AttributeOperationInterface
	inst         InstOperationInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

func (assoc *association) SetProxy(cls ClassificationOperationInterface, obj ObjectOperationInterface, grp GroupOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory) {
	assoc.cls = cls
	assoc.obj = obj
	assoc.attr = attr
	assoc.inst = inst
	assoc.grp = grp
	assoc.modelFactory = targetModel
	assoc.instFactory = targetInst
}

func (assoc *association) SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error) {

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

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , err: %s, rid: %s", objID, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

func (assoc *association) SearchObjectsAssociation(params types.ContextParams, objIDs []string) ([]metadata.Association, error) {

	cond := condition.CreateCondition()
	if 0 != len(objIDs) {
		cond.Field(common.BKObjIDField).In(objIDs)
	}

	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , err: %s, rid: %s", objIDs, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

func (assoc *association) SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error) {
	intput, err := mapstr.NewFromInterface(query.Condition)
	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: intput})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the association info, query: %#v, err: %s, rid: %s", query, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

// CreateCommonAssociation create a common association, in topo model scene, which doesn't include bk_mainline association type
func (assoc *association) CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (*metadata.Association, error) {
	if data.AsstKindID == common.AssociationKindMainline {
		return nil, params.Err.Error(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}
	if len(data.AsstKindID) == 0 || len(data.AsstObjID) == 0 || len(data.ObjectID) == 0 {
		blog.Errorf("[operation-asst] failed to create the association , association kind id associate/object id is required, rid: %s", params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoAssociationMissingParameters)
	}

	// if the on delete action is empty, set none as default.
	if len(data.OnDelete) == 0 {
		data.OnDelete = metadata.NoAction
	}

	// check if this association has already exist,
	// if yes, it's not allowed to create this association

	//  check the association
	cond := condition.CreateCondition()
	cond.Field(common.AssociatedObjectIDField).Eq(data.AsstObjID)
	cond.Field(common.BKObjIDField).Eq(data.ObjectID)
	cond.Field(common.AssociationKindIDField).Eq(data.AsstKindID)

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data.Info) > 0 {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , the associations %s->%s already exist , rid: %s", params.ReqID,
			cond.ToMapStr(), data.ObjectID, data.AsstObjID)
		return nil, params.Err.Errorf(common.CCErrTopoAssociationAlreadyExist, data.ObjectID, data.AsstObjID)
	}

	// check source object exists
	if err := assoc.obj.IsValidObject(params, data.ObjectID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, err: %s, rid: %s", data.ObjectID, err.Error(), params.ReqID)
		return nil, err
	}

	if err := assoc.obj.IsValidObject(params, data.AsstObjID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, err: %s, rid: %s", data.AsstObjID, err.Error(), params.ReqID)
		return nil, err
	}

	// create a new
	rspAsst, err := assoc.clientSet.CoreService().Association().CreateModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: *data})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s, rid: %s", data, rspAsst.ErrMsg, params.ReqID)
		return nil, params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	data.ID = int64(rspAsst.Data.Created.ID)
	return data, nil
}

func (assoc *association) DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error {

	rsp, err := assoc.clientSet.CoreService().Association().DeleteInstAssociation(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to delete the inst association info , err: %s, rid: %s", rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (assoc *association) CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error {
	// create a new
	rspAsst, err := assoc.clientSet.CoreService().Association().CreateInstAssociation(context.Background(), params.Header, &metadata.CreateOneInstanceAssociation{Data: *data})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s, rid: %s", data, rspAsst.ErrMsg, params.ReqID)
		return params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	return nil
}

func (assoc *association) IsMainlineObject(params types.ContextParams, objID string) (bool, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return false, err
	}

	if !asst.Result {
		return false, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return false, fmt.Errorf("model association [%+v] not found", cond)
	}

	for _, mainline := range asst.Data.Info {
		if mainline.ObjectID == objID || mainline.AsstObjID == objID {
			return true, nil
		}
	}

	return false, nil
}

func (assoc *association) DeleteAssociationWithPreCheck(params types.ContextParams, associationID int64) error {
	// if this association has already been instantiated, then this association should not be deleted.
	// get the association with id at first.
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(associationID)
	result, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %v, rid: %s", associationID, err, params.ReqID)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !result.Result {
		blog.Errorf("[operation-asst] delete association with id[%d], but get this association for pre check failed, err: %s, rid: %s", associationID, result.ErrMsg, params.ReqID)
		return params.Err.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("[operation-asst] delete association with id[%d], but can not find this association, return now., rid: %s", associationID, params.ReqID)
		return nil
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("[operation-asst] delete association with id[%d], but got multiple association, rid: %s", associationID, params.ReqID)
		return params.Err.Error(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if result.Data.Info[0].AsstKindID == common.AssociationKindMainline {
		return params.Err.Error(common.CCErrorTopoAssociationKindMainlineUnavailable)
	}

	// find instance(s) belongs to this association
	cond = condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(result.Data.Info[0].AssociationName)
	query := metadata.QueryInput{Condition: cond.ToMapStr()}
	insts, err := assoc.SearchInstAssociation(params, &query)
	if err != nil {
		blog.Errorf("[operation-asst] delete association with id[%d], but association instance(s) failed, err: %v, rid: %s", associationID, err, params.ReqID)
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if len(insts) != 0 {
		// object association has already been instantiated, association can not be deleted.
		blog.Errorf("[operation-asst] delete association with id[%d], but has multiple instances, can not be deleted., rid: %s", associationID, params.ReqID)
		return params.Err.Error(common.CCErrTopoAssociationHasAlreadyBeenInstantiated)
	}

	// TODO: check association on_delete action before really delete this association.
	// all the pre check has finished, delete the association now.
	cond = condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(associationID)
	return assoc.DeleteAssociation(params, cond)
}

func (assoc *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %v, rid: %s", cond.ToMapStr(), err, params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("delete object association, but get association with cond[%v] failed, err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		// we assume this association has already been deleted.
		blog.Warnf("delete object association, but can not get association with cond[%v] , rid: %s", cond.ToMapStr(), params.ReqID)
		return nil
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data.Info[0].IsPre && *rsp.Data.Info[0].IsPre {
		blog.Errorf("delete object association with cond[%v], but it's a pre-defined association, can not be deleted., rid: %s", cond.ToMapStr(), params.ReqID)
		return params.Err.Error(common.CCErrorTopoDeletePredefinedAssociation)
	}

	// delete the object association
	result, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s, rid: %s", cond.ToMapStr(), result.ErrMsg, params.ReqID)
		return params.Err.Error(result.Code)
	}

	return nil
}

func (assoc *association) UpdateAssociation(params types.ContextParams, data mapstr.MapStr, assoID int64) error {
	asst := &metadata.Association{}
	err := data.MarshalJSONInto(asst)
	if err != nil {
		blog.Errorf("[operation-asst] update association with  %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if field, can := asst.CanUpdate(); !can {
		blog.Warnf("update association[%d], but request to update a forbidden update field[%s]., rid: %s", assoID, field, params.ReqID)
		return params.Err.Error(common.CCErrorTopoObjectAssociationUpdateForbiddenFields)
	}

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(assoID)
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(params.SupplierAccount)

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to update the association (%#v) , err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) < 1 {
		blog.Errorf("[operation-asst] failed to update the object association , id %d not found, rid: %s", assoID, params.ReqID)
		return params.Err.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	// a pre-defined association can not be updated.
	if nil != rsp.Data.Info[0].IsPre && *rsp.Data.Info[0].IsPre {
		blog.Errorf("update object association[%d], but it's a pre-defined association, can not be updated., rid: %s", assoID, params.ReqID)
		return params.Err.Error(common.CCErrorTopoUpdatePredefinedAssociation)
	}

	// check object exists
	if err := assoc.obj.IsValidObject(params, rsp.Data.Info[0].ObjectID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s, rid: %s", rsp.Data.Info[0].ObjectID, err.Error(), params.ReqID)
		return err
	}

	if err := assoc.obj.IsValidObject(params, rsp.Data.Info[0].AsstObjID); nil != err {
		blog.Errorf("[operation-asst] the object(%s) is invalid, error info is %s, rid: %s", rsp.Data.Info[0].AsstObjID, err.Error(), params.ReqID)
		return err
	}

	updateopt := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(assoID).ToMapStr(),
		Data:      data,
	}
	rspAsst, err := assoc.clientSet.CoreService().Association().UpdateModelAssociation(context.Background(), params.Header, &updateopt)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , err: %s, rid: %s", data, rspAsst.ErrMsg, params.ReqID)
		return params.Err.Error(rspAsst.Code)
	}

	return nil
}

// CheckBeAssociation and return error if the obj has been bind
func (assoc *association) CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	exists, err := assoc.SearchInstAssociation(params, &metadata.QueryInput{Condition: cond.ToMapStr()})
	if nil != err {
		return err
	}

	if len(exists) > 0 {
		beAsstObject := make([]string, 0)
		for _, asst := range exists {
			instRsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, asst.ObjectID,
				&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKInstIDField: asst.InstID}})
			if err != nil {
				return params.Err.Error(common.CCErrObjectSelectInstFailed)
			}
			if !instRsp.Result {
				return params.Err.New(instRsp.Code, instRsp.ErrMsg)
			}
			if len(instRsp.Data.Info) <= 0 {
				// 作为补充而存在，删除实例主机已经不存在的脏实例关联
				if delErr := assoc.DeleteInstAssociation(params, condition.CreateCondition().
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
func (assoc *association) SearchObjectAssocWithAssocKindList(params types.ContextParams, asstKindIDs []string) (resp *metadata.AssociationList, err error) {
	if len(asstKindIDs) == 0 {
		return &metadata.AssociationList{Associations: make([]metadata.AssociationDetail, 0)}, nil
	}

	asso := make([]metadata.AssociationDetail, 0)
	for _, id := range asstKindIDs {
		cond := condition.CreateCondition()
		cond.Field(common.AssociationKindIDField).Eq(id)

		r, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v, rid: %s", id, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !r.Result {
			blog.Errorf("get object association list with association kind[%s] failed, err: %v, rid: %s", id, r.ErrMsg, params.ReqID)
			return nil, params.Err.Errorf(r.Code, r.ErrMsg)
		}

		asso = append(asso, metadata.AssociationDetail{AssociationKindID: id, Associations: r.Data.Info})
	}

	return &metadata.AssociationList{Associations: asso}, nil
}

func (assoc *association) SearchType(params types.ContextParams, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
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

	return assoc.clientSet.CoreService().Association().ReadAssociationType(context.Background(), params.Header, &input)
}

func (assoc *association) CreateType(params types.ContextParams, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {

	rsp, err := assoc.clientSet.CoreService().Association().CreateAssociationType(params.Context, params.Header, &metadata.CreateAssociationKind{Data: *request})
	if err != nil {
		blog.Errorf("create association type failed, kind id: %s, err: %v, rid: %s", request.AssociationKindID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}
	if rsp.Result == false || rsp.Code != 0 {
		blog.ErrorJSON("create association type failed, request: %s, response: %s, rid: %s", request, rsp, params.ReqID)
		return nil, errors.NewCCError(rsp.Code, rsp.ErrMsg)
	}
	resp = &metadata.CreateAssociationTypeResult{BaseResp: rsp.BaseResp}
	resp.Data.ID = int64(rsp.Data.Created.ID)
	request.ID = resp.Data.ID
	if err := assoc.authManager.RegisterAssociationTypeByID(params.Context, params.Header, resp.Data.ID); err != nil {
		blog.Error("create association type: %s success, but register id: %d to auth failed, err: %v, rid: %s", request.AssociationKindID, resp.Data.ID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
	}

	return resp, nil

}

func (assoc *association) UpdateType(params types.ContextParams, asstTypeID int64, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	if len(request.AsstName) != 0 {
		if err := assoc.authManager.UpdateAssociationTypeByID(params.Context, params.Header, asstTypeID); err != nil {
			blog.Errorf("update association type %s, but got update resource to auth failed, err: %v, rid: %s", request.AsstName, err, params.ReqID)
			return nil, params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
		}
	}

	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstTypeID).ToMapStr(),
		Data:      mapstr.NewFromStruct(request, "json"),
	}

	rsp, err := assoc.clientSet.CoreService().Association().UpdateAssociationType(context.Background(), params.Header, &input)
	if err != nil {
		blog.Errorf("update association type failed, kind id: %d, err: %v, rid: %s", asstTypeID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}
	resp = &metadata.UpdateAssociationTypeResult{BaseResp: rsp.BaseResp}
	return resp, nil
}

func (assoc *association) DeleteType(params types.ContextParams, asstTypeID int64) (resp *metadata.DeleteAssociationTypeResult, err error) {
	if err := assoc.authManager.DeregisterAssociationTypeByIDs(params.Context, params.Header, asstTypeID); err != nil {
		blog.Errorf("delete association type id: %d, but deregister from auth failed, err: %v, rid: %s", asstTypeID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrCommUnRegistResourceToIAMFailed, err.Error())
	}
	cond := condition.CreateCondition()
	cond.Field("id").Eq(asstTypeID)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	query := &metadata.SearchAssociationTypeRequest{
		Condition: cond.ToMapStr(),
	}

	result, err := assoc.SearchType(params, query)
	if err != nil {
		blog.Errorf("delete association kind[%d], but get detailed info failed, err: %v, rid: %s", asstTypeID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete association kind[%d], but get detailed info failed, err: %s, rid: %s", asstTypeID, result.ErrMsg, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(result.Data.Info) > 1 {
		blog.Errorf("delete association kind[%d], but get multiple instance, rid: %s", asstTypeID, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoGetMultipleAssocKindInstWithOneID)
	}

	if len(result.Data.Info) == 0 {
		return &metadata.DeleteAssociationTypeResult{BaseResp: metadata.SuccessBaseResp, Data: common.CCSuccessStr}, nil
	}

	if result.Data.Info[0].IsPre != nil && *result.Data.Info[0].IsPre {
		blog.Errorf("delete association kind[%d], but this is a pre-defined association kind, can not be deleted., rid: %s", asstTypeID, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoDeletePredefinedAssociationKind)
	}

	// a already used association kind can not be deleted.
	cond = condition.CreateCondition()
	cond.Field(common.AssociationKindIDField).Eq(result.Data.Info[0].AssociationKindID)
	asso, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %v, rid: %s", asstTypeID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("delete association kind[%d], but get objects that used this asso kind failed, err: %s, rid: %s", asstTypeID, result.ErrMsg, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(asso.Data.Info) != 0 {
		blog.Warnf("delete association kind[%d], but it has already been used, can not be deleted., rid: %s", asstTypeID, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoAssociationKindHasBeenUsed)
	}

	rsp, err := assoc.clientSet.CoreService().Association().DeleteAssociationType(
		context.Background(), params.Header, &metadata.DeleteOption{
			Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstTypeID).ToMapStr(),
		},
	)
	if err != nil {
		blog.Errorf("delete association type failed, kind id: %d, err: %v, rid: %s", asstTypeID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrTopoCreateAssocKindFailed, err.Error())
	}

	return &metadata.DeleteAssociationTypeResult{BaseResp: rsp.BaseResp}, nil
}

func (assoc *association) SearchObject(params types.ContextParams, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: request.Condition})

	resp = &metadata.SearchAssociationObjectResult{BaseResp: rsp.BaseResp, Data: []*metadata.Association{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, err
}

func (assoc *association) CreateObject(params types.ContextParams, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	rsp, err := assoc.clientSet.CoreService().Association().CreateModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: *request})

	resp = &metadata.CreateAssociationObjectResult{
		BaseResp: rsp.BaseResp,
	}
	resp.Data.ID = int64(rsp.Data.Created.ID)
	return resp, err
}

func (assoc *association) UpdateObject(params types.ContextParams, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstID).ToMapStr(),
		Data:      mapstr.NewFromStruct(request, "json"),
	}

	rsp, err := assoc.clientSet.CoreService().Association().UpdateModelAssociation(context.Background(), params.Header, &input)
	resp = &metadata.UpdateAssociationObjectResult{
		BaseResp: rsp.BaseResp,
	}
	return resp, err
}

func (assoc *association) DeleteObject(params types.ContextParams, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {

	input := metadata.DeleteOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(asstID).ToMapStr(),
	}
	rsp, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(context.Background(), params.Header, &input)
	return &metadata.DeleteAssociationObjectResult{BaseResp: rsp.BaseResp}, err

}

func (assoc *association) SearchInst(params types.ContextParams, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: request.Condition})

	resp = &metadata.SearchAssociationInstResult{BaseResp: rsp.BaseResp, Data: []*metadata.InstAsst{}}
	for index := range rsp.Data.Info {
		resp.Data = append(resp.Data, &rsp.Data.Info[index])
	}

	return resp, err
}

func (assoc *association) CreateInst(params types.ContextParams, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	var bizID int64
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			blog.Errorf("parse business id from request failed, params: %+v, err: %+v, rid: %s", params, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPInputInvalid)
		}
	}

	cond := condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
	result, err := assoc.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}

	if len(result.Data) == 0 {
		blog.Errorf("create instance association, but can not find object association[%s]. rid: %s", request.ObjectAsstID, params.ReqID)
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
		instance, err := assoc.SearchInst(params, &metadata.SearchAssociationInstRequest{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}

		cond = condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
		cond.Field(common.BKAsstInstIDField).Eq(request.AsstInstID)

		instance, err = assoc.SearchInst(params, &metadata.SearchAssociationInstRequest{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToOneAssociation)
		}
	case metadata.OneToManyMapping:
		cond = condition.CreateCondition()
		cond.Field(common.AssociationObjAsstIDField).Eq(request.ObjectAsstID)
		cond.Field(common.BKAsstInstIDField).Eq(request.AsstInstID)

		instance, err := assoc.SearchInst(params, &metadata.SearchAssociationInstRequest{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instance.Result {
			blog.Errorf("create association instance, but check instance with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}
		if len(instance.Data) >= 1 {
			return nil, params.Err.Error(common.CCErrorTopoCreateMultipleInstancesForOneToManyAssociation)
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
	createResult, err := assoc.clientSet.CoreService().Association().CreateInstAssociation(context.Background(), params.Header, &input)
	if err != nil {
		blog.Errorf("create instance association failed, do coreservice create failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	resp = &metadata.CreateAssociationInstResult{BaseResp: createResult.BaseResp}
	instanceAssociationID := int64(createResult.Data.Created.ID)
	resp.Data.ID = instanceAssociationID

	curData := mapstr.NewFromStruct(input.Data, "json")
	curData.Set("name", objectAsst.AssociationAliasName)
	// record audit log
	auditlog := metadata.SaveAuditLogParams{
		ID:    request.InstID,
		Model: objID,
		Content: metadata.Content{
			CurData: curData,
			Headers: InstanceAssociationAuditHeaders,
		},
		OpDesc: "create instance association",
		OpType: auditoplog.AuditOpTypeAdd,
		BizID:  bizID,
	}
	auditresp, err := assoc.clientSet.CoreService().Audit().SaveAuditLog(params.Context, params.Header, auditlog)
	if err != nil {
		blog.Errorf("CreateInst success, but save audit log failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrAuditSaveLogFailed)
	}
	if !auditresp.Result {
		blog.Errorf("CreateInst success, but save audit log failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.New(auditresp.Code, auditresp.ErrMsg)
	}

	return resp, err
}

func (assoc *association) DeleteInst(params types.ContextParams, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error) {
	var bizID int64
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			blog.Errorf("parse business id from request failed, params: %+v, err: %+v, rid: %s", params, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPInputInvalid)
		}
	}

	// record audit log
	searchCondition := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(assoID).ToMapStr(),
	}
	data, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, &searchCondition)
	if err != nil {
		blog.Errorf("DeleteInst failed, get instance association failed, params: %+v, err: %+v, rid: %s", params, err, params.ReqID)
		return nil, err
	}
	if len(data.Data.Info) == 0 {
		blog.Errorf("DeleteInst failed, instance association not found, searchCondition: %+v, err: %+v, rid: %s", searchCondition, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommNotFound)
	}
	if len(data.Data.Info) > 1 {
		blog.Errorf("DeleteInst failed, get instance association with id:%d get multiple, err: %+v, rid: %s", assoID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommNotFound)
	}

	instanceAssociation := data.Data.Info[0]

	cond := condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).Eq(instanceAssociation.ObjectAsstID)
	assInfoResult, err := assoc.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !assInfoResult.Result {
		blog.Errorf("create association instance, but search object association with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}

	input := metadata.DeleteOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(assoID).ToMapStr(),
	}
	rsp, err := assoc.clientSet.CoreService().Association().DeleteInstAssociation(context.Background(), params.Header, &input)
	resp = &metadata.DeleteAssociationInstResult{
		BaseResp: rsp.BaseResp,
	}

	preData := mapstr.NewFromStruct(instanceAssociation, "json")
	if len(assInfoResult.Data) > 0 {
		preData.Set("name", assInfoResult.Data[0].AssociationAliasName)
	}
	auditlog := metadata.SaveAuditLogParams{
		ID:    instanceAssociation.InstID,
		Model: instanceAssociation.ObjectID,
		Content: metadata.Content{
			PreData: preData,
			Headers: InstanceAssociationAuditHeaders,
		},
		OpDesc: "delete instance association",
		OpType: auditoplog.AuditOpTypeDel,
		BizID:  bizID,
	}
	auditresp, err := assoc.clientSet.CoreService().Audit().SaveAuditLog(params.Context, params.Header, auditlog)
	if err != nil {
		blog.Errorf("DeleteInst finished, but save audit log failed, delete inst response: %+v, err: %v, rid: %s", auditresp, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrAuditSaveLogFailed)
	}
	if !auditresp.Result {
		blog.Errorf("DeleteInst finished, but save audit log failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.New(auditresp.Code, auditresp.ErrMsg)
	}

	return resp, err
}

// SearchInstAssociationList 与实例有关系的实例关系数据,以分页的方式返回
func (assoc *association) SearchInstAssociationList(params types.ContextParams, query *metadata.QueryCondition) ([]metadata.InstAsst, uint64, error) {

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, query)
	if nil != err {
		blog.Errorf("ReadInstAssociation http do error, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, 0, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.ErrorJSON("ReadInstAssociation http response error, query: %s, response: %s, rid: %s", query, rsp, params.ReqID)
		return nil, 0, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, rsp.Data.Count, nil
}

// SearchInstAssociationUIList 与实例有关系的实例关系数据,以分页的方式返回
func (assoc *association) SearchInstAssociationUIList(params types.ContextParams, objID string, query *metadata.QueryCondition) (result interface{}, asstCnt uint64, err error) {

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, query)
	if nil != err {
		blog.Errorf("ReadInstAssociation http do error, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, 0, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.ErrorJSON("ReadInstAssociation http response error, query: %s, response: %s, rid: %s", query, rsp, params.ReqID)
		return nil, 0, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	objIDInstIDMap := make(map[string][]int64, 0)
	var objSrcAsstArr []metadata.InstAsst
	var objDstAsstArr []metadata.InstAsst
	for _, instAsst := range rsp.Data.Info {
		objIDInstIDMap[instAsst.ObjectID] = append(objIDInstIDMap[instAsst.ObjectID], instAsst.InstID)
		objIDInstIDMap[instAsst.AsstObjectID] = append(objIDInstIDMap[instAsst.AsstObjectID], instAsst.AsstInstID)
		if instAsst.ObjectID == objID {
			objSrcAsstArr = append(objSrcAsstArr, instAsst)
		} else {
			objDstAsstArr = append(objDstAsstArr, instAsst)

		}
	}

	instInfo := make(map[string][]mapstr.MapStr, 0)
	for instObjID, instIDArr := range objIDInstIDMap {
		idField := metadata.GetInstIDFieldByObjID(instObjID)
		cond := condition.CreateCondition()
		cond.Field(idField).In(instIDArr)
		input := &metadata.QueryCondition{
			Condition: cond.ToMapStr(),
			Limit: metadata.SearchLimit{
				Offset: 0,
				Limit:  common.BKNoLimit,
			},
			Fields: []string{metadata.GetInstNameFieldName(instObjID), idField},
		}
		instResp, err := assoc.clientSet.CoreService().Instance().
			ReadInstance(context.Background(), params.Header, instObjID, input)
		if err != nil {
			blog.Errorf("ReadInstance http do error, err: %s, input:%s, rid: %s", err.Error(), input, params.ReqID)
			return nil, 0, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
		}
		if !instResp.Result {
			blog.ErrorJSON("ReadInstance http response error, query: %s, response: %s, rid: %s", query, rsp, params.ReqID)
			return nil, 0, params.Err.New(rsp.Code, rsp.ErrMsg)
		}
		instInfo[instObjID] = instResp.Data.Info

	}
	instAsstMap := map[string][]metadata.InstAsst{
		"src": objSrcAsstArr,
		"dst": objDstAsstArr,
	}

	result = mapstr.MapStr{
		"association": instAsstMap,
		"instance":    instInfo,
	}

	return result, rsp.Data.Count, nil
}

// SearchInstAssociationUIList 与实例有关系的实例关系数据,以分页的方式返回
// returnInstInfoObjID 根据条件查询出来关联关系，需要返回实例信息（实例名，实例ID）的模型ID
func (assoc *association) SearchInstAssociationSingleObjectInstInfo(params types.ContextParams, returnInstInfoObjID string, query *metadata.QueryCondition) (result []metadata.InstBaseInfo, cnt uint64, err error) {

	rsp, err := assoc.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), params.Header, query)
	if nil != err {
		blog.Errorf("ReadInstAssociation http do error, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, 0, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.ErrorJSON("ReadInstAssociation http response error, query: %s, response: %s, rid: %s", query, rsp, params.ReqID)
		return nil, 0, params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	// association count
	cnt = rsp.Data.Count

	if cnt == 0 {
		return nil, 0, nil
	}

	var objIDInstIDArr []int64

	for _, instAsst := range rsp.Data.Info {
		if instAsst.ObjectID == returnInstInfoObjID {
			objIDInstIDArr = append(objIDInstIDArr, instAsst.InstID)
		} else if instAsst.AsstObjectID == returnInstInfoObjID {
			objIDInstIDArr = append(objIDInstIDArr, instAsst.AsstInstID)

		}
	}

	idField := metadata.GetInstIDFieldByObjID(returnInstInfoObjID)
	nameField := metadata.GetInstNameFieldName(returnInstInfoObjID)
	cond := condition.CreateCondition()
	cond.Field(idField).In(objIDInstIDArr)
	input := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Limit: metadata.SearchLimit{
			Offset: 0,
			Limit:  common.BKNoLimit,
		},
		Fields: []string{nameField, idField},
	}
	instResp, err := assoc.clientSet.CoreService().Instance().
		ReadInstance(context.Background(), params.Header, returnInstInfoObjID, input)
	if err != nil {
		blog.Errorf("ReadInstance http do error, err: %s, input:%s, rid: %s", err.Error(), input, params.ReqID)
		return nil, 0, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !instResp.Result {
		blog.ErrorJSON("ReadInstance http response error, query: %s, response: %s, rid: %s", query, rsp, params.ReqID)
		return nil, 0, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	result = make([]metadata.InstBaseInfo, 0)
	for _, row := range instResp.Data.Info {
		id, err := row.Int64(idField)
		if err != nil {
			blog.ErrorJSON("ReadInstance  convert field(%s) to int error. err:%s, inst:%s,rid:%s", idField, err.Error(), row, params.ReqID)
			// CCErrCommInstFieldConvertFail  convert %s  field %s to %s error %s
			return nil, 0, params.Err.Errorf(common.CCErrCommInstFieldConvertFail, returnInstInfoObjID, idField, "int", err.Error())
		}
		name, err := row.String(nameField)
		if err != nil {
			blog.ErrorJSON("ReadInstance  convert field(%s) to int error. err:%s, inst:%s,rid:%s", nameField, err.Error(), row, params.ReqID)
			// CCErrCommInstFieldConvertFail  convert %s  field %s to %s error %s
			return nil, 0, params.Err.Errorf(common.CCErrCommInstFieldConvertFail, returnInstInfoObjID, idField, "string", err.Error())
		}
		result = append(result, metadata.InstBaseInfo{
			ID:   id,
			Name: name,
		})
	}

	return result, cnt, nil
}
