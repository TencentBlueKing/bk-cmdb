/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
)

// BaseResp common result struct
type BaseResp struct {
	Result bool   `json:"result"`
	Code   int    `json:"bk_error_code"`
	ErrMsg string `json:"bk_error_msg"`
}

// SuccessBaseResp default result
var SuccessBaseResp = BaseResp{Result: true, Code: common.CCSuccess, ErrMsg: common.CCSuccessStr}

// CreatedCount created count struct
type CreatedCount struct {
	Count int64 `json:"created_count"`
}

// UpdatedCount created count struct
type UpdatedCount struct {
	Count int64 `json:"updated_count"`
}

// DeletedCount created count struct
type DeletedCount struct {
	Count int64 `json:"deleted_count"`
}

// Exception exception info
type Exception struct {
	Message string      `json:"message"`
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
}

// CreateManyInfoResult create many function return this result struct
type CreateManyInfoResult struct {
	Created    []int64         `json:"created"`
	Repeated   []mapstr.MapStr `json:"repeated"`
	Exceptions []Exception     `json:"exception"`
}

// CreateManyDataResult the data struct definition in create many function result
type CreateManyDataResult struct {
	Count int64                `json:"count"`
	Info  CreateManyInfoResult `json:"info"`
}

// CreateOneDataResult the data struct definition in create one function result
type CreateOneDataResult struct {
	Count int64        `json:"count"`
	Info  CreatedCount `json:"info"`
}

// SetManyInfoResult set many function return this result struct
type SetManyInfoResult struct {
	Created    []int64         `json:"created"`
	Updated    []mapstr.MapStr `json:"updated"`
	Exceptions []Exception     `json:"exception"`
}

// SetManyDataResult the data struct definition in create many function result
type SetManyDataResult struct {
	Count int64             `json:"count"`
	Info  SetManyInfoResult `json:"info"`
}

// SetOneInfoResult the info struct definition in create one function result's Info field
type SetOneInfoResult struct {
	CreatedCount `json:",inline"`
	UpdatedCount `json:",inline"`
}

// SetOneDataResult the data struct definition in create one function result
type SetOneDataResult struct {
	Count int64        `json:"count"`
	Info  CreatedCount `json:"info"`
}

// SearchDataResult common search data result
type SearchDataResult struct {
	Count int64           `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}
