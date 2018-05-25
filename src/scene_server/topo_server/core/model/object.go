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

var _ Object = (*object)(nil)

type object struct {
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

func (cli *object) IsMainLine() bool {
	return false
}

func (cli *object) Parse(data frtypes.MapStr) error {

	err := frcommon.SetValueToStructByTags(cli, data)
	if nil != err {
		return err
	}

	// TODO 增加校验

	return err
}

func (cli *object) Save() error {
	dataMapStr := frcommon.SetValueToMapStrByTags(cli)

	_ = dataMapStr
	return nil
}

func (cli *object) CreateGroup() Group {
	return &group{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
}

func (cli *object) CreateAttribute() Attribute {
	return &attribute{
		OwnerID:  cli.OwnerID,
		ObjectID: cli.ObjectID,
	}
}

func (cli *object) SetClassification(class Classification) {
	cli.ObjCls = class.GetID()
}

func (cli *object) GetClassification() (Classification, error) {
	return nil, nil
}

func (cli *object) SetIcon(objectIcon string) {
	cli.ObjIcon = objectIcon
}

func (cli *object) GetIcon() string {
	return cli.ObjIcon
}

func (cli *object) SetID(objectID string) {
	cli.ObjectID = objectID
}

func (cli *object) GetID() string {
	return cli.ObjectID
}

func (cli *object) SetName(objectName string) {
	cli.ObjectName = objectName
}

func (cli *object) GetName() string {
	return cli.ObjectName
}

func (cli *object) SetIsPre(isPre bool) {
	cli.IsPre = isPre
}

func (cli *object) GetIsPre() bool {
	return cli.IsPre
}

func (cli *object) SetIsPaused(isPaused bool) {
	cli.IsPaused = isPaused
}

func (cli *object) GetIsPaused() bool {
	return cli.IsPaused
}

func (cli *object) SetPosition(position string) {
	cli.Position = position
}

func (cli *object) GetPosition() string {
	return cli.Position
}

func (cli *object) SetSupplierAccount(supplierAccount string) {
	cli.OwnerID = supplierAccount
}

func (cli *object) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *object) SetDescription(description string) {
	cli.Description = description
}

func (cli *object) GetDescription() string {
	return cli.Description
}

func (cli *object) SetCreator(creator string) {
	cli.Creator = creator
}

func (cli *object) GetCreator() string {
	return cli.Creator
}

func (cli *object) SetModifier(modifier string) {
	cli.Modifier = modifier
}

func (cli *object) GetModifier() string {
	return cli.Modifier
}
