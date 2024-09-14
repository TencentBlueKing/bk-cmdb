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

// Package types defines cmdb data syncer types
package types

import (
	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
)

// ResType is the synchronize resource type
type ResType string

const (
	// Biz is the business synchronize resource type
	Biz ResType = "biz"
	// Set is the set synchronize resource type
	Set ResType = "set"
	// Module is the module synchronize resource type
	Module ResType = "module"
	// Host is the host synchronize resource type
	Host ResType = "host"
	// HostRelation is the host relation synchronize resource type
	HostRelation ResType = "host_relation"
	// ObjectInstance is the object instance synchronize resource type
	ObjectInstance ResType = "object_instance"
	// QuotedInstance is the quoted instance synchronize resource type
	QuotedInstance ResType = "quoted_instance"
	// InstAsst is the instance association synchronize resource type
	InstAsst ResType = "inst_asst"
	// ServiceInstance is the service instance synchronize resource type
	ServiceInstance ResType = "service_instance"
	// Process is the process synchronize resource type
	Process ResType = "process"
	// ProcessRelation is the process relation synchronize resource type
	ProcessRelation ResType = "process_relation"
)

var (
	// allResType is all synchronize resource type in the order of dependency
	allResType = []ResType{Biz, ObjectInstance, Set, Module, Host, HostRelation, InstAsst, ServiceInstance, Process,
		ProcessRelation, QuotedInstance}
	allResTypeMap = make(map[ResType]struct{})
)

func init() {
	for _, resType := range allResType {
		allResTypeMap[resType] = struct{}{}
	}
}

// ListAllResType list all synchronize resource type
func ListAllResType() []ResType {
	return allResType
}

// ListAllResTypeForIncrSync list all synchronize resource type for incremental sync
func ListAllResTypeForIncrSync() []ResType {
	incrResTypes := make([]ResType, 0)
	for _, resType := range allResType {
		if resType == QuotedInstance {
			continue
		}
		incrResTypes = append(incrResTypes, resType)
	}
	return incrResTypes
}

// ResTypeWithSubResMap stores all synchronize resource type with sub resource
var ResTypeWithSubResMap = map[ResType]struct{}{
	ObjectInstance: {},
	InstAsst:       {},
	QuotedInstance: {},
}

// Validate resource type
func (r ResType) Validate(subRes string) ccErr.RawErrorInfo {
	_, exists := allResTypeMap[r]
	if !exists {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{ResTypeField},
		}
	}

	_, withSubRes := ResTypeWithSubResMap[r]
	if (subRes != "") != withSubRes {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{SubResField},
		}
	}

	return ccErr.RawErrorInfo{}
}

const (
	ResTypeField = "resource_type"
	SubResField  = "sub_resource"
)
