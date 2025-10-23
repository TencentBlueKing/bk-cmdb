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
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var opFactory map[OpFactory]Operator

func init() {
	opFactory = make(map[OpFactory]Operator)

	opFactory[Equal.Factory()] = EqualOp(Equal)
	opFactory[NotEqual.Factory()] = NotEqualOp(NotEqual)

	opFactory[GreaterThan.Factory()] = GreaterThanOp(GreaterThan)
	opFactory[GreaterThanEqual.Factory()] = GreaterThanEqualOp(GreaterThanEqual)

	opFactory[LessThan.Factory()] = LessThanOp(LessThan)
	opFactory[LessThanEqual.Factory()] = LessThanEqualOp(LessThanEqual)

	opFactory[In.Factory()] = InOp(In)
	opFactory[NotIn.Factory()] = NotInOp(NotIn)

	opFactory[ContainsSensitive.Factory()] = ContainsSensitiveOp(ContainsSensitive)
	opFactory[ContainsInsensitive.Factory()] = ContainsInsensitiveOp(ContainsInsensitive)

	opFactory[JSONEqual.Factory()] = JSONEqualOp(JSONEqual)
	opFactory[JSONNotEqual.Factory()] = JSONNotEqualOp(JSONEqual)
	opFactory[JSONContains.Factory()] = JSONContainsOp(JSONContains)
	opFactory[JSONHasKey.Factory()] = JSONContainsPathOp(JSONHasKey)
	opFactory[JSONNotHasKey.Factory()] = JSONNotContainsPathOp(JSONNotHasKey)

}

const (
	// And logic operator
	And LogicOperator = "and"
	// Or logic operator
	Or LogicOperator = "or"
	// All logic operator
	All LogicOperator = "all"

	// JSONFieldSeparator is the separator of json field
	JSONFieldSeparator = "."
	// WildcardPlaceholder is the wildcard char in rule field
	WildcardPlaceholder = "*"
)

// LogicOperator defines the logic operator
type LogicOperator string

// Validate the logic operator is valid or not.
func (lo LogicOperator) Validate() error {
	switch lo {
	case And:
	case Or:
	case All:
	default:
		return fmt.Errorf("unsupported expression's logic operator: %s", lo)
	}

	return nil
}

// OpFactory defines the operator's factory type.
type OpFactory string

// Operator return this operator factory's Operator
func (of OpFactory) Operator() Operator {
	op, exist := opFactory[of]
	if !exist {
		unknown := UnknownOp(Unknown)
		return &unknown
	}

	return op
}

// Validate this operator factory is valid or not.
func (of OpFactory) Validate() error {
	typ := OpType(of)
	return typ.Validate()
}

const (
	// Unknown is an unsupported operator
	Unknown OpType = "unknown"
	// Equal operator
	Equal OpType = "eq"
	// NotEqual operator
	NotEqual OpType = "neq"

	// GreaterThan operator
	GreaterThan OpType = "gt"
	// GreaterThanEqual operator
	GreaterThanEqual OpType = "gte"
	// LessThan operator
	LessThan OpType = "lt"
	// LessThanEqual operator
	LessThanEqual OpType = "lte"
	// In operator
	In OpType = "in"
	// NotIn operator
	NotIn OpType = "nin"
	// ContainsSensitive operator
	ContainsSensitive OpType = "cs"
	// ContainsInsensitive operator
	ContainsInsensitive OpType = "cis"
)

// JSONOperatorPrefix json operator prefix
const JSONOperatorPrefix = "json_"

// JSON operators, starts with json_
// reference: https://www.postgresql.org/docs/current/functions-json.html
const (

	// JSONEqual is json field equal operator.
	JSONEqual OpType = "json_eq"
	// JSONNotEqual is json field not equal operator.
	JSONNotEqual OpType = "json_neq"
	// JSONContains is json array contains value
	JSONContains OpType = "json_contains"
	// JSONHasKey is json has value or json array has key
	JSONHasKey OpType = "json_has_keys"
	// JSONNotHasKey is json not has value or json array not has key
	JSONNotHasKey OpType = "json_not_has_keys"
)

// IsJSONOperator check if op is json operator
func IsJSONOperator(op OpType) bool {
	return strings.HasPrefix(string(op), JSONOperatorPrefix)
}

// ArrayOperatorPrefix array operator prefix
const ArrayOperatorPrefix = "array_"

