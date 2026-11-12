/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package operator

import (
	"testing"

	"configcenter/src/common/json"
)

func TestPlan_MarshalJSON(t *testing.T) {
	p := &Plan{
		Kind: ConditionalKind,
		Condition: &AuthCondition{
			Operator: And,
			Element: &Content{
				Content: []*AuthCondition{
					{
						Operator: Equal,
						Element: &FieldValue{
							Field: AttrField("os"),
							Value: "linux",
						},
					},
					{
						Operator: Or,
						Element: &Content{
							Content: []*AuthCondition{
								{
									Operator: In,
									Element: &FieldValue{
										Field: "biz",
										Value: []string{"1"},
									},
								},
								{
									Operator: Equal,
									Element: &FieldValue{
										Field: AttrField("owner"),
										Value: "zhangsan",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	shouldBe := `{"kind":"CONDITIONAL","condition":{"op":"and","content":[{"op":"eq","field":"attr.os","value":"linux"},` +
		`{"op":"or","content":[{"op":"in","field":"biz","value":["1"]},` +
		`{"op":"eq","field":"attr.owner","value":"zhangsan"}]}]}}`
	if string(js) != shouldBe {
		t.Fatalf("invalid marshal, got: %s", js)
	}
}

func TestPlan_UnmarshalJSON(t *testing.T) {
	src := `
{
    "kind": "CONDITIONAL",
    "condition": {
        "op": "and",
        "content": [
            {"field": "attr.os", "op": "eq", "value": "linux"},
            {"op": "or", "content": [
                {"field": "biz", "op": "in", "value": ["1"]},
                {"field": "attr.owner", "op": "eq", "value": "zhangsan"}
            ]}
        ]
    }
}
`
	p := new(Plan)
	if err := json.Unmarshal([]byte(src), p); err != nil {
		t.Fatal(err)
	}

	if p.Kind != ConditionalKind {
		t.Fatal("parse kind failed")
	}
	if p.Condition == nil || p.Condition.Operator != And {
		t.Fatal("parse and operator failed")
	}

	content, ok := p.Condition.Element.(*Content)
	if !ok {
		t.Fatal("parse Content failed")
	}
	if len(content.Content) != 2 {
		t.Fatal("parse content, but got invalid length")
	}

	eqCond := content.Content[0]
	if eqCond.Operator != Equal {
		t.Fatal("parse eq operator failed")
	}
	eqFV, ok := eqCond.Element.(*FieldValue)
	if !ok {
		t.Fatal("parse eq FieldValue failed")
	}
	if eqFV.Field != AttrField("os") || eqFV.Value != "linux" {
		t.Fatal("parse eq condition failed")
	}

	orCond := content.Content[1]
	if orCond.Operator != Or {
		t.Fatal("parse or operator failed")
	}
	orContent, ok := orCond.Element.(*Content)
	if !ok {
		t.Fatal("parse or Content failed")
	}
	if len(orContent.Content) != 2 {
		t.Fatal("parse or content length failed")
	}

	inFV, ok := orContent.Content[0].Element.(*FieldValue)
	if !ok || orContent.Content[0].Operator != In || inFV.Field != "biz" {
		t.Fatal("parse in condition failed")
	}
	inValues, ok := inFV.Value.([]interface{})
	if !ok || len(inValues) != 1 || inValues[0] != "1" {
		t.Fatal("parse in value failed")
	}
	ownerFV, ok := orContent.Content[1].Element.(*FieldValue)
	if !ok || orContent.Content[1].Operator != Equal || ownerFV.Field != AttrField("owner") ||
		ownerFV.Value != "zhangsan" {
		t.Fatal("parse owner eq condition failed")
	}
}

func TestPlan_AlwaysAllowedDenied(t *testing.T) {
	src := `{"kind":"ALWAYS_ALLOWED","condition":null}`
	p := new(Plan)
	if err := json.Unmarshal([]byte(src), p); err != nil {
		t.Fatal(err)
	}
	if p.Kind != AlwaysAllowedKind || p.Condition != nil {
		t.Fatal("parse always allowed plan failed")
	}

	src = `{"kind":"ALWAYS_DENIED","condition":null}`
	p = new(Plan)
	if err := json.Unmarshal([]byte(src), p); err != nil {
		t.Fatal(err)
	}
	if p.Kind != AlwaysDeniedKind || p.Condition != nil {
		t.Fatal("parse always denied plan failed")
	}
}

func TestAuthOp_IsLogical(t *testing.T) {
	if !And.IsLogical() || !Or.IsLogical() || !Not.IsLogical() {
		t.Fatal("logical operators should return true")
	}
	if Equal.IsLogical() || In.IsLogical() || StartWith.IsLogical() {
		t.Fatal("compare operators should return false")
	}
}

func TestFieldKind(t *testing.T) {
	id := Field(IamIDKey)
	if !id.IsID() || id.IsAncestor() || id.IsSelfAttr() {
		t.Fatal("id field classification failed")
	}

	ancestor := Field("biz")
	if !ancestor.IsAncestor() || ancestor.IsID() || ancestor.IsSelfAttr() {
		t.Fatal("ancestor field classification failed")
	}

	attr := AttrField("os")
	if !attr.IsSelfAttr() || attr.IsID() || attr.IsAncestor() || attr.SelfAttr() != "os" {
		t.Fatal("self attribute field classification failed")
	}
	if string(attr) != IamAttrPrefix+"os" {
		t.Fatal("self attribute field prefix failed")
	}
}
