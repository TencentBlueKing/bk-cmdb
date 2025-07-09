/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// RolePriResult TODO
type RolePriResult struct {
	BaseResp `json:",inline"`
	Data     []string `json:"data"`
}

// RoleAppResult TODO
type RoleAppResult struct {
	BaseResp `json:",inline"`
	Data     []map[string]interface{} `json:"data"`
}

// MainLineResult TODO
type MainLineResult struct {
	BaseResp `json:",inline"`
	Data     []map[string]interface{} `json:"data"`
}

// AppQueryResult TODO
type AppQueryResult struct {
	BaseResp `json:",inline"`
	Data     InstResult `json:"data"`
}

// ObjectAttrBatchResult TODO
type ObjectAttrBatchResult struct {
	BaseResp `json:",inline"`
	Data     map[string]ObjectAttr `json:"data"`
}

// ObjectAttr TODO
type ObjectAttr struct {
	Attr []interface{} `json:"attr"`
}

// ObjectAttrResult TODO
type ObjectAttrResult struct {
	BaseResp `json:",inline"`
	Data     []Attribute `json:"data"`
}

// ObjectAttrGroupResult TODO
type ObjectAttrGroupResult struct {
	BaseResp `json:",inline"`
	Data     []AttributeGroup `json:"data"`
}
