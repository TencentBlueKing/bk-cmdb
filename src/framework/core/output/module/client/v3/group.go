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
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

// GroupGetter group getter
type GroupGetter interface {
	Group() GroupInterface
}

// GroupInterface group interface
type GroupInterface interface {
	CreateGroup(data types.MapStr) (int, error)
	DeleteGroup(cond common.Condition) error
	UpdateGroup(data types.MapStr, cond common.Condition) error
	SearchGroups(cond common.Condition) ([]types.MapStr, error)
}

// Group group data struct
type Group struct {
	cli *Client
}

func newGroup(cli *Client) *Group {
	return &Group{
		cli: cli,
	}
}

// CreateGroup create a group
func (g *Group) CreateGroup(data types.MapStr) (int, error) {
	data.Set("bk_supplier_account", g.cli.GetSupplierAccount())
	if !data.Exists("bk_group_name") {
		return 0, errors.New("bk_group_name must set")
	}
	targetURL := fmt.Sprintf("%s/api/v3/objectatt/group/new", g.cli.GetAddress())

	out := data.ToJSON()
	log.Infof("create group %s", out)
	rst, err := g.cli.httpCli.POST(targetURL, nil, out)
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

// DeleteGroup delete a group by condition
func (g *Group) DeleteGroup(cond common.Condition) error {

	data := cond.ToMapStr()
	id, err := data.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/objectatt/group/groupid/%d", g.cli.GetAddress(), id)

	rst, err := g.cli.httpCli.DELETE(targetURL, nil, nil)
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

// UpdateGroup update a group by condition
func (g *Group) UpdateGroup(data types.MapStr, cond common.Condition) error {
	data.Set("bk_supplier_account", g.cli.GetSupplierAccount())

	param := types.MapStr{
		"condition": cond.ToMapStr(),
		"data":      data,
	}

	log.Infof("update group by %s to %s", cond.ToMapStr().ToJSON(), data.ToJSON())
	targetURL := fmt.Sprintf("%s/api/v3/objectatt/group/update", g.cli.GetAddress())
	rst, err := g.cli.httpCli.PUT(targetURL, nil, param.ToJSON())
	if nil != err {
		log.Errorf("post error %v", err)
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		log.Errorf("post error %s", rst)
		return errors.New(gs.Get("bk_error_msg").String())
	}

	return nil
}

// SearchGroups search some group by condition
func (g *Group) SearchGroups(cond common.Condition) ([]types.MapStr, error) {
	objid := cond.ToMapStr().String(ObjectID)

	if len(objid) <= 0 {
		return nil, errors.New("bk_obj_id must set")
	}

	targetURL := fmt.Sprintf("%s/api/v3/objectatt/group/property/owner/%s/object/%s", g.cli.GetAddress(), g.cli.GetSupplierAccount(), objid)
	rst, err := g.cli.httpCli.POST(targetURL, nil, cond.ToMapStr().ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		log.Errorf("falied to search group %s", rst)
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
