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
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
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

func (cli *group) ToMapStr() types.MapStr {
	return types.MapStr{
		GroupID:         cli.GroupID,
		GroupName:       cli.GroupName,
		GroupIndex:      0,
		ObjectID:        cli.ObjectID,
		SupplierAccount: cli.OwnerID,
		IsDefault:       cli.IsDefault,
		IsPre:           cli.IsPre,
	}
}

func (cli *group) Save() error {

	// construct the search condition
	cond := common.CreateCondition().Field(GroupID).Eq(cli.GroupID).Field(ObjectID).Eq(cli.ObjectID)

	// search all group by condition
	dataItems, err := client.GetClient().CCV3().Group().SearchGroups(cond)
	if nil != err {
		return err
	}

	// create a new object
	if 0 == len(dataItems) {
		if _, err = client.GetClient().CCV3().Group().CreateGroup(cli.ToMapStr()); nil != err {
			return err
		}
		return nil
	}

	// update the exists one
	for _, item := range dataItems {

		item.Set(GroupName, cli.GroupName)
		item.Set(GroupIndex, 0)
		item.Set(IsDefault, cli.IsDefault)

		cond := common.CreateCondition().Field(ObjectID).Eq(cli.ObjectID).Field(GroupID).Eq(cli.GroupID)
		if err = client.GetClient().CCV3().Group().UpdateGroup(item, cond); nil != err {
			return err
		}
	}

	// success
	return nil
}

func (cli *group) SetID(id string) {
	cli.GroupID = id
}

func (cli *group) GetID() string {

	return cli.GroupID
}

func (cli *group) SetName(name string) {
	cli.GroupName = name
}

func (cli *group) GetName() string {
	return cli.GroupName
}

func (cli *group) SetIndex(idx int) {
	cli.GroupIndex = idx
}

func (cli *group) GetIndex() int {
	return cli.GroupIndex
}

func (cli *group) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}

func (cli *group) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *group) SetDefault() {
	cli.IsDefault = true
}
func (cli *group) SetNonDefault() {
	cli.IsDefault = false
}

func (cli *group) GetDefault() bool {
	return cli.IsDefault
}

func (cli *group) CreateAttribute() Attribute {
	attr := &attribute{
		PropertyGroup: cli.GroupID,
	}
	return attr
}

func (cli *group) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	cond := common.CreateCondition().Field(PropertyName).Like(attributeName)
	return newAttributeIterator(cond)
}

func (cli *group) FindAttributesByCondition(cond common.Condition) (AttributeIterator, error) {
	return newAttributeIterator(cond)
}