// Array operators
// reference: https://www.postgresql.org/docs/current/functions-array.html
const (
	// ArrayEqual equal for array data
	ArrayEqual OpType = "array_equal"
	// ArrayNotEqual not equal for array data
	ArrayNotEqual OpType = "array_neq"
	// ArrayContains A array_contains B means B is a subset of A.
	ArrayContains OpType = "array_contains"
	// ArraySubset A array_subset B means A is a subset of B.
	ArraySubset OpType = "array_subset"
	// ArrayOverlap A array_overlap B means A and B have any common element.
	ArrayOverlap OpType = "array_overlap"
)

// IsArrayOperator check if op is array operator
func IsArrayOperator(op OpType) bool {
	return strings.HasPrefix(string(op), ArrayOperatorPrefix)
}

// OpType defines the operators supported by mysql.
type OpType string

// Validate test the operator is valid or not.
func (op OpType) Validate() error {
	switch op {
	case Equal, NotEqual,
		GreaterThan, GreaterThanEqual,
		LessThan, LessThanEqual,
		In, NotIn,
		ContainsSensitive, ContainsInsensitive:

	case JSONEqual, JSONNotEqual, JSONContains,
		JSONHasKey, JSONNotHasKey:

	default:
		return fmt.Errorf("unsupported operator: %s", op)
	}

	return nil
}

// Factory return opType's factory type.
func (op OpType) Factory() OpFactory {
	return OpFactory(op)
}

// Operator is a collection of supported query operators.
type Operator interface {
	// Name is the operator's name
	Name() OpType
	// ValidateValue validate the operator's value is valid or not
	ValidateValue(rVal reflect.Value, opt *ExprOption) error
}

// UnknownOp is unknown operator
type UnknownOp OpType

// Name is equal operator
func (uo UnknownOp) Name() OpType {
	return Unknown
}

// ValidateValue validate equal's value
func (uo UnknownOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	return errors.New("unknown operator")
}

// EqualOp is equal operator type
type EqualOp OpType

// Name is equal operator
func (eo EqualOp) Name() OpType {
	return Equal
}

// ValidateValue validate equal's value
func (eo EqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isBasicValue(rVal) {
		return errors.New("invalid value field")
	}
	return nil
}

// NotEqualOp is not equal operator type
type NotEqualOp OpType

// Name is not equal operator
func (ne NotEqualOp) Name() OpType {
	return NotEqual
}

// ValidateValue validate not equal's value
func (ne NotEqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isBasicValue(rVal) {
		return errors.New("invalid ne operator's value field")
	}
	return nil
}

// GreaterThanOp is greater than operator
type GreaterThanOp OpType

// Name is greater than operator
func (gt GreaterThanOp) Name() OpType {
	return GreaterThan
}

// ValidateValue validate greater than value
func (gt GreaterThanOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isComparableValue(rVal) {
		return errors.New("invalid gt operator's value, should be a numeric or time format string value")
	}
	return nil
}

// GreaterThanEqualOp is greater than equal operator
type GreaterThanEqualOp OpType

// Name is greater than operator
func (gte GreaterThanEqualOp) Name() OpType {
	return GreaterThanEqual
}

// ValidateValue validate greater than value
func (gte GreaterThanEqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isComparableValue(rVal) {
		return errors.New("invalid gte operator's value, should be a numeric or time format string value")
	}
	return nil
}

// LessThanOp is less than operator
type LessThanOp OpType

// Name is less than equal operator
func (lt LessThanOp) Name() OpType {
	return LessThan
}

// ValidateValue validate less than equal value
func (lt LessThanOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isComparableValue(rVal) {
		return errors.New("invalid lt operator's value, should be a numeric or time format string value")
	}
	return nil
}

// LessThanEqualOp is less than equal operator
type LessThanEqualOp OpType

// Name is less than equal operator
func (lte LessThanEqualOp) Name() OpType {
	return LessThanEqual
}

// ValidateValue validate less than equal value
func (lte LessThanEqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isComparableValue(rVal) {
		return errors.New("invalid lte operator's value, should be a numeric or time format string value")
	}
	return nil
}

// InOp is in operator
type InOp OpType

// Name is in operator
func (io InOp) Name() OpType {
	return In
}

