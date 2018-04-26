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
	GroupID    string `json:"bk_group_id"`
	GroupName  string `json:"bk_group_name"`
	GroupIndex int    `json:"bk_group_index"`
	ObjectID   string `json:"bk_obj_id"`
	OwnerID    string `json:"bk_supplier_account"`
	IsDefault  bool   `json:"bk_isdefault"`
	IsPre      bool   `json:"ispre"`
}

func (cli *group) Save() error {
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

func (cli *group) Default() bool {
	return cli.IsDefault
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
