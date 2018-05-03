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

// CreateModule create a new module object
func (cli *Client) CreateModule(data types.MapStr) (int, error) {

	appID := data.String(BusinessID)
	if 0 == len(appID) {
		return 0, errors.New("the business id is not set")
	}

	data.Remove(BusinessID)

	setID := data.String(SetID)
	if 0 == len(appID) {
		return 0, errors.New("the set id is not set")
	}

	data.Remove(SetID)

	targetURL := fmt.Sprintf("%s/api/v3/module/%s/%s", cli.GetAddress(), appID, setID)

	rst, err := cli.httpCli.POST(targetURL, nil, data.ToJSON())
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

// DeleteModule delete a module by condition
func (cli *Client) DeleteModule(cond common.Condition) error {

	data := cond.ToMapStr()
	moduleID := data.String(ModuleID)
	if 0 == len(moduleID) {
		return errors.New("the module id is not set")
	}

	data.Remove(ModuleID)

	appID := data.String(BusinessID)
	if 0 == len(appID) {
		return errors.New("the business id is not set")
	}

	data.Remove(BusinessID)

	setID := data.String(SetID)
	if 0 == len(appID) {
		return errors.New("the set id is not set")
	}

	data.Remove(SetID)

	targetURL := fmt.Sprintf("%s/api/v3/module/%s/%s/%s", cli.GetAddress(), appID, setID, moduleID)

	rst, err := cli.httpCli.DELETE(targetURL, nil, nil)
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

// UpdateModule update a module by condition
func (cli *Client) UpdateModule(data types.MapStr, cond common.Condition) error {

	condData := cond.ToMapStr()
	moduleID := condData.String(ModuleID)
	if 0 == len(moduleID) {
		return errors.New("the module id is not set")
	}

	condData.Remove(ModuleID)

	appID := condData.String(BusinessID)
	if 0 == len(appID) {
		return errors.New("the business id is not set")
	}

	condData.Remove(BusinessID)

	setID := condData.String(SetID)
	if 0 == len(appID) {
		return errors.New("the set id is not set")
	}

	condData.Remove(SetID)

	targetURL := fmt.Sprintf("%s/api/v3/module/%s/%s/%s", cli.GetAddress(), appID, setID, moduleID)

	rst, err := cli.httpCli.PUT(targetURL, nil, data.ToJSON())
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

// SearchModules search some modules by condition
func (cli *Client) SearchModules(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	appID := data.String(BusinessID)
	if 0 == len(appID) {
		return nil, errors.New("the business id is not set")
	}

	data.Remove(BusinessID)

	setID := data.String(SetID)
	if 0 == len(appID) {
		return nil, errors.New("the set id is not set")
	}

	data.Remove(SetID)

	targetURL := fmt.Sprintf("%s/api/v3/module/search/%s/%s/%s", cli.GetAddress(), cli.supplierAccount, appID, setID)

	rst, err := cli.httpCli.POST(targetURL, nil, data.ToJSON())
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
