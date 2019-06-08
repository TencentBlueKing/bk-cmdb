/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package datamodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

type CustomField interface {
	GetName() string
	GetKey() string
	GetType() string
	IsRequired() bool
	IsEditable() bool
	// ex: `json:"key" bson:"key" max:"10" min:"1"`
	GetTag() reflect.StructTag
	TypeOf() reflect.Type
}

type FieldBase struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Editable    bool   `json:"editable"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

func (fb *FieldBase) GetKey() string {
	return fb.Key
}

func (fb *FieldBase) GetName() string {
	return fb.Name
}

func (fb *FieldBase) GetType() string {
	return fb.Type
}

func (fb *FieldBase) IsRequired() bool {
	return fb.Required
}

func (fb *FieldBase) IsEditable() bool {
	return fb.Editable
}

type IntField struct {
	FieldBase
	Max *int64 `json:"max"`
	Min *int64 `json:"min"`
}

func (i *IntField) GetTag() reflect.StructTag {
	tag := fmt.Sprintf("json:\"%s\"", i.Key)
	validates := make([]string, 0)
	if i.Max != nil {
		validates = append(validates, fmt.Sprintf("lte:%d", i.Max))
	}
	if i.Max != nil {
		validates = append(validates, fmt.Sprintf("gte:%d", i.Min))
	}
	if len(validates) > 0 {
		tag += fmt.Sprintf(" validate:\"%s\"", strings.Join(validates, ","))
	}
	return reflect.StructTag(tag)
}

func (i *IntField) TypeOf() reflect.Type {
	return reflect.TypeOf(int64(0))
}

type StringField struct {
	FieldBase
	MaxLength *int64 `json:"max_length"`
	MinLength *int64 `json:"min_length"`
	Regex     string `json:"regex"`
}

func (i *StringField) GetTag() reflect.StructTag {
	tag := fmt.Sprintf("json:\"%s\"", i.Key)
	validates := make([]string, 0)
	if i.MaxLength != nil {
		validates = append(validates, fmt.Sprintf("max:%d", i.MaxLength))
	}
	if i.MinLength != nil {
		validates = append(validates, fmt.Sprintf("min:%d", i.MinLength))
	}
	if len(validates) > 0 {
		tag += fmt.Sprintf(" validate:\"%s\"", strings.Join(validates, ","))
	}
	return reflect.StructTag(tag)
}

func (i *StringField) TypeOf() reflect.Type {
	return reflect.TypeOf("")
}

type DynamicStructure struct {
	CustomFields []CustomField
	structValue  interface{}
	structType   reflect.Type
}

func (ds *DynamicStructure) Get(key string) (interface{}, error) {
	rv := reflect.ValueOf(ds.structValue).Elem()
	value := rv.FieldByName(key)
	if value.IsValid() {
		return value.Interface(), nil
	}
	return nil, errors.New("not found")
}

func (ds *DynamicStructure) StructOf() reflect.Type {
	if ds.structType == nil {
		reflectFields := make([]reflect.StructField, 0)
		for _, field := range ds.CustomFields {
			reflectField := reflect.StructField{
				Name: field.GetKey(),
				Type: field.TypeOf(),
				Tag:  field.GetTag(),
			}
			reflectFields = append(reflectFields, reflectField)
		}
		ds.structType = reflect.StructOf(reflectFields)
	}
	return ds.structType
}

func (ds *DynamicStructure) UnmarshalJSON(bs []byte) error {
	structType := ds.StructOf()
	v := reflect.New(structType).Elem()
	ds.structValue = v.Addr().Interface()

	if err := json.Unmarshal(bs, &ds.structValue); err != nil {
		return err
	}

	return nil
}

func (ds *DynamicStructure) MarshalJSON() ([]byte, error) {
	structType := ds.StructOf()
	v := reflect.New(structType).Elem()
	ds.structValue = v.Addr().Interface()

	bs, err := json.Marshal(ds.structValue)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (ds *DynamicStructure) Validate() error {
	validate := validator.New()
	err := validate.Struct(ds.structValue)
	return err
}
