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

	"gorm.io/gorm/schema"
)

// Name of table
type Name string

const (
	// IDGeneratorTable ...
	IDGeneratorTable Name = "id_generator_test"
	// TestModelTable ...
	TestModelTable Name = "test_model"
)

// Validate whether the table name is valid or not.
func (n Name) Validate() error {
	_, valid := tableRegistry[n]
	if valid {
		return nil
	}

	return fmt.Errorf("table name is invalid: %s", n)
}

var tableRegistry = make(map[Name]any)

// Register table
func (n Name) Register(tableStruct Tabler) {
	tableRegistry[n] = tableStruct
}

// String return table name string
func (n Name) String() string {
	return string(n)
}

// Tabler have table name method
type Tabler = schema.Tabler
