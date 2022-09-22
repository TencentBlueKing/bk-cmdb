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

// Package operator TODO
package operator

import (
	"errors"
	"reflect"
	"strings"
)

var factory map[string]Operator

func init() {
	factory = make(map[string]Operator)

	equal := EqualOper("")
	factory[equal.Name()] = &equal

	notEqual := NotEqualOper("")
	factory[notEqual.Name()] = &notEqual

	in := InOper("")
	factory[in.Name()] = &in

	notIn := NotInOper("")
	factory[notIn.Name()] = &notIn

	contains := ContainsOper("")
	factory[contains.Name()] = &contains

	notContains := NotContainsOper("")
	factory[notContains.Name()] = &notContains

	startWith := StartsWithOper("")
	factory[startWith.Name()] = &startWith

	notStartWith := NotStartsWithOper("")
	factory[notStartWith.Name()] = &notStartWith

	endWith := EndsWithOper("")
	factory[endWith.Name()] = &endWith

	notEndWith := NotEndsWithOper("")
	factory[notEndWith.Name()] = &notEndWith

	lessThan := LessThanOper("")
	factory[lessThan.Name()] = &lessThan

	lessThanEqual := LessThanEqualOper("")
	factory[lessThanEqual.Name()] = &lessThanEqual

	greaterThan := GreaterThanOper("")
	factory[greaterThan.Name()] = &greaterThan

	greaterThanEqual := GreaterThanEqualOper("")
	factory[greaterThanEqual.Name()] = &greaterThanEqual

	any := AnyOper("")
	factory[any.Name()] = &any

}

// Operator TODO
type Operator interface {
	// Name of the operator
	Name() string

	// Match is used to check if "match" is "logical equal" to the "with"
	// with different OperType, different OperType has different definition
	// of "logical equal", if "logical equal" then return bool "true" value.

	// Match TODO
	// match: the value to test
	// with: the value to compare to, which is also the template
	Match(match interface{}, with interface{}) (bool, error)
}

const (
	// Unknown TODO
	Unknown = "unknown"
	// Equal TODO
	Equal = "eq"
	// NEqual TODO
	NEqual = "not_eq"
	// Any TODO
	Any = "any"
	// In TODO
	In = "in"
	// Nin TODO
	Nin = "not_in"
	// Contains TODO
	Contains = "contains"
	// NContains TODO
	NContains = "not_contains"
	// StartWith TODO
	StartWith = "starts_with"
	// NStartWith TODO
	NStartWith = "not_starts_with"
	// EndWith TODO
	EndWith = "ends_with"
	// NEndWith TODO
	NEndWith = "not_ends_with"
	// LessThan TODO
	LessThan = "lt"
	// LessThanEqual TODO
	LessThanEqual = "lte"
	// GreaterThan TODO
	GreaterThan = "gt"
	// GreaterThanEqual TODO
	GreaterThanEqual = "gte"
)

// OperType TODO
type OperType string

// Operator TODO
func (o *OperType) Operator() Operator {
	if o == nil {
		unknown := UnknownOper("")
		return &unknown
	}

	oper, support := factory[string(*o)]
	if !support {
		unknown := UnknownOper("")
		return &unknown
	}

	return oper
}

// UnknownOper TODO
type UnknownOper OperType

// Name TODO
func (u *UnknownOper) Name() string {
	return Unknown
}

// Match TODO
func (u *UnknownOper) Match(_ interface{}, _ interface{}) (bool, error) {
	return false, errors.New("unknown type, can not do match")
}

// EqualOper TODO
type EqualOper OperType

// Name TODO
func (e *EqualOper) Name() string {
	return Equal
}

// Match TODO
func (e *EqualOper) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return reflect.DeepEqual(match, with), nil
}

// NotEqualOper TODO
type NotEqualOper OperType

// Name TODO
func (e *NotEqualOper) Name() string {
	return NEqual
}

// Match TODO
func (e *NotEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return !reflect.DeepEqual(match, with), nil
}

// InOper TODO
type InOper OperType

// Name TODO
func (e *InOper) Name() string {
	return In
}

