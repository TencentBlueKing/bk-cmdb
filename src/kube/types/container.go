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

package types

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
)

// ContainerQueryReq container query request
type ContainerQueryReq struct {
	PodID  int64              `json:"bk_pod_id"`
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields,omitempty"`
	Page   metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate ContainerQueryReq
func (p *ContainerQueryReq) Validate() errors.RawErrorInfo {
	if p.PodID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{BKPodIDField},
		}
	}

	if errInfo, err := p.Page.Validate(false); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errInfo},
		}
	}

	// todo validate Filter
	return errors.RawErrorInfo{}
}

// BuildCond build query container condition
func (p *ContainerQueryReq) BuildCond(supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		BKPodIDField: p.PodID,
	}
	if supplierAccount != "" {
		cond[common.BkSupplierAccount] = supplierAccount
	}

	if p.Filter != nil {
		filterCond, err := p.Filter.ToMgo()
		if err != nil {
			return nil, err
		}
		cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, filterCond}}
	}
	return cond, nil
}
