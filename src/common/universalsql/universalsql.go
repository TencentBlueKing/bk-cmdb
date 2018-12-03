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

package universalsql

import (
	"configcenter/src/common/mapstr"
)

type Result interface {
	ToSQL() string
	ToMapStr() mapstr.MapStr
}

// ConditionElement some operators implment this interface, like $eq $neq $nin $in etc.
type ConditionElement interface {
	ToMapStr() mapstr.MapStr
}

type Condition interface {
	Result
	Element(element ConditionElement) Condition
	And(elements ...ConditionElement) Condition
	Or(elements ...ConditionElement) Condition
	Embed(embedName string) Condition
}

type CreateStatement interface {
	Fields(fields ...Field) Result
}

type UpdateStatement interface {
	Set(fields ...Field) UpdateStatement
	Where(cond Condition) Result
}

type DeleteStatement interface {
	Where(cond Condition) Result
}

type SelectStatement interface {
	Fields(fieldName ...string) SelectStatement
	Count() SelectStatement
	Where(cond Condition) Result
}

type TableOperation interface {
	Create() CreateStatement
	Update() UpdateStatement
	Delete() DeleteStatement
	Select() SelectStatement
}


// From used to construct condition and set the table name
func From(tableName string) TableOperation {
	return nil
}

// Where used to construct condtion instance
func Where() Condition {
	return nil
}

// CreateConditionfromstruct construct a new condition instance by a struct.
// It will be ignore, if the field of the struct is not be seted any value.
// It will use the field name as the database collection field name, if the field tag is not seted
// The array field will be convert in rules.
// eg: struct { FieldName []string  `field`:"field_name"} will be converted to  {"field_name":{"$in":[]}}
func CreateConditionFromStruct(targetStruct interface{}) Condition {
	return nil
}
