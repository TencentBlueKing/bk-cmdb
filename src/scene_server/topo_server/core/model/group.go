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
	"configcenter/src/apimachinery"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

var _ Group = (*group)(nil)

type group struct {
	grp       metadata.Group
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *group) Create() error {
	return nil
}

func (cli *group) Update() error {
	return nil
}

func (cli *group) Delete() error {
	return nil
}

func (cli *group) IsExists() (bool, error) {
	return false, nil
}

func (cli *group) Parse(data frtypes.MapStr) (*metadata.Group, error) {

	err := metadata.SetValueToStructByTags(cli, data)

	if nil != err {
		return nil, err
	}

	// TODO 实现校验逻辑

	return nil, err
}
func (cli *group) ToMapStr() (frtypes.MapStr, error) {
	return nil, nil
}

func (cli *group) GetAttributes() ([]Attribute, error) {
	return nil, nil
}

func (cli *group) Save() error {
	dataMapStr := metadata.SetValueToMapStrByTags(cli)

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
	cli.grp.GroupIndex = int(groupIndex)
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
