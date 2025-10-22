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

package structs_test

import (
	"encoding/json/v2"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

var (
	jsonStr    = []byte(`{"str": "test","int": 1,"float": 1.1,"bool": true,"array": [1,2,3],"map": {"a": "a","b": 2,"c": 3.3}}`)
	mapData    = make(map[string]any)
	structData *structs.Struct
)

func init() {
	// init test struct builder
	builder, err := structs.UpsertBuilderByFields("test", []structs.Field{
		{
			Name: "Str",
			Type: structs.StringType,
			Tags: map[string]string{"json": "str"},
		},
		{
			Name: "Int",
			Type: structs.Int64Type,
			Tags: map[string]string{"json": "int"},
		},
		{
			Name: "Float",
			Type: structs.Float64Type,
			Tags: map[string]string{"json": "float"},
			Validator: func(data any) error {
				if data.(float64) < 0 {
					return errors.New("float value cannot be negative")
				}
				return nil
			},
		},
		{
			Name: "Bool",
			Type: structs.BoolType,
			Tags: map[string]string{"json": "bool"},
		},
		{
			Name:    "Array",
			Type:    structs.Int64Type,
			IsSlice: true,
			Tags:    map[string]string{"json": "array"},
		},
		{
			Name: "Map",
			Type: structs.MapType,
			Tags: map[string]string{"json": "map"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// init test struct data
	structData = builder.New()
	if err = json.Unmarshal(jsonStr, structData.Pointer()); err != nil {
		log.Fatal(err)
	}

	// init test map data
	if err = json.Unmarshal(jsonStr, &mapData); err != nil {
		log.Fatal(err)
	}
}

func TestStruct(t *testing.T) {
	builder, _ := structs.GetBuilder("test")

	// test struct builder
	s := builder.New()
	data := s.Pointer()
	if err := json.Unmarshal(jsonStr, data); err != nil {
		t.Fatal(err)
	}

	// test struct get
	str, err := s.Get("Str")
	if err != nil {
		t.Fatal(err)
	}
	if str.String() != "test" {
		t.Fatalf("expect test, got %v", str)
	}

	// test struct set
	if err = s.Set("Int", 2); err != nil {
		t.Fatal(err)
	}
	intVal, err := s.Get("Int")
	if err != nil {
		t.Fatal(err)
	}
	if intVal.Int64() != int64(2) {
		t.Fatalf("expect 2, got %v", intVal)
	}

	// test struct validate
	if err = s.Validate(); err != nil {
		t.Fatal(err)
	}
	if err = s.Set("Float", -1.2); err != nil {
		t.Fatal(err)
	}
	if err = s.Validate(); err == nil {
		t.Fatal("expect validation error, got nil")
	}
}

func TestSlice(t *testing.T) {
	builder, _ := structs.GetBuilder("test")

	// test slice builder
	s := builder.NewSlice(0, 0)
	data := s.Pointer()

	sliceJson := []byte(`[{"str": "test1","int": 1,"float": 1.1,"bool": true,"array": [1,2,3],"map": {"a": "a","b": 2,"c": 3.3}},{"str":"test2"}]`)
	if err := json.Unmarshal(sliceJson, &data); err != nil {
		t.Fatal(err)
	}

	// test slice len
	if s.Len() != 2 {
		t.Fatalf("expect 2, got %v", s.Len())
	}

	// test slice get
	s0Val, err := s.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if intVal := reflect.ValueOf(s0Val).FieldByName("Int").Int(); intVal != 1 {
		t.Fatalf("expect 1, got %v", intVal)
	}

	// test slice get struct
	s0, err := s.GetStruct(0)
	if err != nil {
		t.Fatal(err)
	}
	s0Str, err := s0.Get("Str")
	if err != nil {
		t.Fatal(err)
	}
	if s0Str.String() != "test1" {
		t.Fatalf("expect test, got %v", s0Str)
	}

	s1, err := s.GetStruct(1)
	if err != nil {
		t.Fatal(err)
	}
	s1Str, err := s1.Get("Str")
	if err != nil {
		t.Fatal(err)
	}
	if s1Str.String() != "test2" {
		t.Fatalf("expect test, got %v", s0Str)
	}

	// test slice set
	setVal := structData.Value()
	if err = s.Set(0, setVal); err != nil {
		t.Fatal(err)
	}
	structVal, err := s.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(structVal, setVal) {
		t.Fatalf("expect slice data equal struct data, got %v", structVal)
	}

	// test slice append
	if err = s.Append(s0Val); err != nil {
		t.Fatal(err)
	}

	// test slice validate
	if err = s.Validate(); err != nil {
		t.Fatal(err)
	}

	if err = s1.Set("Float", -1.2); err != nil {
		t.Fatal(err)
	}
	if err = s.SetStruct(1, s1); err != nil {
		t.Fatal(err)
	}
	if err = s.Validate(); err == nil {
		t.Fatal("expect validation error, got nil")
	}
}

func TestDynamicStructAsField(t *testing.T) {
	testBuilder, ok := structs.GetBuilder("test")
	if !ok {
		t.Fatal("test builder not found")
	}
	testWrapperBuilder, err := structs.UpsertBuilderByFields("testWrapper", []structs.Field{
		{
			Name: "Test",
			Type: "test",
			Tags: map[string]string{"json": "test"},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test wrapper builder: %v", err)
	}
	testWrapperInst := testWrapperBuilder.New()
	testInst := testBuilder.New()
	testInst.Set("Str", "testStr")
	testInst.Set("Int", 1)
	testInst.Set("Float", 1.1)
	testInst.Set("Bool", true)
	testInst.Set("Array", []int64{1, 2, 3})
	testInst.Set("Map", map[string]any{
		"strKey": "str_value",
	})
	err = testWrapperInst.Set("Test", testInst.Value())
	if err != nil {
		t.Fatalf("failed to set test wrapper value to test: %v", err)
	}
	actualJson, err := json.Marshal(testWrapperInst.Value())
	if err != nil {
		t.Fatalf("fail to marshal testWrapper to json, err: %w", err)
	}
	data := map[string]any{
		"test": map[string]any{
			"str":   "testStr",
			"int":   int64(1),
			"float": float64(1.1),
			"bool":  true,
			"array": []int64{1, 2, 3},
			"map": map[string]any{
				"strKey": "str_value",
			},
		},
	}
	wantJson, _ := json.Marshal(data)
	assert.JSONEq(t, string(wantJson), string(actualJson))
}

func BenchmarkMapMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(mapData); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStructMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(structData.Pointer()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[string]any)
		if err := json.Unmarshal(jsonStr, &m); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStructUnmarshal(b *testing.B) {
	builder, _ := structs.GetBuilder("test")

	for n := 0; n < b.N; n++ {
		data := builder.New().Pointer()
		if err := json.Unmarshal(jsonStr, data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, key := range []string{"str", "int", "float", "bool", "array", "map"} {
			_, ok := mapData[key]
			if !ok {
				b.Fatalf("key %s not found", key)
			}
		}
	}
}

func BenchmarkStructGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, key := range []string{"Str", "Int", "Float", "Bool", "Array", "Map"} {
			_, err := structData.Get(key)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkMapSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mapData["str"] = "abc"
		mapData["int"] = 2
		mapData["float"] = 3.3
		mapData["bool"] = false
		mapData["array"] = []int64{4, 5}
		mapData["map"] = map[string]any{"d": 1, "e": "f"}
	}
}

func BenchmarkStructSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if err := structData.Set("Str", "abc"); err != nil {
			b.Fatal(err)
		}
		if err := structData.Set("Int", 2); err != nil {
			b.Fatal(err)
		}
		if err := structData.Set("Float", 3.3); err != nil {
			b.Fatal(err)
		}
		if err := structData.Set("Bool", false); err != nil {
			b.Fatal(err)
		}
		if err := structData.Set("Array", []int64{4, 5}); err != nil {
			b.Fatal(err)
		}
		if err := structData.Set("Map", map[string]any{"d": 1, "e": "f"}); err != nil {
			b.Fatal(err)
		}
	}
}
