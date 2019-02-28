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

package mongo

import (
	"encoding/json"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql"
)

type FieldItem struct {
	Key string
	Val mapstr.MapStr
}

func (k *FieldItem) ToSQL() (string, error) {
	sql, err := json.Marshal(*k)
	return string(sql), err
}

func (k *FieldItem) legal() bool {
	if 0 == len(k.Key) {
		return false
	}
	//TODO:any other illegal case

	return true
}

func (k *FieldItem) ToMapStr() mapstr.MapStr {
	rst := mapstr.New()

	if !k.legal() {
		//drop it
		return rst
	}

	rst[k.Key] = k.Val
	return rst
}

//Comparision operator start

//Field create a new field
func Field(k string) *FieldItem {
	//check legality for k in func (k *FieldItem) ToMapStr()
	//if k is illegal, the field should be throw away
	return &FieldItem{Key: k, Val: mapstr.New()}
}

//Eq add an element like { <field> : { $eq: <val> } } for the field
func (k *FieldItem) Eq(val interface{}) *FieldItem {
	k.Val[universalsql.EQ] = val
	return k
}

//Neq add an element like { <field> : { $neq: <val> } } for the field
func (k *FieldItem) Neq(val interface{}) *FieldItem {
	k.Val[universalsql.NEQ] = val
	return k
}

//Gt add an element like { <field>: { $gt: <val> } } for the field
func (k *FieldItem) Gt(val interface{}) *FieldItem {
	k.Val[universalsql.GT] = val
	return k
}

// Regex add an element like { <field>: {$regex:<val>} } for the field
func (k *FieldItem) Regex(val interface{}) *FieldItem {
	k.Val[universalsql.REGEX] = val
	return k
}

//Gte add an element like { <field>: { $gte: <val> } } for the field
func (k *FieldItem) Gte(val interface{}) *FieldItem {
	k.Val[universalsql.GTE] = val
	return k
}

//Lt add an element like { <field>: { $lt: <val> } } for the field
func (k *FieldItem) Lt(val interface{}) *FieldItem {
	k.Val[universalsql.LT] = val
	return k
}

//Lte add an element like { <field>: { $lte: <val> } } for the field
func (k *FieldItem) Lte(val interface{}) *FieldItem {
	k.Val[universalsql.LTE] = val
	return k
}

//In add an element like { <field>: { $in: [ <val1>, <val2>,...<valn> ] } } for the field
func (k *FieldItem) In(val interface{}) *FieldItem {
	k.Val[universalsql.IN] = val
	return k
}

//Nin add an element like { <field>: { $nin: [ <val1>, <val2>,...<valn> ] } } for the field
func (k *FieldItem) Nin(val interface{}) *FieldItem {
	k.Val[universalsql.NIN] = val
	return k
}

//Comparision operator end
//Elements operator start
//Exists add an element like { <field>: { $exists: bool } } for the field
func (k *FieldItem) Exists(val bool) *FieldItem {
	k.Val[universalsql.EXISTS] = val
	return k
}

func (k *FieldItem) Type(val interface{}) *FieldItem {
	//TODO:type is not safe
	return k
}

//Elements operator end
//Array operator start
//All add an element like { <field>: { $all: [ <value1> , <value2> ... ] } } for the field
func (k *FieldItem) All(val interface{}) *FieldItem {
	k.Val[universalsql.ALL] = val
	return k
}
func (k *FieldItem) ElemMatch() *FieldItem {
	//TODO:too complicated
	return k
}

//Size add an element like { <field>: { $size: value } } for the field
//Size matches any array with the number of elements specified by the argument.
func (k *FieldItem) Size(val int) *FieldItem {
	k.Val[universalsql.SIZE] = val
	return k
}

//Array operator end
