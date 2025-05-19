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

// BizBriefTopology is the brief topology of business
type BizBriefTopology struct {
	// basic business info
	Biz *BriefBizInfo `json:"biz"`
	// the idle set nodes info
	Idle []*BriefNode `json:"idle"`
	// the other common nodes
	Nodes []*BriefNode `json:"nds"`
}

// BriefBizInfo is the brief info of business
type BriefBizInfo struct {
	// business id
	ID int64 `json:"id" bson:"bk_biz_id"`
	// business name
	Name string `json:"nm" bson:"bk_biz_name"`
	// describe it's a resource pool business or normal business.
	// 0: normal business
	// >0: special business, like resource pool business.
	Default int `json:"dft" bson:"default"`
}

// BriefNode is the brief biz topo node
type BriefNode struct {
	// the object of this node, like set or module
	Object string `json:"obj"`
	// the node's instance id, like set id or module id
	ID int64 `json:"id"`
	// the node's name, like set name or module name
	Name string `json:"nm"`
	// only set, module has this field.
	// describe what kind of set or module this node is.
	// 0: normal module or set.
	// >1: special set or module
	Default *int `json:"dft,omitempty"`
	// the sub-nodes of current node
	SubNodes []*BriefNode `json:"nds"`
}
