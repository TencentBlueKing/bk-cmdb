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

// classification the model classification definition
type classification struct {
	ClassificationID   string `field:"bk_classification_id"`
	ClassificationName string `field:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon"`
}

func (cli *classification) ToMapStr() types.MapStr {
	return types.MapStr{
		ClassificationID:   cli.ClassificationID,
		ClassificationName: cli.ClassificationName,
		ClassificationType: cli.ClassificationType,
		ClassificationIcon: cli.ClassificationIcon,
	}
}

func (cli *classification) Save() error {

	// construct the search condition
	cond := common.CreateCondition()

	cond.Field(ClassificationID).Eq(cli.ClassificationID)

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
		item.Set(ClassificationName, cli.ClassificationName)
		item.Set(ClassificationIcon, cli.ClassificationIcon)
		item.Set(ClassificationType, cli.ClassificationType)
		cond := common.CreateCondition()
		cond.Field(ClassificationID).Eq(cli.ClassificationID)
		if err = v3.GetClient().UpdateClassification(item, cond); nil != err {
			return err
		}
	}

	// success
	return nil
}

func (cli *classification) GetID() string {
	return cli.ClassificationID
}

func (cli *classification) SetID(id string) {
	cli.ClassificationID = id
}

func (cli *classification) SetName(name string) {
	cli.ClassificationName = name
}

func (cli *classification) SetIcon(iconName string) {
	cli.ClassificationIcon = iconName
}

func (cli *classification) CreateModel() Model {
	m := &model{}
	m.SetClassification(cli.ClassificationID)
	m.SetIcon(cli.ClassificationIcon)
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
