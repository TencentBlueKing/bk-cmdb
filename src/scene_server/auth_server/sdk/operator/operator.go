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

type Operator interface {
	// name of the operator
	Name() string

	// Match is used to check if "match" is "logical equal" to the "with"
	// with different OperType, different OperType has different definition
	// of "logical equal", if "logical equal" then return bool "true" value.

	// match: the value to test
	// with: the value to compare to, which is also the template
	Match(match interface{}, with interface{}) (bool, error)
}

const (
	Unknown          = "unknown"
	Equal            = "eq"
	NEqual           = "not_eq"
	Any              = "any"
	In               = "in"
	Nin              = "not_in"
	Contains         = "contains"
	NContains        = "not_contains"
	StartWith        = "starts_with"
	NStartWith       = "not_starts_with"
	EndWith          = "ends_with"
	NEndWith         = "not_ends_with"
	LessThan         = "lt"
	LessThanEqual    = "lte"
	GreaterThan      = "gt"
	GreaterThanEqual = "gte"
)

type OperType string

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

type UnknownOper OperType

func (u *UnknownOper) Name() string {
	return Unknown
}

func (u *UnknownOper) Match(_ interface{}, _ interface{}) (bool, error) {
	return false, errors.New("unknown type, can not do match")
}

type EqualOper OperType

func (e *EqualOper) Name() string {
	return Equal
}

func (e *EqualOper) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return reflect.DeepEqual(match, with), nil
}

type NotEqualOper OperType

func (e *NotEqualOper) Name() string {
	return NEqual
}

func (e *NotEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return !reflect.DeepEqual(match, with), nil
}

type InOper OperType

func (e *InOper) Name() string {
	return In
}

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

type NotInOper OperType

func (n *NotInOper) Name() string {
	return Nin
}

func (n *NotInOper) Match(match interface{}, with interface{}) (bool, error) {
	inOper := InOper("in")
	hit, err := inOper.Match(match, with)
	if err != nil {
		return false, err
	}

	return !hit, nil
}

type ContainsOper OperType

func (c *ContainsOper) Name() string {
	return Contains
}

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

type NotContainsOper OperType

func (c *NotContainsOper) Name() string {
	return NContains
}

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

type StartsWithOper OperType

func (s *StartsWithOper) Name() string {
	return StartWith
}

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

type NotStartsWithOper OperType

func (n *NotStartsWithOper) Name() string {
	return NStartWith
}

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

type EndsWithOper OperType

func (e *EndsWithOper) Name() string {
	return EndWith
}

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

type NotEndsWithOper OperType

func (e *NotEndsWithOper) Name() string {
	return NEndWith
}

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

type LessThanOper OperType

func (l *LessThanOper) Name() string {
	return LessThan
}

func (l *LessThanOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) < toFloat64(with), nil
}

type LessThanEqualOper OperType

func (l *LessThanEqualOper) Name() string {
	return LessThanEqual
}

func (l *LessThanEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) <= toFloat64(with), nil
}

type GreaterThanOper OperType

func (gt *GreaterThanOper) Name() string {
	return GreaterThan
}

func (gt *GreaterThanOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

type GreaterThanEqualOper OperType

func (gte *GreaterThanEqualOper) Name() string {
	return GreaterThanEqual
}

func (gte *GreaterThanEqualOper) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

type AnyOper OperType

func (a *AnyOper) Name() string {
	return Any
}

func (a *AnyOper) Match(match interface{}, _ interface{}) (bool, error) {
	return true, nil
}
