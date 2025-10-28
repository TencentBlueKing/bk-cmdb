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

// Package filter defines the filter expression
package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/TencentBlueKing/bk-cmdb/pkg/util"
)

const (
	// DefaultMaxInLimit defines the default max in limit
	DefaultMaxInLimit = uint(500)
	// DefaultMaxNotInLimit defines the default max nin limit
	DefaultMaxNotInLimit = uint(500)
	// DefaultMaxRuleLimit defines the default max number of rules limit
	DefaultMaxRuleLimit = uint(10)
	// DefaultMaxDepth defines the default max depth of the expression
	DefaultMaxDepth = uint(5)
	// DefaultMaxArrayElemLimit defines the default max element for array operator
	DefaultMaxArrayElemLimit = uint(500)
)

// ExprOption defines how to validate an
// expression.
type ExprOption struct {
	// RuleFields:
	// 1. used to test if all the expression rule's field
	//    is in the RuleFields' key restricts.
	// 2. all the expression's rule field should be a sub-set
	//    of the RuleFields' key.
	// 3. to support json field wildcard, use single '*' as suffix,
	//    for instance, 'tag.*' matches fields starts with 'tag.',
	//    only last segment of a dot-separate field can be '*'.
	RuleFields map[string]FieldType
	// MaxInLimit defines the max element of the in operator
	// If not set, then use default value: DefaultMaxInLimit
	MaxInLimit uint
	// MaxNotInLimit defines the max element of the nin operator
	// If not set, then use default value: DefaultMaxNotInLimit
	MaxNotInLimit uint
	// MaxRulesLimit defines the max number of rules an expression allows.
	// If not set, then use default value: DefaultMaxRuleLimit
	MaxRulesLimit uint
	// MaxArrayElemLimit defines the max element of the array operator
	// If not set, then use default value: DefaultMaxInLimit
	MaxArrayElemLimit uint
	// MaxDepth defines the max depth of whole expression tree.
	// If not set, then use default value: DefaultMaxDepth
	MaxDepth *uint
	// curDepth stores the current depth of the expression tree, used internally for depth validation.
	curDepth uint
}

// ExprOptionFunc expr option func defines.
type ExprOptionFunc func(opt *ExprOption)

// RuleFields set rule fields func.
func RuleFields(fields map[string]FieldType) ExprOptionFunc {
	return func(opt *ExprOption) {
		opt.RuleFields = fields
	}
}

// MaxInLimit set max in limit func.
func MaxInLimit(limit uint) ExprOptionFunc {
	return func(opt *ExprOption) {
		opt.MaxInLimit = limit
	}
}

// MaxNotInLimit set max not in limit func.
func MaxNotInLimit(limit uint) ExprOptionFunc {
	return func(opt *ExprOption) {
		opt.MaxNotInLimit = limit
	}
}

// MaxRulesLimit set max rule limit func.
func MaxRulesLimit(limit uint) ExprOptionFunc {
	return func(opt *ExprOption) {
		opt.MaxRulesLimit = limit
	}
}

// MaxDepth set max depth func.
func MaxDepth(depth uint) ExprOptionFunc {
	return func(opt *ExprOption) {
		opt.MaxDepth = &depth
	}
}

// NewExprOption new expr option.
// ExprOptionFunc: RuleFields、MaxInLimit、MaxNotInLimit、MaxRulesLimit
func NewExprOption(opts ...ExprOptionFunc) *ExprOption {
	exprOpt := new(ExprOption)
	exprOpt.RuleFields = make(map[string]FieldType)
	for _, opt := range opts {
		opt(exprOpt)
	}

	return exprOpt
}

// Expression is to build a query expression
type Expression struct {
	Op    LogicOperator `json:"op"`
	Rules []RuleFactory `json:"rules"`
}

