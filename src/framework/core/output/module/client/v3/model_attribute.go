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

type AttributeGetter interface {
	Attribute() AttributeInterface
}
type AttributeInterface interface {
	CreateObjectAttribute(data types.MapStr) (int, error)
	DeleteObjectAttribute(cond common.Condition) error
	UpdateObjectAttribute(data types.MapStr, cond common.Condition) error
	SearchObjectAttributes(cond common.Condition) ([]types.MapStr, error)
}

type Attribute struct {
	cli *Client
}

func newAttribute(cli *Client) *Attribute {
	return &Attribute{
		cli: cli,
	}
}

// CreateObjectAttribute create a new model object attribute
func (m *Attribute) CreateObjectAttribute(data types.MapStr) (int, error) {

	targetURL := fmt.Sprintf("%s/api/v3/create/objectattr", m.cli.GetAddress())

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

// DeleteObjectAttribute delete a object attribute by condition
func (m *Attribute) DeleteObjectAttribute(cond common.Condition) error {

	data := cond.ToMapStr()
	id, err := data.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/delete/objectattr/%d", m.cli.GetAddress(), id)

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

// UpdateObjectAttribute update a object attribute by condition
func (m *Attribute) UpdateObjectAttribute(data types.MapStr, cond common.Condition) error {

	dataCond := cond.ToMapStr()
	id, err := dataCond.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/update/objectattr/%d", m.cli.GetAddress(), id)

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

// SearchObjectAttributes search some object attributes by condition
func (m *Attribute) SearchObjectAttributes(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	targetURL := fmt.Sprintf("%s/api/v3/find/objectattr", m.cli.GetAddress())

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
