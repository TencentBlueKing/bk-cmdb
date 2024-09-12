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
	// InstAsst is the instance association synchronize resource type
	InstAsst ResType = "inst_asst"
)

// AllResTypeMap stores all synchronize resource type
var AllResTypeMap = map[ResType]struct{}{
	Biz:            {},
	Set:            {},
	Module:         {},
	Host:           {},
	HostRelation:   {},
	ObjectInstance: {},
	InstAsst:       {},
}

// ResTypeWithSubResMap stores all synchronize resource type with sub resource
var ResTypeWithSubResMap = map[ResType]struct{}{
	ObjectInstance: {},
	InstAsst:       {},
}

// Validate resource type
func (r ResType) Validate(subRes string) ccErr.RawErrorInfo {
	_, exists := AllResTypeMap[r]
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
