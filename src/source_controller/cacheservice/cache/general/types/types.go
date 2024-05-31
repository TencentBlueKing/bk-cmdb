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

// Package types defines general resource cache types
package types

import (
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/watch"
)

const (
	// PageSize is the default page size
	PageSize = common.BKMaxLimitSize
)

// EventType is the event type
type EventType string

const (
	// Upsert is the create or update event type
	Upsert EventType = "upsert"
	// Delete is the delete event type
	Delete EventType = "delete"
	// Init is the initialization event type
	Init EventType = "init"
)

// FullSyncCondEvent is the event of full sync condition
type FullSyncCondEvent struct {
	EventMap map[EventType][]*fullsynccond.FullSyncCond
}

// FullSyncCondInfo is the full sync condition info for the general resource cache
type FullSyncCondInfo struct {
	SubResource     string
	IsAll           bool
	Interval        time.Duration
	Condition       *filter.Expression
	SupplierAccount string
}

// ListDetailByIDsOpt is list general resource detail cache by ids option
type ListDetailByIDsOpt struct {
	SubRes   string
	IsSystem bool
	IDKeys   []string
	Fields   []string
}

// Validate ListDetailByIDsOpt
func (o *ListDetailByIDsOpt) Validate(hasSubRes bool) ccErr.RawErrorInfo {
	if o.SubRes == "" && hasSubRes {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"sub_res"}}
	}

	if len(o.IDKeys) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(o.IDKeys) > PageSize {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit, Args: []interface{}{"ids", PageSize}}
	}

	return ccErr.RawErrorInfo{}
}

// ListDetailByUniqueKeyOpt is list general resource detail cache by unique keys option
type ListDetailByUniqueKeyOpt struct {
	SubRes   string
	IsSystem bool
	Type     general.UniqueKeyType
	Keys     []string
	Fields   []string
}

// Validate ListDetailByUniqueKeyOpt
func (o *ListDetailByUniqueKeyOpt) Validate(hasSubRes bool) ccErr.RawErrorInfo {
	if o.SubRes == "" && hasSubRes {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"sub_res"}}
	}

	if len(o.Type) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"type"}}
	}

	if len(o.Keys) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"keys"}}
	}

	if len(o.Keys) > PageSize {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit, Args: []interface{}{"keys", PageSize}}
	}

	return ccErr.RawErrorInfo{}
}

// RefreshDetailByIDsOpt is refresh general resource detail cache by ids option
type RefreshDetailByIDsOpt struct {
	SubResource string
	IDKeys      []string
}

// Validate RefreshDetailByIDsOpt
func (o *RefreshDetailByIDsOpt) Validate() ccErr.RawErrorInfo {
	if len(o.IDKeys) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(o.IDKeys) > PageSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", PageSize},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListDetailOpt is list general resource cache from start option
type ListDetailOpt struct {
	Fields       []string
	OnlyListID   bool
	IDListFilter *IDListFilterOpt
	Page         *general.PagingOption
}

// Validate ListDetailOpt
func (o *ListDetailOpt) Validate(hasSubRes bool) ccErr.RawErrorInfo {
	if o.IDListFilter == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"filter_opt"},
		}
	}

	if rawErr := o.IDListFilter.Validate(hasSubRes); rawErr.ErrCode != 0 {
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
	StartID    int64
	StartOid   string
	StartIndex int64
	Limit      int64
}

// Validate PagingOption
func (o *PagingOption) Validate() ccErr.RawErrorInfo {
	if o.Limit <= 0 || o.Limit > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"limit"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// IDListFilterOpt is the id list filter option for list general resource operation
type IDListFilterOpt struct {
	IDListKey string
	*BasicFilter

	// full sync cond filter option
	IsAll bool
	Cond  *filter.Expression
}

// Validate IDListFilterOpt
func (o *IDListFilterOpt) Validate(hasSubRes bool) ccErr.RawErrorInfo {
	if o.BasicFilter == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"basic_filter"},
		}
	}

	if rawErr := o.BasicFilter.Validate(hasSubRes); rawErr.ErrCode != 0 {
		return rawErr
	}

	if !o.IsAll && o.Cond == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"cond"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// BasicFilter is the basic filter for getting general resource data from db
type BasicFilter struct {
	SubRes          string
	SupplierAccount string

	// IsSystem defines whether id list is for system use, system resource do not need to be filtered by SupplierAccount
	IsSystem bool
}

// Validate BasicFilter
func (o *BasicFilter) Validate(hasSubRes bool) ccErr.RawErrorInfo {
	if o.SubRes == "" && hasSubRes {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"sub_res"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// WatchEventData is the general resource watch event data
type WatchEventData struct {
	ChainNode *watch.ChainNode
	Data      filter.JsonString
}
