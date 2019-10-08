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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	CreateObjectGroup(params types.ContextParams, data mapstr.MapStr) (model.GroupInterface, error)
	DeleteObjectGroup(params types.ContextParams, groupID int64) error
	FindObjectGroup(params types.ContextParams, cond condition.Condition) ([]model.GroupInterface, error)
	FindGroupByObject(params types.ContextParams, objID string, cond condition.Condition) ([]model.GroupInterface, error)
	UpdateObjectGroup(params types.ContextParams, cond *metadata.UpdateGroupCondition) error
	UpdateObjectAttributeGroup(params types.ContextParams, cond []metadata.PropertyGroupObjectAtt) error
	DeleteObjectAttributeGroup(params types.ContextParams, objID, propertyID, groupID string) error
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface)
}

// NewGroupOperation create a new group operation instance
func NewGroupOperation(client apimachinery.ClientSetInterface) GroupOperationInterface {
	return &group{
		clientSet: client,
	}
}

type group struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	obj          ObjectOperationInterface
}

func (g *group) SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface) {
	g.modelFactory = modelFactory
	g.instFactory = instFactory
	g.obj = obj
}

func (g *group) CreateObjectGroup(params types.ContextParams, data mapstr.MapStr) (model.GroupInterface, error) {

	grp := g.modelFactory.CreateGroup(params)

	_, err := grp.Parse(data)
	if nil != err {
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}

	//  check the object
	if err = g.obj.IsValidObject(params, grp.Group().ObjectID); nil != err {
		blog.Errorf("[operation-grp] the group (%#v) is in valid, rid: %s", data, params.ReqID)
		return nil, params.Err.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	err = grp.Create()
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	return grp, nil
}

func (g *group) DeleteObjectGroup(params types.ContextParams, groupID int64) error {
	cond := condition.CreateCondition().Field(common.BKFieldID).Eq(groupID)
	rsp, err := g.clientSet.CoreService().Model().DeleteAttributeGroupByCondition(context.Background(), params.Header, metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-grp]failed to request object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp]failed to delete the group(%d), err: %s, rid: %s", groupID, rsp.ErrMsg, params.ReqID)
		return params.Err.Error(common.CCErrTopoObjectGroupDeleteFailed)
	}

	return nil
}

func (g *group) FindObjectGroup(params types.ContextParams, cond condition.Condition) ([]model.GroupInterface, error) {
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}
	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroupByCondition(context.Background(), params.Header, metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[opeartion-grp] failed to search the group by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateGroup(params, g.clientSet, rsp.Data.Info), nil
}

func (g *group) FindGroupByObject(params types.ContextParams, objID string, cond condition.Condition) ([]model.GroupInterface, error) {

	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroup(context.Background(), params.Header, objID, metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to search the group of the object(%s) by the condition (%#v), error info is %s, rid: %s", objID, fCond, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateGroup(params, g.clientSet, rsp.Data.Info), nil
}

func (g *group) UpdateObjectAttributeGroup(params types.ContextParams, conds []metadata.PropertyGroupObjectAtt) error {
	for _, cond := range conds {
		// if the target group doesn't exist, don't change the original group
		grpCond := condition.CreateCondition()
		grpCond.Field(metadata.GroupFieldGroupID).Eq(cond.Data.PropertyGroupID)
		grpCond.Field(metadata.GroupFieldObjectID).Eq(cond.Condition.ObjectID)
		grps, err := g.FindObjectGroup(params, grpCond)
		if nil != err {
			blog.Errorf("[operation-grp] failed to get the group  by the condition (%#v), error info is %s , rid: %s", cond, err.Error(), params.ReqID)
			return err
		}
		if len(grps) != 1 {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is group is invalid, rid: %s", cond, params.ReqID)
			return params.Err.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupID)
		}
	}

	for _, cond := range conds {
		input := metadata.UpdateOption{
			Condition: mapstr.NewFromStruct(cond.Condition, "json"),
			Data:      mapstr.NewFromStruct(cond.Data, "json"),
		}

		rsp, err := g.clientSet.CoreService().Model().UpdateModelAttrsByCondition(context.Background(), params.Header, &input)
		if nil != err {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is %s , rid: %s", cond, err.Error(), params.ReqID)
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is %s , rid: %s", cond, rsp.ErrMsg, params.ReqID)
			return params.Err.New(rsp.Code, rsp.ErrMsg)
		}
	}

	return nil
}

func (g *group) DeleteObjectAttributeGroup(params types.ContextParams, objID, propertyID, groupID string) error {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().
			Field(common.BKObjIDField).Eq(objID).
			Field(common.BKPropertyIDField).Eq(propertyID).
			Field(common.BKPropertyGroupField).Eq(groupID).ToMapStr(),
		Data: mapstr.MapStr{
			"bk_property_index":         -1,
			common.BKPropertyGroupField: "default",
		},
	}

	rsp, err := g.clientSet.CoreService().Model().UpdateModelAttrs(context.Background(), params.Header, objID, &input)
	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group , error info is %s , rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to set the group , error info is %s , rid: %s", rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (g *group) UpdateObjectGroup(params types.ContextParams, cond *metadata.UpdateGroupCondition) error {

	if cond.Data.Index == nil && cond.Data.Name == nil {
		return nil
	}
	input := metadata.UpdateOption{
		Condition: mapstr.NewFromStruct(cond.Condition, "json"),
		Data:      mapstr.NewFromStruct(cond.Data, "json"),
	}
	rsp, err := g.clientSet.CoreService().Model().UpdateAttributeGroupByCondition(context.Background(), params.Header, input)
	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
