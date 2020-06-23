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

import "fmt"

type ResourceType string

// Decision describes the authorize decision, have already been authorized(true) or not(false)
type Decision struct {
	Authorized bool `json:"authorized"`
}

// AuthOptions describes a item to be authorized
type AuthOptions struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

type Subject struct {
	Type ResourceType `json:"type"`
	ID   string       `json:"id"`
}

// Action define's the use's action, which is must correspond to the registered action ids in iam.
type Action struct {
	ID string `json:"id"`
}

// Resource defines all the information used to authorize a resource.
type Resource struct {
	System    string             `json:"system"`
	Type      ResourceType       `json:"type"`
	ID        string             `json:"id"`
	Attribute ResourceAttributes `json:"attribute"`
}

// ResourceAttributes is the attributes of resource.
// map key: one of the attribute of this resource.
// map value: the value of this attribute for a resource instance.
// value can only be one of string, int, boolean.
// Note: _bk_iam_path_ key is a special key, which represent the resource's depended auth topology path.
// it's value's protocol should be like this: ["/biz,1/set,2/"].
type ResourceAttributes map[string]interface{}

type AuthError struct {
	// request id, parsed from iam's http response header(X-Request-Id)
	Rid     string
	Code    int64
	Message string
}

func (ae *AuthError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, rid :%s", ae.Code, ae.Message, ae.Rid)
}

type ListInstancesOptions struct {
	Type ResourceType `json:"type"`
	// only support "fetch_instance_info" value for now.
	Method string     `json:"method"`
	Filter ListFilter `json:"filter"`
}

type ListFilter struct {
	// resource instance id list
	IDList []string `json:"ids"`
	// attribute key array
	Attributes []string `json:"attrs"`
}
