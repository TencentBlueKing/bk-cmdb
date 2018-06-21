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
	cccommon "configcenter/src/common"
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"

	"configcenter/src/framework/core/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

// HostGetter host getter
type HostGetter interface {
	Host() HostInterface
}

// HostInterface host operation
type HostInterface interface {
	// SearchHost search host by condition,
	SearchHost(cond common.Condition) ([]types.MapStr, error)
	// CreateHostBatch create host
	CreateHostBatch(bizID int64, moduleIDS []int64, data ...types.MapStr) ([]int, error)
	// update update host by hostID, hostID could be separated by a comma
	UpdateHostBatch(data types.MapStr, hostID string) error
	// DeleteHost delete host by hostID, hostID could be separated by a comma
	DeleteHostBatch(hostID string) error
}

// Host define
type Host struct {
	cli *Client
}

func newHost(cli *Client) *Host {
	return &Host{
		cli: cli,
	}
}

// CreateHostBatch batch to create hosts
func (h *Host) CreateHostBatch(bizID int64, moduleIDS []int64, data ...types.MapStr) ([]int, error) {
	infos := map[int]map[string]interface{}{}
	for index := range data {
		data[index].Set("import_from", "3") // 3 means api import hosts
		infos[index] = data[index]
	}
	param := types.MapStr{
		"bk_biz_id":      bizID,
		"bk_module_id":   moduleIDS,
		"bk_supplier_id": cccommon.BKDefaultSupplierID,
		"host_info":      infos,
	}
	targetURL := fmt.Sprintf("%s/api/v3/hosts/sync/new/host", h.cli.GetAddress())
	rst, err := h.cli.httpCli.POST(targetURL, nil, param.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		fmt.Println("host info:", string(param.ToJSON()))
		fmt.Println("rst:", string(rst))
		return nil, errors.New(gs.Get("bk_error_msg").String())
	}

	ids := []int{}
	gs.Get("data.success").ForEach(func(key, value gjson.Result) bool {
		ids = append(ids, int(value.Int()))
		return true
	})

	return ids, nil
}

// UpdateHostBatch batch to update the hosts
func (h *Host) UpdateHostBatch(data types.MapStr, hostID string) error {

	data.Set("bk_host_id", hostID)
	targetURL := fmt.Sprintf("%s/api/v3/hosts/batch", h.cli.GetAddress())
	rst, err := h.cli.httpCli.PUT(targetURL, nil, data.ToJSON())
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

// DeleteHostBatch batch to delete the host
func (h *Host) DeleteHostBatch(hostID string) error {
	data := common.CreateCondition().Field("bk_host_id").Eq(hostID)

	targetURL := fmt.Sprintf("%s/api/v3/hosts/batch", h.cli.GetAddress())
	rst, err := h.cli.httpCli.DELETE(targetURL, nil, data.ToMapStr().ToJSON())
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

// SearchHost search the host
func (h *Host) SearchHost(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	conditions := []map[string]interface{}{}
	for key, val := range data {
		var operator string
		var value interface{}
		switch v := val.(type) {
		case types.MapStr:
			for k, vv := range v {
				operator = k
				value = vv
			}
		default:
			operator = cccommon.BKDBEQ
			value = val
		}
		condition := map[string]interface{}{
			"field":    key,
			"operator": operator,
			"value":    value,
		}

		conditions = append(conditions, condition)
	}

	params := map[string]interface{}{
		"bk_biz_id": -1,
		"page": types.MapStr{
			"start": cond.GetStart(),
			"limit": cond.GetLimit(),
			"sort":  cond.GetSort(),
		},
	}

	params["condition"] = []map[string]interface{}{
		map[string]interface{}{
			"bk_obj_id": "host",
			"fields":    []string{},
			"condition": conditions,
		},
		map[string]interface{}{
			"bk_obj_id": "set",
			"fields":    []string{},
			"condition": make([]map[string]interface{}, 0),
		},
		map[string]interface{}{
			"bk_obj_id": "module",
			"fields":    []string{},
			"condition": make([]map[string]interface{}, 0),
		},
		map[string]interface{}{
			"bk_obj_id": "biz",
			"fields":    []string{},
			"condition": make([]map[string]interface{}, 0),
		},
	}

	out, err := json.Marshal(params)
	if err != nil {
		log.Errorf("json error: %v", err)
	}

	//fmt.Printf("search host param:%s\n", out)

	targetURL := fmt.Sprintf("%s/api/v3/hosts/search", h.cli.GetAddress())
	rst, err := h.cli.httpCli.POST(targetURL, nil, out)
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

	hostMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &hostMap)
	if nil != err {
		return nil, err
	}
	resultMap := make([]types.MapStr, 0)

	for _, item := range hostMap {

		biz, err := item.MapStrArray("biz")
		if nil != err {
			return nil, err
		}

		set, err := item.MapStrArray("set")
		if nil != err {
			return nil, err
		}

		module, err := item.MapStrArray("module")
		if nil != err {
			return nil, err
		}

		host, err := item.MapStr("host")
		if nil != err {
			return nil, err
		}

		host.Set("biz", biz)
		host.Set("set", set)
		host.Set("module", module)

		item, err = item.MapStr("host")
		if nil != err {
			return nil, err
		}

		resultMap = append(resultMap, item)

	}

	//fmt.Println("host data:", resultMap)

	return resultMap, err
}
