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
	"context"
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Group group opeartion interface declaration
type Group interface {
	Operation

	Parse(data frtypes.MapStr) (*metadata.Group, error)
	CreateAttribute() Attribute

	GetAttributes() ([]Attribute, error)

	Origin() metadata.Group

	SetID(groupID string)
	GetID() string

	SetName(groupName string)
	GetName() string

	SetIndex(groupIndex int64)
	GetIndex() int64

	SetObjectID(objID string)
	GetObjectID() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetDefault(isDefault bool)
	GetDefault() bool

	SetIsPre(isPre bool)
	GetIsPre() bool

	ToMapStr() (frtypes.MapStr, error)
}

var _ Group = (*group)(nil)

type group struct {
	grp       metadata.Group
	isNew     bool
	params    types.ContextParams
	clientSet apimachinery.ClientSetInterface
}

func (g *group) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.grp)
}

func (g *group) Origin() metadata.Group {
	return g.grp
}

func (g *group) SetObjectID(objID string) {
	g.grp.ObjectID = objID
}
func (g *group) GetObjectID() string {
	return g.grp.ObjectID
}

func (g *group) Create() error {

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(g.grp.GroupID)
	cond.Field(metadata.GroupFieldGroupIndex).Eq(g.grp.GroupIndex)
	cond.Field(metadata.GroupFieldGroupName).Eq(g.grp.GroupName)

	grpItems, err := g.search(cond)
	if nil != err {
		blog.Errorf("[model-grp] failed to search the groups by the condition(%#v), error info is %s", cond, err.Error())
		return err
	}

	if 0 != len(grpItems) {
		blog.Errorf("[model-grp] the group(%#v) is repeated", g.grp)
		return g.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	rsp, err := g.clientSet.ObjectController().Meta().CreatePropertyGroup(context.Background(), g.params.Header, &g.grp)

	if nil != err {
		blog.Errorf("[model-grp] failed to request object controller, error info is %s", err.Error())
		return g.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp] failed to create the group(%s), error info is is %s", g.grp.GroupID, rsp.ErrMsg)
		return g.params.Err.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	g.grp.ID = rsp.Data.ID

	return nil
}

func (g *group) Update(data frtypes.MapStr) error {

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(g.grp.GroupID)
	grps, err := g.search(cond)
	if nil != err {
		return err
	}

	for _, grpItem := range grps { // only one item

		cond := &metadata.UpdateGroupCondition{}
		cond.Condition.GroupID = grpItem.GroupID
		cond.Data.Index = g.grp.GroupIndex
		cond.Data.Name = g.grp.GroupName

		rsp, err := g.clientSet.ObjectController().Meta().UpdatePropertyGroup(context.Background(), g.params.Header, cond)
		if nil != err {
			blog.Errorf("[model-grp]failed to request object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-grp]failed to update the group(%s), error info is %s", grpItem.GroupID, err.Error())
			return g.params.Err.Error(common.CCErrTopoObjectAttributeUpdateFailed)
		}

		g.grp.ID = grpItem.ID

	}
	return nil
}

func (g *group) Delete() error {

	rsp, err := g.clientSet.ObjectController().Meta().DeletePropertyGroup(context.Background(), g.grp.GroupID, g.params.Header)
	if nil != err {
		blog.Error("[model-grp]failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp]failed to delte the group(%s), error info is %s", g.grp.GroupID, rsp.ErrMsg)
		return g.params.Err.Error(common.CCErrTopoObjectGroupDeleteFailed)
	}

	return nil
}

func (g *group) IsExists() (bool, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(g.grp.GroupID)
	cond.Field(metadata.GroupFieldGroupName).Eq(g.grp.GroupName)

	rsp, err := g.clientSet.ObjectController().Meta().SelectGroup(context.Background(), g.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request object controller ,error info is %s", err.Error())
		return false, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to query group, error info is  %s", rsp.ErrMsg)
		return false, g.params.Err.Error(common.CCErrTopoObjectGroupSelectFailed)
	}

	return 0 != len(rsp.Data), nil
}

func (g *group) Parse(data frtypes.MapStr) (*metadata.Group, error) {

	err := metadata.SetValueToStructByTags(&g.grp, data)
	return &g.grp, err
}
func (g *group) ToMapStr() (frtypes.MapStr, error) {

	rst := metadata.SetValueToMapStrByTags(&g.grp)
	return rst, nil
}

func (g *group) GetAttributes() ([]Attribute, error) {
	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(g.grp.ObjectID).
		Field(metadata.AttributeFieldPropertyGroup).Eq(g.grp.GroupID).
		Field(metadata.AttributeFieldSupplierAccount).Eq(g.params.SupplierAccount)

	rsp, err := g.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), g.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, g.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", g.grp.ObjectID, rsp.ErrMsg)
		return nil, g.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Attribute, 0)
	for _, item := range rsp.Data {

		attr := &attribute{
			attr:      item,
			params:    g.params,
			clientSet: g.clientSet,
		}

		rstItems = append(rstItems, attr)
	}

	return rstItems, nil
}

func (g *group) search(cond condition.Condition) ([]metadata.Group, error) {

	rsp, err := g.clientSet.ObjectController().Meta().SelectGroup(context.Background(), g.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the classificaiont, error info is %s", rsp.ErrMsg)
		return nil, g.params.Err.Error(rsp.Code)
	}

	return rsp.Data, nil
}
func (g *group) Save() error {

	if exists, err := g.IsExists(); nil != err {
		return err
	} else if !exists {
		return g.Create()
	}
	data := metadata.SetValueToMapStrByTags(g.grp)
	return g.Update(data)
}

func (g *group) CreateAttribute() Attribute {
	return &attribute{
		params:    g.params,
		clientSet: g.clientSet,
		attr: metadata.Attribute{
			OwnerID:  g.grp.OwnerID,
			ObjectID: g.grp.ObjectID,
		},
	}
}

func (g *group) SetID(groupID string) {
	g.grp.GroupID = groupID
}

func (g *group) GetID() string {
	return g.grp.GroupID
}

func (g *group) SetName(groupName string) {
	g.grp.GroupName = groupName
}

func (g *group) GetName() string {
	return g.grp.GroupName
}

func (g *group) SetIndex(groupIndex int64) {
	g.grp.GroupIndex = groupIndex
}

func (g *group) GetIndex() int64 {
	return int64(g.grp.GroupIndex)
}

func (g *group) SetSupplierAccount(supplierAccount string) {
	g.grp.OwnerID = supplierAccount
}

func (g *group) GetSupplierAccount() string {
	return g.grp.OwnerID
}

func (g *group) SetDefault(isDefault bool) {
	g.grp.IsDefault = isDefault
}

func (g *group) GetDefault() bool {
	return g.grp.IsDefault
}

func (g *group) SetIsPre(isPre bool) {
	g.grp.IsPre = isPre
}

func (g *group) GetIsPre() bool {
	return g.grp.IsPre
}
