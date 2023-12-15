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

// Package types defines biz topology cache types
package types

// GetBizTopoOption get business topology options
type GetBizTopoOption struct {
	BizID int64 `json:"bk_biz_id"`
}

// BizTopo is the business topology
type BizTopo struct {
	Biz   *BizInfo `json:"biz"`
	Nodes []Node   `json:"nds"`
}

// BizInfo is the business info
type BizInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"nm"`
	Count *int64 `json:"cnt,omitempty"`
}

// Node is the common topology node info
type Node struct {
	// Kind topology node kind
	Kind string `json:"kind"`
	// ID topology node id
	ID int64 `json:"id"`
	// Name topology node name
	Name string `json:"nm"`
	// Count the resource count of the topology node
	Count *int64 `json:"cnt,omitempty"`
	// SubNodes the sub-nodes of current topology node
	SubNodes []Node `json:"nds,omitempty"`

	// ParentID topology node parent id, is an intermediate value only used to rearrange topology tree
	ParentID int64 `json:"-"`
}

// TopoType is the topology tree type
type TopoType string

const (
	// KubeType is the kube topology tree type
	KubeType TopoType = "kube"
)
