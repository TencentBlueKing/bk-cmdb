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

package table

import (
	"reflect"

	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

// BaseModelName defines base model name.
const BaseModelName = "base_model"

// Base defines base model.
type Base struct {
	ID string `json:"id" gorm:"column:id;size:64;primaryKey,index:,unique"`
}

// SetID sets the ID.
func (m *Base) SetID(id string) {
	m.ID = id
}

func init() {
	structs.RegisterFieldType(BaseModelName, reflect.TypeFor[Base]())
}
