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

// Package filter ...
package filter

import (
	"time"
)

// FieldType define the db table's column data type
type FieldType string

const (
	// Numeric means this field is Numeric data type.
	Numeric FieldType = "numeric"
	// Boolean means this field is Boolean data type.
	Boolean FieldType = "bool"
	// String means this field is String data type.
	String FieldType = "string"
	// Time means this field is Time data type.
	Time FieldType = "time"
	// Any means this field is Any data type, will skip value type check
	Any FieldType = "any"
)

const (

	// TimeStdFormat is the system's standard time format to store or to query, equal to time.RFC3339
	TimeStdFormat = time.RFC3339
	// DateLayout is the date layout with '%Y-%m-%d'
	DateLayout = time.DateOnly
	// DateTimeLayout is the date layout with '%Y-%m-%d %H:%M:%S'
	DateTimeLayout = time.DateTime
)

// RuleType is the expression rule's rule type.
type RuleType string

const (
	// EmptyType means the rules is empty
	EmptyType RuleType = "Empty"
	// AtomType means it's a AtomRule
	AtomType RuleType = "AtomRule"
	// ExpressionType means that it may be a query expression or a sub query expression.
	ExpressionType RuleType = "Expression"
	// UnknownType means it's an unknown type.
	UnknownType RuleType = "Unknown"
)
