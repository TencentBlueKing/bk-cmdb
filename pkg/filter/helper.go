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

package filter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/samber/lo"
)

// isNumeric test if an interface is a numeric value or not.
func isNumeric(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		if value.Type() == reflect.TypeFor[json.Number]() {
			return true
		}
		return false
	}
}

// isComparable test if an interface is comparable or not, like int, float, string, time.Time
func isComparableValue(val reflect.Value) bool {
	if isNumeric(val) {
		return true
	}
	if _, err := parseTime(val); err == nil {
		return true
	}
	return false
}

// isBasicValue test if is the basic supported golang type or not.
func isBasicValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// parseTime test if a value can be parsed as time.Time
func parseTime(value reflect.Value) (time.Time, error) { //nolint:unparam
	if value.Kind() != reflect.String {
		return time.Time{}, fmt.Errorf("value should be a string time format")
	}
	str := value.String()
	return time.Parse(TimeStdFormat, str)
}

// ContainersExpression 生成资源字段包含的过滤条件，即fieldName in (1,2,3)
func ContainersExpression[T any](fieldName string, values []T) *Expression {
	return &Expression{
		Op: And,
		Rules: []RuleFactory{
			&AtomRule{Field: fieldName, Op: In.Factory(), Value: lo.ToAnySlice(values)},
		},
	}
}

// AllExpression 生成全量查询filter。
func AllExpression() *Expression {
	return &Expression{
		Op: All,
	}
}

// MergeWithAnd merge expressions using 'and' operation.
func MergeWithAnd(rules ...RuleFactory) (*Expression, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("rules are not set")
	}

	andRules := make([]RuleFactory, 0)
	for _, rule := range rules {
		switch rule.WithType() {
		case AtomType:
			andRules = append(andRules, rule)
		case ExpressionType:
			expr, ok := rule.(*Expression)
			if !ok {
				return nil, fmt.Errorf("rule type is not expression")
			}
			if expr.Op == And {
				andRules = append(andRules, expr.Rules...)
				continue
			}
			andRules = append(andRules, expr)
		default:
			return nil, fmt.Errorf("rule type %s is invalid", rule.WithType())
		}
	}

	return &Expression{
		Op:    And,
		Rules: andRules,
	}, nil
}

// RuleEqual 生成资源字段等于查询的AtomRule，即fieldName=value
func RuleEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: Equal.Factory(), Value: value}
}

// RuleNotEqual 生成资源字段等于查询的AtomRule，即fieldName!= value
func RuleNotEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: NotEqual.Factory(), Value: value}
}

// RuleIn 生成资源字段等于查询的AtomRule，即fieldName in values
func RuleIn[T any](fieldName string, values []T) *AtomRule {
	return &AtomRule{Field: fieldName, Op: In.Factory(), Value: values}
}

// RuleNotIn 生成资源字段等于查询的AtomRule，即fieldName nin values
func RuleNotIn[T any](fieldName string, values []T) *AtomRule {
	return &AtomRule{Field: fieldName, Op: NotIn.Factory(), Value: values}
}

// RuleCis 生成资源字段不区分大小写匹配查询的AtomRule，即LOWER(fieldName) like value
func RuleCis[T any](fieldName string, value T) *AtomRule {
	return &AtomRule{Field: fieldName, Op: ContainsInsensitive.Factory(), Value: value}
}

// RuleGreaterThan 生成资源字段大于查询的AtomRule，即fieldName > values
func RuleGreaterThan(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: GreaterThan.Factory(), Value: value}
}

// RuleLessThan 生成资源字段小于查询的AtomRule，即fieldName < values
func RuleLessThan(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: LessThan.Factory(), Value: value}
}

// RuleGreaterThanEqual 生成资源字段大于等于给定值的AtomRule，即fieldName >= values
func RuleGreaterThanEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: GreaterThanEqual.Factory(), Value: value}
}

// RuleLessThanEqual 生成资源字段小于等于给定值的AtomRule，即fieldName <= values
func RuleLessThanEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: LessThanEqual.Factory(), Value: value}
}

// RuleJSONEqual 生成资源字段等于查询的AtomRule，即fieldName=value
func RuleJSONEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: JSONEqual.Factory(), Value: value}
}

// RuleJSONNotEqual 生成资源字段等于查询的AtomRule，即fieldName!=value
func RuleJSONNotEqual(fieldName string, value any) *AtomRule {
	return &AtomRule{Field: fieldName, Op: JSONNotEqual.Factory(), Value: value}
}

// RuleJSONContains 生成资源字段等于查询的AtomRule，即values in fieldName
func RuleJSONContains[T any](fieldName string, values T) *AtomRule {
	return &AtomRule{Field: fieldName, Op: JSONContains.Factory(), Value: values}
}

// RuleJSONHasKey 生成资源字段等于查询的AtomRule，即field ? v
func RuleJSONHasKey(fieldName string, value string) *AtomRule {
	return &AtomRule{
		Field: fieldName,
		Op:    JSONHasKey.Factory(),
		Value: value,
	}
}

// ExpressionAnd expression with op and
func ExpressionAnd(rules ...RuleFactory) *Expression {
	return &Expression{
		Op:    And,
		Rules: rules,
	}
}

// ExpressionOr expression with op or
func ExpressionOr(rules ...RuleFactory) *Expression {
	return &Expression{
		Op:    Or,
		Rules: rules,
	}
}
