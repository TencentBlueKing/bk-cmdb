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

	"configcenter/src/framework/core/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type HostGetter interface {
	Host() HostInterface
}

type HostInterface interface {
	// SearchHost search host by condition,
	SearchHost(cond common.Condition) ([]types.MapStr, error)
	// CreateHostBatch create host
	CreateHostBatch(data ...types.MapStr) ([]int, error)
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
func (h *Host) CreateHostBatch(data ...types.MapStr) ([]int, error) {
	infos := map[int]map[string]interface{}{}
	for index := range data {
		data[index].Set("import_from", "3")
		infos[index] = data[index]
	}
	param := types.MapStr{
		"bk_supplier_id": cccommon.BKDefaultSupplierID,
		"host_info":      infos,
	}
	targetURL := fmt.Sprintf("%s/api/v3/hosts/add", h.cli.GetAddress())
	rst, err := h.cli.httpCli.POST(targetURL, nil, param.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return nil, errors.New(gs.Get("bk_error_msg").String())
	}

	ids := []int{}
	gs.Get("data.success").ForEach(func(key, value gjson.Result) bool {
		ids = append(ids, int(value.Int()))
		return true
	})

	return ids, nil
}

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

func (h *Host) SearchHost(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	targetURL := fmt.Sprintf("%s/api/v3/hosts/search", h.cli.GetAddress())
	rst, err := h.cli.httpCli.POST(targetURL, nil, data.ToJSON())
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
