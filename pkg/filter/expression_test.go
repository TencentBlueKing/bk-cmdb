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
package filter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	exprJson := `
{
	"op": "and",
	"rules": [{
			"field": "deploy_type",
			"op": "eq",
			"value": "common"
		},
		{
			"field": "creator",
			"op": "eq",
			"value": "tom"
		}
	]
}
`
	expr := new(Expression)
	err := expr.UnmarshalJSON([]byte(exprJson))
	if err != nil {
		t.Error(err)
		return
	}
	expected := ExpressionAnd(
		RuleEqual("deploy_type", "common"),
		RuleEqual("creator", "tom"),
	)
	assert.Equal(t, expected, expr, "expression is not expected")
}

func TestExpressionValidateOption(t *testing.T) {
	expr := ExpressionAnd(
		RuleEqual("name", "cmdb"),
		RuleGreaterThan("age", 18),
		RuleLessThan("age", 30),
		RuleIn("servers", []string{"api", "web"}),
		RuleEqual("asDefault", true),
		RuleGreaterThan("created_at", "2006-01-02T15:04:05Z"),
	)

	fields := map[string]FieldType{
		"name":       String,
		"age":        Numeric,
		"servers":    String,
		"asDefault":  Boolean,
		"created_at": Time,
	}
	opt := NewExprOption(RuleFields(fields), MaxInLimit(1), MaxRulesLimit(5))

	if err := expr.Validate(opt); !strings.Contains(err.Error(), "rules elements number exceeded, it at most have 5 rules") {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}
	opt.MaxRulesLimit = 6

	if err := expr.Validate(opt); !strings.Contains(err.Error(), "invalid in operator's value, at most have 1 elements") {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}

	opt.MaxInLimit = 0
	if err := expr.Validate(opt); err != nil {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}

	// test invalidate scenario
	opt.RuleFields["name"] = Numeric
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a numeric") {
		t.Errorf("validate numeric type failed, err: %v", err)
		return
	}
	opt.RuleFields["name"] = String

	opt.RuleFields["age"] = String
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a string") {
		t.Errorf("validate string type failed, err: %v", err)
		return
	}
	opt.RuleFields["age"] = Numeric

	opt.RuleFields["asDefault"] = Time
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a string time format") {
		t.Errorf("validate time type failed, err: %v", err)
		return
	}
	opt.RuleFields["asDefault"] = Boolean

	opt.RuleFields["created_at"] = Boolean
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a boolean") {
		t.Errorf("validate boolean type failed, err: %v", err)
		return
	}
}

func TestAllExpression(t *testing.T) {
	expr := ExpressionAnd(
		AllExpression(),
		ExpressionOr(RuleEqual("name", "cmdb"), RuleGreaterThan("age", 18)),
	)
	fields := map[string]FieldType{
		"name":       String,
		"age":        Numeric,
		"servers":    String,
		"asDefault":  Boolean,
		"created_at": Time,
	}
	opt := NewExprOption(RuleFields(fields), MaxInLimit(1), MaxRulesLimit(5))
	if err := expr.Validate(opt); err != nil {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}

}

func TestWildcardExpressionValidateOpt(t *testing.T) {
	fieldMap := map[string]FieldType{
		"extension.name":     String,
		"extension.*.field2": String,
	}
	opt := NewExprOption(RuleFields(fieldMap))
	t.Run("single_dot", func(t *testing.T) {
		expr := ExpressionAnd(RuleJSONEqual("extension.field1", "cmdb"))

		if err := expr.Validate(opt); !strings.Contains(err.Error(),
			"rule field(extension.field1) should not exist(not supported)") {
			t.Errorf("validate field failed, err: %v", err)
			return
		}
		opt.RuleFields["extension.*"] = String
		if err := expr.Validate(opt); err != nil {
			t.Errorf("validate wildcard expression field, err: %v", err)
			return
		}
	})

	t.Run("second_dot_wildcard", func(t *testing.T) {
		expr := RuleEqual("extension.field1.field2", "cmdb")
		if err := expr.Validate(opt); !strings.Contains(err.Error(),
			"rule field(extension.field1.field2) should not exist(not supported)") {
			t.Errorf("validate field failed, err: %v", err)
			return
		}

		opt.RuleFields["extension.field1.*"] = String
		if err := expr.Validate(opt); err != nil {
			t.Errorf("validate wildcard expression failed, err: %v", err)
			return
		}
	})

}

