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

// check the interface
var _ Attribute = (*attribute)(nil)

// attribute the metadata structure definition of the model attribute
type attribute struct {
	ownerID       string `field:"bk_supplier_account"`
	objectID      string `field:"bk_obj_id"`
	propertyID    string `field:"bk_property_id"`
	propertyName  string `field:"bk_property_name"`
	propertyGroup string `field:"bk_property_group"`
	propertyIndex int    `field:"bk_property_index"`
	unit          string `field:"unit"`
	placeholder   string `field:"placeholder"`
	isEditable    bool   `field:"editable"`
	isPre         bool   `field:"ispre"`
	isRequired    bool   `field:"isrequired"`
	isReadOnly    bool   `field:"isreadonly"`
	isOnly        bool   `field:"isonly"`
	isSystem      bool   `field:"bk_issystem"`
	isAPI         bool   `field:"bk_isapi"`
	propertyType  string `field:"bk_property_type"`
	option        string `field:"option"`
	description   string `field:"description"`
	creator       string `field:"creator"`
}

func (cli *attribute) Save() error {
	return nil
}

func (cli *attribute) SetID(id string) {
	cli.propertyID = id
}

func (cli *attribute) SetName(name string) {
	cli.propertyName = name
}

func (cli *attribute) SetUnit(unit string) {
	cli.unit = unit
}

func (cli *attribute) SetPlaceholer(placeHoler string) {
	cli.placeholder = placeHoler
}

func (cli *attribute) SetEditable() {
	cli.isEditable = true
}

func (cli *attribute) SetNonEditable() {
	cli.isEditable = false
}

func (cli *attribute) Editable() bool {
	return cli.isEditable
}

func (cli *attribute) SetRequired() {
	cli.isRequired = true
}

func (cli *attribute) SetNonRequired() {
	cli.isRequired = false
}

func (cli *attribute) Required() bool {
	return cli.isRequired
}

func (cli *attribute) SetKey(isKey bool) {
	cli.isOnly = isKey
}

func (cli *attribute) SetOption(option string) {
	cli.option = option
}

func (cli *attribute) SetDescrition(des string) {
	cli.description = des
}
