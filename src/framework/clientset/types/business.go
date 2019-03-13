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

import (
    "configcenter/src/framework/core/types"
)

type BusinessResponse struct {
	BaseResp `json:",inline"`
	Data     types.MapStr `json:"data"`
}

type CreateBusinessCtx struct {
    BaseCtx
	Tenancy      string
	BusinessInfo types.MapStr
}

type UpdateBusinessCtx struct {
    BaseCtx
	Tenancy      string
	BusinessID   int64
	BusinessInfo types.MapStr
}

type DeleteBusinessCtx struct {
    BaseCtx
	Tenancy    string
	BusinessID int64
}

type ListBusinessCtx struct {
    BaseCtx
	Tenancy   string
	QueryInfo Query
}

type ListBusinessResult struct {
    BaseResp `json:",inline"`
    Data     ListInfo `json:"data"`
}


