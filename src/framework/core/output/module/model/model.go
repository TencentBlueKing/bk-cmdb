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

// Package model TODO
package model

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
)

var _ Model = (*model)(nil)

// model the metadata structure definition of the model
type model struct {
	ObjCls      string `field:"bk_classification_id"`
	ObjIcon     string `field:"bk_obj_icon"`
	ObjectID    string `field:"bk_obj_id"`
	ObjectName  string `field:"bk_obj_name"`
	IsPre       bool   `field:"ispre"`
	IsPaused    bool   `field:"bk_ispaused"`
	Position    string `field:"position"`
	OwnerID     string `field:"bk_supplier_account"`
	Description string `field:"description"`
	Creator     string `field:"creator"`
	Modifier    string `field:"modifier"`
	id          int64
}

// ToMapStr TODO
func (cli *model) ToMapStr() types.MapStr {
	return common.SetValueToMapStrByTags(cli)
}

// Attributes TODO
func (cli *model) Attributes() ([]Attribute, error) {

	cond := common.CreateCondition().Field(ObjectID).Like(cli.ObjectID).Field(SupplierAccount).Eq(cli.OwnerID)

	dataMap, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Attribute().SearchObjectAttributes(cond)

	if nil != err {
		return nil, err
	}

	attrs := make([]Attribute, 0)

	for _, item := range dataMap {
		tmpItem := &attribute{}

		if err := common.SetValueToStructByTags(tmpItem, item); nil != err {
			log.Errorf("failed to convert, %s", err.Error())
		}

		attrs = append(attrs, tmpItem)
	}

	return attrs, nil

}
func (cli *model) search() ([]types.MapStr, error) {

	cond := common.CreateCondition().Field(ObjectID).Eq(cli.ObjectID).Field(SupplierAccount).Eq(cli.OwnerID)

	// search all objects by condition
	return client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Model().SearchObjects(cond)
}

// IsExists TODO
func (cli *model) IsExists() (bool, error) {

	items, err := cli.search()
	if nil != err {
		return false, err
	}
	return 0 != len(items), nil
}

// Create TODO
func (cli *model) Create() error {

	id, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Model().CreateObject(cli.ToMapStr())
	if nil != err {
		return err
	}

	cli.id = id
	return nil
}

// Update TODO
func (cli *model) Update() error {

	dataItems, err := cli.search()
	if nil != err {
		return err
	}

	// update the exists one
	for _, item := range dataItems {

		item.Set(ObjectIcon, cli.ObjIcon)
		item.Set(ClassificationID, cli.ObjCls)
		item.Set(ObjectName, cli.ObjectName)
		item.Set(IsPre, cli.IsPre)
		item.Set(IsPaused, cli.IsPaused)
		item.Set(Position, cli.Position)
		item.Set(Description, cli.Description)
		item.Set(Modifier, cli.Modifier)

		item.Remove(ObjectID)

		id, err := item.Int("id")
		if nil != err {
			return err
		}

		cond := common.CreateCondition()
		cond.Field(ObjectID).Eq(cli.ObjectID).Field(SupplierAccount).Eq(cli.OwnerID).Field("id").Eq(id)
		if err = client.GetClient().CCV3(client.Params{SupplierAccount: cli.OwnerID}).Model().UpdateObject(item,
			cond); nil != err {
			return err
		}
	}
	return nil
}

// Save TODO
func (cli *model) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}
	return cli.Create()
}

// CreateAttribute TODO
func (cli *model) CreateAttribute() Attribute {
	attr := &attribute{
		ObjectID:      cli.ObjectID,
		OwnerID:       cli.OwnerID,
		Creator:       cli.Creator,
		PropertyGroup: "default",
	}
	return attr
}

// SetClassification TODO
func (cli *model) SetClassification(classificationID string) {
	cli.ObjCls = classificationID
}

// GetClassification TODO
func (cli *model) GetClassification() string {
	return cli.ObjCls
}

// SetIcon TODO
func (cli *model) SetIcon(iconName string) {
	cli.ObjIcon = iconName
}

// GetIcon TODO
func (cli *model) GetIcon() string {
	return cli.ObjIcon
}

// SetID TODO
func (cli *model) SetID(id string) {
	cli.ObjectID = id
}

// GetID TODO
func (cli *model) GetID() string {
	return cli.ObjectID
}

// SetName TODO
func (cli *model) SetName(name string) {
	cli.ObjectName = name
}

// GetName TODO
func (cli *model) GetName() string {
	return cli.ObjectName
}

// SetPaused TODO
func (cli *model) SetPaused() {
	cli.IsPaused = true
}

// SetNonPaused TODO
func (cli *model) SetNonPaused() {
	cli.IsPaused = false
}

// Paused TODO
func (cli *model) Paused() bool {
	return cli.IsPaused
}

// SetPosition TODO
func (cli *model) SetPosition(position string) {
	cli.Position = position
}

// GetPosition TODO
func (cli *model) GetPosition() string {
	return cli.Position
}

// SetSupplierAccount TODO
func (cli *model) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}

// GetSupplierAccount TODO
func (cli *model) GetSupplierAccount() string {
	return cli.OwnerID
}

// SetDescription TODO
func (cli *model) SetDescription(desc string) {
	cli.Description = desc
}

// GetDescription TODO
func (cli *model) GetDescription() string {
	return cli.Description
}

// SetCreator TODO
func (cli *model) SetCreator(creator string) {
	cli.Creator = creator
}

// GetCreator TODO
func (cli *model) GetCreator() string {
	return cli.Creator
}

// SetModifier TODO
func (cli *model) SetModifier(modifier string) {
	cli.Modifier = modifier
}

// GetModifier TODO
func (cli *model) GetModifier() string {
	return cli.Modifier
}

// CreateGroup TODO
func (cli *model) CreateGroup() Group {
	g := &group{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
	return g
}

// FindAttributesLikeName TODO
func (cli *model) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	cond := common.CreateCondition().Field(PropertyName).Like(attributeName)
	return newAttributeIterator(cli.OwnerID, cond)
}

// FindAttributesByCondition TODO
func (cli *model) FindAttributesByCondition(cond common.Condition) (AttributeIterator, error) {
	return newAttributeIterator(cli.OwnerID, cond)
}

// FindGroupsLikeName TODO
func (cli *model) FindGroupsLikeName(groupName string) (GroupIterator, error) {
	cond := common.CreateCondition().Field(GroupName).Like(groupName).Field(ObjectID).Eq(cli.GetID())
	return newGroupIterator(cli.OwnerID, cond)
}

// FindGroupsByCondition TODO
func (cli *model) FindGroupsByCondition(cond common.Condition) (GroupIterator, error) {
	return newGroupIterator(cli.OwnerID, cond.Field(ObjectID).Eq(cli.GetID()))
}
