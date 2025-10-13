/*
 * Tencent is pleased to support the open source community by making
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

// Package codec provides encoding and decoding utilities across various formats
package codec

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func decodeTo(r *http.Request, val any) error {
	rt := reflect.TypeOf(val).Elem()
	rv := reflect.ValueOf(val).Elem()

	fields, err := getStructFields(rt, rv)
	if err != nil {
		return err
	}

	pathCodec := NewPathCodec(r)
	queryCodec := NewQueryCodec(r)
	for _, f := range fields {
		switch f.tag.In {
		case pathOptName:
			if err := pathCodec.Decode(f.field, f.fv, f.tag); err != nil {
				return fmt.Errorf("field[%s] decode path: %w", f.field.Name, err)
			}
		case queryOptName:
			if err := queryCodec.Decode(f.field, f.fv, f.tag); err != nil {
				return fmt.Errorf("field[%s] decode query: %w", f.field.Name, err)
			}
		case "":
			return fmt.Errorf("field[%s] in option is required", f.field.Name)
		default:
			return fmt.Errorf("field[%s] in[%s] option not valid", f.field.Name, f.tag.In)
		}
	}

	// json 整个解析
	jsonCodec := NewJsonCodec(r)
	if err := jsonCodec.Decode(val); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

// Decode 按结构体反序列化Request
func Decode[T any](r *http.Request) (*T, error) {
	rt := reflect.TypeFor[T]()
	if rt.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct type")
	}

	t := new(T)
	err := decodeTo(r, t)
	if err != nil {
		return nil, fmt.Errorf("decode req: %w", err)
	}

	return t, nil
}

type structField struct {
	field reflect.StructField // 结构体字段
	fv    reflect.Value       // 字段的值
	tag   *Tag                // 字段解析后的req tag
}

// getStructFields 获取字段列表, 校验json/req的唯一性
func getStructFields(rt reflect.Type, rv reflect.Value) ([]structField, error) {
	fields := []structField{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// 非导出需要跳过, 无法设置值
		if !field.IsExported() {
			continue
		}

		reqTagStr := field.Tag.Get(tagName)
		if reqTagStr == "" {
			continue
		}
		tag, err := parseTag(reqTagStr)
		if err != nil {
			return nil, fmt.Errorf("field[%s] %w", field.Name, err)
		}
		// tag name为空或者-忽略
		if tag.Name == "" || tag.Name == "-" {
			continue
		}

		jsonTagName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonTagName != "" && jsonTagName != "-" {
			return nil, fmt.Errorf("field[%s] req and json tag are mutually exclusive", field.Name)
		}

		fields = append(fields, structField{field: field, fv: rv.Field(i), tag: tag})
	}
	return fields, nil
}

// getFieldValue 获取字段值
func getFieldValue(field reflect.Type, tag *Tag, values []string) (reflect.Value, error) {
	// 指针类型
	if field.Kind() == reflect.Pointer {
		typ := field.Elem()
		rv, err := getFieldValue(typ, tag, values)
		if err != nil {
			return reflect.Value{}, err
		}

		newPtr := reflect.New(typ)
		newPtr.Elem().Set(rv)
		return newPtr, nil
	}

	// slice类型
	if field.Kind() == reflect.Slice {
		typ := field.Elem()

		// []byte 特殊处理
		if typ == byteType {
			return ParseValue(field, values[0], tag.Option)
		}

		val := reflect.MakeSlice(field, 0, len(values))
		for _, v := range values {
			rv, err := getFieldValue(typ, tag, []string{v})
			if err != nil {
				return reflect.Value{}, err
			}
			val = reflect.Append(val, rv)
		}
		return val, nil
	}

	return ParseValue(field, values[0], tag.Option)
}
