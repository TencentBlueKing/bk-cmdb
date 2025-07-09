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

// CreatedManyOptionResult create many api http response return this result struct
type CreatedManyOptionResult struct {
	BaseResp `json:",inline"`
	Data     CreateManyDataResult `json:"data"`
}

// CreatedOneOptionResult create One api http response return this result struct
type CreatedOneOptionResult struct {
	BaseResp `json:",inline"`
	Data     CreateOneDataResult `json:"data"`
}

// BatchCreateResp batch create response
type BatchCreateResp struct {
	BaseResp `json:",inline"`
	Data     BatchCreateResult `json:"data"`
}

// BatchCreateResult batch create result
type BatchCreateResult struct {
	IDs []uint64 `json:"ids"`
}