// Validate the expression is valid or not.
func (exp *Expression) Validate(opt *ExprOption) (hitErr error) {
	defer func() {
		if hitErr != nil {
			hitErr = fmt.Errorf("expression validate failed: %w", hitErr)
		}
	}()

	if exp.IsEmpty() {
		return fmt.Errorf("expression rules is empty")
	}

	if err := exp.Op.Validate(); err != nil {
		return err
	}

	if len(exp.Rules) == 0 {
		return nil
	}

	maxRules := DefaultMaxRuleLimit
	maxDepth := DefaultMaxDepth

	if opt != nil {
		if opt.MaxRulesLimit > 0 {
			maxRules = opt.MaxRulesLimit
		}
		if opt.MaxDepth != nil {
			maxDepth = *opt.MaxDepth
		}
	} else {
		opt = NewExprOption()
	}

	opt.curDepth++
	defer func() {
		opt.curDepth--
	}()

	if opt.curDepth > maxDepth {
		return fmt.Errorf("expression depth exceeded, please reduce the depth less than %d", maxDepth)
	}

	if len(exp.Rules) > int(maxRules) {
		return fmt.Errorf("rules elements number exceeded, it at most have %d rules", maxRules)
	}

	for _, r := range exp.Rules {
		switch r.WithType() {
		case AtomType:
		case ExpressionType:
		default:
			return fmt.Errorf("unknown rule type: %s", r.WithType())
		}
	}

	for _, one := range exp.Rules {
		if err := one.Validate(opt); err != nil {
			return err
		}
	}

	return nil
}

// IsEmpty when rules is empty or filter is null
func (exp *Expression) IsEmpty() bool {
	return exp == nil || (exp.Op != All && len(exp.Rules) == 0)
}

// WithType return this expression rule's tye.
func (exp *Expression) WithType() RuleType {
	return ExpressionType
}

// UnmarshalJSON unmarshal a json raw to this expression
func (exp *Expression) UnmarshalJSON(raw []byte) error {
	parsed := gjson.GetManyBytes(raw, "op", "rules")
	op := LogicOperator(parsed[0].String())
	rules := parsed[1]
	rules.Raw = strings.TrimSpace(rules.Raw)

	if len(op) == 0 {
		// both op and raw is empty, then it's an empty expression json.
		if len(rules.Raw) == 0 {
			return nil
		}

		return errors.New("invalid expression, operator field is empty, but have none empty rules")
	}

	exp.Op = op
	if err := op.Validate(); err != nil {
		return err
	}

	if rules.Raw == "null" {
		return nil
	}

	if !rules.IsArray() {
		return errors.New("rules should be an array")
	}

	if rules.Raw == "[]" {
		return nil
	}

	if strings.TrimSpace(rules.Raw) == "[]" {
		return nil
	}

	for _, value := range rules.Array() {
		if isAtomType(value) {
			atom := new(AtomRule)
			if err := json.Unmarshal([]byte(value.Raw), atom); err != nil {
				return err
			}

			exp.Rules = append(exp.Rules, atom)
			continue
		}

		if isExpressionType(value) {
			expr := new(Expression)
			if err := json.Unmarshal([]byte(value.Raw), &expr); err != nil {
				return err
			}

			exp.Rules = append(exp.Rules, expr)
			continue
		}

		return fmt.Errorf("unknown expression rule type: %s", value.Raw)
	}

	return nil
}

func isAtomType(value gjson.Result) bool {
	parsed := gjson.GetMany(value.Raw, "field", "op", "value")
	if !parsed[0].Exists() || !parsed[1].Exists() || !parsed[2].Exists() {
		return false
	}

	return true
}

func isExpressionType(value gjson.Result) bool {
	parsed := gjson.GetMany(value.Raw, "op", "rules")
	if !parsed[0].Exists() || !parsed[1].Exists() {
		return false
	}

	return true
}

// RuleFactory defines an expression's basic rule.
// which is used to filter the resources.
type RuleFactory interface {
	// WithType get a rule's type
	WithType() RuleType
	// Validate this rule is valid or not
	Validate(opt *ExprOption) error
}

var _ RuleFactory = new(AtomRule)

var _ RuleFactory = new(Expression)

// AtomRule is the basic query rule.
type AtomRule struct {
	Field string    `json:"field"`
	Op    OpFactory `json:"op"`
	Value any       `json:"value"`
}

