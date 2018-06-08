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

var _ Attribute = (*attribute)(nil)

// attribute the metadata structure definition of the model attribute
type attribute struct {
	attr      metadata.Attribute
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *attribute) Parse(data frtypes.MapStr) (*metadata.Attribute, error) {
	return cli.attr.Parse(data)
}

func (cli *attribute) ToMapStr() (frtypes.MapStr, error) {
	return nil, nil
}

func (cli *attribute) Create() error {
	return nil
}

func (cli *attribute) Update() error {
	return nil
}

func (cli *attribute) IsExists() (bool, error) {
	return false, nil
}

func (cli *attribute) Delete() error {
	return nil
}

func (cli *attribute) Save() error {

	return nil
}

func (cli *attribute) SetSupplierAccount(supplierAccount string) {

	cli.attr.OwnerID = supplierAccount

}

func (cli *attribute) GetSupplierAccount() string {
	return cli.attr.OwnerID
}

func (cli *attribute) SetObjectID(objectID string) {
	cli.attr.ObjectID = objectID
}

func (cli *attribute) GetObjectID() string {
	return cli.attr.ObjectID
}

func (cli *attribute) SetID(attributeID string) {
	cli.attr.PropertyID = attributeID
}

func (cli *attribute) GetID() string {
	return cli.attr.PropertyID
}

func (cli *attribute) SetName(attributeName string) {
	cli.attr.PropertyName = attributeName
}

func (cli *attribute) GetName() string {
	return cli.attr.PropertyName
}

func (cli *attribute) SetGroup(grp Group) {
	cli.attr.PropertyGroup = grp.GetID()
}

func (cli *attribute) GetGroup() (Group, error) {
	return nil, nil
}

func (cli *attribute) SetGroupIndex(attGroupIndex int64) {
	cli.attr.PropertyIndex = int(attGroupIndex)
}

func (cli *attribute) GetGroupIndex() int64 {
	return int64(cli.attr.PropertyIndex)
}

func (cli *attribute) SetUnint(unit string) {
	cli.attr.Unit = unit
}

func (cli *attribute) GetUnint() string {
	return cli.attr.Unit
}

func (cli *attribute) SetPlaceholder(placeHolder string) {
	cli.attr.Placeholder = placeHolder
}

func (cli *attribute) GetPlaceholder() string {
	return cli.attr.Placeholder
}

func (cli *attribute) SetIsEditable(isEditable bool) {
	cli.attr.IsEditable = isEditable
}

func (cli *attribute) GetIsEditable() bool {
	return cli.attr.IsEditable
}

func (cli *attribute) SetIsPre(isPre bool) {
	cli.attr.IsPre = isPre
}

func (cli *attribute) GetIsPre() bool {
	return cli.attr.IsPre
}

func (cli *attribute) SetIsReadOnly(isReadOnly bool) {
	cli.attr.IsReadOnly = isReadOnly
}

func (cli *attribute) GetIsReadOnly() bool {
	return cli.attr.IsReadOnly
}

func (cli *attribute) SetIsOnly(isOnly bool) {
	cli.attr.IsOnly = isOnly
}

func (cli *attribute) GetIsOnly() bool {
	return cli.attr.IsOnly
}

func (cli *attribute) SetIsSystem(isSystem bool) {
	cli.attr.IsSystem = isSystem
}

func (cli *attribute) GetIsSystem() bool {
	return cli.attr.IsSystem
}

func (cli *attribute) SetIsAPI(isAPI bool) {
	cli.attr.IsAPI = isAPI
}

func (cli *attribute) GetIsAPI() bool {
	return cli.attr.IsAPI
}

func (cli *attribute) SetType(attributeType string) {
	cli.attr.PropertyType = attributeType
}

func (cli *attribute) GetType() string {
	return cli.attr.PropertyType
}

func (cli *attribute) SetOption(attributeOption interface{}) {
	cli.attr.Option = attributeOption
}

func (cli *attribute) GetOption() interface{} {
	return cli.attr.Option
}

func (cli *attribute) SetDescription(attributeDescription string) {
	cli.attr.Description = attributeDescription
}

func (cli *attribute) GetDescription() string {
	return cli.attr.Description
}

func (cli *attribute) SetCreator(attributeCreator string) {
	cli.attr.Creator = attributeCreator
}

func (cli *attribute) GetCreator() string {
	return cli.attr.Creator
}
