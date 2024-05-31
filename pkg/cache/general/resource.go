/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package general

import (
	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
)

// ResType is the resource type for general resource cache
type ResType string

const (
	// Host is the resource type for host cache
	Host ResType = "host"
	// ModuleHostRel is the resource type for host relation cache
	ModuleHostRel ResType = "host_relation"
	// Biz is the resource type for business cache
	Biz ResType = "biz"
	// Set is the resource type for set cache
	Set ResType = "set"
	// Module is the resource type for module cache
	Module ResType = "module"
	// Process is the resource type for process cache
	Process ResType = "process"
	// ProcessRelation is the resource type for process instance relation cache
	ProcessRelation ResType = "process_relation"
	// BizSet is the resource type for  cache
	BizSet ResType = "biz_set"
	// Plat is the resource type for cloud area cache
	Plat ResType = "plat"
	// Project is the resource type for project cache
	Project ResType = "project"
	// ObjectInstance is the resource type for common object instance cache, its sub resource specifies the object id
	ObjectInstance ResType = "object_instance"
	// MainlineInstance is the resource type for mainline instance cache, its sub resource specifies the object id
	MainlineInstance ResType = "mainline_instance"
	// InstAsst is the resource type for instance association cache, its sub resource specifies the associated object id
	InstAsst ResType = "inst_asst"
	// KubeCluster is the resource type for kube cluster cache
	KubeCluster ResType = "kube_cluster"
	// KubeNode is the resource type for kube node cache
	KubeNode ResType = "kube_node"
	// KubeNamespace is the resource type for kube namespace cache
	KubeNamespace ResType = "kube_namespace"
	// KubeWorkload is the resource type for kube workload cache,  its sub resource specifies the workload type
	KubeWorkload ResType = "kube_workload"
	// KubePod is the resource type for kube pod cache, its event detail is pod info with containers in it
	KubePod ResType = "kube_pod"
)

// SupportedResTypeMap is a map whose key is resource type that is supported by general resource cache
// not all resource types are supported now, add related logics if other resource type needs cache.
var SupportedResTypeMap = map[ResType]struct{}{
	Host:             {},
	Biz:              {},
	Set:              {},
	Module:           {},
	BizSet:           {},
	Plat:             {},
	ObjectInstance:   {},
	MainlineInstance: {},
}

// ResTypeHasSubResMap is a map of supported resource type -> whether it has sub resource
var ResTypeHasSubResMap = map[ResType]struct{}{
	ObjectInstance:   {},
	MainlineInstance: {},
	InstAsst:         {},
	KubeWorkload:     {},
}

// ValidateWithSubRes validate ResType with sub resource
func (r ResType) ValidateWithSubRes(subRes string) ccErr.RawErrorInfo {
	_, exists := SupportedResTypeMap[r]
	if !exists {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{ResourceField},
		}
	}

	_, hasSubRes := ResTypeHasSubResMap[r]
	if (subRes != "") != hasSubRes {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{SubResField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ResTypeNeedOidMap is a map whose key is resource type that needs oid to generate id key
var ResTypeNeedOidMap = map[ResType]struct{}{
	ModuleHostRel:   {},
	ProcessRelation: {},
}
