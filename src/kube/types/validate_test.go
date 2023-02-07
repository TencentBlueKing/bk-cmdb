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

package types

import "testing"

// TestValidateNumeric validation function unit test for numeric types
func TestValidateNumeric(t *testing.T) {

	a := 5
	param := NumericSettings{
		Min: 0,
		Max: 8,
	}
	if err := ValidateNumeric(a, param); err != nil {
		t.Fatalf("test validate numeric is failed， err: %v", err)
	}

	b := "aaaa"
	if err := ValidateNumeric(b, param); err != nil {
		t.Logf("test data type cannot be string success")
	}

	c := 10
	if err := ValidateNumeric(c, param); err != nil {
		t.Logf("test that the numeric type cannot exceed the specified range successfully")
	}

	var d interface{}
	if err := ValidateNumeric(d, param); err != nil {
		t.Logf("test data cannot be nil, err: %v", err)
	}
	return
}

//TestValidateBoolen validation function unit test for numeric bool
func TestValidateBoolen(t *testing.T) {
	a := false
	if err := ValidateBoolen(a); err != nil {
		t.Fatalf("data type is bool, test failed, err: %v", err)
	}
	b := 5
	if err := ValidateBoolen(b); err != nil {
		t.Logf("the data type cannot be a number, the verification is successful, err: %v", err)
	}
	c := "aaa"
	if err := ValidateBoolen(c); err != nil {
		t.Logf("the data type cannot be a string, the verification is successful, err: %v", err)
	}
	d := []int{1, 2, 3}
	if err := ValidateBoolen(d); err != nil {
		t.Logf("the data type cannot be an array, the verification is successful, err: %v", err)
	}
	var e interface{}
	if err := ValidateBoolen(e); err != nil {
		t.Logf("test data cannot be nil, err: %v", err)
	}
}

// TestValidateString validation function unit test for numeric string
func TestValidateString(t *testing.T) {
	a := "a"
	param := StringSettings{
		MaxLength:    10,
		RegularCheck: "",
	}
	if err := ValidateString(a, param); err != nil {
		t.Fatalf("data type is string, test failed, err: %v", err)
	}
	b := 5
	if err := ValidateString(b, param); err != nil {
		t.Logf("the data type cannot be a numeric, the verification is successful, err: %v", err)
	}
	c := "1234567890"
	if err := ValidateString(c, param); err != nil {
		t.Logf("the length of the string must conform to the specified range, and the verification fails, err: %v", err)
	}
	d := "baa"
	param.RegularCheck = "^a"
	if err := ValidateString(d, param); err != nil {
		t.Logf("in the presence of a regular expression, the string calibration satisfies the regular, err: %v", err)
	}
	var e interface{}
	if err := ValidateString(e, param); err != nil {
		t.Logf("test data cannot be nil, err: %v", err)
	}
}

// TestValidateMapString validation function unit test for map string, like map[string]string
func TestValidateMapString(t *testing.T) {
	a := map[string]string{"1": "a"}
	if err := ValidateMapString(a, 10); err != nil {
		t.Fatalf("both the key and value of the map must be strings, and the verification failed, err: %v", err)
	}

	b := 5
	if err := ValidateMapString(b, 10); err != nil {
		t.Logf("the data type cannot be a number, the verification failed, err: %v", err)
	}

	c := "5"
	if err := ValidateMapString(c, 10); err != nil {
		t.Logf("the data type cannot be a string, the verification failed, err: %v", err)
	}

	d := map[string]int{"1": 1}
	if err := ValidateMapString(d, 10); err != nil {
		t.Logf("the value in the map must be a string, err: %v", err)
	}

	e := map[string]string{"1": "a", "2": "b"}
	if err := ValidateMapString(e, 1); err != nil {
		t.Logf("the number of key-value pairs in the map cannot exceed the specified maximum number, err: %v", err)
	}

	f := map[int]string{1: "a", 2: "b"}
	if err := ValidateMapString(f, 4); err != nil {
		t.Logf("the key in the map must be a string, err: %v", err)
	}

}

