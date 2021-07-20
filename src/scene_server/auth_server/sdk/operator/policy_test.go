/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operator

import (
	"testing"

	"configcenter/src/common/json"
)

func TestPolicy_MarshalJSON(t *testing.T) {
	p := &Policy{
		Operator: "AND",
		Element: &Content{
			Content: []*Policy{
				{
					Operator: "eq",
					Element: &FieldValue{
						Field: Field{
							Resource:  "host",
							Attribute: "os",
						},
						Value: "linux",
					},
				},
				{
					Operator: "OR",
					Element: &Content{
						Content: []*Policy{
							{
								Operator: "neq",
								Element: &FieldValue{
									Field: Field{
										Resource:  "host",
										Attribute: "os",
									},
									Value: "windows",
								},
							},
						},
					},
				},
			},
		},
	}

	js, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	shouldBe := "{\"op\":\"AND\",\"content\":[{\"op\":\"eq\",\"field\":\"host.os\",\"value\":\"linux\"},{\"op\":\"OR\",\"content\":[{\"op\":\"neq\",\"field\":\"host.os\",\"value\":\"windows\"}]}]}"

	if string(js) != shouldBe {
		t.Fatal("invalid unmarshal")
	}

}

func TestPolicy_UnmarshalJSON(t *testing.T) {
	src := `
{
    "op":"AND",
    "content":[
        {
            "op":"eq",
            "field":"host.os",
            "value":"linux"
        },
        {
            "op":"OR",
            "content":[
                {
                    "op":"neq",
                    "field":"host.os",
                    "value":"windows"
                }
            ]
        }
    ]
}
`
	p := new(Policy)
	if err := json.Unmarshal([]byte(src), p); err != nil {
		t.Fatal(err)
		return
	}

	if p.Operator != And {
		t.Fatal("parse AND oper failed")
	}

	content, ok := p.Element.(*Content)
	if !ok {
		t.Fatal("parse Content failed")
	}

	if len(content.Content) != 2 {
		t.Fatal("parse Content, but got invalid length")
	}

	// check element 0, eq operator
	if content.Content[0].Operator != "eq" {
		t.Fatal("parse content.eq operator failed")
	}

	eqPolicy, ok := content.Content[0].Element.(*FieldValue)
	if !ok {
		t.Fatal("parse content.FieldValue failed")
	}

	if eqPolicy.Field.Resource != "host" ||
		eqPolicy.Field.Attribute != "os" ||
		eqPolicy.Value != "linux" {
		t.Fatal("parse eq policy failed")
	}

	// check element 1, neq operator
	if content.Content[1].Operator != Or {
		t.Fatal("parse or operator failed")
	}

	contentOR, ok := content.Content[1].Element.(*Content)
	if !ok {
		t.Fatal("parse content OR failed")
	}

	if len(contentOR.Content) != 1 {
		t.Fatal("parse content OR length failed")
	}

	if contentOR.Content[0].Operator != "neq" {
		t.Fatal("parse content.neq operator failed")
	}

	neqPolicy, ok := contentOR.Content[0].Element.(*FieldValue)
	if !ok {
		t.Fatal("parse content.FieldValue failed")
	}

	if neqPolicy.Field.Resource != "host" ||
		neqPolicy.Field.Attribute != "os" ||
		neqPolicy.Value != "windows" {
		t.Fatal("parse neq policy failed")
	}

}

func TestField_UnmarshalJSON(t *testing.T) {
	normal := `{"field":"host.os", "value":"windows"}`
	f := new(FieldValue)
	if err := json.Unmarshal([]byte(normal), f); err != nil {
		t.Fatal(err)
	}

	if f.Field.Resource != "host" {
		t.Fatal("parse host resource failed")
	}

	if f.Field.Attribute != "os" {
		t.Fatal("parse host os attribute failed")
	}

	if f.Value != "windows" {
		t.Fatal("parse host value failed")
	}

	abnormalAttribute := `{"field":"host.", "value":"windows"}`
	if err := json.Unmarshal([]byte(abnormalAttribute), f); err == nil {
		t.Fatal("should parse failed with empty attribute failed.")
	}

	abnormalResource := `{"field":".os", "value":"windows"}`
	if err := json.Unmarshal([]byte(abnormalResource), f); err == nil {
		t.Fatal("should parse failed with empty resource failed.")
	}

	abnormalDot := `{"field":"host:os", "value":"windows"}`
	if err := json.Unmarshal([]byte(abnormalDot), f); err == nil {
		t.Fatal("should parse failed with empty dot failed.")
	}

}

func TestField_MarshalJSON(t *testing.T) {
	f := &FieldValue{
		Field: Field{
			Resource:  "host",
			Attribute: "os",
		},
		Value: "windows",
	}
	js, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal FieldVaile failed, err: %v", err)
	}

	shouldBe := `{"field":"host.os","value":"windows"}`
	if string(js) != shouldBe {
		t.Fatal("filed value marshaled failed")
	}
}
