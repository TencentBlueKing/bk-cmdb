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

// CreateCommonInst create a common inst instance
func (cli *Client) CreateCommonInst(data types.MapStr) (int, error) {

	objID := data.String(ObjectID)
	if 0 == len(objID) {
		return 0, errors.New("the object id is not set")
	}

	data.Remove(ObjectID)

	targetURL := fmt.Sprintf("%s/api/v3/inst/%s/%s", cli.GetAddress(), cli.supplierAccount, objID)

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

// DeleteCommonInst delete a common inst instance
func (cli *Client) DeleteCommonInst(cond common.Condition) error {

	condData := cond.ToMapStr()
	instID, err := condData.Int(CommonInstID)
	if nil != err {
		return err
	}

	objID := condData.String(ObjectID)
	if 0 == len(objID) {
		return errors.New("the object id is not set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/inst/%s/%s/%d", cli.GetAddress(), cli.supplierAccount, objID, instID)

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

// UpdateCommonInst update a common inst instance
func (cli *Client) UpdateCommonInst(data types.MapStr, cond common.Condition) error {

	condData := cond.ToMapStr()
	instID, err := condData.Int(CommonInstID)
	if nil != err {
		return err
	}

	objID := condData.String(ObjectID)
	if 0 == len(objID) {
		return errors.New("the object id is not set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/inst/%s/%s/%d", cli.GetAddress(), cli.supplierAccount, objID, instID)

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

// SearchInst search all inst by condition
func (cli *Client) SearchInst(cond common.Condition) ([]types.MapStr, error) {

	condData := cond.ToMapStr()

	objID := condData.String(ObjectID)
	if 0 == len(objID) {
		return nil, errors.New("the object id is not set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/inst/search/owner/%s/object/%s", cli.GetAddress(), cli.supplierAccount, objID)

	// convert to the condition
	condInner := types.MapStr{
		"fields":    []string{},
		"condition": condData,
		"page": types.MapStr{
			"start": cond.GetStart(),
			"limit": 10,
			"sort":  cond.GetSort(),
		},
	}

	//fmt.Println("inner cond:", string(condInner.ToJSON()))

	rst, err := cli.httpCli.POST(targetURL, nil, condInner.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return nil, errors.New(gs.Get("bk_error_msg").String())
	}

	dataStr := gs.Get("data.info").String()
	if 0 == len(dataStr) {
		return nil, errors.New("data is empty")
	}

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)
	return resultMap, err
}
