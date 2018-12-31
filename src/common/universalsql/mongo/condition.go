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
	"encoding/json"
)

var _ universalsql.Condition = (*mongoCondition)(nil)

type mongoCondition struct {
	elements []universalsql.ConditionElement
	and      []universalsql.ConditionElement
	or       []universalsql.ConditionElement
	not      []universalsql.ConditionElement
	nor      []universalsql.ConditionElement
	embed    map[string]*mongoCondition // key is the embed filed name
}

func newCondition() *mongoCondition {
	return &mongoCondition{
		embed: map[string]*mongoCondition{},
	}
}

// NewCondition create a condition instance
func NewCondition() universalsql.Condition {
	return newCondition()
}

// NewConditionFromMapStr create a condition by MapStr
func NewConditionFromMapStr(inputCond mapstr.MapStr) (outputCond universalsql.Condition, err error) {

	return parseConditionFromMapStr(newCondition(), "", inputCond)
}

func (m *mongoCondition) ToSQL() (string, error) {
	sql, err := json.Marshal(m.ToMapStr())
	return string(sql), err
}

func (m *mongoCondition) ToMapStr() mapstr.MapStr {

	result := mapstr.New()

	// merge elements, default
	for _, ele := range m.elements {
		tmp := ele.ToMapStr()
		if 0 != len(tmp) {
			result.Merge(tmp)
		}
	}

	// merge and elements
	andElements := mapstr.NewArray()
	for _, ele := range m.and {
		tmp := ele.ToMapStr()
		if 0 != len(tmp) {
			andElements = append(andElements, tmp)
		}
	}
	if 0 != len(andElements) {
		result.Set(universalsql.AND, andElements)
	}

	// merge or elements
	orElements := mapstr.NewArray()
	for _, ele := range m.or {
		tmp := ele.ToMapStr()
		if 0 != len(tmp) {
			orElements = append(orElements, tmp)
		}
	}
	if 0 != len(orElements) {
		result.Set(universalsql.OR, orElements)
	}

	//merge not elements
	notElements := mapstr.NewArray()
	for _, ele := range m.not {
		tmp := ele.ToMapStr()
		if 0 != len(tmp) {
			notElements = append(notElements, tmp)
		}
	}
	if 0 != len(notElements) {
		result.Set(universalsql.NOT, notElements)
	}

	//merge nor elements
	norElements := mapstr.NewArray()
	for _, ele := range m.nor {
		tmp := ele.ToMapStr()
		if 0 != len(tmp) {
			norElements = append(norElements, tmp)
		}
	}
	if 0 != len(norElements) {
		result.Set(universalsql.NOR, norElements)
	}

	// merge embed elements
	for embedName, embedCondition := range m.embed {
		result.Set(embedName, embedCondition.ToMapStr())
	}

	return result
}

func (m *mongoCondition) Element(elements ...universalsql.ConditionElement) universalsql.Condition {
	m.elements = append(m.elements, elements...)
	return m
}

func (m *mongoCondition) And(elements ...universalsql.ConditionElement) universalsql.Condition {
	m.and = append(m.and, elements...)
	return m
}

func (m *mongoCondition) Or(elements ...universalsql.ConditionElement) universalsql.Condition {
	m.or = append(m.or, elements...)
	return m
}

func (m *mongoCondition) Not(elements ...universalsql.ConditionElement) universalsql.Condition {
	m.not = append(m.or, elements...)
	return m
}

func (m *mongoCondition) Nor(elements ...universalsql.ConditionElement) universalsql.Condition {
	m.nor = append(m.or, elements...)
	return m
}

func (m *mongoCondition) Embed(embedName string) (origin, embed universalsql.Condition) {

	tmp := newCondition()
	m.embed[embedName] = tmp
	origin = m
	embed = tmp
	return origin, embed
}

func (m *mongoCondition) merge(sourceCond *mongoCondition) *mongoCondition {

	m.and = append(m.and, sourceCond.and...)
	m.or = append(m.or, sourceCond.or...)
	m.elements = append(m.elements, sourceCond.elements...)

	for sourceEmbedName, sourceCond := range sourceCond.embed {

		targetEmbedCond, ok := m.embed[sourceEmbedName]
		if !ok {
			m.embed[sourceEmbedName] = sourceCond
			continue
		}

		targetEmbedCond = targetEmbedCond.merge(sourceCond)
	}

	return m
}
