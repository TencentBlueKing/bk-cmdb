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

import "configcenter/src/framework/core/types"

// CreateInstanceCtx TODO
type CreateInstanceCtx struct {
	BaseCtx
	Tenancy  string
	ObjectID string
	Instance types.MapStr
}

// CreateInstanceResult TODO
type CreateInstanceResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		ID int64 `json:"bk_inst_id"`
	} `json:"data"`
}

// ListInstanceCtx TODO
type ListInstanceCtx struct {
	BaseCtx
	Tenancy  string
	ObjectID string
	Filter   Query
}

// QueryInstance TODO
type QueryInstance struct {
	Page   Page     `json:"page"`
	Fields []string `json:"fields"`
	// Format: map[object id][]QueryVerb
	Condition map[string][]QueryVerb `json:"condition"`
}

// ListInstanceResult TODO
type ListInstanceResult struct {
	BaseResp `json:",inline"`
	Data     ListInfo `json:"data"`
}

// UpdateObjectCtx TODO
type UpdateObjectCtx struct {
	BaseCtx
	Tenancy    string
	ObjectID   string
	InstanceID int64
	Object     types.MapStr
}

// DeleteObjectCtx TODO
type DeleteObjectCtx struct {
	BaseCtx
	Tenancy    string
	ObjectID   string
	InstanceID int64
}
