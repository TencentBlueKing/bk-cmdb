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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql"
)

type element struct {
	Key string
	Val interface{}
}

type KV element

func (k *KV) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		k.Key: k.Val,
	}
}

// Comparison Operator Start

// Eq mongodb operator $eq
type Eq element

var _ universalsql.ConditionElement = (*Eq)(nil)

// ToMapStr return map[string]interface{}
func (e *Eq) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		e.Key: e.Val,
	}
}

// Neq mongodb operator $neq
type Neq element

var _ universalsql.ConditionElement = (*Neq)(nil)

// ToMapStr return the format result
func (n *Neq) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		n.Key: mapstr.MapStr{
			universalsql.NEQ: n.Val,
		},
	}
}

// Gt mongodb operator $gt
type Gt element

var _ universalsql.ConditionElement = (*Gt)(nil)

// ToMapStr return the format result
func (g *Gt) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		g.Key: mapstr.MapStr{
			universalsql.GT: g.Val,
		},
	}
}

// Lt mongodb operator $lt
type Lt element

var _ universalsql.ConditionElement = (*Lt)(nil)

// ToMapStr return the format result
func (l *Lt) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		l.Key: mapstr.MapStr{
			universalsql.LT: l.Val,
		},
	}
}

// Gte mongodb operator $gte
type Gte element

var _ universalsql.ConditionElement = (*Gte)(nil)

// ToMapStr return the format result
func (g *Gte) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		g.Key: mapstr.MapStr{
			universalsql.GTE: g.Val,
		},
	}
}

// Lte mongodb oeprator $lte
type Lte element

var _ universalsql.ConditionElement = (*Lte)(nil)

// ToMapStr return the format result
func (l *Lte) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		l.Key: mapstr.MapStr{
			universalsql.LTE: l.Val,
		},
	}
}

// In mongodb oeprator $in
type In element

var _ universalsql.ConditionElement = (*In)(nil)

// ToMapStr return the format result
func (i *In) ToMapStr() mapstr.MapStr {
	if nil == i.Val {
		i.Val = []interface{}{}
	}
	return mapstr.MapStr{
		i.Key: mapstr.MapStr{
			universalsql.IN: i.Val,
		},
	}
}

// Nin mongodb operator $nin
type Nin element

var _ universalsql.ConditionElement = (*Nin)(nil)

// ToMapStr return the format result
func (n *Nin) ToMapStr() mapstr.MapStr {
	if nil == n.Val {
		n.Val = []interface{}{}
	}
	return mapstr.MapStr{
		n.Key: mapstr.MapStr{
			universalsql.NIN: n.Val,
		},
	}
}

// Regex mongodb operator $regex
type Regex element

var _ universalsql.ConditionElement = (*Regex)(nil)

// ToMapStr return the format result
func (r *Regex) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		r.Key: mapstr.MapStr{
			universalsql.REGEX: r.Val,
		},
	}
}

// Comparison Operator End

// Exists Operator Start
type Exists element

var _ universalsql.ConditionElement = (*Exists)(nil)

// Exists mongodb operator $exists
func (e *Exists) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{
		e.Key: mapstr.MapStr{
			universalsql.EXISTS: e.Val,
		}}
}

// Elements Operator End
