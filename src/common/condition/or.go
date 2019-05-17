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

package condition

import (
	types "configcenter/src/common/mapstr"
)

// Field create a field
type OR interface {
	Item(val interface{}) Condition
	Array(val []interface{}) Condition
	MapStrArr(valArr []types.MapStr) Condition
	ToMapStr() types.MapStr
}

// orField the field object
type orField struct {
	condition  Condition
	fieldValue []interface{}
	fields     []Field
}

// ToMapStr conver to serch condition
func (cli *orField) ToMapStr() types.MapStr {

	tmpResult := types.MapStr{}
	for _, item := range cli.fields {
		tmpResult.Merge(item.ToMapStr())
	}

	tmpResult.Merge(types.MapStr{BKDBOR: cli.fieldValue})
	return tmpResult
}

// Item add or  query one of the conditions or
func (cli *orField) Item(val interface{}) Condition {
	cli.fieldValue = append(cli.fieldValue, val)
	return cli.condition
}

// Array add or  multiple query conditions
func (cli *orField) Array(val []interface{}) Condition {
	cli.fieldValue = append(cli.fieldValue, val...)
	return cli.condition
}

// MapStrArr  add or  multiple query conditions
func (cli *orField) MapStrArr(valArr []types.MapStr) Condition {
	for _, val := range valArr {
		cli.fieldValue = append(cli.fieldValue, val)
	}
	return cli.condition
}
