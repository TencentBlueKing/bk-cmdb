/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package general defines general resource cache types and utils
package general

import (
	"bytes"
	"encoding/json"
	"strings"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// general resource related field names
const (
	IDField       = "id"
	ResourceField = "resource"
	SubResField   = "sub_resource"
)

// ListDetailByIDsOpt is list general resource detail cache by ids option
type ListDetailByIDsOpt struct {
	Resource    ResType  `json:"resource"`
	SubResource string   `json:"sub_resource"`
	IDs         []int64  `json:"ids"`
	Fields      []string `json:"fields"`
}

// Validate ListDetailByIDsOpt
func (o *ListDetailByIDsOpt) Validate() ccErr.RawErrorInfo {
	if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if len(o.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(o.IDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListGeneralCacheResp is the general resource cache list response
type ListGeneralCacheResp struct {
	metadata.BaseResp `json:",inline"`
	Data              *ListGeneralCacheRes `json:"data"`
}

// ListGeneralCacheRes is the general resource cache list result
type ListGeneralCacheRes struct {
	Info StringArrRes `json:"info"`
}

// StringArrRes is string array result for cache value
type StringArrRes []string

// MarshalJSON marshal json
func (s StringArrRes) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	buf.WriteString(strings.Join(s, ","))
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshal json
func (s *StringArrRes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	arr := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*s = make([]string, len(arr))
	for i, val := range arr {
		(*s)[i] = string(val)
	}
	return nil
}

// DecodeStringArrRes decode StringArrRes into specified type array
func DecodeStringArrRes[T any](s StringArrRes) ([]T, error) {
	result := make([]T, len(s))
	for i, str := range s {
		err := json.Unmarshal([]byte(str), &result[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ListDetailByUniqueKeyOpt is the option to list general cache with key.
type ListDetailByUniqueKeyOpt struct {
	Resource    ResType       `json:"resource"`
	SubResource string        `json:"sub_resource"`
	Type        UniqueKeyType `json:"type"`
	Keys        []string      `json:"keys"`
	Fields      []string      `json:"fields"`
}

// Validate ListDetailByUniqueKeyOpt
func (o ListDetailByUniqueKeyOpt) Validate() ccErr.RawErrorInfo {
	if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if len(o.Type) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"type"},
		}
	}

	if len(o.Keys) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"keys"},
		}
	}

	if len(o.Keys) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"keys", common.BKMaxLimitSize},
		}
	}

	return ccErr.RawErrorInfo{}
}

// RefreshIDListOpt is refresh general resource detail cache by ids option
type RefreshIDListOpt struct {
	Resource ResType `json:"resource"`
	SubRes   string  `json:"sub_resource"`

	// full sync cond filter option
	CondID int64 `json:"cond_id"`
}

// Validate RefreshIDListOpt
func (o *RefreshIDListOpt) Validate() ccErr.RawErrorInfo {
	_, exists := SupportedResTypeMap[o.Resource]
	if !exists {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{ResourceField},
		}
	}

	_, exists = ResTypeHasSubResMap[o.Resource]
	if exists && o.SubRes == "" {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{SubResField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// RefreshDetailByIDsOpt is refresh general resource detail cache by ids option
type RefreshDetailByIDsOpt struct {
	Resource    ResType `json:"resource"`
	SubResource string  `json:"sub_resource"`
	IDs         []int64 `json:"ids"`
}

// Validate RefreshDetailByIDsOpt
func (o *RefreshDetailByIDsOpt) Validate() ccErr.RawErrorInfo {
	if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if len(o.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(o.IDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListDetailOpt is list general resource cache option
type ListDetailOpt struct {
	Resource    ResType       `json:"resource"`
	SubResource string        `json:"sub_resource"`
	Fields      []string      `json:"fields"`
	Page        *PagingOption `json:"page"`
}

// Validate ListDetailOpt
func (o *ListDetailOpt) Validate() ccErr.RawErrorInfo {
	if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if o.Page == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"paging_opt"},
		}
	}

	return o.Page.Validate()
}

// PagingOption is the paging option for list general resource operation
type PagingOption struct {
	StartID     int64  `json:"start_id,omitempty"`
	StartOid    string `json:"start_oid,omitempty"`
	StartIndex  int64  `json:"start_index,omitempty"`
	Limit       int64  `json:"limit"`
	EnableCount bool   `json:"enable_count,omitempty"`
}

// Validate PagingOption
func (o *PagingOption) Validate() ccErr.RawErrorInfo {
	if o.EnableCount {
		if o.StartID != 0 || o.StartOid != "" || o.StartIndex != 0 || o.Limit != 0 {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"page.enable_count"},
			}
		}
		return ccErr.RawErrorInfo{}
	}

	if o.Limit <= 0 || o.Limit > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"limit"},
		}
	}

	return ccErr.RawErrorInfo{}
}