// Match TODO
func (e *InOper) Match(match interface{}, with interface{}) (bool, error) {
	if match == nil || with == nil {
		return false, errors.New("invalid parameter")
	}

	if !reflect.ValueOf(match).IsValid() || !reflect.ValueOf(with).IsValid() {
		return false, errors.New("invalid parameter value")
	}

	mKind := reflect.TypeOf(match).Kind()
	if mKind == reflect.Slice || mKind == reflect.Array {
		return false, errors.New("invalid type, can not be array or slice")
	}

	wKind := reflect.TypeOf(with).Kind()
	if !(wKind == reflect.Slice || wKind == reflect.Array) {
		return false, errors.New("invalid type, should be array or slice")
	}

	// compare string if it's can
	if m, ok := match.(string); ok {
		valWith := reflect.ValueOf(with)
		for i := 0; i < valWith.Len(); i++ {
			v, ok := valWith.Index(i).Interface().(string)
			if !ok {
				return false, errors.New("unsupported compare with type")
			}
			if m == v {
				return true, nil
			}
		}
		return false, nil
	}

	// compare bool if it's can
	if m, ok := match.(bool); ok {
		valWith := reflect.ValueOf(with)
		for i := 0; i < valWith.Len(); i++ {
			v, ok := valWith.Index(i).Interface().(bool)
			if !ok {
				return false, errors.New("unsupported compare with type")
			}
			if m == v {
				return true, nil
			}
		}
		return false, nil
	}

	// compare numeric value if it's can
	if !isNumeric(match) {
		return false, errors.New("unsupported compare type")
	}

	// with value is slice or array, so we need to compare it one by one.
	hit := false
	valWith := reflect.ValueOf(with)

	for i := 0; i < valWith.Len(); i++ {
		if !isNumeric(valWith.Index(i).Interface()) {
			return false, errors.New("unsupported compare with type")
		}
		if toFloat64(match) == toFloat64(valWith.Index(i).Interface()) {
			hit = true
			break
		}
	}

	return hit, nil

}

// NotInOper TODO
type NotInOper OperType

// Name TODO
func (n *NotInOper) Name() string {
	return Nin
}

// Match TODO
func (n *NotInOper) Match(match interface{}, with interface{}) (bool, error) {
	inOper := InOper("in")
	hit, err := inOper.Match(match, with)
	if err != nil {
		return false, err
	}

	return !hit, nil
}

// ContainsOper TODO
type ContainsOper OperType

// Name TODO
func (c *ContainsOper) Name() string {
	return Contains
}

// Match TODO
func (c *ContainsOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return strings.Contains(m, w), nil
}

// NotContainsOper TODO
type NotContainsOper OperType

// Name TODO
func (c *NotContainsOper) Name() string {
	return NContains
}

// Match TODO
func (c *NotContainsOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return !strings.Contains(m, w), nil
}

// StartsWithOper TODO
type StartsWithOper OperType

// Name TODO
func (s *StartsWithOper) Name() string {
	return StartWith
}

// Match TODO
func (s *StartsWithOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return strings.HasPrefix(m, w), nil
}

// NotStartsWithOper TODO
type NotStartsWithOper OperType

// Name TODO
func (n *NotStartsWithOper) Name() string {
	return NStartWith
}

// Match TODO
func (n *NotStartsWithOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return !strings.HasPrefix(m, w), nil
}

// EndsWithOper TODO
type EndsWithOper OperType

// Name TODO
func (e *EndsWithOper) Name() string {
	return EndWith
}

// Match TODO
func (e *EndsWithOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return strings.HasSuffix(m, w), nil
}

// NotEndsWithOper TODO
type NotEndsWithOper OperType

// Name TODO
func (e *NotEndsWithOper) Name() string {
	return NEndWith
}

// Match TODO
func (e *NotEndsWithOper) Match(match interface{}, with interface{}) (bool, error) {
	m, ok := match.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	w, ok := with.(string)
	if !ok {
		return false, errors.New("invalid parameter")
	}

	return !strings.HasSuffix(m, w), nil
}

// LessThanOper TODO
type LessThanOper OperType

// Name TODO
func (l *LessThanOper) Name() string {
	return LessThan
}

// Match TODO
func (l *LessThanOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) < toFloat64(with), nil
}

// LessThanEqualOper TODO
type LessThanEqualOper OperType

// Name TODO
func (l *LessThanEqualOper) Name() string {
	return LessThanEqual
}

// Match TODO
func (l *LessThanEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) <= toFloat64(with), nil
}

// GreaterThanOper TODO
type GreaterThanOper OperType

// Name TODO
func (gt *GreaterThanOper) Name() string {
	return GreaterThan
}

// Match TODO
func (gt *GreaterThanOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

// GreaterThanEqualOper TODO
type GreaterThanEqualOper OperType

// Name TODO
func (gte *GreaterThanEqualOper) Name() string {
	return GreaterThanEqual
}

// Match TODO
func (gte *GreaterThanEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

// AnyOper TODO
type AnyOper OperType

// Name TODO
func (a *AnyOper) Name() string {
	return Any
}

// Match TODO
func (a *AnyOper) Match(match interface{}, _ interface{}) (bool, error) {
	return true, nil
}
