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
	"encoding/json"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/coccyx/timeparser"
)

const defaultError = "{\"result\": false, \"bk_error_code\": 1199000, \"bk_error_msg\": %s}"

// RespError
type RespError struct {
	// error message
	Msg error
	// error code
	ErrCode int
	Data    interface{}
}

func (r *RespError) Error() string {
	br := new(Response)
	br.Code = r.ErrCode
	if nil != r.Msg {
		if ccErr, ok := (r.Msg).(errors.CCErrorCoder); ok {
			br.Code = ccErr.GetCode()
			br.ErrMsg = ccErr.Error()
		} else {
			br.ErrMsg = r.Msg.Error()
		}
	}
	br.Data = r.Data

	js, err := json.Marshal(br)
	if err != nil {
		return fmt.Sprintf(defaultError, err.Error())
	}

	return string(js)
}

// data is the data you want to return to client.
func NewSuccessResp(data interface{}) *Response {
	return &Response{
		BaseResp: BaseResp{Result: true, Code: common.CCSuccess, ErrMsg: common.CCSuccessStr},
		Data:     data,
	}
}

type Response struct {
	BaseResp `json:",inline"`
	Data     interface{} `json:"data"`
}

type BoolResponse struct {
	BaseResp `json:",inline"`
	Data     bool `json:"data"`
}

type Uint64Response struct {
	BaseResp `json:",inline"`
	Count    uint64 `json:"count"`
}

type CoreUint64Response struct {
	BaseResp `json:",inline"`
	Data     uint64 `json:"data"`
}

type MapArrayResponse struct {
	BaseResp `json:",inline"`
	Data     []mapstr.MapStr `json:"data"`
}

// ResponseInstData
type ResponseInstData struct {
	BaseResp `json:",inline"`
	Data     InstDataInfo `json:"data"`
}

// InstDataInfo response instance data result Data field
type InstDataInfo struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type ResponseDataMapStr struct {
	BaseResp `json:",inline"`
	Data     mapstr.MapStr `json:"data"`
}

type QueryInput struct {
	Condition interface{} `json:"condition"`
	Fields    string      `json:"fields,omitempty"`
	Start     int         `json:"start,omitempty"`
	Limit     int         `json:"limit,omitempty"`
	Sort      string      `json:"sort,omitempty"`
}

// ConvTime cc_type key
func (o *QueryInput) ConvTime() error {
	conds, ok := o.Condition.(map[string]interface{})
	if true != ok && nil != conds {
		return nil
	}
	for key, item := range conds {
		convItem, err := o.convTimeItem(item)
		if nil != err {
			continue
		}
		conds[key] = convItem
	}

	return nil
}

// convTimeItem cc_time_type
func (o *QueryInput) convTimeItem(item interface{}) (interface{}, error) {

	switch item.(type) {
	case map[string]interface{}:

		arrItem, ok := item.(map[string]interface{})
		if true == ok {
			_, timeTypeOk := arrItem[common.BKTimeTypeParseFlag]
			if timeTypeOk {
				delete(arrItem, common.BKTimeTypeParseFlag)
			}

			for key, value := range arrItem {
				switch value.(type) {

				case []interface{}:
					var err error
					arrItem[key], err = o.convTimeItem(value)
					if nil != err {
						return nil, err
					}
				case map[string]interface{}:
					arrItemVal, ok := value.(map[string]interface{})
					if ok {
						for key, value := range arrItemVal {
							var err error
							arrItemVal[key], err = o.convTimeItem(value)
							if nil != err {
								return nil, err
							}
						}
						arrItem[key] = value
					}

				default:
					if timeTypeOk {
						var err error
						arrItem[key], err = o.convInterfaceToTime(value)
						if nil != err {
							return nil, err
						}
					}

				}
			}
			item = arrItem
		}
	case []interface{}:
		arrItem, ok := item.([]interface{})
		if true == ok {
			for index, value := range arrItem {
				newValue, err := o.convTimeItem(value)
				if nil != err {
					return nil, err

				}
				arrItem[index] = newValue
			}
			item = arrItem

		}

	}

	return item, nil
}

func (o *QueryInput) convInterfaceToTime(val interface{}) (interface{}, error) {
	switch val.(type) {
	case string:
		ts, err := timeparser.TimeParserInLocation(val.(string), time.Local)
		if nil != err {
			return nil, err
		}
		return ts.Local(), nil
	default:
		ts, err := util.GetInt64ByInterface(val)
		if nil != err {
			return 0, err
		}
		t := time.Unix(ts, 0).Local()
		return t, nil
	}

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
	ApplicationID int64    `json:"bk_biz_id"`
	HostIDs       []int64  `json:"bk_host_id"`
	Metadata      Metadata `field:"metadata" json:"metadata" bson:"metadata"`
}

// common search struct
type SearchParams struct {
	Condition map[string]interface{} `json:"condition"`
	Page      map[string]interface{} `json:"page,omitempty"`
	Fields    []string               `json:"fields,omitempty"`
}

// PropertyGroupCondition used to reflect the property group json
type PropertyGroupCondition struct {
	Condition map[string]interface{} `json:"condition"`
	Data      map[string]interface{} `json:"data"`
}

type UpdateParams struct {
	Condition map[string]interface{} `json:"condition"`
	Data      map[string]interface{} `json:"data"`
}
type ListHostWithoutAppResponse struct {
	BaseResp `json:",inline"`
	Data     ListHostResult `json:"data"`
}
