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

	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"

	"configcenter/src/framework/core/errors"
	"configcenter/src/framework/core/types"
)

// BusinessGetter business getter interface
type BusinessGetter interface {
	Business() BusinessInterface
}

// BusinessInterface business operation interface
type BusinessInterface interface {
	// SearchBusiness search host by condition,
	SearchBusiness(cond common.Condition) ([]types.MapStr, error)
	// CreateBusiness create host
	CreateBusiness(data types.MapStr) (int, error)
	// update update host by bizID, bizID could be separated by a comma
	UpdateBusiness(data types.MapStr, bizID int) error
	// DeleteBusiness delete host by bizID, bizID could be separated by a comma
	DeleteBusiness(bizID int) error
}

// Business define
type Business struct {
	cli *Client
}

func newBusiness(cli *Client) *Business {
	return &Business{
		cli: cli,
	}
}

// CreateBusiness create business
func (h *Business) CreateBusiness(data types.MapStr) (int, error) {
	targetURL := fmt.Sprintf("%s/api/v3/biz/%s", h.cli.GetAddress(), h.cli.GetSupplierAccount())
	rst, err := h.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return 0, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		if gs.Get(cccommon.HTTPBKAPIErrorCode).Int() == cccommon.CCErrCommDuplicateItem {
			return 0, errors.ErrDuplicateDataExisted
		}
		return 0, errors.New(gs.Get(cccommon.HTTPBKAPIErrorMessage).String())
	}

	id := gs.Get("data.bk_biz_id").Int()

	return int(id), nil
}

// UpdateBusiness update business by iD
func (h *Business) UpdateBusiness(data types.MapStr, bizID int) error {
	targetURL := fmt.Sprintf("%s/api/v3/biz/%s/%d", h.cli.GetAddress(), h.cli.GetSupplierAccount(), bizID)
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

// DeleteBusiness delete business by ID
func (h *Business) DeleteBusiness(bizID int) error {

	targetURL := fmt.Sprintf("%s/api/v3/biz/%s/%d", h.cli.GetAddress(), h.cli.GetSupplierAccount(), bizID)
	rst, err := h.cli.httpCli.DELETE(targetURL, nil, nil)
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

// SearchBusiness search bussiness operation
func (h *Business) SearchBusiness(cond common.Condition) ([]types.MapStr, error) {
	data := cond.ToMapStr()

	param := types.MapStr{
		"native":    1,
		"condition": data,
		"page": types.MapStr{
			"start": cond.GetStart(),
			"limit": cond.GetLimit(),
			"sort":  cond.GetSort(),
		},
	}

	out := param.ToJSON()
	//log.Infof("search business param %s", out)

	targetURL := fmt.Sprintf("%s/api/v3/biz/search/%s", h.cli.GetAddress(), h.cli.GetSupplierAccount())
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

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)
	return resultMap, err
}
