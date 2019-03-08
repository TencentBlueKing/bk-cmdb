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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
)

// SetGetter the set getter interface
type SetGetter interface {
	Set() SetInterface
}

// SetInterface the set interface
type SetInterface interface {
	CreateSet(bizID int64, data types.MapStr) (int, error)
	DeleteSet(cond common.Condition) error
	UpdateSet(bizID int64, data types.MapStr, cond common.Condition) error
	SearchSets(cond common.Condition) ([]types.MapStr, error)
}

func newSet(cli *Client) *Set {
	return &Set{
		cli: cli,
	}
}

// Set the set interface implement
type Set struct {
	cli *Client
}

// CreateSet create a new Set
func (cli *Set) CreateSet(bizID int64, data types.MapStr) (int, error) {

	data.Set(SupplierAccount, cli.cli.GetSupplierAccount())
	targetURL := fmt.Sprintf("%s/api/v3/set/%d", cli.cli.GetAddress(), bizID)

	rst, err := cli.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return 0, err
	}

	//fmt.Println("the set id:", string(rst))
	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return 0, errors.New(gs.Get("bk_error_msg").String())
	}

	// parse id
	id := gs.Get("data.bk_set_id").Int()

	return int(id), nil
}

// DeleteSet delete a set by condition
func (cli *Set) DeleteSet(cond common.Condition) error {
	data := cond.ToMapStr()

	appID := data.String(BusinessID)
	if 0 == len(appID) {
		return errors.New("the business id is not set")
	}

	setID := data.String(SetID)
	if 0 == len(appID) {
		return errors.New("the set id is not set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/set/%s/%s", cli.cli.GetAddress(), appID, setID)

	rst, err := cli.cli.httpCli.DELETE(targetURL, nil, nil)
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

// UpdateSet update a set by condition
func (cli *Set) UpdateSet(bizID int64, data types.MapStr, cond common.Condition) error {

	condData := cond.ToMapStr()

	setID := condData.String(SetID)

	targetURL := fmt.Sprintf("%s/api/v3/set/%d/%s", cli.cli.GetAddress(), bizID, setID)

	rst, err := cli.cli.httpCli.PUT(targetURL, nil, data.ToJSON())
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

// SearchSets search some sets by condition
func (cli *Set) SearchSets(cond common.Condition) ([]types.MapStr, error) {
	data := cond.ToMapStr()

	appID := data.String(BusinessID)
	if 0 == len(appID) {
		return nil, errors.New("the business id is not set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/set/search/%s/%s", cli.cli.GetAddress(), cli.cli.supplierAccount, appID)
	// convert to the condition
	condInner := types.MapStr{
		"fields":    []string{},
		"condition": data,
		"page": types.MapStr{
			"start": cond.GetStart(),
			"limit": cond.GetLimit(),
			"sort":  cond.GetSort(),
		},
	}
	rst, err := cli.cli.httpCli.POST(targetURL, nil, condInner.ToJSON())
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

	//fmt.Println("data:", dataStr)

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)
	return resultMap, err
}