// ValidateValue validate in operator's value
func (io InOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	switch rVal.Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return errors.New("in operator's value should be an array")
	}

	length := rVal.Len()
	if length == 0 {
		return errors.New("invalid in operator's value, at least have one element")
	}

	maxInV := DefaultMaxInLimit
	if opt != nil {
		if opt.MaxInLimit > 0 {
			maxInV = opt.MaxInLimit
		}
	}

	if length > int(maxInV) {
		return fmt.Errorf("invalid in operator's value, at most have %d elements", maxInV)
	}

	// each element in the array or slice should be a basic type.
	for i := range length {
		if !isBasicValue(rVal.Index(i)) {
			return fmt.Errorf("invalid in operator's value: %v, each element's value should be a basic type",
				rVal.Index(i).Interface())
		}
	}

	return nil
}

// NotInOp is not in operator
type NotInOp OpType

// Name is not in operator
func (nio NotInOp) Name() OpType {
	return NotIn
}

// ValidateValue validate not in value
func (nio NotInOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	switch rVal.Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return errors.New("nin operator's value should be an array")
	}

	length := rVal.Len()
	if length == 0 {
		return errors.New("invalid nin operator's value, at least have one element")
	}

	maxNotInV := DefaultMaxNotInLimit
	if opt != nil {
		if opt.MaxNotInLimit > 0 {
			maxNotInV = opt.MaxNotInLimit
		}
	}

	if length > int(maxNotInV) {
		return fmt.Errorf("invalid nin operator's value, at most have %d elements", maxNotInV)
	}

	// each element in the array or slice should be a basic type.
	for i := range length {
		if !isBasicValue(rVal.Index(i)) {
			return fmt.Errorf("invalid nin operator's value: %v, each element's value should be a basic type",
				rVal.Index(i).Interface())
		}
	}

	return nil
}

// ContainsSensitiveOp is contains sensitive operator
type ContainsSensitiveOp OpType

// Name is 'like' expression with camel sensitive operator
func (cso ContainsSensitiveOp) Name() OpType {
	return ContainsSensitive
}

// ValidateValue validate 'like' operator's value
func (cso ContainsSensitiveOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if rVal.Kind() != reflect.String {
		return errors.New("cs operator's value should be an string")
	}

	if rVal.Len() == 0 {
		return errors.New("cs operator's value can not be a empty string")
	}

	return nil
}

// ContainsInsensitiveOp is contains insensitive operator
type ContainsInsensitiveOp OpType

// Name is 'like' expression with camel insensitive operator
func (cio ContainsInsensitiveOp) Name() OpType {
	return ContainsInsensitive
}

// ValidateValue validate 'like' operator's value
func (cio ContainsInsensitiveOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if rVal.Kind() != reflect.String {
		return errors.New("cis operator's value should be an string")
	}

	if rVal.Len() == 0 {
		return errors.New("cis operator's value can not be a empty string")
	}

	return nil
}

// JSONEqualOp is json field equal operator
type JSONEqualOp OpType

// Name is json field equal operator
func (op JSONEqualOp) Name() OpType {
	return JSONEqual
}

// ValidateValue validate json field equal's value
func (op JSONEqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isBasicValue(rVal) {
		return errors.New("invalid value field")
	}
	return nil
}

// JSONNotEqualOp is json field equal operator
type JSONNotEqualOp OpType

// Name is json field equal operator
func (op JSONNotEqualOp) Name() OpType {
	return JSONNotEqual
}

// ValidateValue validate json field equal's value
func (op JSONNotEqualOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isBasicValue(rVal) {
		return errors.New("invalid value field")
	}
	return nil
}

// JSONContainsOp is json array field contain operator
type JSONContainsOp OpType

// Name is json field in operator
func (op JSONContainsOp) Name() OpType {
	return JSONContains
}

// ValidateValue validate json field in's value
func (op JSONContainsOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if !isBasicValue(rVal) {
		return errors.New("invalid value field")
	}
	return nil
}

// JSONContainsPathOp is json field json contain path operator
type JSONContainsPathOp OpType

// Name is json field json contain path operator
func (op JSONContainsPathOp) Name() OpType {
	return JSONHasKey
}

// ValidateValue validate json field equal's value
func (op JSONContainsPathOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if rVal.Kind() != reflect.String {
		return errors.New("invalid value field")
	}

	return nil
}

// JSONNotContainsPathOp is json field json contain path operator
type JSONNotContainsPathOp OpType

// Name is json field json contain path operator
func (op JSONNotContainsPathOp) Name() OpType {
	return JSONNotHasKey
}

// ValidateValue validate json field equal's value
func (op JSONNotContainsPathOp) ValidateValue(rVal reflect.Value, opt *ExprOption) error {
	if rVal.Kind() != reflect.String {
		return errors.New("invalid value field")
	}

	return nil
}
