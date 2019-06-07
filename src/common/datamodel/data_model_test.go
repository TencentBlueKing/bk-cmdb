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
	"testing"
)

type SimpleModelBase struct {
	ID   int64
	Name string
}
type SimpleModel struct {
	DynamicStructure
	SimpleModelBase
}

func (sm *SimpleModel) UnmarshalJSON(bs []byte) error {
	data := SimpleModel{}
	if err := json.Unmarshal(bs, &data.SimpleModelBase); err != nil {
		return err
	}
	if err := json.Unmarshal(bs, &sm.DynamicStructure); err != nil {
		return err
	}
	sm.Name = data.Name
	sm.ID = data.ID
	return nil
}

func NewSimpleModel(fields []CustomField) *SimpleModel {
	sm := SimpleModel{}
	sm.DynamicStructure.CustomFields = fields
	return &sm
}

func TestLoad(t *testing.T) {
	structFields := make([]CustomField, 0)

	// define a int field
	maxValue := int64(100)
	minValue := int64(1)
	field1 := &IntField{
		FieldBase: FieldBase{
			Key:         "Custom1",
			Name:        "Field 1",
			Type:        "int",
			Required:    false,
			Editable:    false,
			Unit:        "test..",
			Description: "this is a description",
		},
		Max: &maxValue,
		Min: &minValue,
	}
	structFields = append(structFields, field1)

	// define a string field
	maxLength := int64(20)
	minLength := int64(1)
	field2 := &StringField{
		FieldBase: FieldBase{
			Key:         "Custom2",
			Name:        "field 2",
			Type:        "int",
			Required:    false,
			Editable:    false,
			Unit:        "test..",
			Description: "this is a description",
		},
		MaxLength: &maxLength,
		MinLength: &minLength,
	}
	structFields = append(structFields, field2)

	sm := NewSimpleModel(structFields)
	blob := `{"id":666,"name":"a simple test case","Custom1":10,"custom2":"field a value"}`
	if err := json.Unmarshal([]byte(blob), &sm); err != nil {
		t.Fatalf("Unmarshal SimpleModel failed, err: %+v", err)
		return
	}
	t.Logf("ID field: %d", sm.ID)
	t.Logf("Name field: %s", sm.Name)

	custom1, err := sm.Get("Custom1")
	if err != nil {
		t.Fatalf("get field Custom1 failed, err: %+v", err)
	}
	t.Logf("Custom1 field: %d", custom1.(int64))

	custom2, err := sm.Get("Custom2")
	if err != nil {
		t.Fatalf("get field Custom2 failed, err: %+v", err)
	}
	t.Logf("Custom2 field: %s", custom2.(string))
}