// TestValidateKVObject kay value type test case.
func TestValidateKVObject(t *testing.T) {
	a := map[string]string{
		"1": "1",
		"2": "2",
	}
	param := MapObjectSettings{
		MaxDeep:   5,
		MaxLength: 5,
	}
	if err := ValidateKVObject(a, param, 1); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	b := map[string]map[int]map[string]string{
		"1": {
			1: {"1": "1"},
		},
		"2": {
			2: {"2": "2"},
		},
	}
	param = MapObjectSettings{
		MaxDeep:   3,
		MaxLength: 5,
	}
	if err := ValidateKVObject(b, param, 1); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	param.MaxLength = 1
	if err := ValidateKVObject(b, param, 1); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	c := map[string]map[int]map[string]map[int]map[int]int64{
		"1": {
			2: {"3": {4: {5: 1}}},
		},
		"2": {
			2: {"2": {2: {2: 2}}},
		},
	}
	param = MapObjectSettings{
		MaxDeep:   5,
		MaxLength: 5,
	}
	if err := ValidateKVObject(c, param, 1); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	d := map[string]map[string]string{
		"1": {"1": "1"},
		"2": {"1": "1"},
		"3": {"1": "1"},
	}
	param = MapObjectSettings{
		MaxDeep:   2,
		MaxLength: 3,
	}

	if err := ValidateKVObject(d, param, 1); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}
}

// TestValidateArray array type test case
func TestValidateArray(t *testing.T) {

	a := []string{"11", "22"}
	param := &ArraySettingsParam{
		ArrayMaxLength: 5,
		StringParam: StringSettings{
			MaxLength:    5,
			RegularCheck: "",
		},
	}
	if err := ValidateArray(a, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 1,
		StringParam: StringSettings{
			MaxLength:    5,
			RegularCheck: "",
		},
	}
	if err := ValidateArray(a, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 5,
		StringParam: StringSettings{
			MaxLength:    1,
			RegularCheck: "",
		},
	}
	if err := ValidateArray(a, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 5,
		StringParam: StringSettings{
			MaxLength:    5,
			RegularCheck: "^a",
		},
	}
	if err := ValidateArray(a, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 3,
		StringParam: StringSettings{
			MaxLength:    1,
			RegularCheck: "",
		},
	}
	if err := ValidateArray(a, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	b := []map[string]string{
		{"1": "1"},
		{"2": "2"},
		{"3": "3"},
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 3,
		MapObjectParam: MapObjectSettings{
			MaxLength: 10,
			MaxDeep:   5,
		},
	}
	if err := ValidateArray(b, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	d := []int{1, 2, 3, 4}
	param = &ArraySettingsParam{
		ArrayMaxLength: 5,
		NumericParam: NumericSettings{
			Min: 0,
			Max: 9,
		},
	}
	if err := ValidateArray(d, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	param.NumericParam.Max = 3
	if err := ValidateArray(d, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	e := []map[string]map[string]int{
		{"1": {"1": 1, "a": 1, "b": 1, "c": 1}},
		{"2": {"2": 2, "b": 2, "a": 1, "c": 1}},
		{"3": {"3": 3, "c": 3, "b": 1, "a": 1}},
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 3,
		MapObjectParam: MapObjectSettings{
			MaxLength: 10,
			MaxDeep:   1,
		},
	}
	if err := ValidateArray(e, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 2,
		MapObjectParam: MapObjectSettings{
			MaxLength: 10,
			MaxDeep:   1,
		},
	}
	if err := ValidateArray(e, param); err != nil {
		t.Logf("the element requires an array type, the validation failed, err: %v", err)
	}

	param = &ArraySettingsParam{
		ArrayMaxLength: 3,
		MapObjectParam: MapObjectSettings{
			MaxLength: 4,
			MaxDeep:   3,
		},
	}
	if err := ValidateArray(e, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	f := []map[string]map[string]map[int]int{
		{"1": {"1": {1: 1}, "2": {2: 2}}},
	}
	param = &ArraySettingsParam{
		ArrayMaxLength: 3,
		MapObjectParam: MapObjectSettings{
			MaxLength: 3,
			MaxDeep:   3,
		},
	}
	if err := ValidateArray(f, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}

	g := []map[string]map[string]int{
		{"1": {"1": 1, "2": 1}},
		{"2": {"2": 2, "1": 1}},
	}
	param = &ArraySettingsParam{
		ArrayMaxLength: 2,
		MapObjectParam: MapObjectSettings{
			MaxLength: 2,
			MaxDeep:   2,
		},
	}
	if err := ValidateArray(g, param); err != nil {
		t.Fatalf("the element requires an array type, the validation failed, err: %v", err)
	}
}
