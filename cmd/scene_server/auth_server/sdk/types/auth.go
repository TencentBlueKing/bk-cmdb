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
	operator2 "configcenter/cmd/scene_server/auth_server/sdk/operator"
	"errors"
	"fmt"

	"configcenter/api/util"
	"github.com/prometheus/client_golang/prometheus"
)

// Config TODO
type Config struct {
	Iam     IamConfig
	Options Options
}

// IamConfig TODO
type IamConfig struct {
	// blueking's auth center addresses
	Address []string
	// app code is used for authorize used.
	AppCode string
	// app secret is used for authorized
	AppSecret string
	// the system id which used in auth center.
	SystemID string
	// http TLS config
	TLS util.TLSClientConfig
}

// Validate TODO
func (a IamConfig) Validate() error {
	if len(a.Address) == 0 {
		return errors.New("no iam address")
	}

	if len(a.AppCode) == 0 {
		return errors.New("no iam app code")
	}

	if len(a.AppSecret) == 0 {
		return errors.New("no iam app secret")
	}
	return nil
}

// Options TODO
type Options struct {
	Metric prometheus.Registerer
}

// AuthorizeList Defines the list structure of authorized instance ids. If the permission type is unlimited, the
// "IsAny" field is true and the "IDS" is empty. Otherwise, the "IsAny" field is false and the "ids" is the specific
// instance ID.
type AuthorizeList struct {
	// ids with permission.
	Ids []string `json:"ids"`
	// is the permission type unrestricted.
	IsAny bool `json:"isAny"`
}

// ResourceType TODO
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

// Validate TODO
func (a AuthOptions) Validate() error {
	if len(a.System) == 0 {
		return errors.New("system is empty")
	}

	if len(a.Subject.Type) == 0 {
		return errors.New("subject.type is empty")
	}

	if len(a.Subject.ID) == 0 {
		return errors.New("subject.id is empty")
	}

	if len(a.Action.ID) == 0 {
		return errors.New("action.id is empty")
	}

	return nil
}

// AuthBatchOptions TODO
type AuthBatchOptions struct {
	System  string       `json:"system"`
	Subject Subject      `json:"subject"`
	Batch   []*AuthBatch `json:"batch"`
}

// Validate TODO
func (a AuthBatchOptions) Validate() error {
	if len(a.System) == 0 {
		return errors.New("system is empty")
	}

	if len(a.Subject.Type) == 0 {
		return errors.New("subject.type is empty")
	}

	if len(a.Subject.ID) == 0 {
		return errors.New("subject.id is empty")
	}

	if len(a.Batch) == 0 {
		return nil
	}

	for _, b := range a.Batch {
		if len(b.Action.ID) == 0 {
			return errors.New("empty action id")
		}
	}
	return nil
}

// AuthBatch TODO
type AuthBatch struct {
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

// Subject TODO
type Subject struct {
	Type ResourceType `json:"type"`
	ID   string       `json:"id"`
}

// Action define's the use's action, which is must correspond to the registered action ids in iam.
type Action struct {
	ID string `json:"id"`
}

// ActionPolicy TODO
type ActionPolicy struct {
	Action Action            `json:"action"`
	Policy *operator2.Policy `json:"condition"`
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

// AuthError TODO
type AuthError struct {
	// request id, parsed from iam's http response header(X-Request-Id)
	Rid     string
	Code    int64
	Message string
}

// Error 用于错误处理
func (ae *AuthError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, rid :%s", ae.Code, ae.Message, ae.Rid)
}

// ListWithAttributes TODO
type ListWithAttributes struct {
	Operator operator2.OperType `json:"op"`
	// resource instance id list, this list is not required, it also
	// one of the query filter with Operator.
	IDList       []string            `json:"ids"`
	AttrPolicies []*operator2.Policy `json:"attr_policies"`
	Type         ResourceType        `json:"type"`
}
