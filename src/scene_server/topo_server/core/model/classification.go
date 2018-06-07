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

	"configcenter/src/scene_server/topo_server/core/types"
)

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	ClassificationID   string `field:"bk_classification_id"`
	ClassificationName string `field:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon"`

	params types.LogicParams
}

func (cli *classification) Parse(data frtypes.MapStr) error {

	err := frcommon.SetValueToStructByTags(cli, data)

	if nil != err {
		return err
	}

	// TODO 增加校验逻辑

	return err
}

func (cli *classification) ToMapStr() (frtypes.MapStr, error) {
	return nil, nil
}

func (cli *classification) GetObjects() ([]Object, error) {
	return nil, nil
}

func (cli *classification) Save() error {
	dataMapStr := frcommon.SetValueToMapStrByTags(cli)

	_ = dataMapStr
	return nil
}

func (cli *classification) SetID(classificationID string) {
	cli.ClassificationID = classificationID
}

func (cli *classification) GetID() string {
	return cli.ClassificationID
}

func (cli *classification) SetName(classificationName string) {
	cli.ClassificationName = classificationName
}

func (cli *classification) GetName() string {
	return cli.ClassificationName
}

func (cli *classification) SetType(classificationType string) {
	cli.ClassificationType = classificationType
}

func (cli *classification) GetType() string {
	return cli.ClassificationType
}

func (cli *classification) SetSupplierAccount(supplierAccount string) {
	// TODO: need to add owner field
}

func (cli *classification) GetSupplierAccount() string {
	// TODO: need to add owner field
	return ""
}

func (cli *classification) SetIcon(classificationIcon string) {
	cli.ClassificationIcon = classificationIcon
}

func (cli *classification) GetIcon() string {
	return cli.ClassificationIcon
}
