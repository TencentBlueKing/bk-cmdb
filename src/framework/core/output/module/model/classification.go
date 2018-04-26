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

import "configcenter/src/framework/common"

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	classificationID   string `field:"bk_classification_id"`
	classificationName string `field:"bk_classification_name"`
	classificationType string `field:"bk_classification_type"`
	classificationIcon string `field:"bk_classification_icon"`
}

func (cli *classification) Save() error {
	return nil
}

func (cli *classification) GetID() string {
	return cli.classificationID
}

func (cli *classification) SetID(id string) {
	cli.classificationID = id
}

func (cli *classification) SetName(name string) {
	cli.classificationName = name
}

func (cli *classification) SetIcon(iconName string) {
	cli.classificationIcon = iconName
}

func (cli *classification) CreateModel() Model {
	m := &model{}
	m.ObjCls = cli.classificationID
	m.ObjIcon = cli.classificationIcon
	return m
}

func (cli *classification) FindModelsLikeName(modelName string) (Iterator, error) {
	// TODO: 按照名字正则查找，返回已经包含一定数量的Model数据的迭代器。
	return nil, nil
}

func (cli *classification) FindModelsByCondition(condition *common.Condition) (Iterator, error) {
	// TODO: 按照条件查找，返回一定数量的Model
	return nil, nil
}
