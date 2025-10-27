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

package conv

import (
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
)

var optModifier = map[filter.OpType]converter{
	filter.ArrayEqual:    &ArrayValueConverter{OP: "="},
	filter.ArrayNotEqual: &ArrayValueConverter{OP: "<>"},
	filter.ArrayContains: &ArrayValueConverter{OP: "@>"},
	filter.ArraySubset:   &ArrayValueConverter{OP: "<@"},
	filter.ArrayOverlap:  &ArrayValueConverter{OP: "&&"},
	filter.ArrayIsEmpty: &RawValueWithFunction{
		OP:            "IS NULL",
		ColumnWrapper: &functionWrapper{Name: "array_length", ExtraArgs: []any{1}},
		ValueWrapper:  &emptyWrapper{}},
	filter.ArrayNotEmpty: &RawValueWithFunction{
		OP:            "IS NOT NULL",
		ColumnWrapper: &functionWrapper{Name: "array_length", ExtraArgs: []any{1}},
		ValueWrapper:  &emptyWrapper{}},
}

// arrayRuleToClauseExpr convert to postgresql array expression
// ref: https://www.postgresql.org/docs/current/functions-array.html
func arrayRuleToClauseExpr(rule *filter.AtomRule) (clause.Expression, error) {
	opModifier, exists := optModifier[filter.OpType(rule.Op)]
	if !exists {
		return nil, fmt.Errorf("not support array OP %s", rule.Op)
	}
	return opModifier.GetExpression(rule.Field, rule.Value)
}

// SimpleExpression simple expression
type SimpleExpression struct {
	Column        any
	OP            string
	Value         any
	ColumnWrapper wrapper
	ValueWrapper  wrapper
}

// Build sql
func (e *SimpleExpression) Build(builder clause.Builder) {
	if e.ColumnWrapper != nil {
		e.ColumnWrapper.Wrap(builder, e.Column)
	} else {
		builder.WriteQuoted(e.Column)
	}
	_ = builder.WriteByte(' ')
	_, _ = builder.WriteString(e.OP)
	_ = builder.WriteByte(' ')
	if e.ValueWrapper != nil {
		e.ValueWrapper.Wrap(builder, e.Value)
	} else {
		builder.AddVar(builder, e.Value)
	}
}

type converter interface {
	GetExpression(column, value any) (clause.Expression, error)
}

// ArrayValueConverter default array value converter
type ArrayValueConverter struct {
	OP string
}

// GetExpression build column
func (a *ArrayValueConverter) GetExpression(column, value any) (clause.Expression, error) {
	// 尝试转换value为array sql
	valueArraySQL, err := buildArraySQL(value)
	if err != nil {
		return nil, err
	}
	exp := &SimpleExpression{
		Column: column,
		OP:     a.OP,
		Value:  valueArraySQL,
	}
	return exp, nil
}

type functionWrapper struct {
	Name      string
	ExtraArgs []any
}

// Wrap wrap
func (f *functionWrapper) Wrap(builder clause.Builder, val any) {
	_, _ = builder.WriteString(f.Name)
	_ = builder.WriteByte('(')
	builder.WriteQuoted(val)
	for _, arg := range f.ExtraArgs {
		_ = builder.WriteByte(',')
		builder.AddVar(builder, arg)
	}
	_ = builder.WriteByte(')')
}

type wrapper interface {
	Wrap(builder clause.Builder, val any)
}

type emptyWrapper struct{}

// Wrap write nothing
func (e *emptyWrapper) Wrap(builder clause.Builder, val any) {
}

// RawValueWithFunction raw value function
type RawValueWithFunction struct {
	OP            string
	ColumnWrapper wrapper
	ValueWrapper  wrapper
}

// GetExpression ...
func (a *RawValueWithFunction) GetExpression(column, value any) (clause.Expression, error) {
	exp := &SimpleExpression{
		Column:        column,
		OP:            a.OP,
		Value:         value,
		ColumnWrapper: a.ColumnWrapper,
		ValueWrapper:  a.ValueWrapper,
	}
	return exp, nil
}

func buildArraySQL(value any) (sql string, err error) {
	if value == nil {
		return "", nil
	}
	var encoder = pq.Array(value)
	val, err := encoder.Value()
	if err != nil {
		return "", fmt.Errorf("encode value to array sql failed: %w", err)
	}
	if val == nil {
		// empty array
		return "{}", nil
	}

	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("valuer.Value() is not string: %T", val)
	}
	return s, nil
}
