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

type OperType string

type Equal OperType

func (e *Equal) Name() string {
	return "eq"
}

func (e *Equal) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return reflect.DeepEqual(match, with), nil
}

type NotEqual OperType

func (e *NotEqual) Name() string {
	return "neq"
}

func (e *NotEqual) Match(match interface{}, with interface{}) (bool, error) {
	mType := reflect.TypeOf(match)
	wType := reflect.TypeOf(with)
	if mType.Kind() != wType.Kind() {
		return false, errors.New("mismatch type")
	}

	return !reflect.DeepEqual(match, with), nil
}

type In OperType

func (e *In) Name() string {
	return "in"
}

func (e *In) Match(match interface{}, with interface{}) (bool, error) {
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
		if w, ok := with.([]string); ok {
			hit := false
			for _, to := range w {
				if to == m {
					hit = true
					break
				}
			}
			return hit, nil
		}
	}

	// compare bool if it's can
	if m, ok := match.(bool); ok {
		if w, ok := with.([]bool); ok {
			hit := false
			for _, to := range w {
				if to == m {
					hit = true
					break
				}
			}
			return hit, nil
		}
	}

	// compare numeric value if it's can
	if !isNumeric(match) {
		return false, errors.New("unsupported compare type")
	}

	// with value is slice or array, so we need to compare it one by one.
	hit := false
	valWith := reflect.ValueOf(with)
	if !isNumeric(valWith.Interface()) {
		return false, errors.New("unsupported compare with type")
	}

	for i := 0; i < valWith.Len(); i++ {
		if toFloat64(match) == toFloat64(valWith.Index(i).Interface()) {
			hit = true
		}
	}

	return hit, nil

}

type NotIn OperType

func (n *NotIn) Name() string {
	return "nin"
}

func (n *NotIn) Match(match interface{}, with interface{}) (bool, error) {
	inOper := In("in")
	hit, err := inOper.Match(match, with)
	if err != nil {
		return false, err
	}

	return !hit, nil
}

type Contains OperType

func (c *Contains) Name() string {
	return "contains"
}

func (c *Contains) Match(match interface{}, with interface{}) (bool, error) {
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

type NotContains OperType

func (c *NotContains) Name() string {
	return "not_contains"
}

func (c *NotContains) Match(match interface{}, with interface{}) (bool, error) {
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

type StartsWith OperType

func (s *StartsWith) Name() string {
	return "starts_with"
}

func (s *StartsWith) Match(match interface{}, with interface{}) (bool, error) {
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

type NotStartsWith OperType

func (n *NotStartsWith) Name() string {
	return "not_starts_with"
}

func (n *NotStartsWith) Match(match interface{}, with interface{}) (bool, error) {
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

type EndsWith OperType

func (e *EndsWith) Name() string {
	return "ends_with"
}

func (e *EndsWith) Match(match interface{}, with interface{}) (bool, error) {
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

type NotEndsWith OperType

func (e *NotEndsWith) Name() string {
	return "not_ends_with"
}

func (e *NotEndsWith) Match(match interface{}, with interface{}) (bool, error) {
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

type LessThan OperType

func (l *LessThan) Name() string {
	return "lt"
}

func (l *LessThan) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) < toFloat64(with), nil
}

type LessThanEqual OperType

func (l *LessThanEqual) Name() string {
	return "lte"
}

func (l *LessThanEqual) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) <= toFloat64(with), nil
}

type GreaterThan OperType

func (gt *GreaterThan) Name() string {
	return "gt"
}

func (gt *GreaterThan) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

type GreaterThanEqual OperType

func (gte *GreaterThanEqual) Name() string {
	return "gte"
}

func (gte *GreaterThanEqual) Match(match interface{}, with interface{}) (bool, error) {
	if !isNumeric(match) || !isNumeric(with) {
		return false, errors.New("invalid parameter")
	}

	return toFloat64(match) > toFloat64(with), nil
}

type Any OperType

func (a *Any) Name() string {
	return "any"
}

func (a *Any) Match(match interface{}, with interface{}) (bool, error) {
	if !reflect.ValueOf(match).IsValid() {
		return false, errors.New("invalid parameter")
	}
	return true, nil
}
