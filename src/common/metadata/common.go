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

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/basetype"
	"configcenter/src/common/errors"
	"fmt"
	"github.com/gin-gonic/gin/json"
)

type BaseResp struct {
	Result bool   `json:"result"`
	Code   int    `json:"bk_error_code"`
	ErrMsg string `json:"bk_error_msg"`
}

var SuccessBaseResp = BaseResp{Result: true, Code: common.CCSuccess, ErrMsg: common.CCSuccessStr}

type RespError struct {
	// error message
	Msg error
	// error code
	ErrCode int
}

const defaultError = "{\"result\": false, \"bk_error_code\": 1199000, \"bk_error_msg\": %s}"

func (r *RespError) Error() string {
	br := new(Response)
	br.Code = r.ErrCode
	br.ErrMsg = r.Msg.Error()
	if nil != r.Msg {
		if ccErr, ok := (r.Msg).(errors.CCErrorCoder); ok {
			br.Code = ccErr.GetCode()
			br.ErrMsg = ccErr.Error()
		}
	}

	js, err := json.Marshal(br)
	if err != nil {
		return fmt.Sprintf(defaultError, err.Error())
	}

	return string(js)
}

// data is the data you want to return to client.
func NewSuccessResp(data interface{}) *Response {
	return &Response{
		BaseResp: BaseResp{true, common.CCSuccess, common.CCSuccessStr},
		Data:     data,
	}
}

type Response struct {
	BaseResp `json:",inline"`
	Data     interface{} `json:"data"`
}

type MapResponse struct {
	BaseResp `json:",inline"`
	Data     map[string]*basetype.Type `json:"data"`
}

type RecursiveMapResponse struct {
	BaseResp `json:",inline"`
	Data     map[string]map[string]*basetype.Type `json:"data"`
}

type QueryInput struct {
	Condition interface{} `json:"condition"`
	Fields    string      `json:"fields,omitempty"`
	Start     int         `json:"start,omitempty"`
	Limit     int         `json:"limit,omitempty"`
	Sort      string      `json:"sort",omitempty`
}

type CloudHostModuleParams struct {
	ApplicationID int64        `json:"bk_biz_id"`
	HostInfoArr   []BkHostInfo `json:"host_info"`
	ModuleID      int64        `json:"bk_module_id"`
}

type BkHostInfo struct {
	IP      string `json:"bk_host_innerip"`
	CloudID int    `json:"bk_cloud_id"`
}

type DefaultModuleHostConfigParams struct {
	ApplicationID int   `json:"bk_biz_id"`
	HostID        []int `json:"bk_host_id"`
}

//common search struct
type SearchParams struct {
	Condition map[string]interface{} `json:"condition"`
	Page      map[string]interface{} `json:"page,omitempty"`
	Fields    []string               `json:"fields,omitempty"`
	Native    int                    `json:"native,omitempty"`
}

// PropertyGroupCondition used to reflect the property group json
type PropertyGroupCondition struct {
	Condition map[string]interface{} `json:"condition"`
	Data      map[string]interface{} `json:"data"`
}
