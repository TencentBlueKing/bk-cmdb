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
	"context"
	"net/http"
)

type BaseResp struct {
	Result bool   `json:"result"`
	Code   int    `json:"bk_error_code"`
	ErrMsg string `json:"bk_error_msg"`
}

type Response struct {
	BaseResp `json:",inline"`
	Data     interface{} `json:"data"`
}

type Query struct {
	Page      Page         `json:"page"`
	Fields    []string     `json:"fields"`
	Condition types.MapStr `json:"condition"`
}

type Page struct {
	Limit uint64 `json:"limit,omitempty"`
	Start uint64 `json:"start,omitempty"`
	Sort  string `json:"sort,omitempty"`
}

type ListInfo struct {
	Count int64          `json:"count"`
	Info  []types.MapStr `json:"info"`
}

type BaseCtx struct {
	Ctx    context.Context
	Header http.Header
}

type OperatorKind string

const (
	Regex OperatorKind = "$regex"
	Eq    OperatorKind = "$eq"
)

type QueryVerb struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
