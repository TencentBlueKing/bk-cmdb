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

package dal

import (
	"errors"
	"time"
)

// FieldType define the table's field data type.
type FieldType string

const (
	// Numeric means this field is numeric data type.
	Numeric FieldType = "numeric"
	// Boolean means this field is boolean data type.
	Boolean FieldType = "bool"
	// String means this field is string data type.
	String FieldType = "string"
	// MapString means this field is map string type.there is a special map for
	// container scenarios, in which both key and value must be of string type,
	// such as label, taints, etc. At this time, these variables need to be set
	// to this type to facilitate judgment and verification.
	MapString FieldType = "mapString"
	// Array means this field is array data type.
	Array FieldType = "array"
	// Object means this field is object data type.
	Object FieldType = "object"
	// Enum means this field is enum enum type.
	Enum FieldType = "enum"
	// Note: subsequent support for other types can be added here.
	// after adding a type, pay attention to adding a verification
	// function for this type synchronously. special attention is
	// paid to whether the array elements also need to synchronize support for this type.
)

// Fields table's fields details.
type Fields struct {
	// descriptors specific description of the field.
	descriptors []FieldDescriptor
	// fields defines all the table's fields.
	fields []string
	// fieldType the type corresponding to the field.
	fieldType map[string]FieldType
}

// FieldsDescriptors table of field descriptor.
type FieldsDescriptors []FieldDescriptor

func mergeFields(all ...FieldsDescriptors) *Fields {
	result := &Fields{
		descriptors: make([]FieldDescriptor, 0),
		fields:      make([]string, 0),
		fieldType:   make(map[string]FieldType),
	}

	if len(all) == 0 {
		return result
	}

	for _, col := range all {
		for _, f := range col {
			result.descriptors = append(result.descriptors, f)
			result.fieldType[f.Field] = f.Type
			result.fields = append(result.fields, f.Field)
		}
	}
	return result
}

// FieldsType returns the corresponding type of all fields.
func (f Fields) FieldsType() map[string]FieldType {
	copied := make(map[string]FieldType)
	for k, v := range f.fieldType {
		copied[k] = v
	}

	return copied
}

// OneFieldType returns the type corresponding to the specified field.
func (f Fields) OneFieldType(field string) FieldType {
	return f.fieldType[field]
}

// FieldsDescriptor returns table's all fields descriptor.
func (f Fields) FieldsDescriptor() []FieldDescriptor {
	return f.descriptors
}

// OneFieldDescriptor returns one field's descriptor.
func (f Fields) OneFieldDescriptor(field string) FieldDescriptor {
	if field == "" {
		return FieldDescriptor{}
	}

	for idx := range f.descriptors {
		if f.descriptors[idx].Field == field {
			return f.descriptors[idx]
		}
	}
	return FieldDescriptor{}
}

// Fields returns all the table's fields.
func (f Fields) Fields() []string {
	copied := make([]string, len(f.fields))
	for idx := range f.fields {
		copied[idx] = f.fields[idx]
	}
	return copied
}

// mergeFieldDescriptors merge all fields of a table together.
func mergeFieldDescriptors(resources ...FieldsDescriptors) FieldsDescriptors {
	if len(resources) == 0 {
		return make([]FieldDescriptor, 0)
	}

	merged := make([]FieldDescriptor, 0)
	for _, one := range resources {
		merged = append(merged, one...)
	}

	return merged
}

// FieldDescriptor defines a table's field related information.
type FieldDescriptor struct {
	// Field is field's name.
	Field string
	// Type is this field's data type.
	Type FieldType
	// IsRequired is it required.
	IsRequired bool
	// IsEditable is it editable.
	IsEditable bool
	// Option additional information for the field.
	// the content corresponding to different fields may be different.
	Option interface{}
	_      struct{}
}

// Revision resource revision information.
type Revision struct {
	Creator    string `json:"creator" bson:"creator"`
	Modifier   string `json:"modifier" bson:"modifier"`
	CreateTime int64  `json:"create_time" bson:"create_time"`
	LastTime   int64  `json:"last_time" bson:"last_time"`
}

// lagSeconds fault tolerance for ntp errors of different devices.
const lagSeconds = 5 * 60

// ValidateCreate validation of parameters in the creation scene.
func (r Revision) ValidateCreate() error {

	if len(r.Creator) == 0 || len(r.Modifier) == 0 {
		return errors.New("creator can not be empty")
	}

	if r.Creator != r.Modifier {
		return errors.New("creator can not be empty")
	}

	if r.CreateTime == 0 {
		return errors.New("create time must be set")
	}

	now := time.Now().Unix()
	if (r.CreateTime <= (now - lagSeconds)) || (r.CreateTime >= (now + lagSeconds)) {
		return errors.New("invalid create time")
	}

	return nil
}

// ValidateUpdate validate revision when updated.
func (r Revision) ValidateUpdate() error {
	if len(r.Modifier) == 0 {
		return errors.New("reviser can not be empty")
	}

	if len(r.Creator) != 0 {
		return errors.New("creator can not be updated")
	}

	now := time.Now().Unix()
	if (r.LastTime <= (now - lagSeconds)) || (r.LastTime >= (now + lagSeconds)) {
		return errors.New("invalid update time")
	}

	if r.LastTime < r.CreateTime-lagSeconds {
		return errors.New("update time must be later than create time")
	}
	return nil
}
