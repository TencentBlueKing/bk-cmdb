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

package v3

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type ClassificationGetter interface {
	Classification() ClassificationInterface
}
type ClassificationInterface interface {
	CreateClassification(data types.MapStr) (int, error)
	DeleteClassification(cond common.Condition) error
	SearchClassifications(cond common.Condition) ([]types.MapStr, error)
	SearchClassificationWithObjects(cond common.Condition) ([]types.MapStr, error)
	UpdateClassification(data types.MapStr, cond common.Condition) error
}

type Classification struct {
	cli *Client
}

func newClassification(cli *Client) *Classification {
	return &Classification{
		cli: cli,
	}
}

// CreateClassification create a new classification
func (m *Classification) CreateClassification(data types.MapStr) (int, error) {

	targetURL := fmt.Sprintf("%s/api/v3/create/objectclassification", m.cli.GetAddress())

	rst, err := m.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return 0, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return 0, errors.New(gs.Get("bk_error_msg").String())
	}

	// parse id
	id := gs.Get("data.id").Int()

	return int(id), nil
}

// DeleteClassification delete some classification by condition
func (m *Classification) DeleteClassification(cond common.Condition) error {

	data := cond.ToMapStr()
	id, err := data.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/delete/objectclassification/%d", m.cli.GetAddress(), id)

	rst, err := m.cli.httpCli.DELETE(targetURL, nil, nil)
	if nil != err {
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return errors.New(gs.Get("bk_error_msg").String())
	}

	return nil
}

// SearchClassifications search some classification by condition
func (m *Classification) SearchClassifications(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	targetURL := fmt.Sprintf("%s/api/v3/find/objectclassification", m.cli.GetAddress())

	rst, err := m.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return nil, errors.New(gs.Get("message").String())
	}

	dataStr := gs.Get("data").String()
	if 0 == len(dataStr) {
		return nil, errors.New("data is empty")
	}

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)

	return resultMap, err
}

// SearchClassificationWithObjects search some classification with objects
func (m *Classification) SearchClassificationWithObjects(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	targetURL := fmt.Sprintf("%s/api/v3/find/classificationobject", m.cli.GetAddress())

	rst, err := m.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return nil, errors.New(gs.Get("bk_error_msg").String())
	}

	dataStr := gs.Get("data").String()
	if 0 == len(dataStr) {
		return nil, errors.New("data is empty")
	}

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)
	return resultMap, err
}

// UpdateClassification update the classification by condition
func (m *Classification) UpdateClassification(data types.MapStr, cond common.Condition) error {

	dataCond := cond.ToMapStr()
	id, err := dataCond.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/update/objectclassification/%d", m.cli.GetAddress(), id)

	rst, err := m.cli.httpCli.PUT(targetURL, nil, data.ToJSON())
	if nil != err {
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return errors.New(gs.Get("bk_error_msg").String())
	}

	return nil
}
