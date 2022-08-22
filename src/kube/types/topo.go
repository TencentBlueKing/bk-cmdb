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
)

// HostPathReq node path for hosts request
type HostPathReq struct {
	HostIDs []int64 `json:"ids"`
}

// Validate validate HostPathReq
func (h *HostPathReq) Validate() errors.RawErrorInfo {
	if len(h.HostIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(h.HostIDs) > common.BKMaxLimitSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}
	return errors.RawErrorInfo{}
}

// HostPathResp node path for hosts response
type HostPathResp struct {
	metadata.BaseResp `json:",inline"`
	Data              HostPathData `json:"data"`
}

// HostPathData node path for hosts data
type HostPathData struct {
	Info []HostNodePath `json:"info"`
}

// HostNodePath node path for host
type HostNodePath struct {
	HostID int64      `json:"bk_host_id"`
	Paths  []NodePath `json:"paths"`
}

// NodePath node path
type NodePath struct {
	Path      string `json:"path"`
	ClusterID int64  `json:"bk_cluster_id"`
	BizID     int64  `json:"bk_biz_id"`
}

// HostNodeRelation get host and node relation message
type HostNodeRelation struct {
	BizIDs              []int64
	HostWithNode        map[int64][]mapstr.MapStr
	NodeIDWithBizID     map[int64]int64
	NodeIDWithClusterID map[int64]int64
	ClusterIDWithUid    map[int64]string
}

// PodPathReq pod path request
type PodPathReq struct {
	PodIDs []int64 `json:"ids"`
}

// Validate validate PodPathReq
func (p *PodPathReq) Validate() errors.RawErrorInfo {
	if len(p.PodIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(p.PodIDs) > common.BKMaxLimitSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}
	return errors.RawErrorInfo{}
}

// PodPathResp pod container topological path response
type PodPathResp struct {
	metadata.BaseResp `json:",inline"`
	Data              PodPathData `json:"data"`
}

// PodPathData pod container topological path data
type PodPathData struct {
	Info []PodPath `json:"info"`
}

// PodPath pod container topological path
type PodPath struct {
	PodID      int64        `json:"bk_pod_id"`
	Path       string       `json:"path"`
	WorkloadID int64        `json:"bk_workload_id"`
	Kind       WorkloadType `json:"kind"`
}
