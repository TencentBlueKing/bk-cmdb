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
	"configcenter/src/common/querybuilder"
	cctime "configcenter/src/common/time"
	"configcenter/src/common/util"

	"github.com/coccyx/timeparser"
)

const defaultError = "{\"result\": false, \"bk_error_code\": 1199000, \"bk_error_msg\": %s}"

// RespError TODO
type RespError struct {
	// error message
	Msg error
	// error code
	ErrCode int
	Data    interface{}
}

// Error 用于错误处理
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

// NewSuccessResp TODO
// data is the data you want to return to client.
func NewSuccessResp(data interface{}) *Response {
	return &Response{
		BaseResp: BaseResp{Result: true, Code: common.CCSuccess, ErrMsg: common.CCSuccessStr},
		Data:     data,
	}
}

// Response TODO
type Response struct {
	BaseResp `json:",inline"`
	Data     interface{} `json:"data" mapstructure:"data"`
}

// CountResponseContent count action response content.
type CountResponseContent struct {
	// Count count num.
	Count uint64 `json:"count"`
}

// CountResponse count action response.
type CountResponse struct {
	BaseResp `json:",inline"`
	Data     CountResponseContent `json:"data"`
}

// BoolResponse TODO
type BoolResponse struct {
	BaseResp `json:",inline"`
	Data     bool `json:"data"`
}

// Uint64Response TODO
type Uint64Response struct {
	BaseResp `json:",inline"`
	Count    uint64 `json:"count"`
}

// CoreUint64Response TODO
type CoreUint64Response struct {
	BaseResp `json:",inline"`
	Data     uint64 `json:"data"`
}

// ArrayResponse TODO
type ArrayResponse struct {
	BaseResp `json:",inline"`
	Data     []interface{} `json:"data"`
}

// HostCountResponse host count
type HostCountResponse struct {
	BaseResp `json:",inline"`
	Data     int64 `json:"data"`
}

// MapArrayResponse TODO
type MapArrayResponse struct {
	BaseResp `json:",inline"`
	Data     []mapstr.MapStr `json:"data"`
}

// ResponseInstData TODO
type ResponseInstData struct {
	BaseResp `json:",inline"`
	Data     InstDataInfo `json:"data"`
}

// InstDataInfo response instance data result Data field
type InstDataInfo struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

// ResponseDataMapStr TODO
type ResponseDataMapStr struct {
	BaseResp `json:",inline"`
	Data     mapstr.MapStr `json:"data"`
}

// QueryInput TODO
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

// TimeConditionItem TODO
type TimeConditionItem struct {
	Field string       `json:"field" bson:"field"`
	Start *cctime.Time `json:"start" bson:"start"`
	End   *cctime.Time `json:"end" bson:"end"`
}

// TimeCondition TODO
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

// ConditionWithTime TODO
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

// CloudHostModuleParams TODO
type CloudHostModuleParams struct {
	ApplicationID int64        `json:"bk_biz_id"`
	HostInfoArr   []BkHostInfo `json:"host_info"`
	ModuleID      int64        `json:"bk_module_id"`
}

// BkHostInfo TODO
type BkHostInfo struct {
	IP      string `json:"bk_host_innerip"`
	CloudID int    `json:"bk_cloud_id"`
}

// DefaultModuleHostConfigParams TODO
type DefaultModuleHostConfigParams struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostIDs       []int64 `json:"bk_host_id"`
	ModuleID      int64   `json:"bk_module_id"`
}

// Condition is common simple condition parameter struct.
type Condition struct {
	// Condition conditions.
	Condition map[string]interface{} `json:"condition"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
}

// SearchParams TODO
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

// UpdateParams TODO
type UpdateParams struct {
	Condition map[string]interface{} `json:"condition"`
	Data      map[string]interface{} `json:"data"`
}

// ListHostWithoutAppResponse TODO
type ListHostWithoutAppResponse struct {
	BaseResp `json:",inline"`
	Data     ListHostResult `json:"data"`
}

// SearchInstBatchOption TODO
type SearchInstBatchOption struct {
	IDs    []int64  `json:"bk_ids"`
	Fields []string `json:"fields"`
}

// Validate TODO
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

// BKResponse TODO
type BKResponse struct {
	BkBaseResp `json:",inline"`
	Data       interface{} `json:"data"`
}

// CommonSearchResult is common search action result.
type CommonSearchResult struct {
	// Info search result.
	Info []interface{} `json:"info"`
}

// BatchCreateSetRequest batch create set request struct
type BatchCreateSetRequest struct {
	Sets []map[string]interface{} `json:"sets"`
}

// OneSetCreateResult create one set return result
type OneSetCreateResult struct {
	Index    int         `json:"index"`
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"error_message"`
}

