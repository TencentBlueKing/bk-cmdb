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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
)

var opFactory map[OpFactory]Operator

func init() {
	opFactory = make(map[OpFactory]Operator)

	eq := EqualOp(Equal)
	opFactory[OpFactory(eq.Name())] = &eq
	ne := NotEqualOp(NotEqual)
	opFactory[OpFactory(ne.Name())] = &ne
	in := InOp(In)
	opFactory[OpFactory(in.Name())] = &in
	nin := NotInOp(NotIn)
	opFactory[OpFactory(nin.Name())] = &nin
	lt := LessOp(Less)
	opFactory[OpFactory(lt.Name())] = &lt
	lte := LessOrEqualOp(LessOrEqual)
	opFactory[OpFactory(lte.Name())] = &lte
	gt := GreaterOp(Greater)
	opFactory[OpFactory(gt.Name())] = &gt
	gte := GreaterOrEqualOp(GreaterOrEqual)
	opFactory[OpFactory(gte.Name())] = &gte
	datetimeLt := DatetimeLessOp(DatetimeLess)
	opFactory[OpFactory(datetimeLt.Name())] = &datetimeLt
	datetimeLte := DatetimeLessOrEqualOp(DatetimeLessOrEqual)
	opFactory[OpFactory(datetimeLte.Name())] = &datetimeLte
	datetimeGt := DatetimeGreaterOp(DatetimeGreater)
	opFactory[OpFactory(datetimeGt.Name())] = &datetimeGt
	datetimeGte := DatetimeGreaterOrEqualOp(DatetimeGreaterOrEqual)
	opFactory[OpFactory(datetimeGte.Name())] = &datetimeGte
	beginsWith := BeginsWithOp(BeginsWith)
	opFactory[OpFactory(beginsWith.Name())] = &beginsWith
	beginsWithInsensitive := BeginsWithInsensitiveOp(BeginsWithInsensitive)
	opFactory[OpFactory(beginsWithInsensitive.Name())] = &beginsWithInsensitive
	notBeginsWith := NotBeginsWithOp(NotBeginsWith)
	opFactory[OpFactory(notBeginsWith.Name())] = &notBeginsWith
	notBeginsWithInsensitive := NotBeginsWithInsensitiveOp(NotBeginsWithInsensitive)
	opFactory[OpFactory(notBeginsWithInsensitive.Name())] = &notBeginsWithInsensitive
	contains := ContainsOp(Contains)
	opFactory[OpFactory(contains.Name())] = &contains
	containsSensitive := ContainsSensitiveOp(ContainsSensitive)
	opFactory[OpFactory(containsSensitive.Name())] = &containsSensitive
	notContains := NotContainsOp(NotContains)
	opFactory[OpFactory(notContains.Name())] = &notContains
	notContainsInsensitive := NotContainsInsensitiveOp(NotContainsInsensitive)
	opFactory[OpFactory(notContainsInsensitive.Name())] = &notContainsInsensitive
	endsWith := EndsWithOp(EndsWith)
	opFactory[OpFactory(endsWith.Name())] = &endsWith
	endsWithInsensitive := EndsWithInsensitiveOp(EndsWithInsensitive)
	opFactory[OpFactory(endsWithInsensitive.Name())] = &endsWithInsensitive
	notEndsWith := NotEndsWithOp(NotEndsWith)
	opFactory[OpFactory(notEndsWith.Name())] = &notEndsWith
	notEndsWithInsensitive := NotEndsWithInsensitiveOp(NotEndsWithInsensitive)
	opFactory[OpFactory(notEndsWithInsensitive.Name())] = &notEndsWithInsensitive
	isEmpty := IsEmptyOp(IsEmpty)
	opFactory[OpFactory(isEmpty.Name())] = &isEmpty
	isNotEmpty := IsNotEmptyOp(IsNotEmpty)
	opFactory[OpFactory(isNotEmpty.Name())] = &isNotEmpty
	size := SizeOp(Size)
	opFactory[OpFactory(size.Name())] = &size
	isNull := IsNullOp(IsNull)
	opFactory[OpFactory(isNull.Name())] = &isNull
	isNotNull := IsNotNullOp(IsNotNull)
	opFactory[OpFactory(isNotNull.Name())] = &isNotNull
	exist := ExistOp(Exist)
	opFactory[OpFactory(exist.Name())] = &exist
	notExist := NotExistOp(NotExist)
	opFactory[OpFactory(notExist.Name())] = &notExist
	obj := ObjectOp(Object)
	opFactory[OpFactory(obj.Name())] = &obj
	filterArr := ArrayOp(Array)
	opFactory[OpFactory(filterArr.Name())] = &filterArr
}

