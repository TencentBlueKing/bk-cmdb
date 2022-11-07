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

package orm

import (
	"encoding/json"
	"reflect"
	"testing"
)

// InlineSpec spec struct which can be inline
type InlineSpec struct {
	SpecIntVal  *int    `json:"spec_int_val" bson:"spec_int_val"`
	SpecStrVal  *string `json:"spec_str_val" bson:"spec_str_val"`
	SpecBoolVal *bool   `json:"spec_bool_val" bson:"spec_bool_val"`
}

// innerStruct inner struct
type innerStruct struct {
	IntVal int    `json:"int_val" bson:"int_val"`
	StrVal string `json:"str_val" bson:"str_val"`
}

// cases test case struct
type cases struct {
	InlineSpec  `json:",inline" bson:",inline"`
	IntVal      *int               `json:"int_val" bson:"int_val"`
	NoPointVal  int                `json:"no_point_val" bson:"no_point_val"`
	StrVal      *string            `json:"str_val" bson:"str_val"`
	BoolVal     *bool              `json:"bool_val" bson:"bool_val"`
	MapVal      *map[string]string `json:"map_val" bson:"map_val"`
	ArrayVal    *[]string          `json:"array_val" bson:"array_val"`
	InnerStruct *innerStruct       `json:"inner_struct" bson:"inner_struct"`
	IgnoreField *string            `json:"ignore_field" bson:"ignore_field"`
}

// TestGetUpdateFieldWithOption test GetUpdateFieldWithOption func
func TestGetUpdateFieldWithOption(t *testing.T) {
	// build test data
	intVal := 1
	strVal := "stringVal"
	boolVal := true
	ignoreFiled := "ignoreFiled"
	arrayVal := []string{
		"array1",
		"array2",
	}
	structVal := &innerStruct{
		IntVal: 1,
	}
	mapVal := map[string]string{
		"testMap1": "1",
		"testMap2": "2",
	}

	requestData := &cases{
		IntVal:      &intVal,
		NoPointVal:  intVal,
		StrVal:      &strVal,
		BoolVal:     &boolVal,
		ArrayVal:    &arrayVal,
		MapVal:      &mapVal,
		InnerStruct: structVal,
		InlineSpec: InlineSpec{
			SpecIntVal: &intVal,
		},
		IgnoreField: &ignoreFiled,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Errorf("marshal request data failed, data: %v, err: %v", requestData, err)
		return
	}
	t.Logf("request data: %s\n", string(jsonData))

	data := &cases{}
	if err := json.Unmarshal(jsonData, data); err != nil {
		t.Errorf("unmarshal data failed, data: %s, err: %v", string(jsonData), err)
		return
	}
	t.Logf("unmarshal data: %v\n", data)

	// get update field with option
	opts := NewFieldOptions().AddIgnoredFields("ignore_filed")
	values, err := GetUpdateFieldsWithOption(data, opts)
	if err != nil {
		t.Errorf("get update field failed, err: %v", err)
		return
	}

	// check whether the returned value is the data we want to update.
	if !reflect.DeepEqual(values["int_val"], 1) {
		t.Errorf("get int_val field failed, not equal")
		return
	}

	if !reflect.DeepEqual(values["no_point_val"], nil) {
		t.Errorf("get no point value")
		return
	}

	if !reflect.DeepEqual(values["str_val"], "stringVal") {
		t.Errorf("get str_val field failed, not equal")
		return
	}

	if !reflect.DeepEqual(values["inner_struct"], *structVal) {
		t.Errorf("get inner_struct field failed, not equal")
		return
	}

	if !reflect.DeepEqual(values["array_val"], arrayVal) {
		t.Errorf("get array_val field failed, not equal")
		return
	}

	for key, val := range values {
		t.Logf("key: %v, val: %v\n", key, val)
	}

	// test nil
	values, err = GetUpdateFieldsWithOption(nil, opts)
	if err == nil {
		t.Errorf("get nil update field should have error")
		return
	}
}
