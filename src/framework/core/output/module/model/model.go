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
	"fmt"
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
}

func (cli *model) Save() error {
	fmt.Println("test model")
	return nil
}

func (cli *model) CreateAttribute() Attribute {
	attr := &attribute{}
	return attr
}

func (cli *model) SetClassification(classificationID string) {
	cli.ObjCls = classificationID
}

func (cli *model) GetClassification() string {
	return cli.ObjCls
}

func (cli *model) SetIcon(iconName string) {
	cli.ObjIcon = iconName
}

func (cli *model) GetIcon() string {
	return cli.ObjIcon
}

func (cli *model) SetID(id string) {
	cli.ObjectID = id
}

func (cli *model) GetID() string {
	return cli.ObjectID
}

func (cli *model) SetName(name string) {
	cli.ObjectName = name
}
func (cli *model) GetName() string {
	return cli.ObjectName
}

func (cli *model) SetPaused() {
	cli.IsPaused = true
}

func (cli *model) SetNonPaused() {
	cli.IsPaused = false
}

func (cli *model) Paused() bool {
	return cli.IsPaused
}

func (cli *model) SetPosition(position string) {
	cli.Position = position
}

func (cli *model) GetPosition() string {
	return cli.Position
}

func (cli *model) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}
func (cli *model) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *model) SetDescription(desc string) {
	cli.Description = desc
}
func (cli *model) GetDescription() string {
	return cli.Description
}
func (cli *model) SetCreator(creator string) {
	cli.Creator = creator
}
func (cli *model) GetCreator() string {
	return cli.Creator
}
func (cli *model) SetModifier(modifier string) {
	cli.Modifier = modifier
}
func (cli *model) GetModifier() string {
	return cli.Modifier
}
func (cli *model) CreateGroup() Group {
	g := &group{}
	return g
}

func (cli *model) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
func (cli *model) FindGroupsLikeName(groupName string) (GroupIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindGroupsByCondition(condition *common.Condition) (GroupIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
