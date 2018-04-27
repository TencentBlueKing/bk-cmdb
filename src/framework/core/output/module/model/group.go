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
)

var _ Group = (*group)(nil)

type group struct {
	groupID    string `field:"bk_group_id"`
	groupName  string `field:"bk_group_name"`
	groupIndex int    `field:"bk_group_index"`
	objectID   string `field:"bk_obj_id"`
	ownerID    string `field:"bk_supplier_account"`
	isDefault  bool   `field:"bk_isdefault"`
	isPre      bool   `field:"ispre"`
}

func (cli *group) Save() error {
	return nil
}

func (cli *group) SetID(id string) {
	cli.groupID = id
}

func (cli *group) GetID() string {

	return cli.groupID
}

func (cli *group) SetName(name string) {
	cli.groupName = name
}

func (cli *group) SetIndex(idx int) {
	cli.groupIndex = idx
}

func (cli *group) GetIndex() int {
	return cli.groupIndex
}

func (cli *group) SetSupplierAccount(ownerID string) {
	cli.ownerID = ownerID
}

func (cli *group) GetSupplierAccount() string {
	return cli.ownerID
}

func (cli *group) SetDefault() {
	cli.isDefault = true
}
func (cli *group) SetNonDefault() {
	cli.isDefault = false
}

func (cli *group) Default() bool {
	return cli.isDefault
}

func (cli *group) CreateAttribute() Attribute {
	attr := &attribute{}
	return attr
}

func (cli *group) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	// TODO: 按照名字正则查找
	return nil, nil
}

func (cli *group) FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error) {
	// TODO: 按照条件查找
	return nil, nil
}
