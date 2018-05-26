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
	frcommon "configcenter/src/framework/common"
	frtypes "configcenter/src/framework/core/types"
)

var _ Group = (*group)(nil)

type group struct {
	GroupID    string `field:"bk_group_id"`
	GroupName  string `field:"bk_group_name"`
	GroupIndex int    `field:"bk_group_index"`
	ObjectID   string `field:"bk_obj_id"`
	OwnerID    string `field:"bk_supplier_account"`
	IsDefault  bool   `field:"bk_isdefault"`
	IsPre      bool   `field:"ispre"`
}

func (cli *group) Parse(data frtypes.MapStr) error {

	err := frcommon.SetValueToStructByTags(cli, data)

	if nil != err {
		return err
	}

	// TODO 实现校验逻辑

	return err
}

func (cli *group) GetAttributes() ([]Attribute, error) {
	return nil, nil
}

func (cli *group) Save() error {
	dataMapStr := frcommon.SetValueToMapStrByTags(cli)

	_ = dataMapStr
	return nil
}

func (cli *group) CreateAttribute() Attribute {
	return &attribute{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
}

func (cli *group) SetID(groupID string) {
	cli.GroupID = groupID
}

func (cli *group) GetID() string {
	return cli.GroupID
}

func (cli *group) SetName(groupName string) {
	cli.GroupName = groupName
}

func (cli *group) GetName() string {
	return cli.GroupName
}

func (cli *group) SetIndex(groupIndex int64) {
	cli.GroupIndex = int(groupIndex)
}

func (cli *group) GetIndex() int64 {
	return int64(cli.GroupIndex)
}

func (cli *group) SetSupplierAccount(supplierAccount string) {
	cli.OwnerID = supplierAccount
}

func (cli *group) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *group) SetDefault(isDefault bool) {
	cli.IsDefault = isDefault
}

func (cli *group) GetDefault() bool {
	return cli.IsDefault
}

func (cli *group) SetIsPre(isPre bool) {
	cli.IsPre = isPre
}

func (cli *group) GetIsPre() bool {
	return cli.IsPre
}
