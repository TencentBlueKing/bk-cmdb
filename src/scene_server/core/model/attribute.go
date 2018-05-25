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

var _ Attribute = (*attribute)(nil)

// attribute the metadata structure definition of the model attribute
type attribute struct {
	OwnerID       string      `field:"bk_supplier_account"`
	ObjectID      string      `field:"bk_obj_id"`
	PropertyID    string      `field:"bk_property_id"`
	PropertyName  string      `field:"bk_property_name"`
	PropertyGroup string      `field:"bk_property_group"`
	PropertyIndex int         `field:"bk_property_index"`
	Unit          string      `field:"unit"`
	Placeholder   string      `field:"placeholder"`
	IsEditable    bool        `field:"editable"`
	IsPre         bool        `field:"ispre"`
	IsRequired    bool        `field:"isrequired"`
	IsReadOnly    bool        `field:"isreadonly"`
	IsOnly        bool        `field:"isonly"`
	IsSystem      bool        `field:"bk_issystem"`
	IsAPI         bool        `field:"bk_isapi"`
	PropertyType  string      `field:"bk_property_type"`
	Option        interface{} `field:"option"`
	Description   string      `field:"description"`
	Creator       string      `field:"creator"`
}

func (cli *attribute) Parse(data frtypes.MapStr) error {
	err := frcommon.SetValueToStructByTags(cli, data)

	if nil != err {
		return err
	}

	// TODO 实现数据校验

	return err
}

func (cli *attribute) Save() error {
	return nil
}

func (cli *attribute) SetSupplierAccount(supplierAccount string) {

	cli.OwnerID = supplierAccount

}

func (cli *attribute) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *attribute) SetObjectID(objectID string) {
	cli.ObjectID = objectID
}

func (cli *attribute) GetObjectID() string {
	return cli.ObjectID
}

func (cli *attribute) SetID(attributeID string) {
	cli.PropertyID = attributeID
}

func (cli *attribute) GetID() string {
	return cli.PropertyID
}

func (cli *attribute) SetName(attributeName string) {
	cli.PropertyName = attributeName
}

func (cli *attribute) GetName() string {
	return cli.PropertyName
}

func (cli *attribute) SetGroup(grp Group) {
	cli.PropertyGroup = grp.GetID()
}

func (cli *attribute) GetGroup() (Group, error) {
	return nil, nil
}

func (cli *attribute) SetGroupIndex(attGroupIndex int64) {
	cli.PropertyIndex = int(attGroupIndex)
}

func (cli *attribute) GetGroupIndex() int64 {
	return int64(cli.PropertyIndex)
}

func (cli *attribute) SetUnint(unit string) {
	cli.Unit = unit
}

func (cli *attribute) GetUnint() string {
	return cli.Unit
}

func (cli *attribute) SetPlaceholder(placeHolder string) {
	cli.Placeholder = placeHolder
}

func (cli *attribute) GetPlaceholder() string {
	return cli.Placeholder
}

func (cli *attribute) SetIsEditable(isEditable bool) {
	cli.IsEditable = isEditable
}

func (cli *attribute) GetIsEditable() bool {
	return cli.IsEditable
}

func (cli *attribute) SetIsPre(isPre bool) {
	cli.IsPre = isPre
}

func (cli *attribute) GetIsPre() bool {
	return cli.IsPre
}

func (cli *attribute) SetIsReadOnly(isReadOnly bool) {
	cli.IsReadOnly = isReadOnly
}

func (cli *attribute) GetIsReadOnly() bool {
	return cli.IsReadOnly
}

func (cli *attribute) SetIsOnly(isOnly bool) {
	cli.IsOnly = isOnly
}

func (cli *attribute) GetIsOnly() bool {
	return cli.IsOnly
}

func (cli *attribute) SetIsSystem(isSystem bool) {
	cli.IsSystem = isSystem
}

func (cli *attribute) GetIsSystem() bool {
	return cli.IsSystem
}

func (cli *attribute) SetIsAPI(isAPI bool) {
	cli.IsAPI = isAPI
}

func (cli *attribute) GetIsAPI() bool {
	return cli.IsAPI
}

func (cli *attribute) SetType(attributeType string) {
	cli.PropertyType = attributeType
}

func (cli *attribute) GetType() string {
	return cli.PropertyType
}

func (cli *attribute) SetOption(attributeOption interface{}) {
	cli.Option = attributeOption
}

func (cli *attribute) GetOption() interface{} {
	return cli.Option
}

func (cli *attribute) SetDescription(attributeDescription string) {
	cli.Description = attributeDescription
}

func (cli *attribute) GetDescription() string {
	return cli.Description
}

func (cli *attribute) SetCreator(attributeCreator string) {
	cli.Creator = attributeCreator
}

func (cli *attribute) GetCreator() string {
	return cli.Creator
}

// Save update the data in the database
func (cli *attribute) Save() error {

	return nil
}
