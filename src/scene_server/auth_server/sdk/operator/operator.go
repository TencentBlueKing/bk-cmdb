/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package operator defines the authorization plan expression operators and evaluation helpers.
package operator

import (
	"errors"
	"reflect"
	"strings"
)

var factory map[OperType]Operator

func init() {
	factory = make(map[OperType]Operator)

	equal := EqualOper(Equal)
	factory[Equal] = &equal

	in := InOper(In)
	factory[In] = &in

	startWith := StartsWithOper(StartWith)
	factory[StartWith] = &startWith
}

// Operator is used to evaluate a compare node against a resource attribute.
type Operator interface {
	// Name of the operator
	Name() OperType

	// Match is used to check if "match" is "logical equal" to the "with".
	// Different OperType has different definition of "logical equal".
	// match: the value to test
	// with: the value to compare to, which is also the template
	Match(match interface{}, with interface{}) (bool, error)
}

const (
	// And is the n-ary logical and operator, whose operands are in the content field.
	And OperType = "and"
	// Or is the n-ary logical or operator, whose operands are in the content field.
	Or OperType = "or"
	// Not is the unary logical not operator, whose content has exactly one element.
	Not OperType = "not"

	// Equal compares if the attribute equals to the value.
	Equal OperType = "eq"
	// In compares if the attribute is one of the value list.
	In OperType = "in"
	// StartWith compares if the attribute has the value prefix.
	StartWith OperType = "starts_with"
)

// OperType is the operator of the authorization plan expression.
type OperType string

// Operator returns the evaluator of this compare operator.
func (o OperType) Operator() Operator {
	oper, support := factory[o]
	if !support {
		unknown := UnknownOper("")
		return &unknown
	}

	return oper
}

// IsLogical returns if the operator is a logical one.
func (o OperType) IsLogical() bool {
	switch o {
	case And, Or, Not:
		return true
	default:
		return false
	}
}

// UnknownOper is returned when the operator is not a registered compare operator.
type UnknownOper OperType

// Name TODO
func (u *UnknownOper) Name() OperType {
	return "unknown"
}

// Match TODO
func (u *UnknownOper) Match(_ interface{}, _ interface{}) (bool, error) {
	return false, errors.New("unknown type, can not do match")
}

// EqualOper TODO
type EqualOper OperType

// Name TODO
func (e *EqualOper) Name() OperType {
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

// InOper TODO
type InOper OperType

// Name TODO
func (e *InOper) Name() OperType {
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

// StartsWithOper TODO
type StartsWithOper OperType

// Name TODO
func (s *StartsWithOper) Name() OperType {
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
