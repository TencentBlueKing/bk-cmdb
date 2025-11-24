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

// Package structs defines struct related utilities.
package structs

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
)

var (
	// builderRegistry is the mapping of struct name to struct builder.
	builderRegistry = make(map[string]*Builder)
	// builderLock protects concurrent access to builderRegistry.
	builderLock sync.RWMutex

	// handlers for builder change.
	handlers = make(map[uint64]Handler)
	// nextHandlerIDCounter is the id of the next handler.
	nextHandlerIDCounter = uint64(0)
)

// Handler builder change handler
type Handler func(name string, newBuilder *Builder)

// RegisterBuilder registers struct name and its builder to the registry.
func RegisterBuilder(name string, b *Builder) {
	builderLock.Lock()
	builderRegistry[name] = b
	builderLock.Unlock()
	for _, handler := range handlers {
		handler(name, b)
	}
}

// GetBuilder get struct builder by name.
func GetBuilder(name string) (*Builder, bool) {
	builderLock.RLock()
	builder, exists := builderRegistry[name]
	builderLock.RUnlock()
	return builder, exists
}

// RegisterChangeHandler register for builder change.
func RegisterChangeHandler(handlerFunc Handler) uint64 {
	builderLock.Lock()
	defer builderLock.Unlock()
	nextHandlerIDCounter++
	id := nextHandlerIDCounter
	handlers[id] = handlerFunc
	return id
}

// DeregisterChangeHandler stop listening for builder change.
func DeregisterChangeHandler(id uint64) {
	builderLock.Lock()
	defer builderLock.Unlock()
	_, ok := handlers[id]
	if !ok {
		return
	}
	delete(handlers, id)
}

// Builder is used to build a dynamic struct by its fields.
type Builder struct {
	name string
	// structType is the struct reflection type.
	structType reflect.Type
	// validators is the mapping of struct field name to validate functions,
	// validate functions only check first level field, which means embedded field will be checked as a single field.
	validators map[string]func(any) error
	// fieldIndexMap is the mapping of struct field name to their index in the struct.
	// Outer fields override inner fields with the same name,
	// fields with larger index override those with smaller index.
	fieldIndexMap map[string][]int
}

// UpsertBuilderByFields creates or updates a struct builder in the registry by its name and fields.
func UpsertBuilderByFields(name string, fields []Field) (*Builder, error) {
	builderFields := make([]reflect.StructField, len(fields))

	builder := &Builder{
		name:          name,
		validators:    make(map[string]func(any) error),
		fieldIndexMap: make(map[string][]int),
	}
	for i, field := range fields {
		if field.Name == "" {
			return nil, fmt.Errorf("%s struct %d field name is empty", name, i)
		}
		field.Name = strings.ToUpper(string(field.Name[0])) + field.Name[1:]

		// get the reflection type of the field type.
		typ, exists := GetFieldType(field.Type)
		if !exists {
			// get the struct type from builder registry if field type is not pre-defined type.
			s, ok := GetBuilder(string(field.Type))
			if !ok {
				return nil, fmt.Errorf("%s field %s type %s not supported", name, field.Name, field.Type)
			}
			typ = s.structType
		}

		if field.IsSlice {
			typ = reflect.SliceOf(typ)
		}

		// generate struct tag from the field's tag map.
		tags := make([]string, 0)
		for k, v := range field.Tags {
			tags = append(tags, fmt.Sprintf(`%s:"%s"`, k, v))
		}
		tag := reflect.StructTag(strings.Join(tags, " "))

		// generate the reflection struct field.
		builderFields[i] = reflect.StructField{
			Name:      field.Name,
			Type:      typ,
			Tag:       tag,
			Anonymous: field.Anonymous,
		}

		if field.Validator != nil {
			builder.validators[field.Name] = field.Validator
		}
		// record the field index.
		builder.fieldIndexMap[field.Name] = []int{i}
		// travel the embedded field and record the field index.
		travelEmbeddedField(builder.fieldIndexMap, typ, []int{i})
	}

	// create the dynamic struct type from the fields.
	builder.structType = reflect.StructOf(builderFields)

	// register the builder.
	RegisterBuilder(name, builder)

	return builder, nil
}

func travelEmbeddedField(indexMap map[string][]int, typ reflect.Type, curPath []int) {
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := range typ.NumField() {
		subField := typ.Field(i)
		subPath := append(slices.Clip(curPath), i)
		// three cases to update the index:
		// 1. new field
		// 2. old field is deeper
		// 3. old field is same level but index is smaller
		oldIndex, ok := indexMap[subField.Name]
		switch {
		case !ok:
			fallthrough
		case len(oldIndex) > len(subPath):
			fallthrough
		case len(oldIndex) == len(subPath) && oldIndex[len(oldIndex)-1] < subPath[len(subPath)-1]:
			indexMap[subField.Name] = subPath
		default:
			// do not update, keep the old index
		}

		if subField.Anonymous {
			travelEmbeddedField(indexMap, subField.Type, subPath)
		}
	}
}

// New creates a new struct instance.
func (b *Builder) New() *Struct {
	ptr := reflect.New(b.structType)

	return &Struct{
		name:          b.name,
		data:          ptr.Interface(),
		val:           ptr.Elem(),
		validators:    b.validators,
		fieldIndexMap: b.fieldIndexMap,
	}
}

// NewSlice creates a new slice instance.
func (b *Builder) NewSlice(len, cap int) *Slice {
	sliceType := reflect.SliceOf(b.structType)
	sliceValue := reflect.MakeSlice(sliceType, len, cap)

	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(sliceValue)

	return &Slice{
		name:          b.name,
		data:          slicePtr.Interface(),
		val:           slicePtr.Elem(),
		validators:    b.validators,
		fieldIndexMap: b.fieldIndexMap,
	}
}

// Name get the name of the struct.
func (b *Builder) Name() string {
	return b.name
}

// OfStruct checks if the struct is an instance of the builder type.
func (b *Builder) OfStruct(s *Struct) bool {
	if s == nil {
		return false
	}

	if b.name != s.name {
		return false
	}

	if b.structType != s.val.Type() {
		return false
	}

	return true
}

// OfSlice checks if the slice is an instance of the builder type.
func (b *Builder) OfSlice(s *Slice) bool {
	if s == nil {
		return false
	}
	if b.name != s.name {
		return false
	}
	if s.val.Kind() != reflect.Slice {
		return false
	}
	if b.structType != s.val.Type().Elem() {
		return false
	}
	return true
}
