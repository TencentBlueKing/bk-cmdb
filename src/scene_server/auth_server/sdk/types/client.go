/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import "configcenter/src/scene_server/auth_server/sdk/operator"

// BaseResp TODO
type BaseResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// GetPolicyOption TODO
type GetPolicyOption AuthOptions

// GetPolicyResp TODO
type GetPolicyResp struct {
	BaseResp `json:",inline"`
	Data     *operator.Policy `json:"data"`
}

// ListPolicyOptions TODO
type ListPolicyOptions struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Actions   []Action   `json:"actions"`
	Resources []Resource `json:"resources"`
}

// ListPolicyResp TODO
type ListPolicyResp struct {
	BaseResp `json:",inline"`
	Data     []*ActionPolicy `json:"data"`
}
