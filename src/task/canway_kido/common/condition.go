/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except 
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and 
 * limitations under the License.
 */
 
package common

import (
	"configcenter/src/framework/core/types"
)

// CreateCondition create a condition object
func CreateCondition() Condition {
	return &condition{}
}

// Condition condition interface
type Condition interface {
	SetStart(start int)
	GetStart() int
	SetLimit(limit int)
	GetLimit() int
	SetSort(sort string)
	GetSort() string
	Field(fieldName string) Field
	ToMapStr() types.MapStr
}

// Condition the condition definition
type condition struct {
	start  int
	limit  int
	sort   string
	fields []Field
}

// SetStart set the start
func (cli *condition) SetStart(start int) {
	cli.start = start
}

// GetStart return the start
func (cli *condition) GetStart() int {
	return cli.start
}

// SetLimit set the limit num
func (cli *condition) SetLimit(limit int) {
	cli.limit = limit
}

// GetLimit return the limit num
func (cli *condition) GetLimit() int {
	return cli.limit
}

// SetSort set the sort field
func (cli *condition) SetSort(sort string) {
	cli.sort = sort
}

// GetSort return the sort field
func (cli *condition) GetSort() string {
	return cli.sort
}

// CreateField create a field
func (cli *condition) Field(fieldName string) Field {
	field := &field{
		fieldName: fieldName,
		condition: cli,
	}
	cli.fields = append(cli.fields, field)
	return field
}

// ToMapStr to MapStr object
func (cli *condition) ToMapStr() types.MapStr {
	tmpResult := types.MapStr{}
	for _, item := range cli.fields {
		tmpResult.Merge(item.ToMapStr())
	}
	return tmpResult
}
