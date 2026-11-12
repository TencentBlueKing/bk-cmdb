/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package types defines the authorization sdk types.
package types

import (
	"configcenter/src/scene_server/auth_server/sdk/operator"
	apigwiam "configcenter/src/thirdparty/apigw/iam"
)

// AuthorizeList Defines the list structure of authorized instance ids. If the permission type is unlimited, the
// "IsAny" field is true and the "IDS" is empty. Otherwise, the "IsAny" field is false and the "ids" is the specific
// instance ID.
type AuthorizeList struct {
	// ids with permission.
	Ids []string `json:"ids"`
	// is the permission type unrestricted.
	IsAny bool `json:"isAny"`
}

// Decision describes the authorize decision, have already been authorized(true) or not(false)
type Decision struct {
	Authorized bool `json:"authorized"`
}

// ListWithAttributes TODO
type ListWithAttributes struct {
	Operator operator.OperType `json:"op"`
	// resource instance id list, this list is not required, it also
	// one of the query filter with Operator.
	IDList       []string                  `json:"ids"`
	AttrPolicies []*operator.AuthCondition `json:"attr_policies"`
	Type         apigwiam.IamResourceType  `json:"type"`
}
