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
	"errors"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

const (
	TopoCountMaxNum = 100
)

// HostPathOption find host path request
type HostPathOption struct {
	HostIDs []int64 `json:"ids"`
}

// Validate validate HostPathOption
func (h *HostPathOption) Validate() ccErr.RawErrorInfo {
	if len(h.HostIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(h.HostIDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}
	return ccErr.RawErrorInfo{}
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
	BizID       int64  `json:"bk_biz_id"`
	BizName     string `json:"biz_name"`
	ClusterID   int64  `json:"bk_cluster_id"`
	ClusterName string `json:"cluster_name"`
}

// HostNodeRelation get host and node relation message
type HostNodeRelation struct {
	BizIDs            []int64
	HostWithNode      map[int64][]Node
	ClusterIDWithName map[int64]string
}

// PodPathOption pod path request
type PodPathOption struct {
	PodIDs []int64 `json:"ids"`
}

// Validate validate PodPathReq
func (p *PodPathOption) Validate() ccErr.RawErrorInfo {
	if len(p.PodIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(p.PodIDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", common.BKMaxLimitSize},
		}
	}
	return ccErr.RawErrorInfo{}
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
	BizName      string       `json:"biz_name"`
	ClusterID    int64        `json:"bk_cluster_id"`
	ClusterName  string       `json:"cluster_name"`
	NamespaceID  int64        `json:"bk_namespace_id"`
	Namespace    string       `json:"namespace"`
	Kind         WorkloadType `json:"kind"`
	WorkloadID   int64        `json:"bk_workload_id"`
	WorkloadName string       `json:"workload_name"`
	PodID        int64        `json:"bk_pod_id"`
}

// KubeResourceInfo the type of the requested resource and the corresponding resource ID.
// it should be noted that when the kind is folder, the host cannot be obtained through
// the pod table. In this case, the node table needs to be used to find the corresponding
// number of hosts. Since the node is only associated with the cluster, the id in this
// scenario needs to pass the corresponding clusterID.
type KubeResourceInfo struct {
	Kind string `json:"kind"`
	ID   int64  `json:"id"`
}

// KubeTopoCountOption calculate the number of hosts or pods under the container resource node.
type KubeTopoCountOption struct {
	ResourceInfos []KubeResourceInfo `json:"resource_info"`
}

// Validate validate the KubeTopoCountOption
func (option *KubeTopoCountOption) Validate() ccErr.RawErrorInfo {

	if len(option.ResourceInfos) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"resource_info"},
		}
	}
	if len(option.ResourceInfos) > TopoCountMaxNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"resource_info", TopoCountMaxNum},
		}
	}
	for _, info := range option.ResourceInfos {
		if !IsKubeResourceKind(info.Kind) {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"non-kube objects", info.Kind},
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// KubeTopoCountRsp the response of the node host or the number of pods
type KubeTopoCountRsp struct {
	Kind  string `json:"kind"`
	ID    int64  `json:"id"`
	Count int64  `json:"count"`
}

// KubeTopoPathOption get container topology path request.
type KubeTopoPathOption struct {
	ReferenceObjID string            `json:"bk_reference_obj_id"`
	ReferenceID    int64             `json:"bk_reference_id"`
	Page           metadata.BasePage `json:"page"`
}

// Validate validate the KubeTopoPathOption
func (option *KubeTopoPathOption) Validate() ccErr.RawErrorInfo {

	if option.ReferenceID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{errors.New("bk_reference_id must be set")},
		}
	}

	// is the resource type legal
	if !IsKubeTopoResource(option.ReferenceObjID) {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{option.ReferenceObjID},
		}
	}

	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}
	return ccErr.RawErrorInfo{}
}

// KubeObjectInfo container object information.
type KubeObjectInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// KubeTopoPathRsp get topology path response.
type KubeTopoPathRsp struct {
	Info  []KubeObjectInfo `json:"info"`
	Count int              `json:"count"`
}

// SearchHostOption search host request
type SearchHostOption struct {
	BizID       int64                    `json:"bk_biz_id"`
	ClusterID   int64                    `json:"bk_cluster_id"`
	Folder      bool                     `json:"folder"`
	NamespaceID int64                    `json:"bk_namespace_id"`
	WorkloadID  int64                    `json:"bk_workload_id"`
	WlKind      WorkloadType             `json:"kind"`
	NodeCond    *NodeCondition           `json:"node_cond"`
	Ip          metadata.IPInfo          `json:"ip"`
	HostCond    metadata.SearchCondition `json:"host_condition"`
	Page        metadata.BasePage        `json:"page"`
}
