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

// check the interface
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

	id int
}

// ToMapStr TODO
func (cli *attribute) ToMapStr() types.MapStr {
	return common.SetValueToMapStrByTags(cli)
}

func (cli *attribute) search() ([]types.MapStr, error) {
	// construct the search condition
	cond := common.CreateCondition().Field(PropertyID).Eq(cli.PropertyID)
	cond.Field(ObjectID).Eq(cli.ObjectID)
	cond.Field(SupplierAccount).Eq(cli.OwnerID)

	// search all objects by condition
	dataItems, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Attribute().SearchObjectAttributes(cond)
	return dataItems, err
}

// IsExists TODO
func (cli *attribute) IsExists() (bool, error) {
	items, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(items), nil
}

// Create TODO
func (cli *attribute) Create() error {

	id, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Attribute().CreateObjectAttribute(cli.ToMapStr())
	if nil != err {
		return err
	}
	cli.id = id
	return nil
}

// Update TODO
func (cli *attribute) Update() error {

	dataItems, err := cli.search()

	if nil != err {
		return err
	}

	// update the exists one
	for _, item := range dataItems {

		item.Set(PropertyName, cli.PropertyName)
		item.Set(PropertyGroup, cli.PropertyGroup)
		item.Set(PropertyIndex, cli.PropertyIndex)
		item.Set(Unit, cli.Unit)
		item.Set(PlaceHolder, cli.Placeholder)
		item.Set(IsEditable, cli.IsEditable)
		item.Set(IsRequired, cli.IsRequired)
		item.Set(IsReadOnly, cli.IsReadOnly)
		item.Set(IsOnly, cli.IsOnly)
		item.Set(IsApi, cli.IsAPI)
		item.Set(PropertyType, cli.PropertyType)
		item.Set(Option, cli.Option)
		item.Set(Description, cli.Description)

		id, err := item.Int("id")
		if nil != err {
			return err
		}

		cond := common.CreateCondition()
		cond.Field(ObjectID).Eq(cli.ObjectID).Field(SupplierAccount).Eq(cli.OwnerID).Field(PropertyID).Eq(cli.PropertyID).Field("id").Eq(id)
		if err = client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Attribute().UpdateObjectAttribute(item, cond); nil != err {
			return err
		}
	}

	return nil
}

// Save TODO
func (cli *attribute) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()

}

// SetObjectID TODO
func (cli *attribute) SetObjectID(objectID string) {
	cli.ObjectID = objectID
}

// GetObjectID TODO
func (cli *attribute) GetObjectID() string {
	return cli.ObjectID
}

// SetID TODO
func (cli *attribute) SetID(id string) {
	cli.PropertyID = id
}

// GetRecordID TODO
func (cli *attribute) GetRecordID() int {
	return cli.id
}

// GetID TODO
func (cli *attribute) GetID() string {
	return cli.PropertyID
}

// SetName TODO
func (cli *attribute) SetName(name string) {
	cli.PropertyName = name
}

// GetName TODO
func (cli *attribute) GetName() string {
	return cli.PropertyName
}

// SetUnit TODO
func (cli *attribute) SetUnit(unit string) {
	cli.Unit = unit
}

// GetUnit TODO
func (cli *attribute) GetUnit() string {
	return cli.Unit
}

// SetPlaceholder TODO
func (cli *attribute) SetPlaceholder(placeHolder string) {
	cli.Placeholder = placeHolder
}

// GetPlaceholder TODO
func (cli *attribute) GetPlaceholder() string {
	return cli.Placeholder
}

// SetEditable TODO
func (cli *attribute) SetEditable() {
	cli.IsEditable = true
}

// GetEditable TODO
func (cli *attribute) GetEditable() bool {
	return cli.IsEditable
}

// SetNonEditable TODO
func (cli *attribute) SetNonEditable() {
	cli.IsEditable = false
}

// SetRequired TODO
func (cli *attribute) SetRequired() {
	cli.IsRequired = true
}

// SetNonRequired TODO
func (cli *attribute) SetNonRequired() {
	cli.IsRequired = false
}

// GetRequired TODO
func (cli *attribute) GetRequired() bool {
	return cli.IsRequired
}

// SetKey TODO
func (cli *attribute) SetKey(isKey bool) {
	cli.IsOnly = isKey
}

// GetKey TODO
func (cli *attribute) GetKey() bool {
	return cli.IsOnly
}

// SetOption TODO
func (cli *attribute) SetOption(option interface{}) {
	cli.Option = option
}

// GetOption TODO
func (cli *attribute) GetOption() interface{} {
	return cli.Option
}

// SetDescrition TODO
func (cli *attribute) SetDescrition(des string) {
	cli.Description = des
}

// GetDescription TODO
func (cli *attribute) GetDescription() string {
	return cli.Description
}

// SetType TODO
func (cli *attribute) SetType(dataType FieldDataType) {
	cli.PropertyType = string(dataType)
}

// GetType TODO
func (cli *attribute) GetType() FieldDataType {
	return FieldDataType(cli.PropertyType)
}