// WithType return this atom rule's tye.
func (ar *AtomRule) WithType() RuleType {
	return AtomType
}

// Validate this atom rule is valid or not
// Note: opt can be nil, check it before use it.
func (ar *AtomRule) Validate(opt *ExprOption) error {
	if len(ar.Field) == 0 {
		return errors.New("field is empty")
	}

	// validate operator
	if err := ar.Op.Validate(); err != nil {
		return err
	}

	if ar.Value == nil {
		return errors.New("rule value can not be nil")
	}

	rValue := reflect.ValueOf(ar.Value)
	if opt != nil {
		typ, exist := opt.RuleFields[ar.Field]
		if !exist {
			// try match wildcard field again
			typ, exist = ar.getWildcardFieldType(opt)
		}

		if !exist {
			return fmt.Errorf("rule field(%s) should not exist(not supported)", ar.Field)
		}
		if err := validateFieldValue(rValue, typ); err != nil {
			return fmt.Errorf("invalid %s's value, %v", ar.Field, err)
		}
	}

	// validate the operator's value
	if err := ar.Op.Operator().ValidateValue(rValue, opt); err != nil {
		return fmt.Errorf("%s validate failed, %v", ar.Field, err)
	}

	return nil
}

// getWildcardFieldType get column type by wildcard,
// only supports wildcard(single '*') in last part of a dot-separate field
func (ar *AtomRule) getWildcardFieldType(opt *ExprOption) (typ FieldType, exist bool) {
	dotIdx := strings.LastIndex(ar.Field, JSONFieldSeparator)
	if dotIdx == -1 {
		return "", false
	}
	wildcardField := ar.Field[:dotIdx+1] + WildcardPlaceholder
	typ, exist = opt.RuleFields[wildcardField]
	return typ, exist
}

func validateFieldValue(rVal reflect.Value, typ FieldType) error {
	switch rVal.Kind() {
	case reflect.Array, reflect.Slice:
		return validateSliceElements(rVal, typ)
	default:
	}
	rVal = util.UnpackAny(rVal)
	switch typ {
	case String:
		if rVal.Kind() != reflect.String {
			return errors.New("value should be a string")
		}

	case Numeric:
		if !isNumeric(rVal) {
			return errors.New("value should be a numeric")
		}

	case Boolean:
		if rVal.Kind() != reflect.Bool {
			return errors.New("value should be a boolean")
		}

	case Time:
		if _, err := parseTime(rVal); err != nil {
			return fmt.Errorf("parse as time value failed: %w", err)
		}
	case Any:
		// any字段的类型任意都行，不进行校验

	default:
		return fmt.Errorf("unsupported value type format: %s", typ)
	}

	return nil
}

func validateSliceElements(rVal reflect.Value, typ FieldType) error {
	length := rVal.Len()
	if length == 0 {
		return nil
	}

	// validate each slice's element data type
	for i := range length {
		if err := validateFieldValue(rVal.Index(i), typ); err != nil {
			return err
		}
	}

	return nil
}

type broker struct {
	Field string          `json:"field"`
	Op    OpFactory       `json:"op"`
	Value json.RawMessage `json:"value"`
}

// UnmarshalJSON unmarshal the json raw to AtomRule
func (ar *AtomRule) UnmarshalJSON(raw []byte) error {
	br := new(broker)
	err := json.Unmarshal(raw, br)
	if err != nil {
		return err
	}

	ar.Field = br.Field
	ar.Op = br.Op
	if br.Op == OpFactory(In) || br.Op == OpFactory(NotIn) {
		// in and nin operator's value should be an array.
		array := make([]any, 0)
		if err := json.Unmarshal(br.Value, &array); err != nil {
			return fmt.Errorf("unmarshal in/not_in value to []any failed, err: %v", err)
		}

		ar.Value = array

		return nil
	}

	to := new(any)
	if err := json.Unmarshal(br.Value, to); err != nil {
		return err
	}
	ar.Value = *to

	return nil
}
