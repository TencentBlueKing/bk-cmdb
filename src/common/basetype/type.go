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

package basetype

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func NewType(val interface{}) (*Type, error) {
	if val == nil {
		return nil, errors.New("value can not be nil")
	}

	valueof := reflect.ValueOf(val)
	switch valueof.Kind() {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return &Type{
			valueType: Float64,
			float64:   valueof.Convert(reflect.TypeOf(float64(1))).Float(),
		}, nil
	case reflect.String:
		return &Type{
			valueType: String,
			string:    valueof.String(),
		}, nil
	case reflect.Bool:
		return &Type{
			valueType: Bool,
			bool:      val.(bool),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", reflect.ValueOf(val).String())
	}
}

func NewMustType(val interface{}) *Type {
	if val == nil {
		panic("value can not be nil")
	}

	valueof := reflect.ValueOf(val)
	switch valueof.Kind() {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return &Type{
			valueType: Float64,
			float64:   valueof.Convert(reflect.TypeOf(float64(1))).Float(),
		}
	case reflect.String:
		return &Type{
			valueType: String,
			string:    valueof.String(),
		}
	case reflect.Bool:
		return &Type{
			valueType: Bool,
			bool:      val.(bool),
		}
	default:
		panic(fmt.Sprintf("unsupported data type: %s", reflect.ValueOf(val).String()))
	}
}

func NewTimeType(val interface{}) (*Type, error) {
	if val == nil {
		return nil, errors.New("value can not be nil")
	}

	valueof := reflect.ValueOf(val)
	switch valueof.Kind() {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return &Type{
			valueType: Time,
			float64:   valueof.Convert(reflect.TypeOf(float64(1))).Float(),
		}, nil
	case reflect.String:
		rfcTime, err := time.Parse(time.RFC3339Nano, val.(string))
		if err != nil {
			return nil, err
		}
		return &Type{
			valueType: Time,
			float64:   float64(rfcTime.Unix()),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported time type: %s", reflect.ValueOf(val).String())
	}
}

func NewMustTimeType(val interface{}) *Type {
	if val == nil {
		panic("time can not be nil")
	}

	valueof := reflect.ValueOf(val)
	switch valueof.Kind() {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return &Type{
			valueType: Time,
			float64:   valueof.Convert(reflect.TypeOf(float64(1))).Float(),
		}
	case reflect.String:
		rfcTime, err := time.Parse(TimeLayout, val.(string))
		if err != nil {
			panic(fmt.Sprintf("parse time failed, err: %v", err))
		}
		return &Type{
			valueType: Time,
			float64:   float64(rfcTime.Unix()),
		}

	default:
		panic(fmt.Errorf("unsupported time type: %s", reflect.ValueOf(val).String()))
	}
}

type Type struct {
	valueType ValueType
	bool      bool
	string    string
	float64   float64
}

type ValueType string

const (
	TimeLayout = time.RFC3339

	Bool    ValueType = "bool"
	String  ValueType = "string"
	Float64 ValueType = "int64"
	Time    ValueType = "time"
)

func (bt *Type) MarshalJSON() ([]byte, error) {
	switch bt.valueType {
	case Float64:
		return json.Marshal(bt.float64)
	case String:
		return json.Marshal(bt.string)
	case Bool:
		return json.Marshal(bt.bool)
	case Time:
		return json.Marshal(time.Unix(int64(bt.float64), 0).Format(TimeLayout))
	default:
		return []byte{}, fmt.Errorf("unsupported type: %s", bt.valueType)
	}
}

func (bt *Type) UnmarshalJSON(b []byte) error {
	if string(b) == "true" {
		bt.valueType = Bool
		bt.bool = true
	} else if string(b) == "false" {
		bt.valueType = Bool
		bt.bool = false
	}

	f, err := strconv.ParseFloat(string(b), 10)
	if nil == err {
		bt.valueType = Float64
		bt.float64 = f
		return nil
	}

	bt.valueType = String
	bt.string = string(b)
	return nil
}

func (bt *Type) Type() ValueType {
	return bt.valueType
}

func (bt *Type) IsBool() bool {
	return bt.valueType == Bool
}

func (bt *Type) IsNumeric() bool {
	return bt.valueType == Float64
}

func (bt *Type) IsString() bool {
	return bt.valueType == String
}

func (bt *Type) Bool() bool {
	return bt.bool
}

func (bt *Type) String() string {
	return bt.string
}

func (bt *Type) Int() int {
	return int(bt.float64)
}

func (bt *Type) Int8() int8 {
	return int8(bt.float64)
}

func (bt *Type) Int16() int16 {
	return int16(bt.float64)
}

func (bt *Type) Int32() int32 {
	return int32(bt.float64)
}

func (bt *Type) Int64() int64 {
	return int64(bt.float64)
}

func (bt *Type) Uint() uint {
	return uint(bt.float64)
}

func (bt *Type) Uint8() uint8 {
	return uint8(bt.float64)
}

func (bt *Type) Uint16() uint16 {
	return uint16(bt.float64)
}

func (bt *Type) Uint32() uint32 {
	return uint32(bt.float64)
}

func (bt *Type) Uint64() uint64 {
	return uint64(bt.float64)
}

func (bt *Type) Float32() float32 {
	return float32(bt.float64)
}

func (bt *Type) Float64() float64 {
	return bt.float64
}

func (bt *Type) Time() (time.Time, error) {
	if len(bt.string) != 0 {
		return time.Parse(TimeLayout, bt.string)
	}
	return time.Unix(int64(bt.float64), 0), nil
}
