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
	"context"
	"net/http"

	"configcenter/src/framework/core/types"
)

type ListHostsCtx struct {
    BaseCtx
	Filter ListHostFilter
}

type ListHostFilter struct {
	Page Page `json:"page"`
	// if you list hosts from the host resource pool,
	// the value should be -1, this is a special one.
	BusinessID      int64           `json:"bk_biz_id"`
	IPCtx           IPCtx           `json:"ip"`
	SearchCondition SearchCondition `json:"condition"`
}

type IPCtxFlag string

const (
	// search the inner ip or outer ip.
	InnerIP IPCtxFlag = "bk_host_innerip"
	OuterIP IPCtxFlag = "bk_host_outerip"
)

type IPCtx struct {
	Flag IPCtxFlag `json:"flag"`
	// whether match the ip exactly:
	// 0: fuzzy match the ip .
	// 1: match the ip exactly.
	Exact int64 `json:"exact"`
	// ip list you want to search.
	IPList []string `json:"data"`
}

type SearchCondition struct {
	// can be "host", "module", "set", "biz", "object"
	ObjectName string   `json:"bk_obj_id"`
	Fields     []string `json:"fields"`
	Condition  []Filter `json:"condition"`
}

type Filter struct {
	// attribute's name
	Field string `json:"field"`
	// can be one of $eq, $neq, $in, $nin.
	Operator string `json:"operator"`
	// the value of this Field
	Value string `json:"value"`
}

type ListHostResult struct {
	BaseResp `json:",inline"`
	Data     HostsInfo `json:"data"`
}

type HostsInfo struct {
	Count int64 `json:"count"`
	// info map format:
	// map["module name"][module info]
	// module name: biz, host, module, set, object.
	Info []map[string]types.MapStr `json:"info"`
}

type GetHostCtx struct {
	Ctx     context.Context
	Header  http.Header
	Tenancy string
	HostID  int64
}

type GetHostResult struct {
	BaseResp `json:",inline"`
	Data     []HostAttribute `json:"data"`
}

type HostAttribute struct {
	ID    string `json:"bk_property_id"`
	Name  string `json:"bk_property_name"`
	Value string `json:"bk_property_value"`
}

type GetHostSnapshotCtx struct {
    BaseCtx
	HostID int64
}

type GetHostSnapshotResult struct {
	BaseResp `json:",inline"`
	Data     types.MapStr `json:"data"`
}

type UpdateHostsAttributesCtx struct {
    BaseCtx
	Attributes HostsAttributes
}

type HostsAttributes struct {
	// host ids, comma separated.
	// like: "1,2,4"
	HostIDs    string       `json:"bk_host_id"`
	Attributes types.MapStr `json:",inline"`
}

type DeleteHostsCtx struct {
    BaseCtx
	Hosts  DeletedHostsInfo
}

type DeletedHostsInfo struct {
	// host ids, comma separated.
	// like: "1,2,4"
	HostIDs string `json:"bk_host_id"`
	Tenancy string `json:"bk_supplier_account"`
}
