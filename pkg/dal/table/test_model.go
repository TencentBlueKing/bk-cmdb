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
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/datatype"
)

// TestModel ...
type TestModel struct {
	Base     `gorm:"embedded" json:",inline"`
	Name     string                  `gorm:"column:name" json:"name,omitempty"`
	Size     int                     `gorm:"column:size" json:"size,omitempty"`
	Weight   float64                 `gorm:"column:weight" json:"weight,omitempty"`
	Int64s   datatype.Array[int64]   `gorm:"column:int64s" json:"int64s,omitempty"`
	Strings  datatype.Array[string]  `gorm:"column:strings" json:"strings,omitempty"`
	Strings2 *datatype.Array[string] `gorm:"column:strings2" json:"strings2,omitempty"`
}

// TableName ...
func (*TestModel) TableName() string {
	return TestModelTable.String()
}

func init() {
	TestModelTable.Register(&TestModel{})
}
