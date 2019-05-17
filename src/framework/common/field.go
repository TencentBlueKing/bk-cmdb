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
	cccommon "configcenter/src/common"
	"configcenter/src/framework/core/types"
)

// Field create a field
type Field interface {
	Eq(val interface{}) Condition
	NotEq(val interface{}) Condition
	Like(val interface{}) Condition
	In(val interface{}) Condition
	NotIn(val interface{}) Condition
	Lt(val interface{}) Condition
	Lte(val interface{}) Condition
	Gt(val interface{}) Condition
	Gte(val interface{}) Condition
	ToMapStr() types.MapStr
}

// Field the field object
type field struct {
	condition  Condition
	fieldName  string
	opeartor   string
	fieldValue interface{}
	fields     []Field
}

// ToMapStr conver to serch condition
func (cli *field) ToMapStr() types.MapStr {

	tmpResult := types.MapStr{}
	for _, item := range cli.fields {
		tmpResult.Merge(item.ToMapStr())
	}

	if cccommon.BKDBEQ == cli.opeartor {
		tmpResult.Merge(types.MapStr{cli.fieldName: cli.fieldValue})
	} else {
		tmpResult.Merge(types.MapStr{
			cli.fieldName: types.MapStr{
				cli.opeartor: cli.fieldValue,
			},
		})
	}

	return tmpResult
}

// Eqset a filed equal a value
func (cli *field) Eq(val interface{}) Condition {

	switch v := val.(type) {
	default:
		cli.opeartor = cccommon.BKDBEQ
		cli.fieldValue = v
		cli.fieldValue = v
	case string:
		cli.opeartor = cccommon.BKDBIN
		cli.fieldValue = []string{v}
	}

	return cli.condition
}

// NotEq set a filed equal a value
func (cli *field) NotEq(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBNE
	cli.fieldValue = val
	return cli.condition
}

// Like field like value
func (cli *field) Like(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBLIKE
	cli.fieldValue = val
	return cli.condition
}

// In in a array
func (cli *field) In(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBIN
	cli.fieldValue = val
	return cli.condition
}

// NotIn not in a array
func (cli *field) NotIn(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBNIN
	cli.fieldValue = val
	return cli.condition
}

// Lt lower than a  value
func (cli *field) Lt(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBLT
	cli.fieldValue = val
	return cli.condition
}

// Lte lower or equal than a value
func (cli *field) Lte(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBLTE
	cli.fieldValue = val
	return cli.condition
}

// Gt greater than a value
func (cli *field) Gt(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBGT
	cli.fieldValue = val
	return cli.condition
}

// Gte greater or euqal than a value
func (cli *field) Gte(val interface{}) Condition {
	cli.opeartor = cccommon.BKDBGTE
	cli.fieldValue = val
	return cli.condition
}
