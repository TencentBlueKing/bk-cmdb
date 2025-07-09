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

// AuthBathVerifyRequest TODO
type AuthBathVerifyRequest struct {
	Resources []AuthResource `json:"resources"`
}

// AuthResource TODO
type AuthResource struct {
	BizID        int64  `json:"bk_biz_id"`
	ResourceType string `json:"resource_type"`
	ResourceID   int64  `json:"resource_id"`
	ResourceIDEx string `json:"resource_id_ex"`
	Action       string `json:"action"`
	ParentLayers []struct {
		ResourceType string `json:"resource_type"`
		ResourceID   int64  `json:"resource_id"`
		ResourceIDEx string `json:"resource_id_ex"`
	} `json:"parent_layers"`
}

// AuthBathVerifyResult TODO
type AuthBathVerifyResult struct {
	AuthResource
	// the authorize decision, whether a user has been authorized or not.
	Passed bool `json:"is_pass"`
}
