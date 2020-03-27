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

package metric

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func FormFloatOrString(val interface{}) (*FloatOrString, error) {
	valueof := reflect.ValueOf(val)
	switch valueof.Kind() {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return &FloatOrString{
			Type:  Float,
			Float: valueof.Convert(reflect.TypeOf(float64(1))).Float(),
		}, nil
	case reflect.String:
		return &FloatOrString{
			Type:   String,
			String: valueof.String(),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", reflect.ValueOf(val).String())
	}
}

type ValueType string

const (
	Float  ValueType = "Float"
	String ValueType = "String"
)

type FloatOrString struct {
	Type   ValueType
	Float  float64
	String string
}

func (fs FloatOrString) MarshalJSON() ([]byte, error) {
	switch fs.Type {
	case Float:
		return json.Marshal(fs.Float)
	case String:
		return json.Marshal(fs.String)
	default:
		return []byte{}, fmt.Errorf("unsupported type: %s", fs.Type)
	}
}

func (fs *FloatOrString) UnmarshalJSON(b []byte) error {
	f, err := strconv.ParseFloat(string(b), 10)
	if nil == err {
		fs.Type = Float
		fs.Float = f
		return nil
	}

	fs.Type = String
	fs.String = string(b)
	return nil
}
