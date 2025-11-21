/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package meta defines the metadata for authorization.
package meta

// ResourceType is the resource type for authorization.
type ResourceType string

// Action is the action for authorization.
type Action string

// Basic defines the basic info of an auth resource.
type Basic struct {
	// Type is the resource type.
	Type ResourceType `json:"type"`
	// Action is the action that user want to perform on this resource.
	Action Action `json:"action"`
	// Name is the resource name.
	Name string `json:"name"`
	// ID is the resource id.
	ID string `json:"id"`
}

// ResourceAttribute represents one iam resource.
type ResourceAttribute struct {
	*Basic
	// Layers defines the topology layer items before this object in a topology.
	Layers []Basic `json:"layers"`
}

// Decision is the authorize decision.
type Decision struct {
	// Authorized defines whether the user has the permission to the resource.
	Authorized bool `json:"decision"`
}

// ListAuthResOptions is the list authorized resource options.
type ListAuthResOptions struct {
	ResourceType ResourceType `json:"resource_type"`
	Action       Action       `json:"action"`
}

// AuthResInfo is the authorized resource info.
type AuthResInfo struct {
	// IDs is the resource ids that user has permission to.
	IDs []string `json:"ids"`
	// IsAny defines whether the user has the permission to all resources.
	IsAny bool `json:"is_any"`
}
