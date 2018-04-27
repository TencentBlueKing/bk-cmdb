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
	objCls      string `field:"bk_classification_id"`
	objIcon     string `field:"bk_obj_icon"`
	objectID    string `field:"bk_obj_id"`
	objectName  string `field:"bk_obj_name"`
	isPre       bool   `field:"ispre"`
	isPaused    bool   `field:"bk_ispaused"`
	position    string `field:"position"`
	ownerID     string `field:"bk_supplier_account"`
	description string `field:"description"`
	creator     string `field:"creator"`
	modifier    string `field:"modifier"`
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
	cli.objCls = classificationID
}

func (cli *model) GetClassification() string {
	return cli.objCls
}

func (cli *model) SetIcon(iconName string) {
	cli.objIcon = iconName
}

func (cli *model) GetIcon() string {
	return cli.objIcon
}

func (cli *model) SetID(id string) {
	cli.objectID = id
}

func (cli *model) GetID() string {
	return cli.objectID
}

func (cli *model) SetName(name string) {
	cli.objectName = name
}
func (cli *model) GetName() string {
	return cli.objectName
}

func (cli *model) SetPaused() {
	cli.isPaused = true
}

func (cli *model) SetNonPaused() {
	cli.isPaused = false
}

func (cli *model) Paused() bool {
	return cli.isPaused
}

func (cli *model) SetPosition(position string) {
	cli.position = position
}

func (cli *model) GetPosition() string {
	return cli.position
}

func (cli *model) SetSupplierAccount(ownerID string) {
	cli.ownerID = ownerID
}
func (cli *model) GetSupplierAccount() string {
	return cli.ownerID
}

func (cli *model) SetDescription(desc string) {
	cli.description = desc
}
func (cli *model) GetDescription() string {
	return cli.description
}
func (cli *model) SetCreator(creator string) {
	cli.creator = creator
}
func (cli *model) GetCreator() string {
	return cli.creator
}
func (cli *model) SetModifier(modifier string) {
	cli.modifier = modifier
}
func (cli *model) GetModifier() string {
	return cli.modifier
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
