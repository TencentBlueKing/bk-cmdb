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
	cctime "configcenter/src/common/time"
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
	Data     interface{} `json:"data" mapstructure:"data"`
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

type ArrayResponse struct {
	BaseResp `json:",inline"`
	Data     []interface{} `json:"data"`
}

// HostCountResponse host count
type HostCountResponse struct {
	BaseResp `json:",inline"`
	Data     int64 `json:"data"`
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
	Condition map[string]interface{} `json:"condition"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition  *TimeCondition `json:"time_condition,omitempty"`
	Fields         string         `json:"fields,omitempty"`
	Start          int            `json:"start,omitempty"`
	Limit          int            `json:"limit,omitempty"`
	Sort           string         `json:"sort,omitempty"`
	DisableCounter bool           `json:"disable_counter,omitempty"`
}

type TimeConditionItem struct {
	Field string       `json:"field" bson:"field"`
	Start *cctime.Time `json:"start" bson:"start"`
	End   *cctime.Time `json:"end" bson:"end"`
}

type TimeCondition struct {
	Operator string              `json:"oper" bson:"oper"`
	Rules    []TimeConditionItem `json:"rules" bson:"rules"`
}

// MergeTimeCondition parse time condition and merge with common condition to construct a DB condition, only used by DB
func (tc *TimeCondition) MergeTimeCondition(condition map[string]interface{}) (map[string]interface{}, error) {
	if tc == nil {
		return nil, nil
	}

	if tc.Operator != "and" {
		return nil, errors.New(common.CCErrCommParamsInvalid, "time condition oper is invalid")
	}

	if len(tc.Rules) == 0 {
		return nil, errors.New(common.CCErrCommParamsNeedSet, "time condition rules not set")
	}

	timeCondition := make(map[string]interface{})
	for _, cond := range tc.Rules {
		if len(cond.Field) == 0 {
			return nil, errors.New(common.CCErrCommParamsNeedSet, "time condition field not set")
		}

		if cond.Start == nil && cond.End == nil {
			return nil, errors.New(common.CCErrCommParamsInvalid, "time condition start and end both not set")
		}

		if cond.Start == nil {
			timeCondition[cond.Field] = map[string]interface{}{common.BKDBLTE: cond.End}
			continue
		}

		if cond.End == nil {
			timeCondition[cond.Field] = map[string]interface{}{common.BKDBGTE: cond.Start}
			continue
		}

		if *cond.Start == *cond.End {
			timeCondition[cond.Field] = map[string]interface{}{common.BKDBEQ: cond.Start}
			continue
		}

		timeCondition[cond.Field] = map[string]interface{}{common.BKDBGTE: cond.Start, common.BKDBLTE: cond.End}
	}

	if len(condition) == 0 {
		return timeCondition, nil
	}

	return map[string]interface{}{common.BKDBAND: []map[string]interface{}{condition, timeCondition}}, nil
}

type ConditionWithTime struct {
	Condition []ConditionItem `json:"condition"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
}

// ConvTime cc_type key
func (o *QueryInput) ConvTime() error {
	for key, item := range o.Condition {
		convItem, err := o.convTimeItem(item)
		if nil != err {
			continue
		}
		o.Condition[key] = convItem
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
	ApplicationID int64   `json:"bk_biz_id"`
	HostIDs       []int64 `json:"bk_host_id"`
	ModuleID      int64   `json:"bk_module_id"`
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

type SearchInstBatchOption struct {
	IDs    []int64  `json:"bk_ids"`
	Fields []string `json:"fields"`
}

func (s *SearchInstBatchOption) Validate() (rawError errors.RawErrorInfo) {
	if len(s.IDs) == 0 || len(s.IDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_ids", common.BKMaxInstanceLimit},
		}
	}

	if len(s.Fields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"fields"},
		}
	}

	return errors.RawErrorInfo{}
}

// BkBaseResp base response defined in blueking api protocol
type BkBaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BKResponse struct {
	BkBaseResp `json:",inline"`
	Data       interface{} `json:"data"`
}
