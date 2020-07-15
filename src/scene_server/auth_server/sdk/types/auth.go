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
	"errors"
	"fmt"

	"configcenter/src/apimachinery/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"

	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Iam     IamConfig
	Options Options
}

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

type Options struct {
	Metric prometheus.Registerer
}

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

type AuthBatchOptions struct {
	System  string       `json:"system"`
	Subject Subject      `json:"subject"`
	Batch   []*AuthBatch `json:"batch"`
}

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

type AuthBatch struct {
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

type ActionPolicy struct {
	Action Action           `json:"action"`
	Policy *operator.Policy `json:"condition"`
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

type ListWithAttributes struct {
	Operator operator.OperType `json:"op"`
	// resource instance id list, this list is not required, it also
	// one of the query filter with Operator.
	IDList     []string               `json:"ids"`
	Attributes []*operator.FieldValue `json:"attributes"`
	Type       ResourceType           `json:"type"`
}
