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
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
)

// IDGenerator id generator model
type IDGenerator struct {
	// Resource identify id, commonly be table name.
	// Note: the length limit of table name on PostgreSQL is 63 characters, on MySQL it is 64 characters.
	Resource types.Name `json:"resource" gorm:"resource;primaryKey;size:64"`
	MaxID    uint64     `json:"max_id" gorm:"max_id;size:64;default:0"`
}

// TableName id generator table name
func (ig IDGenerator) TableName() string {
	return IDGeneratorTable.String()
}
