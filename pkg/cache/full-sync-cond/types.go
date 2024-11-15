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

// Package fullsynccond defines full sync cond related types
package fullsynccond

import (
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

const (
	// BKTableNameFullSyncCond is the full synchronization cache condition table
	BKTableNameFullSyncCond = "cc_FullSyncCond"
)

// FullSyncCond is the full synchronization cache condition
type FullSyncCond struct {
	ID          int64              `json:"id" bson:"id"`
	Resource    general.ResType    `json:"resource" bson:"resource"`
	SubResource string             `json:"sub_resource" bson:"sub_resource"`
	IsAll       bool               `json:"is_all" bson:"is_all"`
	Interval    int                `json:"interval" bson:"interval"`
	Condition   *filter.Expression `json:"condition,omitempty" bson:"condition,omitempty"`
	TenantID    string             `json:"tenant_id" bson:"tenant_id"`
}

// full sync cond field names
const (
	IDField        = "id"
	ResourceField  = "resource"
	SubResField    = "sub_resource"
	IsAllField     = "is_all"
	IntervalField  = "interval"
	ConditionField = "condition"
)

// NotAllCondLimit is the limit of full sync cond whose is_all is false
const NotAllCondLimit = 100

// CreateFullSyncCondOpt is the full synchronization cache condition create option
type CreateFullSyncCondOpt struct {
	Resource    general.ResType    `json:"resource"`
	SubResource string             `json:"sub_resource"`
	IsAll       bool               `json:"is_all"`
	Condition   *filter.Expression `json:"condition,omitempty"`
	Interval    int                `json:"interval"`
}

// Validate CreateFullSyncCondOpt
func (o *CreateFullSyncCondOpt) Validate() ccErr.RawErrorInfo {
	if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if !o.IsAll {
		if o.Condition == nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{fmt.Sprintf("%s and %s", IsAllField, ConditionField)},
			}
		}

		validateOpt := filter.NewDefaultExprOpt(make(map[string]enumor.FieldType))
		validateOpt.IgnoreRuleFields = true
		if err := o.Condition.Validate(validateOpt); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{fmt.Sprintf("%s is invalid, err: %v", ConditionField, err)},
			}
		}
	}

	if o.Interval < 6 || o.Interval > 7*24 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{IntervalField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// UpdateFullSyncCondOpt is the full synchronization cache condition update option
type UpdateFullSyncCondOpt struct {
	ID   int64                   `json:"id"`
	Data *UpdateFullSyncCondData `json:"data"`
}

// Validate UpdateFullSyncCondOpt
func (o *UpdateFullSyncCondOpt) Validate() ccErr.RawErrorInfo {
	if o.ID <= 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{IDField},
		}
	}

	if o.Data == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if o.Data.Interval < 6 || o.Data.Interval > 7*24 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{IntervalField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// UpdateFullSyncCondData is the full synchronization cache condition update data
type UpdateFullSyncCondData struct {
	Interval int `json:"interval"`
}

// DeleteFullSyncCondOpt is the full synchronization cache condition delete option
type DeleteFullSyncCondOpt struct {
	ID int64 `json:"id"`
}

// Validate DeleteFullSyncCondOpt
func (o *DeleteFullSyncCondOpt) Validate() ccErr.RawErrorInfo {
	if o.ID <= 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{IDField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListFullSyncCondOpt is the full synchronization cache condition list option
type ListFullSyncCondOpt struct {
	Resource    general.ResType `json:"resource"`
	SubResource string          `json:"sub_resource"`
	IDs         []int64         `json:"ids"`
}

// Validate ListFullSyncCondOpt
func (o *ListFullSyncCondOpt) Validate() ccErr.RawErrorInfo {
	if o.Resource == "" && len(o.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{fmt.Sprintf("%s and %s", ResourceField, "ids")},
		}
	}

	if o.Resource != "" {
		if rawErr := o.Resource.ValidateWithSubRes(o.SubResource); rawErr.ErrCode != 0 {
			return rawErr
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListFullSyncCondResp is the full synchronization cache condition list response
type ListFullSyncCondResp struct {
	metadata.BaseResp `json:",inline"`
	Data              *ListFullSyncCondRes `json:"data"`
}

// ListFullSyncCondRes is the full synchronization cache condition list result
type ListFullSyncCondRes struct {
	Info []FullSyncCond `json:"info"`
}

// ListCacheByFullSyncCondOpt is the list general resource cache by full sync cond option
type ListCacheByFullSyncCondOpt struct {
	CondID int64    `json:"cond_id"`
	Cursor int64    `json:"cursor"`
	Limit  int64    `json:"limit"`
	Fields []string `json:"fields"`
}

// Validate CreateFullSyncCondOpt
func (o *ListCacheByFullSyncCondOpt) Validate() ccErr.RawErrorInfo {
	if o.CondID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"cond_id"},
		}
	}

	if o.Cursor < 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"cursor"},
		}
	}

	if o.Limit <= 0 || o.Limit > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"limit"},
		}
	}

	return ccErr.RawErrorInfo{}
}