// CommonSearchFilter is a common search action filter struct,
// such like search instance or instance associations.
// And the conditions must abide by query filter.
type CommonSearchFilter struct {
	// ObjectID is target model object id.
	ObjectID string `json:"bk_obj_id"`

	// Conditions is target search conditions that make up by the query filter.
	Conditions *querybuilder.QueryFilter `json:"conditions"`

	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`

	// Fields indicates which fields should be returns, it's would be ignored if not exists.
	Fields []string `json:"fields"`

	// Page batch query action page.
	Page BasePage `json:"page"`
}

// Validate validates the common search filter struct,
// return the key and error if any one of keys is invalid.
func (f *CommonSearchFilter) Validate() (string, error) {
	// validates object id parameter.
	if len(f.ObjectID) == 0 {
		return "bk_obj_id", fmt.Errorf("empty bk_obj_id")
	}

	// validate page parameter.
	if err := f.Page.ValidateLimit(common.BKMaxInstanceLimit); err != nil {
		return "page.limit", err
	}

	// validate conditions parameter.
	if f.Conditions == nil {
		// empty conditions to match all.
		return "", nil
	}

	option := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if invalidKey, err := f.Conditions.Validate(option); err != nil {
		return fmt.Sprintf("conditions.%s", invalidKey), err
	}

	if f.Conditions.GetDeep() > querybuilder.MaxDeep {
		return "conditions.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
	}

	return "", nil
}

// GetConditions returns a database type conditions base on the query filter.
func (f *CommonSearchFilter) GetConditions() (map[string]interface{}, error) {
	if f.Conditions == nil {
		// empty conditions to match all.
		return map[string]interface{}{}, nil
	}

	// convert to mongo conditions.
	mgoFilter, invalidKey, err := f.Conditions.ToMgo()
	if err != nil {
		return nil, fmt.Errorf("invalid key, conditions.%s, err: %s", invalidKey, err)
	}

	return mgoFilter, nil
}

// CommonCountResult is common count action result.
type CommonCountResult struct {
	// Count count result.
	Count uint64 `json:"count"`
}

// CommonCountFilter is a common count action filter struct,
// such like search instance count or instance associations count.
// And the conditions must abide by query filter.
type CommonCountFilter struct {
	// ObjectID is target model object id.
	ObjectID string `json:"bk_obj_id"`

	// Conditions is target search conditions that make up by the query filter.
	Conditions *querybuilder.QueryFilter `json:"conditions"`

	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
}

// Validate validates the common count filter struct,
// return the key and error if any one of keys is invalid.
func (f *CommonCountFilter) Validate() (string, error) {
	// validates object id parameter.
	if len(f.ObjectID) == 0 {
		return "bk_obj_id", fmt.Errorf("empty bk_obj_id")
	}

	// validate conditions parameter.
	if f.Conditions == nil {
		// empty conditions to match all.
		return "", nil
	}

	option := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if invalidKey, err := f.Conditions.Validate(option); err != nil {
		return fmt.Sprintf("conditions.%s", invalidKey), err
	}

	if f.Conditions.GetDeep() > querybuilder.MaxDeep {
		return "conditions.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
	}

	return "", nil
}

// GetConditions returns a database type conditions base on the query filter.
func (f *CommonCountFilter) GetConditions() (map[string]interface{}, error) {
	if f.Conditions == nil {
		// empty conditions to match all.
		return map[string]interface{}{}, nil
	}

	// convert to mongo conditions.
	mgoFilter, invalidKey, err := f.Conditions.ToMgo()
	if err != nil {
		return nil, fmt.Errorf("invalid key, conditions.%s, err: %s", invalidKey, err)
	}

	return mgoFilter, nil
}