func TestExprConstructionFunc(t *testing.T) {
	cases := []struct {
		name string
		expr *Expression
		want *Expression
	}{
		{
			name: "simple and",
			expr: ExpressionAnd(
				RuleEqual("name", "cmdb"),
				RuleGreaterThan("age", 18),
			),
			want: &Expression{
				Op: And,
				Rules: []RuleFactory{
					&AtomRule{
						Field: "name",
						Op:    Equal.Factory(),
						Value: "cmdb",
					},
					&AtomRule{
						Field: "age",
						Op:    GreaterThan.Factory(),
						Value: 18,
					},
				},
			},
		},
		{
			name: "simple or with sub expression",
			expr: ExpressionAnd(
				RuleEqual("name", "cmdb"),
				ExpressionOr(
					RuleGreaterThan("age", 18),
					RuleLessThan("height", 1.8),
				),
			),
			want: &Expression{
				Op: And,
				Rules: []RuleFactory{
					&AtomRule{
						Field: "name",
						Op:    Equal.Factory(),
						Value: "cmdb",
					},
					&Expression{
						Op: Or,
						Rules: []RuleFactory{
							&AtomRule{
								Field: "age",
								Op:    GreaterThan.Factory(),
								Value: 18,
							},
							&AtomRule{
								Field: "height",
								Op:    LessThan.Factory(),
								Value: 1.8,
							},
						},
					},
				},
			},
		},
		{
			name: "simple or with json equal",
			expr: ExpressionAnd(
				RuleEqual("name", "cmdb"),
				RuleJSONEqual("extension.vpc.id", 3),
			),
			want: &Expression{
				Op: And,
				Rules: []RuleFactory{
					&AtomRule{
						Field: "name",
						Op:    Equal.Factory(),
						Value: "cmdb",
					},
					&AtomRule{
						Field: "extension.vpc.id",
						Op:    JSONEqual.Factory(),
						Value: 3,
					},
				},
			},
		},
		{
			name: "json_contains",
			expr: ExpressionAnd(
				RuleJSONContains("managers", "cmdb"),
			),
			want: &Expression{
				Op: And,
				Rules: []RuleFactory{
					&AtomRule{
						Field: "managers",
						Op:    JSONContains.Factory(),
						Value: "cmdb",
					},
				},
			},
		},
	}
	for _, c := range cases {
		assert.Equalf(t, c.want, c.expr, "expression is not expected, name: %s", c.name)
	}
}

func genWithDepth(depth int) *Expression {
	rule := RuleEqual("depth", depth)
	if depth == 1 {
		return ExpressionAnd(rule)
	}
	if depth&1 == 1 {
		return ExpressionAnd(
			genWithDepth(depth-1),
			rule,
		)
	} else {
		return ExpressionOr(
			genWithDepth(depth-1),
			rule,
		)
	}
}

func TestExpressionDepth(t *testing.T) {
	exp1 := genWithDepth(1)
	assert.Equal(t, ExpressionAnd(RuleEqual("depth", 1)), exp1, "depth1")

	exp2 := genWithDepth(2)
	assert.Equal(t,
		ExpressionOr(
			ExpressionAnd(RuleEqual("depth", 1)),
			RuleEqual("depth", 2)),
		exp2)

	exp := genWithDepth(int(DefaultMaxDepth))

	fields := map[string]FieldType{"name": String, "depth": Numeric}
	opt := NewExprOption(MaxDepth(DefaultMaxDepth-1), RuleFields(fields))
	err := exp.Validate(opt)
	assert.NotNil(t, err, "max depth-1")
	assert.ErrorContains(t, err, "expression depth exceeded")

	opt = opt.CopyWith(MaxDepth(DefaultMaxDepth))
	err = exp.Validate(opt)
	assert.Nil(t, err, "max depth")

	opt = NewExprOption(RuleFields(fields))
	err = exp.Validate(opt)
	assert.Nil(t, err, "default max depth without set depth")

	err = exp.Validate(nil)
	assert.ErrorContains(t, err, "rule field(depth) should not exist(not supported)")
}

func TestEmptyExpression(t *testing.T) {
	exp := &Expression{}
	opt := NewExprOption()
	if err := exp.Validate(opt); err == nil {
		t.Error("empty expression should be invalid")
		return
	}
	exp = ExpressionAnd()
	if err := exp.Validate(opt); err == nil {
		t.Error("empty and expression should be invalid")
		return
	}

	exp = ExpressionOr()
	if err := exp.Validate(opt); err == nil {
		t.Error("empty or expression should be invalid")
		return
	}

	exp = AllExpression()
	if err := exp.Validate(opt); err != nil {
		t.Error("all expression is ok for empty rule")
		return
	}

}
