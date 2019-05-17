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

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	ClassificationID   string `field:"bk_classification_id"`
	ClassificationName string `field:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon"`

	id int
}

func (cli *classification) ToMapStr() types.MapStr {
	return common.SetValueToMapStrByTags(cli)
}
func (cli *classification) search() ([]types.MapStr, error) {

	// construct the search condition
	cond := common.CreateCondition().Field(ClassificationID).Eq(cli.ClassificationID)

	// search all classifications by condition
	return client.GetClient().CCV3(client.Params{}).Classification().SearchClassifications(cond)
}
func (cli *classification) IsExists() (bool, error) {
	items, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(items), nil
}

func (cli *classification) Create() error {

	id, err := client.GetClient().CCV3(client.Params{}).Classification().CreateClassification(cli.ToMapStr())
	if nil != err {
		return err
	}

	cli.id = id
	return nil
}
func (cli *classification) Update() error {

	dataItems, err := cli.search()
	if nil != err {
		return err
	}

	// update the exists one
	for _, item := range dataItems {
		item.Set(ClassificationName, cli.ClassificationName)
		item.Set(ClassificationIcon, cli.ClassificationIcon)
		item.Set(ClassificationType, cli.ClassificationType)
		item.Remove(ClassificationID)

		id, err := item.Int("id")
		if nil != err {
			return err
		}

		cond := common.CreateCondition()
		cond.Field(ClassificationID).Eq(cli.ClassificationID).Field("id").Eq(id)
		if err = client.GetClient().CCV3(client.Params{}).Classification().UpdateClassification(item, cond); nil != err {
			return err
		}
	}
	return nil
}
func (cli *classification) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}

func (cli *classification) GetRecordID() int {
	return cli.id
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
func (cli *classification) GetName() string {
	return cli.ClassificationName
}

func (cli *classification) SetIcon(iconName string) {
	cli.ClassificationIcon = iconName
}

func (cli *classification) GetIcon() string {
	return cli.ClassificationIcon
}

func (cli *classification) CreateModel() Model {
	m := &model{}
	m.SetID("obj_" + common.UUID())
	m.SetClassification(cli.ClassificationID)
	m.SetIcon(objectIconDefault)
	return m
}

func (cli *classification) FindModelsLikeName(supplierAccount string, modelName string) (Iterator, error) {
	cond := common.CreateCondition().Field(ObjectName).Like(modelName)
	return newModelIterator(supplierAccount, cond)
}

func (cli *classification) FindModelsByCondition(supplierAccount string, cond common.Condition) (Iterator, error) {
	return newModelIterator(supplierAccount, cond)
}
