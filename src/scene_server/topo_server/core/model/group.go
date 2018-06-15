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

var _ Group = (*group)(nil)

type group struct {
	grp       metadata.Group
	isNew     bool
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *group) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.grp)
}

func (cli *group) Create() error {

	rsp, err := cli.clientSet.ObjectController().Meta().CreatePropertyGroup(context.Background(), cli.params.Header.ToHeader(), &cli.grp)

	if nil != err {
		blog.Errorf("[model-grp] failed to request object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp] failed to create the group(%s), error info is is %s", cli.grp.GroupID, rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	cli.grp.ID = rsp.Data.ID

	return nil
}

func (cli *group) Update() error {

	cond := &metadata.UpdateGroupCondition{}
	cond.Condition.ID = cli.grp.GroupID
	cond.Data.Index = cli.grp.GroupIndex
	cond.Data.Name = cli.grp.GroupName

	rsp, err := cli.clientSet.ObjectController().Meta().UpdatePropertyGroup(context.Background(), cli.params.Header.ToHeader(), cond)
	if nil != err {
		blog.Errorf("[model-grp]failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp]failed to update the group(%s), error info is %s", cli.grp.GroupID, err.Error())
		return cli.params.Err.Error(common.CCErrTopoObjectAttributeUpdateFailed)
	}
	return nil
}

func (cli *group) Delete() error {

	rsp, err := cli.clientSet.ObjectController().Meta().DeletePropertyGroup(context.Background(), cli.grp.GroupID, cli.params.Header.ToHeader())
	if nil != err {
		blog.Error("[model-grp]failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp]failed to delte the group(%s), error info is %s", cli.grp.GroupID, rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoObjectGroupDeleteFailed)
	}

	return nil
}

func (cli *group) IsExists() (bool, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(cli.grp.GroupID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectGroup(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request object controller ,error info is %s", err.Error())
		return false, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to query group, error info is  %s", rsp.ErrMsg)
		return false, cli.params.Err.Error(common.CCErrTopoObjectGroupSelectFailed)
	}

	return 0 != len(rsp.Data), nil
}

func (cli *group) Parse(data frtypes.MapStr) (*metadata.Group, error) {

	err := metadata.SetValueToStructByTags(&cli.grp, data)
	return &cli.grp, err
}
func (cli *group) ToMapStr() (frtypes.MapStr, error) {

	rst := metadata.SetValueToMapStrByTags(&cli.grp)
	return rst, nil
}

func (cli *group) GetAttributes() ([]Attribute, error) {
	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(cli.grp.ObjectID).
		Field(metadata.AttributeFieldPropertyGroup).Eq(cli.grp.GroupID).
		Field(metadata.AttributeFieldSupplierAccount).Eq(cli.params.Header.OwnerID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.grp.ObjectID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Attribute, 0)
	for _, item := range rsp.Data {

		attr := &attribute{
			attr:      item,
			params:    cli.params,
			clientSet: cli.clientSet,
		}

		rstItems = append(rstItems, attr)
	}

	return rstItems, nil
}

func (cli *group) Save() error {
	dataMapStr := metadata.SetValueToMapStrByTags(&cli.grp)

	_ = dataMapStr
	return nil
}

func (cli *group) CreateAttribute() Attribute {
	return &attribute{
		attr: metadata.Attribute{
			OwnerID:  cli.grp.OwnerID,
			ObjectID: cli.grp.ObjectID,
		},
	}
}

func (cli *group) SetID(groupID string) {
	cli.grp.GroupID = groupID
}

func (cli *group) GetID() string {
	return cli.grp.GroupID
}

func (cli *group) SetName(groupName string) {
	cli.grp.GroupName = groupName
}

func (cli *group) GetName() string {
	return cli.grp.GroupName
}

func (cli *group) SetIndex(groupIndex int64) {
	cli.grp.GroupIndex = groupIndex
}

func (cli *group) GetIndex() int64 {
	return int64(cli.grp.GroupIndex)
}

func (cli *group) SetSupplierAccount(supplierAccount string) {
	cli.grp.OwnerID = supplierAccount
}

func (cli *group) GetSupplierAccount() string {
	return cli.grp.OwnerID
}

func (cli *group) SetDefault(isDefault bool) {
	cli.grp.IsDefault = isDefault
}

func (cli *group) GetDefault() bool {
	return cli.grp.IsDefault
}

func (cli *group) SetIsPre(isPre bool) {
	cli.grp.IsPre = isPre
}

func (cli *group) GetIsPre() bool {
	return cli.grp.IsPre
}
