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

const (
	// PodQueryLimit limit on the number of pod query
	PodQueryLimit = 500
)

// PodQueryReq pod query request
type PodQueryReq struct {
	WorkloadSpec `json:",inline" bson:",inline"`
	HostID       int64              `json:"bk_host_id"`
	NodeID       int64              `json:"bk_node_id"`
	NodeName     string             `json:"node_name"`
	Filter       *filter.Expression `json:"filter"`
	Fields       []string           `json:"fields,omitempty"`
	Page         metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate PodQueryReq
func (p *PodQueryReq) Validate() errors.RawErrorInfo {
	if (p.ClusterID != nil || p.NamespaceID != nil || (p.Ref != nil && p.Ref.ID != nil) || p.NodeID != 0) &&
		(p.ClusterUID != nil || p.Namespace != nil || (p.Ref != nil && p.Ref.Name != nil) || p.NodeName != "") {

		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if p.Ref != nil && ((p.Ref.Name == nil && p.Ref.ID == nil) || p.Ref.Kind == nil ||
		!IsInnerWorkload(WorkloadType(*p.Ref.Kind))) {

		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{RefField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, PodQueryLimit); err.ErrCode != 0 {
		return err
	}

	// todo validate Filter
	return errors.RawErrorInfo{}
}

// BuildCond build query pod condition
func (p *PodQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
	}

	if p.ClusterID != nil {
		cond[BKClusterIDFiled] = p.ClusterID
	}

	if p.ClusterUID != nil {
		cond[ClusterUIDField] = p.ClusterUID
	}

	if p.NamespaceID != nil {
		cond[BKNamespaceIDField] = p.NamespaceID
	}

	if p.Namespace != nil {
		cond[NamespaceField] = p.Namespace
	}

	if p.Ref != nil {
		if p.Ref.Kind != nil {
			cond[RefKindField] = p.Ref.Kind
		}

		if p.Ref.Name != nil {
			cond[RefNameField] = p.Ref.Name
		}

		if p.Ref.ID != nil {
			cond[RefIDField] = p.Ref.ID
		}
	}

	if p.HostID != 0 {
		cond[common.BKHostIDField] = p.HostID
	}

	if p.NodeID != 0 {
		cond[BKNodeIDField] = p.NodeID
	}

	if p.NodeName != "" {
		cond[NodeField] = p.NodeName
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
