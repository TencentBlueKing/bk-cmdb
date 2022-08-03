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
	"configcenter/src/framework/core/log"
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
	id         int
}

// ToMapStr TODO
func (cli *group) ToMapStr() types.MapStr {
	return common.SetValueToMapStrByTags(cli)
}

func (cli *group) search() ([]types.MapStr, error) {

	// construct the search condition
	cond := common.CreateCondition().Field(ObjectID).Eq(cli.ObjectID)

	// search all group by condition
	return client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Group().SearchGroups(cond)

}

// IsExists TODO
func (cli *group) IsExists() (bool, error) {

	items, err := cli.search()
	return 0 != len(items), err
}

// Create TODO
func (cli *group) Create() error {

	id, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Group().CreateGroup(cli.ToMapStr())
	cli.id = id
	return err
}

// Update TODO
func (cli *group) Update() error {

	dataItems, err := cli.search()
	if nil != err {
		return err
	}

	var updateitem types.MapStr
	lastIndex := 1
	var updateBy string
	for _, item := range dataItems {
		index, err := item.Int(GroupIndex)
		if err != nil {
			log.Errorf("get bk_group_index error %v ", err)
		}
		if index > lastIndex {
			lastIndex = index + 1
		}
		if cli.GetName() == item.String(GroupName) {
			updateBy = GroupName
			updateitem = item
		}
		if cli.GetID() == item.String(GroupID) {
			updateBy = GroupID
			updateitem = item
		}
	}
	if len(cli.GetID()) <= 0 {
		cli.SetID(common.UUID())
	}

	if cli.GetIndex() <= 0 {
		cli.SetIndex(lastIndex)
	}

	if nil != updateitem {
		// update exists one
		updateitem.Set(GroupName, cli.GetName())
		updateitem.Set(GroupIndex, cli.GetIndex())
		updateitem.Set(GroupID, cli.GetID())
		cond := common.CreateCondition().Field(ObjectID).Eq(cli.ObjectID)
		if updateBy == GroupID {
			cond = cond.Field(GroupID).Eq(cli.GroupID)
		} else {
			cond = cond.Field(GroupName).Eq(cli.GroupName)
		}
		return client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Group().UpdateGroup(updateitem, cond)
	}
	return nil
}

// Save TODO
func (cli *group) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}

// GetRecordID TODO
func (cli *group) GetRecordID() int {
	return cli.id
}

// SetID TODO
func (cli *group) SetID(id string) {
	cli.GroupID = id
}

// GetID TODO
func (cli *group) GetID() string {

	return cli.GroupID
}

// SetName TODO
func (cli *group) SetName(name string) {
	cli.GroupName = name
}

// GetName TODO
func (cli *group) GetName() string {
	return cli.GroupName
}

// SetIndex TODO
func (cli *group) SetIndex(idx int) {
	cli.GroupIndex = idx
}

// GetIndex TODO
func (cli *group) GetIndex() int {
	return cli.GroupIndex
}

// SetSupplierAccount TODO
func (cli *group) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}

// GetSupplierAccount TODO
func (cli *group) GetSupplierAccount() string {
	return cli.OwnerID
}

// SetDefault TODO
func (cli *group) SetDefault() {
	cli.IsDefault = true
}

// SetNonDefault TODO
func (cli *group) SetNonDefault() {
	cli.IsDefault = false
}

// GetDefault TODO
func (cli *group) GetDefault() bool {
	return cli.IsDefault
}

// CreateAttribute TODO
func (cli *group) CreateAttribute() Attribute {
	attr := &attribute{
		PropertyGroup: cli.GroupID,
	}
	return attr
}

// FindAttributesLikeName TODO
func (cli *group) FindAttributesLikeName(supplierAccount string, attributeName string) (AttributeIterator, error) {
	cond := common.CreateCondition().Field(PropertyName).Like(attributeName)
	return newAttributeIterator(supplierAccount, cond)
}

// FindAttributesByCondition TODO
func (cli *group) FindAttributesByCondition(supplierAccount string, cond common.Condition) (AttributeIterator, error) {
	return newAttributeIterator(supplierAccount, cond)
}
