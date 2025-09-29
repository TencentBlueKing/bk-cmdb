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
)

func decodeTo(r *http.Request, val any) error {
	rt := reflect.TypeOf(val).Elem()
	rv := reflect.ValueOf(val).Elem()

	// json 整个解析
	jsonCodec := NewJsonCodec(r)
	if err := jsonCodec.Decode(val); err != nil {
		return err
	}

	formCodec, err := NewFormCodec(r)
	if err != nil {
		return err
	}

	pathCodec := NewPathCodec(r)
	queryCodec := NewQueryCodec(r)
	headerCodec := NewHeaderCodec(r)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// 非导出需要跳过, 无法设置值
		if !field.IsExported() {
			continue
		}

		tagStr := field.Tag.Get(tagName)
		if tagStr == "" {
			continue
		}
		tag, err := parseTag(tagStr)
		if err != nil {
			return err
		}

		fv := rv.Field(i)
		if err := formCodec.Decode(field, fv, tag); err != nil {
			return err
		}
		if err := queryCodec.Decode(field, fv, tag); err != nil {
			return err
		}
		if err := headerCodec.Decode(field, fv, tag); err != nil {
			return err
		}
		if err := pathCodec.Decode(field, fv, tag); err != nil {
			return err
		}
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
		return nil, fmt.Errorf("codec decode: %w", err)
	}

	return t, nil
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
	if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
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
