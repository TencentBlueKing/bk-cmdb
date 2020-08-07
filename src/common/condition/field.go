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
	"configcenter/src/common/util"
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
	NotGt(val interface{}) Condition
	Gte(val interface{}) Condition
	Or(val interface{}) Condition
	Exists(val interface{}) Condition
	ToMapStr() types.MapStr
	GetFieldName() string
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

	switch cli.opeartor {
	case BKDBEQ:
		tmpResult.Merge(types.MapStr{cli.fieldName: cli.fieldValue})
	case BKDBOR:
		tmpResult.Merge(types.MapStr{BKDBOR: cli.fieldValue})
	default:
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
	cli.opeartor = BKDBEQ
	cli.fieldValue = val
	return cli.condition
}

// NotEq set a filed equal a value
func (cli *field) NotEq(val interface{}) Condition {
	cli.opeartor = BKDBNE
	cli.fieldValue = val
	return cli.condition
}

// Like field like value
func (cli *field) Like(val interface{}) Condition {
	cli.opeartor = "$regex"
	cli.fieldValue = val
	return cli.condition
}

// In in a array
func (cli *field) In(val interface{}) Condition {
	cli.opeartor = BKDBIN
	cli.fieldValue = util.ConverToInterfaceSlice(val)
	return cli.condition
}

// NotIn not in a array
func (cli *field) NotIn(val interface{}) Condition {
	cli.opeartor = BKDBNIN
	cli.fieldValue = val
	return cli.condition
}

// NotGt not in a array
func (cli *field) NotGt(val interface{}) Condition {
	cli.opeartor = BKDBNot
	cli.fieldValue = map[string]interface{}{
		BKDBGT: val,
	}
	return cli.condition
}

// Lt lower than a  value
func (cli *field) Lt(val interface{}) Condition {
	cli.opeartor = BKDBLT
	cli.fieldValue = val
	return cli.condition
}

// Lte lower or equal than a value
func (cli *field) Lte(val interface{}) Condition {
	cli.opeartor = BKDBLTE
	cli.fieldValue = val
	return cli.condition
}

// Gt greater than a value
func (cli *field) Gt(val interface{}) Condition {
	cli.opeartor = BKDBGT
	cli.fieldValue = val
	return cli.condition
}

// Gte greater or euqal than a value
func (cli *field) Gte(val interface{}) Condition {
	cli.opeartor = BKDBGTE
	cli.fieldValue = val
	return cli.condition
}

// Or if any expression of value are true
func (cli *field) Or(val interface{}) Condition {
	cli.opeartor = BKDBOR
	cli.fieldValue = val
	return cli.condition
}

// Exists set the key exists value as val
func (cli *field) Exists(val interface{}) Condition {
	cli.opeartor = BKDBEXISTS
	cli.fieldValue = val
	return cli.condition
}

func (cli *field) GetFieldName() string {
	return cli.fieldName
}
