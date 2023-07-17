/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package filter

import (
	"fmt"
	"reflect"

	"configcenter/pkg/filter"
)

// GenAtomFilter generate atom rule filter
func GenAtomFilter(field string, op filter.OpType, value interface{}) *filter.Expression {
	return &filter.Expression{
		RuleFactory: &filter.AtomRule{
			Field:    field,
			Operator: op.Factory(),
			Value:    value,
		},
	}
}

// And merge expressions using 'and' operation.
func And(rules ...filter.RuleFactory) (*filter.Expression, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("rules are not set")
	}

	andRules := make([]filter.RuleFactory, 0)
	for _, rule := range rules {
		if rule == nil || reflect.ValueOf(rule).IsNil() {
			continue
		}

		for expr, ok := rule.(*filter.Expression); ok; expr, ok = rule.(*filter.Expression) {
			rule = expr.RuleFactory
		}

		switch rule.WithType() {
		case filter.AtomType:
			andRules = append(andRules, rule)
		case filter.CombinedType:
			combinedRule, ok := rule.(*filter.CombinedRule)
			if !ok {
				return nil, fmt.Errorf("combined rule type is invalid")
			}
			if combinedRule.Condition == filter.And {
				andRules = append(andRules, combinedRule.Rules...)
				continue
			}
			andRules = append(andRules, combinedRule)
		default:
			return nil, fmt.Errorf("rule type %s is invalid", rule.WithType())
		}
	}

	if len(andRules) == 0 {
		return nil, fmt.Errorf("rules are all nil")
	}

	if len(andRules) == 1 {
		return &filter.Expression{
			RuleFactory: andRules[0],
		}, nil
	}

	return &filter.Expression{
		RuleFactory: &filter.CombinedRule{
			Condition: filter.And,
			Rules:     andRules,
		},
	}, nil
}
