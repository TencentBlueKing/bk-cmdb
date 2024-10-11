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

// Package idgen defines id generator types
package idgen

import "configcenter/src/common"

// IDGenType is the id generator type whose rule can be updated
type IDGenType string

const (
	// Biz is the business id generator type
	Biz IDGenType = "biz"
	// Set is the set id generator type
	Set IDGenType = "set"
	// Module is the module id generator type
	Module IDGenType = "module"
	// Host is the host id generator type
	Host IDGenType = "host"
	// ObjectInstance is the object instance id generator type
	ObjectInstance IDGenType = "object_instance"
	// InstAsst is the instance association id generator type
	InstAsst IDGenType = "inst_asst"
	// ServiceInstance is the service instance id generator type
	ServiceInstance IDGenType = "service_instance"
	// Process is the process id generator type
	Process IDGenType = "process"
)

// GetIDGenSequenceName get id generator sequence name by id generator type
func GetIDGenSequenceName(typ IDGenType) (string, bool) {
	sequenceName, exists := idGenTypeSeqNameMap[typ]
	return sequenceName, exists
}

// GetAllIDGenTypes get all id generator types
func GetAllIDGenTypes() []IDGenType {
	return allIDGenTypes
}

// GetAllIDGenSeqNames get all id generator sequence names
func GetAllIDGenSeqNames() []string {
	return allIDGenSeqNames
}

// IsIDGenSeqName checks if the sequence name's id generator rule can be changed
func IsIDGenSeqName(seqName string) bool {
	_, ok := seqNameMap[seqName]
	return ok
}

var idGenTypeSeqNameMap = map[IDGenType]string{
	Biz:             common.BKTableNameBaseApp,
	Set:             common.BKTableNameBaseSet,
	Module:          common.BKTableNameBaseModule,
	Host:            common.BKTableNameBaseHost,
	ObjectInstance:  common.BKTableNameBaseInst,
	InstAsst:        common.BKTableNameInstAsst,
	ServiceInstance: common.BKTableNameServiceInstance,
	Process:         common.BKTableNameBaseProcess,
}

var (
	allIDGenTypes    []IDGenType
	allIDGenSeqNames []string
	seqNameMap       = make(map[string]struct{})
)

func init() {
	for typ, seqName := range idGenTypeSeqNameMap {
		allIDGenTypes = append(allIDGenTypes, typ)
		allIDGenSeqNames = append(allIDGenSeqNames, seqName)
		seqNameMap[seqName] = struct{}{}
	}
}
