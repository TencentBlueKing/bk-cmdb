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

package metadata

type PullResourceParam struct {
	Collection string                 `json:"collection"`
	Condition  map[string]interface{} `json:"condition"`
	Fields     []string               `json:"fields"`
	Limit      int64                  `json:"limit"`
	Offset     int64                  `json:"offset"`
}

type PullResourceResponse struct {
	BaseResp `json:",inline"`
	Data     PullResourceResult `json:"data"`
}

type PullResourceResult struct {
	Count int64                    `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

type QueryPolicyParam struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

type Subject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Action struct {
	ID string `json:"id"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type QueryPolicyResponse struct {
	BaseResponse `json:",inline"`
	Data         *PolicyExpression `json:"data"`
}

type Operator string

const (
	OperatorAnd                Operator = "AND"
	OperatorOr                 Operator = "OR"
	OperatorEqual              Operator = "eq"
	OperatorNotEqual           Operator = "not_eq"
	OperatorIn                 Operator = "in"
	OperatorNotIn              Operator = "not_in"
	OperatorContains           Operator = "contains"
	OperatorNotContains        Operator = "not_contains"
	OperatorStartsWith         Operator = "starts_with"
	OperatorNotStartsWith      Operator = "not_starts_with"
	OperatorEndsWith           Operator = "ends_with"
	OperatorNotEndsWith        Operator = "not_ends_with"
	OperatorLessThan           Operator = "lt"
	OperatorLessThanOrEqual    Operator = "lte"
	OperatorGreaterThan        Operator = "gt"
	OperatorGreaterThanOrEqual Operator = "gte"
	OperatorAny                Operator = "any"
)

type PolicyExpression struct {
	Operator Operator            `json:"op"`
	Content  []*PolicyExpression `json:"content,omitempty"`
	Field    string              `json:"field,omitempty"`
	Value    interface{}         `json:"value,omitempty"`
}

type QueryPolicyByActionParam struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Action    []Action   `json:"actions"`
	Resources []Resource `json:"resources"`
}

type QueryPolicyByActionResponse struct {
	BaseResponse `json:",inline"`
	Data         []QueryPolicyByActionResult `json:"data"`
}

type QueryPolicyByActionResult struct {
	ActionID  string            `json:"action_id"`
	Condition *PolicyExpression `json:"condition"`
}
