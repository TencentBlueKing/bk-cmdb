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
	DeleteMainlineAssociaton(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error)
	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) error
	DeleteAssociation(params types.ContextParams, cond condition.Condition) error
	UpdateAssociation(params types.ContextParams, data frtypes.MapStr, cond condition.Condition) error
	SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error)
	SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error)
	CheckBeAssociation(params types.ContextParams, obj model.Object, cond condition.Condition) error
	CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error
	DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error

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
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , error info is %s", objID, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}

func (a *association) SearchInstAssociation(params types.ContextParams, query *metadata.QueryInput) ([]metadata.InstAsst, error) {

	rsp, err := a.clientSet.ObjectController().Instance().SearchObjects(context.Background(), common.BKTableNameInstAsst, params.Header, query)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the association info, query: %#v, error info is %s", query, rsp.ErrMsg)
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

	//  check the association
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(data.AsstObjID)
	cond.Field(metadata.AssociationFieldObjectID).Eq(data.ObjectID)
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(params.SupplierAccount)

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	// create a new

	rspAsst, err := a.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), params.Header, data)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , error info is %s", data, rspAsst.ErrMsg)
		return params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	return nil
}

func (a *association) DeleteInstAssociation(params types.ContextParams, cond condition.Condition) error {

	rsp, err := a.clientSet.ObjectController().Instance().DelObject(context.Background(), common.BKTableNameInstAsst, params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to delete the inst association info , error info is %s", rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (a *association) CreateCommonInstAssociation(params types.ContextParams, data *metadata.InstAsst) error {
	// create a new

	rspAsst, err := a.clientSet.ObjectController().Instance().CreateObject(context.Background(), common.BKTableNameInstAsst, params.Header, data)
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rspAsst.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , error info is %s", data, rspAsst.ErrMsg)
		return params.Err.New(rspAsst.Code, rspAsst.ErrMsg)
	}

	return nil
}
func (a *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {

	// delete the object association
	rsp, err := a.clientSet.ObjectController().Meta().DeleteObjectAssociation(context.Background(), 0, params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to create the association (%#v) , error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
func (a *association) UpdateAssociation(params types.ContextParams, data frtypes.MapStr, cond condition.Condition) error {
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
