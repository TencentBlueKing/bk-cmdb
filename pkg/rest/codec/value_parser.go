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

package codec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	// ErrUnsupportedType unsupported parse type
	ErrUnsupportedType = errors.New("unsupported type")

	parserRegistry = make(map[reflect.Type]Parser)
	byteType       = reflect.TypeFor[byte]()
)

// Parser defines the interface for converting a string to a reflect.Value
// Implementations should parse a string representation into a Go value
type Parser interface {
	Parse(s string) (reflect.Value, error)
}

// NewParser defines a factory interface for creating Parser instances
// Implementations should return a new Parser configured with the provided options
type NewParser interface {
	// New creates a new Parser instance with the given configuration options
	// The options map contains key-value pairs for parser configuration
	New(opt map[string]string) Parser
}

// ParserFunc is a function type that implements the Parser interface
// Allows regular functions to be used as Parser implementations
type ParserFunc func(s string) (reflect.Value, error)

// Parse implements the Parser interface for ParserFunc
// This adapter allows any function with the signature func(string) (reflect.Value, error)
// to be used as a Parser without defining a new type
func (p ParserFunc) Parse(s string) (reflect.Value, error) {
	return p(s)
}

// RegisterParser registers a parser implementation for a specific type T
func RegisterParser[T any](p Parser) {
	parserRegistry[reflect.TypeFor[T]()] = p
}

// ParseValue converts a string to a value of the specified type using registered parsers
// Returns an error if parsing fails or the value type is unsupported
func ParseValue(rt reflect.Type, s string, opt map[string]string) (reflect.Value, error) {
	parser, ok := parserRegistry[rt]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%w: %v", ErrUnsupportedType, rt)
	}

	// 实现自定义初始化
	if v, ok := parser.(NewParser); ok {
		parser = v.New(opt)
	}

	return parser.Parse(s)
}

// StringParser ...
func StringParser(s string) (reflect.Value, error) {
	return reflect.ValueOf(s), nil
}

// BoolParser ...
func BoolParser(s string) (reflect.Value, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(v), nil
}

// Int generic int parser
type Int[T int | int8 | int16 | int32 | int64] struct {
	bitSize int
}

// Parse int parser
func (i Int[T]) Parse(s string) (reflect.Value, error) {
	v, err := strconv.ParseInt(s, 10, i.bitSize)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(T(v)), nil
}

// Uint generic uint parser
type Uint[T uint | uint8 | uint16 | uint32 | uint64] struct {
	bitSize int
}

// Parse uint parser
func (i Uint[T]) Parse(s string) (reflect.Value, error) {
	v, err := strconv.ParseUint(s, 10, i.bitSize)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(T(v)), nil
}

// Float generic float parser
type Float[T float32 | float64] struct {
	bitSize int
}

// Parse float parser
func (f Float[T]) Parse(s string) (reflect.Value, error) {
	v, err := strconv.ParseFloat(s, f.bitSize)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(T(v)), nil
}

// ByteSlice is a wrapper of []byte to implement Parser
type ByteSlice []byte

// Parse ...
func (bs ByteSlice) Parse(s string) (reflect.Value, error) {
	v := []byte(s)

	return reflect.ValueOf(v), nil
}

// Time is a parser with format option
type Time struct {
	option map[string]string
}

// Parse ...
func (t Time) Parse(s string) (reflect.Value, error) {
	format := t.option["format"]
	if format == "" {
		format = time.DateTime
	}

	v, err := time.Parse(format, s)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(v), nil
}

// New ...
func (t *Time) New(opt map[string]string) Parser {
	newParser := &Time{option: opt}
	return newParser
}

func init() {
	// buildin parser
	RegisterParser[string](ParserFunc(StringParser))
	RegisterParser[bool](ParserFunc(BoolParser))
	RegisterParser[int](Int[int]{0})
	RegisterParser[int8](Int[int8]{8})
	RegisterParser[int16](Int[int16]{16})
	RegisterParser[int32](Int[int32]{32})
	RegisterParser[int64](Int[int64]{64})
	RegisterParser[uint](Uint[uint]{0})
	RegisterParser[uint8](Uint[uint8]{8})
	RegisterParser[uint16](Uint[uint16]{16})
	RegisterParser[uint32](Uint[uint32]{32})
	RegisterParser[uint64](Uint[uint64]{64})
	RegisterParser[float32](Float[float32]{32})
	RegisterParser[float64](Float[float64]{64})
	RegisterParser[[]byte](ByteSlice{})
	RegisterParser[time.Time](&Time{})
}
