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

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
)

// GroupRelResByIDsOption group related resource by ids option
type GroupRelResByIDsOption struct {
	IDs       []int64       `json:"ids"`
	IDField   string        `json:"id_field"`
	RelField  string        `json:"rel_field"`
	ExtraCond mapstr.MapStr `json:"extra_cond,omitempty"`
}

// Validate GroupRelResByIDsOption
func (c *GroupRelResByIDsOption) Validate() errors.RawErrorInfo {
	if len(c.IDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(c.IDs) > common.BKMaxUpdateOrCreatePageSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxUpdateOrCreatePageSize},
		}
	}

	if len(c.IDField) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"id_field"}}
	}

	if len(c.RelField) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"rel_field"}}
	}

	return errors.RawErrorInfo{}
}

// GroupRelResByIDsResp group by ids response
type GroupRelResByIDsResp struct {
	BaseResp `json:",inline"`
	Data     map[int64][]interface{} `json:"data"`
}

// GroupByResKind is resource kind for group by operation
type GroupByResKind string

const (
	// ProcInstRelGroupByRes is process instance relation resource kind
	ProcInstRelGroupByRes GroupByResKind = "process_instance_relation"
	// ModuleGroupByRes is module resource kind
	ModuleGroupByRes GroupByResKind = "module"
	// ModuleHostRelGroupByRes is module host relation resource kind
	ModuleHostRelGroupByRes GroupByResKind = "module_host_config"
)

// CountResKindTableMap is GroupByResKind to mongodb table name map
var CountResKindTableMap = map[GroupByResKind]string{
	ProcInstRelGroupByRes:   common.BKTableNameProcessInstanceRelation,
	ModuleGroupByRes:        common.BKTableNameBaseModule,
	ModuleHostRelGroupByRes: common.BKTableNameModuleHostConfig,
}
