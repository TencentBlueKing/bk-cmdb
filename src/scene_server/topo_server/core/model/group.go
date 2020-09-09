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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"github.com/rs/xid"
)

// Group group opeartion interface declaration
type GroupInterface interface {
	Operation
	Parse(data mapstr.MapStr) (*metadata.Group, error)
	CreateAttribute() AttributeInterface
	GetAttributes() ([]AttributeInterface, error)
	Group() metadata.Group
	SetGroup(grp metadata.Group)
	ToMapStr() mapstr.MapStr
}

var _ GroupInterface = (*group)(nil)

func NewGroup(param *rest.Kit, cli apimachinery.ClientSetInterface, bizID int64) GroupInterface {
	return &group{
		grp:       metadata.Group{},
		kit:       param,
		bizID:     bizID,
		clientSet: cli,
		ownerID:   param.SupplierAccount,
	}
}

func NewGroupID(isDefault bool) string {
	if isDefault {
		return "default"
	} else {
		return xid.New().String()
	}
}

type group struct {
	FieldValid
	grp       metadata.Group
	isNew     bool
	kit       *rest.Kit
	bizID     int64
	clientSet apimachinery.ClientSetInterface
	ownerID   string
}

func (g *group) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.grp)
}

func (g *group) Group() metadata.Group {
	return g.grp
}

func (g *group) SetGroup(grp metadata.Group) {
	g.grp = grp
	g.grp.OwnerID = g.ownerID
}

func (g *group) SetObjectID(objID string) {
	g.grp.ObjectID = objID
}
func (g *group) GetObjectID() string {
	return g.grp.ObjectID
}

func (g *group) IsValid(isUpdate bool, data mapstr.MapStr) error {

	if !isUpdate || data.Exists(metadata.GroupFieldGroupID) {
		if _, err := g.FieldValid.Valid(g.kit, data, metadata.GroupFieldGroupID); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.GroupFieldGroupName) {
		val, err := g.FieldValid.Valid(g.kit, data, metadata.GroupFieldGroupName)
		if nil != err {
			return err
		}
		if err = g.FieldValid.ValidName(g.kit, val); nil != err {
			return err
		}
	}

	return nil
}

func (g *group) Create() error {
	if err := g.IsValid(false, g.grp.ToMapStr()); nil != err {
		return err
	}

	rsp, err := g.clientSet.CoreService().Model().CreateAttributeGroup(context.Background(), g.kit.Header, g.GetObjectID(), metadata.CreateModelAttributeGroup{Data: g.grp})
	if nil != err {
		blog.Errorf("[model-grp] failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return g.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[model-grp] failed to create the group(%s), err: is %s, rid: %s", g.grp.GroupID, rsp.ErrMsg, g.kit.Rid)
		return g.kit.CCError.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	g.grp.ID = int64(rsp.Data.Created.ID)

	return nil
}

func (g *group) Update(data mapstr.MapStr) error {
	if err := g.IsValid(true, data); nil != err {
		return err
	}

	exists, err := g.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return g.kit.CCError.Errorf(common.CCErrCommDuplicateItem, g.Group().GroupName)
	}

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(g.grp.GroupID)
	grps, err := g.search(cond)
	if nil != err {
		return err
	}

	for _, grpItem := range grps { // only one item

		input := metadata.UpdateOption{
			Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(grpItem.GroupID).ToMapStr(),
			Data: mapstr.MapStr{
				common.BKPropertyGroupIndexField: g.grp.GroupIndex,
				common.BKPropertyGroupNameField:  g.grp.GroupName,
			},
		}

		rsp, err := g.clientSet.CoreService().Model().UpdateAttributeGroup(context.Background(), g.kit.Header, g.GetObjectID(), input)
		if nil != err {
			blog.Errorf("[model-grp]failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
			return err
		}

		if !rsp.Result {
			blog.Errorf("[model-grp]failed to update the group(%s), err: %s, rid: %s", grpItem.GroupID, err.Error(), g.kit.Rid)
			return g.kit.CCError.Error(common.CCErrTopoObjectAttributeUpdateFailed)
		}

		g.grp.ID = grpItem.ID

	}
	return nil
}

func (g *group) IsExists() (bool, error) {
	// check id
	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(g.grp.GroupID)
	cond.Field(metadata.ModelFieldObjectID).Eq(g.grp.ObjectID)
	cond.Field(metadata.GroupFieldID).NotIn([]int64{g.grp.ID})
	grps, err := g.search(cond)
	if nil != err {
		return false, err
	}
	if 0 != len(grps) {
		return true, nil
	}

	// check name
	cond = condition.CreateCondition()
	cond.Field(metadata.GroupFieldID).NotIn([]int64{g.grp.ID})
	cond.Field(metadata.ModelFieldObjectID).Eq(g.grp.ObjectID)
	cond.Field(metadata.GroupFieldGroupName).Eq(g.grp.GroupName)
	grps, err = g.search(cond)
	if nil != err {
		return false, err
	}
	if 0 != len(grps) {
		return true, nil
	}

	return false, nil
}

func (g *group) Parse(data mapstr.MapStr) (*metadata.Group, error) {

	err := mapstr.SetValueToStructByTags(&g.grp, data)
	return &g.grp, err
}
func (g *group) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(&g.grp)
}

func (g *group) GetAttributes() ([]AttributeInterface, error) {
	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(g.grp.ObjectID).
		Field(metadata.AttributeFieldPropertyGroup).Eq(g.grp.GroupID)

	rsp, err := g.clientSet.CoreService().Model().ReadModelAttr(context.Background(), g.kit.Header, g.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return nil, g.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), err: %s, rid: %s", g.grp.ObjectID, rsp.ErrMsg, g.kit.Rid)
		return nil, g.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	rstItems := make([]AttributeInterface, 0)
	for _, item := range rsp.Data.Info {

		attr := &attribute{
			attr:      item,
			kit:       g.kit,
			clientSet: g.clientSet,
		}

		rstItems = append(rstItems, attr)
	}

	return rstItems, nil
}

func (g *group) search(cond condition.Condition) ([]metadata.Group, error) {
	if g.bizID > 0 {
		cond.Field(common.BKAppIDField).Eq(g.bizID)
	}
	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroup(context.Background(), g.kit.Header, g.GetObjectID(), metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("failed to search the classification, err: %s, rid: %s", rsp.ErrMsg, g.kit.Rid)
		return nil, g.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

func (g *group) Save(data mapstr.MapStr) error {
	if nil != data {
		if _, err := g.grp.Parse(data); nil != err {
			return err
		}
	}

	if exists, err := g.IsExists(); nil != err {
		return err
	} else if !exists {
		return g.Create()
	}

	if nil != data {
		return g.Update(data)
	}

	return g.Update(g.grp.ToMapStr())
}

func (g *group) CreateAttribute() AttributeInterface {
	return &attribute{
		kit:       g.kit,
		clientSet: g.clientSet,
		attr: metadata.Attribute{
			OwnerID:  g.grp.OwnerID,
			ObjectID: g.grp.ObjectID,
		},
	}
}
