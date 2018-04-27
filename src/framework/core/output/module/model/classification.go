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
	"configcenter/src/framework/core/output/module/v3"
	"configcenter/src/framework/core/types"
)

var _ Classification = (*classification)(nil)

const (
	// ClassificationID the const definition
	ClassificationID = "bk_classification_id"
	// ClassificationName the const definition
	ClassificationName = "bk_classification_name"
	// ClassificationType the const definition
	ClassificationType = "bk_classification_type"
	// ClassificationIcon the const definition
	ClassificationIcon = "bk_classification_icon"
)

// classification the model classification definition
type classification struct {
	classificationID   string `field:"bk_classification_id"`
	classificationName string `field:"bk_classification_name"`
	classificationType string `field:"bk_classification_type"`
	classificationIcon string `field:"bk_classification_icon"`
}

func (cli *classification) ToMapStr() types.MapStr {
	return types.MapStr{
		ClassificationID:   cli.classificationID,
		ClassificationName: cli.classificationName,
		ClassificationType: cli.classificationType,
		ClassificationIcon: cli.classificationIcon,
	}
}

func (cli *classification) Save() error {

	// construct the search condition
	cond := common.CreateCondition()

	cond.Field(ClassificationID).Eq(cli.classificationID)

	// search all classifications by condition
	dataItems, err := v3.GetClient().SearchClassifications(cond)
	if nil != err {
		return err
	}

	// create a new classification
	if 0 == len(dataItems) {
		if _, err = v3.GetClient().CreateClassification(cli.ToMapStr()); nil != err {
			return err
		}
		return nil
	}

	// update the exists one
	for _, item := range dataItems {
		item.Set(ClassificationName, cli.classificationName)
		item.Set(ClassificationIcon, cli.classificationIcon)
		item.Set(ClassificationType, cli.classificationType)
		cond := common.CreateCondition()
		cond.Field(ClassificationID).Eq(cli.classificationID)
		if err = v3.GetClient().UpdateClassification(item, cond); nil != err {
			return err
		}
	}

	// success
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
	m.SetClassification(cli.classificationID)
	m.SetIcon(cli.classificationIcon)
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
