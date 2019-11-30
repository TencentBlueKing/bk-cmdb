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

package meta

const (
	JobIntervalMillisecond        = 100
	MinSyncIntervalMinutes        = 45
	IamPageLimit                  = 1000
	IamRequestIntervalMillisecond = 50
)

var (
	// HostResource represent host resource
	HostBizResource  = ResourceType("host")
	HostResourcePool = ResourceType("resourcePoolHost")
	// BusinessResource represent business resource
	BusinessResource = ResourceType("business")
	SetResource      = ResourceType("set")
	ModuleResource   = ResourceType("module")
	ModelResource    = ResourceType("model")
	InstanceResource = ResourceType("instance")
	// AuditCategory          = ResourceType("audit")
	ProcessResource        = ResourceType("process")
	DynamicGroupResource   = ResourceType("dynamicGroup")
	ClassificationResource = ResourceType("classification")
	UserGroupSyncResource  = ResourceType("userGroupSync")

	ServiceTemplateResource = ResourceType("serviceTemplateSync")
	PlatResource            = ResourceType("platSync")
	SetTemplateResource     = ResourceType("setTemplateSync")
)

// ResourceType represent a resource type that will be enqueue to WorkerQueue
type ResourceType string

// WorkRequest represent a task
type WorkRequest struct {
	ResourceType ResourceType
	Data         interface{}
	Header       interface{}
}

// SyncHandler is an interface implemented for sync data to iam
type SyncHandler interface {
	HandleHostSync(task *WorkRequest) error
	HandleHostResourcePoolSync(task *WorkRequest) error
	HandleBusinessSync(task *WorkRequest) error
	HandleSetSync(task *WorkRequest) error
	HandleModuleSync(task *WorkRequest) error
	HandleModelSync(task *WorkRequest) error
	HandleInstanceSync(task *WorkRequest) error
	HandleAuditSync(task *WorkRequest) error
	HandleProcessSync(task *WorkRequest) error
	HandleDynamicGroupSync(task *WorkRequest) error
	HandleClassificationSync(task *WorkRequest) error
	HandleUserGroupSync(task *WorkRequest) error
	HandleServiceTemplateSync(task *WorkRequest) error
	HandlePlatSync(task *WorkRequest) error
	HandleSetTemplateSync(task *WorkRequest) error
}