const (
	// And logic operator
	And LogicOperator = "AND"
	// Or logic operator
	Or LogicOperator = "OR"
)

// LogicOperator defines the logic operator
type LogicOperator string

// Validate the logic operator is valid or not.
func (lo LogicOperator) Validate() error {
	switch lo {
	case And:
	case Or:
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

	// generic operator

	// Equal operator
	Equal OpType = "equal"
	// NotEqual operator
	NotEqual OpType = "not_equal"

	// set operator that is used to filter element using the value array

	// In operator
	In OpType = "in"
	// NotIn operator
	NotIn OpType = "not_in"

	// numeric compare operator

	// Less operator
	Less OpType = "less"
	// LessOrEqual operator
	LessOrEqual OpType = "less_or_equal"
	// Greater operator
	Greater OpType = "greater"
	// GreaterOrEqual operator
	GreaterOrEqual OpType = "greater_or_equal"

	// datetime operator, ** need to be parsed to mongo in coreservice to avoid json marshaling **

	// DatetimeLess operator
	DatetimeLess OpType = "datetime_less"
	// DatetimeLessOrEqual operator
	DatetimeLessOrEqual OpType = "datetime_less_or_equal"
	// DatetimeGreater operator
	DatetimeGreater OpType = "datetime_greater"
	// DatetimeGreaterOrEqual operator
	DatetimeGreaterOrEqual OpType = "datetime_greater_or_equal"

	// string operator

	// BeginsWith operator with case-sensitive
	BeginsWith OpType = "begins_with"
	// BeginsWithInsensitive operator with case-insensitive
	BeginsWithInsensitive OpType = "begins_with_i"
	// NotBeginsWith operator with case-sensitive
	NotBeginsWith OpType = "not_begins_with"
	// NotBeginsWithInsensitive operator with case-insensitive
	NotBeginsWithInsensitive OpType = "not_begins_with_i"
	// Contains operator with case-insensitive, compatible for the query builder's same operator that's case-insensitive
	Contains OpType = "contains"
	// ContainsSensitive operator with case-sensitive
	ContainsSensitive OpType = "contains_s"
	// NotContains operator with case-sensitive
	NotContains OpType = "not_contains"
	// NotContainsInsensitive operator with case-insensitive
	NotContainsInsensitive OpType = "not_contains_i"
	// EndsWith operator with case-sensitive
	EndsWith OpType = "ends_with"
	// EndsWithInsensitive operator with case-insensitive
	EndsWithInsensitive OpType = "ends_with_i"
	// NotEndsWith operator with case-sensitive
	NotEndsWith OpType = "not_ends_with"
	// NotEndsWithInsensitive operator with case-insensitive
	NotEndsWithInsensitive OpType = "not_ends_with_i"

	// array operator

	// IsEmpty operator
	IsEmpty OpType = "is_empty"
	// IsNotEmpty operator
	IsNotEmpty OpType = "is_not_empty"
	// Size operator
	Size OpType = "size"

	// null check operator

	// IsNull operator
	IsNull OpType = "is_null"
	// IsNotNull operator
	IsNotNull OpType = "is_not_null"

	// existence check operator

	// Exist operator
	Exist OpType = "exist"
	// NotExist operator
	NotExist OpType = "not_exist"

	// filter embedded elements operator

	// Object filter object fields operator
	Object OpType = "filter_object"
	// Array filter array elements operator
	Array OpType = "filter_array"
)

// OpType defines the operators supported by cc.
type OpType string

// Validate test the operator is valid or not.
func (op OpType) Validate() error {
	switch op {
	case Equal, NotEqual, In, NotIn, Less, LessOrEqual, Greater, GreaterOrEqual, DatetimeLess, DatetimeLessOrEqual,
		DatetimeGreater, DatetimeGreaterOrEqual, BeginsWith, BeginsWithInsensitive, NotBeginsWith,
		NotBeginsWithInsensitive, Contains, ContainsSensitive, NotContains, NotContainsInsensitive, EndsWith,
		EndsWithInsensitive, NotEndsWith, NotEndsWithInsensitive, IsEmpty, IsNotEmpty, Size, IsNull,
		IsNotNull, Exist, NotExist, Object, Array:
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
	ValidateValue(v interface{}, opt *ExprOption) error
	// ToMgo generate an operator's mongo condition with its field and value.
	ToMgo(field string, value interface{}) (map[string]interface{}, error)
	// Match checks if the first data matches the second data by this operator
	Match(value1, value2 interface{}) (bool, error)
}

// UnknownOp is unknown operator
type UnknownOp OpType

// Name is equal operator
func (o UnknownOp) Name() OpType {
	return Unknown
}

// ValidateValue validate equal's value
func (o UnknownOp) ValidateValue(_ interface{}, _ *ExprOption) error {
	return errors.New("unknown operator")
}

// ToMgo convert this operator's field and value to a mongo query condition.
func (o UnknownOp) ToMgo(_ string, _ interface{}) (map[string]interface{}, error) {
	return nil, errors.New("unknown operator, can not gen mongo expression")
}

// Match checks if the first data matches the second data by this operator
func (o UnknownOp) Match(_, _ interface{}) (bool, error) {
	return false, errors.New("unknown operator, can not check if two value matches this operator")
}

// EqualOp is equal operator type
type EqualOp OpType

// Name is equal operator name
func (o EqualOp) Name() OpType {
	return Equal
}

// ValidateValue validate equal operator's value
func (o EqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsBasicValue(v) {
		return fmt.Errorf("invalid eq value(%+v)", v)
	}
	return nil
}

// ToMgo convert the equal operator's field and value to a mongo query condition.
func (o EqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBEQ: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EqualOp) Match(value1, value2 interface{}) (bool, error) {
	switch t := value1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		val1, val2, err := parseNumericValues(value1, value2)
		if err != nil {
			return false, err
		}
		return val1 == val2, nil
	case string:
		val, ok := value2.(string)
		if !ok {
			return false, fmt.Errorf("rule value type(%T) not matches input type(%T)", value2, value1)
		}
		return val == t, nil
	case bool:
		val, ok := value2.(bool)
		if !ok {
			return false, fmt.Errorf("rule value type(%T) not matches input type(%T)", value2, value1)
		}
		return val == t, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("value(%+v) is not of basic type", value1)
	}
}

// NotEqualOp is not equal operator type
type NotEqualOp OpType

// Name is not equal operator name
func (ne NotEqualOp) Name() OpType {
	return NotEqual
}

// ValidateValue validate not equal operator's value
func (ne NotEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsBasicValue(v) {
		return fmt.Errorf("invalid ne value(%+v)", v)
	}
	return nil
}

// ToMgo convert the not equal operator's field and value to a mongo query condition.
func (ne NotEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBNE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (ne NotEqualOp) Match(value1, value2 interface{}) (bool, error) {
	matched, err := OpFactory(Equal).Operator().Match(value1, value2)
	if err != nil {
		return false, err
	}
	return !matched, nil
}

// InOp is in operator
type InOp OpType

// Name is in operator name
func (o InOp) Name() OpType {
	return In
}

// ValidateValue validate in operator's value
func (o InOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if opt == nil {
		return errors.New("validate option must be set")
	}

	err := valid.ValidateSliceOfBasicType(v, opt.MaxInLimit)
	if err != nil {
		return fmt.Errorf("in operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the in operator's field and value to a mongo query condition.
func (o InOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBIN: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o InOp) Match(value1, value2 interface{}) (bool, error) {
	var itemType string

	switch value1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		itemType = "numeric"
	case string:
		itemType = "string"
	case bool:
		itemType = "bool"
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("value(%+v) is not of basic type", value1)
	}

	if value2 == nil {
		return false, errors.New("rule value is nil")
	}

	switch reflect.TypeOf(value2).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value2)
	}

	v := reflect.ValueOf(value2)
	length := v.Len()
	if length == 0 {
		return false, errors.New("value is empty")
	}

	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()

		switch itemType {
		case "numeric":
			val1, val2, err := parseNumericValues(value1, item)
			if err != nil {
				return false, err
			}
			if val1 == val2 {
				return true, nil
			}
		case "string":
			val, ok := item.(string)
			if !ok {
				return false, fmt.Errorf("array ele index(%d) type(%T) not matches input type(%s)", i, item, itemType)
			}
			if val == value1 {
				return true, nil
			}
		case "bool":
			val, ok := item.(bool)
			if !ok {
				return false, fmt.Errorf("array ele index(%d) type(%T) not matches input type(%s)", i, item, itemType)
			}
			if val == value1 {
				return true, nil
			}
		}
	}

	return false, nil
}

// NotInOp is not in operator
type NotInOp OpType

// Name is not in operator name
func (o NotInOp) Name() OpType {
	return NotIn
}

// ValidateValue validate not in value
func (o NotInOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if opt == nil {
		return errors.New("validate option must be set")
	}

	err := valid.ValidateSliceOfBasicType(v, opt.MaxNotInLimit)
	if err != nil {
		return fmt.Errorf("nin operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not in operator's field and value to a mongo query condition.
func (o NotInOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBNIN: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotInOp) Match(value1, value2 interface{}) (bool, error) {
	matched, err := OpFactory(In).Operator().Match(value1, value2)
	if err != nil {
		return false, err
	}
	return !matched, nil
}

// LessOp is less than operator
type LessOp OpType

// Name is less than operator name
func (o LessOp) Name() OpType {
	return Less
}

// ValidateValue validate less than operator value
func (o LessOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsNumeric(v) {
		return fmt.Errorf("invalid lt operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the less than  operator's field and value to a mongo query condition.
func (o LessOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBLT: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o LessOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 < val2, nil
}

// LessOrEqualOp is less than or equal operator
type LessOrEqualOp OpType

// Name is less than or equal operator name
func (o LessOrEqualOp) Name() OpType {
	return LessOrEqual
}

// ValidateValue validate less than or equal operator value
func (o LessOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsNumeric(v) {
		return errors.New("invalid lte operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the less than or equal operator's field and value to a mongo query condition.
func (o LessOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBLTE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o LessOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 <= val2, nil
}

// GreaterOp is greater than operator
type GreaterOp OpType

// Name is greater than operator name
func (o GreaterOp) Name() OpType {
	return Greater
}

// ValidateValue validate greater than operator value
func (o GreaterOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsNumeric(v) {
		return errors.New("invalid gt operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the greater than operator's field and value to a mongo query condition.
func (o GreaterOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBGT: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o GreaterOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 > val2, nil
}

// GreaterOrEqualOp is greater than or equal operator
type GreaterOrEqualOp OpType

// Name is greater than or equal operator name
func (o GreaterOrEqualOp) Name() OpType {
	return GreaterOrEqual
}

// ValidateValue validate greater than or equal operator value
func (o GreaterOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !util.IsNumeric(v) {
		return errors.New("invalid gte operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the greater than or equal operator's field and value to a mongo query condition.
func (o GreaterOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBGTE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o GreaterOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 >= val2, nil
}

// DatetimeLessOp is datetime less than operator
type DatetimeLessOp OpType

// Name is datetime less than operator name
func (o DatetimeLessOp) Name() OpType {
	return DatetimeLess
}

// ValidateValue validate datetime less than operator value
func (o DatetimeLessOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime less than operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime less than operator's field and value to a mongo query condition.
func (o DatetimeLessOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := util.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBLT: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeLessOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.Before(val2), nil
}

// DatetimeLessOrEqualOp is datetime less than or equal operator
type DatetimeLessOrEqualOp OpType

// Name is datetime less than or equal operator name
func (o DatetimeLessOrEqualOp) Name() OpType {
	return DatetimeLessOrEqual
}

// ValidateValue validate datetime less than or equal operator value
func (o DatetimeLessOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime less than or equal operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime less than or equal operator's field and value to a mongo query condition.
func (o DatetimeLessOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := util.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBLTE: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeLessOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.Before(val2) || val1.Equal(val2), nil
}

// DatetimeGreaterOp is datetime greater than operator
type DatetimeGreaterOp OpType

// Name is datetime greater than operator name
func (o DatetimeGreaterOp) Name() OpType {
	return DatetimeGreater
}

// ValidateValue validate datetime greater than operator value
func (o DatetimeGreaterOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime greater than operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the datetime greater than operator's field and value to a mongo query condition.
func (o DatetimeGreaterOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := util.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBGT: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeGreaterOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.After(val2), nil
}

// DatetimeGreaterOrEqualOp is datetime greater than or equal operator
type DatetimeGreaterOrEqualOp OpType

// Name is datetime greater than or equal operator name
func (o DatetimeGreaterOrEqualOp) Name() OpType {
	return DatetimeGreaterOrEqual
}

// ValidateValue validate datetime greater than or equal operator value
func (o DatetimeGreaterOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime greater than or equal operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime greater than or equal operator's field and value to a mongo query condition.
func (o DatetimeGreaterOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := util.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{common.BKDBGTE: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeGreaterOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.After(val2) || val1.Equal(val2), nil
}

// BeginsWithOp is begins with operator
type BeginsWithOp OpType

// Name is begins with operator name
func (o BeginsWithOp) Name() OpType {
	return BeginsWith
}

// ValidateValue validate begins with operator's value
func (o BeginsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("begins with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the begins with operator's field and value to a mongo query condition.
func (o BeginsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf("^%s", value),
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o BeginsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(val1, val2), nil
}

// BeginsWithInsensitiveOp is begins with insensitive operator
type BeginsWithInsensitiveOp OpType

// Name is begins with insensitive operator name
func (o BeginsWithInsensitiveOp) Name() OpType {
	return BeginsWithInsensitive
}

// ValidateValue validate begins with insensitive operator's value
func (o BeginsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("begins with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the begins with insensitive operator's field and value to a mongo query condition.
func (o BeginsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE:    fmt.Sprintf("^%s", value),
			common.BKDBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o BeginsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// NotBeginsWithOp is not begins with operator
type NotBeginsWithOp OpType

// Name is not begins with operator name
func (o NotBeginsWithOp) Name() OpType {
	return NotBeginsWith
}

// ValidateValue validate not begins with operator's value
func (o NotBeginsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not begins with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not begins with operator's field and value to a mongo query condition.
func (o NotBeginsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf("^%s", value)},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotBeginsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(val1, val2), nil
}

// NotBeginsWithInsensitiveOp is not begins with insensitive operator
type NotBeginsWithInsensitiveOp OpType

// Name is not begins with insensitive operator name
func (o NotBeginsWithInsensitiveOp) Name() OpType {
	return NotBeginsWithInsensitive
}

// ValidateValue validate not begins with insensitive operator's value
func (o NotBeginsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not begins with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not begins with insensitive operator's field and value to a mongo query condition.
func (o NotBeginsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{
				common.BKDBLIKE:    fmt.Sprintf("^%s", value),
				common.BKDBOPTIONS: "i",
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotBeginsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// ContainsOp is contains operator
type ContainsOp OpType

// Name is contains operator name
func (o ContainsOp) Name() OpType {
	return Contains
}

// ValidateValue validate contains operator's value
func (o ContainsOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("contains operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the contains operator's field and value to a mongo query condition.
func (o ContainsOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE:    value,
			common.BKDBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ContainsOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// ContainsSensitiveOp is contains sensitive operator
type ContainsSensitiveOp OpType

// Name is contains sensitive operator name
func (o ContainsSensitiveOp) Name() OpType {
	return ContainsSensitive
}

// ValidateValue validate contains sensitive operator's value
func (o ContainsSensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("contains sensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the contains sensitive operator's field and value to a mongo query condition.
func (o ContainsSensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE: value,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ContainsSensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.Contains(val1, val2), nil
}

// NotContainsOp is not contains operator
type NotContainsOp OpType

// Name is not contains operator name
func (o NotContainsOp) Name() OpType {
	return NotContains
}

// ValidateValue validate not contains operator's value
func (o NotContainsOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not contains operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not contains operator's field and value to a mongo query condition.
func (o NotContainsOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: value},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotContainsOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.Contains(val1, val2), nil
}

// NotContainsInsensitiveOp is not contains insensitive operator
type NotContainsInsensitiveOp OpType

// Name is not contains insensitive operator name
func (o NotContainsInsensitiveOp) Name() OpType {
	return NotContainsInsensitive
}

// ValidateValue validate not contains insensitive operator's value
func (o NotContainsInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not contains insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not contains insensitive operator's field and value to a mongo query condition.
func (o NotContainsInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{
				common.BKDBLIKE:    value,
				common.BKDBOPTIONS: "i",
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotContainsInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.Contains(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// EndsWithOp is ends with operator
type EndsWithOp OpType

// Name is ends with operator name
func (o EndsWithOp) Name() OpType {
	return EndsWith
}

// ValidateValue validate ends with operator's value
func (o EndsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("ends with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the ends with operator's field and value to a mongo query condition.
func (o EndsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf("%s$", value),
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EndsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(val1, val2), nil
}

// EndsWithInsensitiveOp is ends with insensitive operator
type EndsWithInsensitiveOp OpType

// Name is ends with insensitive operator name
func (o EndsWithInsensitiveOp) Name() OpType {
	return EndsWithInsensitive
}

// ValidateValue validate ends with insensitive operator's value
func (o EndsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("ends with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the ends with insensitive operator's field and value to a mongo query condition.
func (o EndsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBLIKE:    fmt.Sprintf("%s$", value),
			common.BKDBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EndsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// NotEndsWithOp is not ends with operator
type NotEndsWithOp OpType

// Name is not ends with operator name
func (o NotEndsWithOp) Name() OpType {
	return NotEndsWith
}

// ValidateValue validate not ends with operator's value
func (o NotEndsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not ends with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not ends with operator's field and value to a mongo query condition.
func (o NotEndsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf("%s$", value)},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotEndsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasSuffix(val1, val2), nil
}

// NotEndsWithInsensitiveOp is not ends with insensitive operator
type NotEndsWithInsensitiveOp OpType

// Name is not ends with insensitive operator name
func (o NotEndsWithInsensitiveOp) Name() OpType {
	return NotEndsWithInsensitive
}

// ValidateValue validate not ends with insensitive operator's value
func (o NotEndsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := valid.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not ends with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// Match checks if the first data matches the second data by this operator
func (o NotEndsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasSuffix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// ToMgo convert the not ends with insensitive operator's field and value to a mongo query condition.
func (o NotEndsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNot: map[string]interface{}{
				common.BKDBLIKE:    fmt.Sprintf("%s$", value),
				common.BKDBOPTIONS: "i",
			},
		},
	}, nil
}

// IsEmptyOp is empty operator
type IsEmptyOp OpType

// Name is empty operator name
func (o IsEmptyOp) Name() OpType {
	return IsEmpty
}

// ValidateValue validate empty operator's value
func (o IsEmptyOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the empty operator's field and value to a mongo query condition.
func (o IsEmptyOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBSize: 0,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsEmptyOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	return reflect.ValueOf(value1).Len() == 0, nil
}

// IsNotEmptyOp is not empty operator
type IsNotEmptyOp OpType

// Name is not empty operator name
func (o IsNotEmptyOp) Name() OpType {
	return IsNotEmpty
}

// ValidateValue validate is not empty operator's value
func (o IsNotEmptyOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not empty operator's field and value to a mongo query condition.
func (o IsNotEmptyOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBSize: map[string]interface{}{common.BKDBGT: 0},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNotEmptyOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	return reflect.ValueOf(value1).Len() > 0, nil
}

// SizeOp size operator
type SizeOp OpType

// Name size operator name
func (o SizeOp) Name() OpType {
	return Size
}

// ValidateValue validate size operator's value
func (o SizeOp) ValidateValue(v interface{}, opt *ExprOption) error {
	intVal, err := util.GetInt64ByInterface(v)
	if err != nil {
		return fmt.Errorf("invalid size operator's value, should be a numeric value, err: %v", err)
	}

	if intVal < 0 {
		return fmt.Errorf("invalid size operator's value, should not be negative")
	}
	return nil
}

// ToMgo convert the size operator's field and value to a mongo query condition.
func (o SizeOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBSize: value,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o SizeOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	intVal, err := util.GetIntByInterface(value2)
	if err != nil {
		return false, fmt.Errorf("invalid size operator's value, should be a numeric value, err: %v", err)
	}

	return reflect.ValueOf(value1).Len() == intVal, nil
}

// IsNullOp is null operator
type IsNullOp OpType

// Name is null operator name
func (o IsNullOp) Name() OpType {
	return IsNull
}

// ValidateValue validate null operator's value
func (o IsNullOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the null operator's field and value to a mongo query condition.
func (o IsNullOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBEQ: nil,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNullOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 == nil, nil
}

// IsNotNullOp is not null operator
type IsNotNullOp OpType

// Name is not null operator name
func (o IsNotNullOp) Name() OpType {
	return IsNotNull
}

// ValidateValue validate is not null operator's value
func (o IsNotNullOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not null operator's field and value to a mongo query condition.
func (o IsNotNullOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBNE: nil,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNotNullOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 != nil, nil
}

// ExistOp is 'exist' operator
type ExistOp OpType

// Name is 'exist' operator name
func (o ExistOp) Name() OpType {
	return Exist
}

// ValidateValue validate 'exist' operator's value
func (o ExistOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the 'exist' operator's field and value to a mongo query condition.
func (o ExistOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBExists: true,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ExistOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 == nil, nil
}

// NotExistOp is not exist operator
type NotExistOp OpType

// Name is not exist operator name
func (o NotExistOp) Name() OpType {
	return NotExist
}

// ValidateValue validate is not exist operator's value
func (o NotExistOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not exist operator's field and value to a mongo query condition.
func (o NotExistOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			common.BKDBExists: false,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotExistOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 != nil, nil
}

// ObjectOp is filter object operator
type ObjectOp OpType

// Name is filter object operator name
func (o ObjectOp) Name() OpType {
	return Object
}

// ValidateValue validate filter object operator value
func (o ObjectOp) ValidateValue(v interface{}, opt *ExprOption) error {
	// filter object operator's value is the sub-rule to filter the object's field.
	subRule, ok := v.(RuleFactory)
	if !ok {
		return fmt.Errorf("filter object operator's value(%+v) is not a rule type", v)
	}

	// validate filter array rule depth, then continues to validate children rule depth
	if opt == nil {
		return errors.New("validate option must be set")
	}
	if opt.MaxRulesDepth <= 1 {
		return fmt.Errorf("expression rules depth exceeds maximum")
	}

	childOpt := cloneExprOption(opt)
	childOpt.MaxRulesDepth = opt.MaxRulesDepth - 1

	if err := subRule.Validate(childOpt); err != nil {
		return fmt.Errorf("invalid value(%+v), err: %v", v, err)
	}

	return nil
}

// ToMgo convert the filter object operator's field and value to a mongo query condition.
func (o ObjectOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	subRule, ok := value.(RuleFactory)
	if !ok {
		return nil, fmt.Errorf("filter object operator's value(%+v) is not a rule type", value)
	}

	parentOpt := &RuleOption{
		Parent:     field,
		ParentType: enumor.Object,
	}

	return subRule.ToMgo(parentOpt)
}

// Match checks if the first data matches the second data by this operator
func (o ObjectOp) Match(value1, value2 interface{}) (bool, error) {
	subRule, ok := value2.(RuleFactory)
	if !ok {
		return false, fmt.Errorf("filter object operator's value(%+v) is not a rule type", value2)
	}

	switch t := value1.(type) {
	case MatchedData:
		return subRule.Match(t)
	case map[string]interface{}:
		return subRule.Match(MapStr(t))
	case mapstr.MapStr:
		return subRule.Match(MapStr(t))
	case string:
		return subRule.Match(JsonString(t))
	case json.RawMessage:
		return subRule.Match(JsonString(t))
	default:
		return false, fmt.Errorf("filter object operator's input value(%+v) is not an object type", value1)
	}
}

const (
	ArrayElement = "element"
)

// ArrayOp is filter array operator
type ArrayOp OpType

// Name is filter array operator name
func (o ArrayOp) Name() OpType {
	return Array
}

// ValidateValue validate filter array operator value
func (o ArrayOp) ValidateValue(v interface{}, opt *ExprOption) error {

	// filter array operator's value is the sub-rule to filter the array's field.
	subRule, ok := v.(RuleFactory)
	if !ok {
		return fmt.Errorf("filter array operator's value(%+v) is not a rule type", v)
	}

	// validate filter array rule depth, then continues to validate children rule depth
	if opt == nil {
		return errors.New("validate option must be set")
	}

	if opt.MaxRulesDepth <= 1 {
		return fmt.Errorf("expression rules depth exceeds maximum")
	}

	childOpt := cloneExprOption(opt)
	childOpt.MaxRulesDepth = opt.MaxRulesDepth - 1

	if err := subRule.Validate(childOpt); err != nil {
		return fmt.Errorf("invalid value(%+v), err: %v", v, err)
	}

	return nil
}

// ToMgo convert the filter array operator's field and value to a mongo query condition.
func (o ArrayOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	subRule, ok := value.(RuleFactory)
	if !ok {
		return nil, fmt.Errorf("filter array operator's value(%+v) is not a rule type", value)
	}

	parentOpt := &RuleOption{
		Parent:     field,
		ParentType: enumor.Array,
	}

	return subRule.ToMgo(parentOpt)
}

// Match checks if the first data matches the second data by this operator
func (o ArrayOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	subRule, ok := value2.(RuleFactory)
	if !ok {
		return false, fmt.Errorf("filter array operator's value(%+v) is not a rule type", value2)
	}

	val := MapStr{
		ArrayElement: value1,
	}

	parentOpt := &RuleOption{
		ParentType: enumor.Array,
	}

	return subRule.Match(val, parentOpt)
}
