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

package enumor

// ColumnType define the db table's column data type
type ColumnType string

const (
	// Numeric means this column is Numeric data type.
	Numeric ColumnType = "numeric"
	// Boolean means this column is Boolean data type.
	Boolean ColumnType = "bool"
	// String means this column is String data type.
	String ColumnType = "string"
	// Time means this column is Time data type.
	Time ColumnType = "time"
	// Object means this column is Object data type.
	Object ColumnType = "object"
	// Array means this column is Array data type.
	Array ColumnType = "array"
)
