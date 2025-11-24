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

// Package table ...
package table

import (
	"fmt"
	"maps"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

var staticRegistry = make(map[types.Name]any)

// Validate whether the table name is valid or not.
func Validate(n types.Name) error {
	// 1. check static registry
	_, valid := staticRegistry[n]
	if valid {
		return nil
	}
	// 2. check dynamic registry
	_, valid = structs.GetBuilder(string(n))
	if valid {
		return nil
	}

	return fmt.Errorf("table name unregistered: %s", n)
}

// Register for static table
func Register(tableStruct types.Tabler) {
	staticRegistry[types.Name(tableStruct.TableName())] = tableStruct
	tableFields.loadTableStruct(types.Name(tableStruct.TableName()), tableStruct)
}

// RegisterWithName for static table with differ table name
func RegisterWithName(name types.Name, tableStruct types.Tabler) {
	staticRegistry[name] = tableStruct
	tableFields.loadTableStruct(name, tableStruct)
}

// GetAllStaticTables get all static tables
func GetAllStaticTables() map[types.Name]any {
	return maps.Clone(staticRegistry)
}
