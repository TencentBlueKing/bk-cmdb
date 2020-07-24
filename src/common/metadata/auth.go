/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

type AuthBathVerifyRequest struct {
	Resources []AuthResource `json:"resources"`
}

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

type AuthBathVerifyResult struct {
	AuthResource
	// the authorize decision, whether a user has been authorized or not.
	Passed bool `json:"is_pass"`
}
