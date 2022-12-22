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

package filter

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"configcenter/src/common/criteria/enumor"

	"go.mongodb.org/mongo-driver/bson"
)

func TestJsonMarshal(t *testing.T) {
	ruleJson, err := json.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr := Expression{
		RuleFactory: exampleRule,
	}
	exprJson, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprJson, ruleJson)
		return
	}
}

func TestJsonMarshalNil(t *testing.T) {
	// check if nil expression json equals nil combined rule json
	var rule *CombinedRule
	ruleJson, err := json.Marshal(rule)
	if err != nil {
		t.Error(err)
		return
	}

	var expr *Expression
	exprJson, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprJson, ruleJson)
		return
	}

	// check if expression with nil combined rule json equals nil combined rule json
	expr = &Expression{
		RuleFactory: rule,
	}
	exprJson, err = json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprJson, ruleJson)
		return
	}

	// check if expression with nil atom rule json equals nil atom rule json
	var atomRule *AtomRule
	ruleJson, err = json.Marshal(atomRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr = &Expression{
		RuleFactory: atomRule,
	}
	exprJson, err = json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprJson, ruleJson)
		return
	}
}

func TestJsonUnmarshal(t *testing.T) {
	exampleExpr := Expression{
		RuleFactory: exampleRule,
	}
	exprJson, err := json.Marshal(exampleExpr)
	if err != nil {
		t.Error(err)
		return
	}

	expr := new(Expression)
	err = json.Unmarshal(exprJson, expr)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, expr.RuleFactory)
}

func TestBsonMarshal(t *testing.T) {
	// TODO test bson marshal Expression value as well as pointer value if bson supports nil pointer with MarshalBSON
	ruleBson, err := bson.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr := &Expression{
		RuleFactory: exampleRule,
	}
	exprBson, err := bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprBson, ruleBson)
		return
	}
}

func TestBsonMarshalNil(t *testing.T) {
	// check if nil expression bson equals nil combined rule bson
	var rule *CombinedRule
	ruleBson, err := bson.Marshal(rule)
	if err != nil {
		t.Error(err)
		return
	}

	var expr *Expression
	exprBson, err := bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprBson, ruleBson)
		return
	}

	// check if expression with nil combined rule bson equals nil combined rule bson
	expr = &Expression{
		RuleFactory: rule,
	}
	exprBson, err = bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprBson, ruleBson)
		return
	}

	// check if expression with nil atom rule bson equals nil atom rule bson
	var atomRule *AtomRule
	ruleBson, err = bson.Marshal(atomRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr = &Expression{
		RuleFactory: atomRule,
	}
	exprBson, err = bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf("expression marshal result %s is not equal to rule marshal result %s", exprBson, ruleBson)
		return
	}
}

func TestBsonUnmarshal(t *testing.T) {
	exampleExpr := &Expression{
		RuleFactory: exampleRule,
	}

	exprBson, err := bson.Marshal(exampleExpr)
	if err != nil {
		t.Error(err)
		return
	}

	expr := new(Expression)
	err = bson.Unmarshal(exprBson, expr)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, expr.RuleFactory)
}

func TestExpressionValidateOption(t *testing.T) {
	expr := &Expression{
		RuleFactory: &CombinedRule{
			Condition: And,
			Rules: []RuleFactory{
				&AtomRule{
					Field:    "string",
					Operator: Equal.Factory(),
					Value:    "a",
				},
				&CombinedRule{
					Condition: Or,
					Rules: []RuleFactory{
						&AtomRule{
							Field:    "int",
							Operator: Greater.Factory(),
							Value:    123,
						},
						&AtomRule{
							Field:    "enum_array",
							Operator: In.Factory(),
							Value:    []string{"b", "c"},
						},
						&AtomRule{
							Field:    "int_array",
							Operator: NotIn.Factory(),
							Value:    []int64{1, 3, 5},
						},
						&AtomRule{
							Field:    "bool",
							Operator: NotEqual.Factory(),
							Value:    false,
						},
					},
				},
				&AtomRule{
					Field:    "time",
					Operator: DatetimeLessOrEqual.Factory(),
					Value:    time.Now().Unix(),
				},
				&AtomRule{
					Field:    "time",
					Operator: DatetimeGreater.Factory(),
					Value:    "2006-01-02 15:04:05",
				},
				// TODO confirm how to deal with filter object & array
			},
		},
	}

	opt := NewDefaultExprOpt(map[string]enumor.FieldType{
		"string":     enumor.String,
		"int":        enumor.Numeric,
		"enum_array": enumor.Enum,
		"int_array":  enumor.Numeric,
		"bool":       enumor.Boolean,
		"time":       enumor.Time,
	})

	if err := expr.Validate(opt); err != nil {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}

	// test invalidate scenario
	opt.RuleFields["string"] = enumor.Numeric
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a numeric") {
		t.Errorf("validate numeric type failed, err: %v", err)
		return
	}
	opt.RuleFields["string"] = enumor.String

	opt.RuleFields["int"] = enumor.String
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a string") {
		t.Errorf("validate string type failed, err: %v", err)
		return
	}
	opt.RuleFields["int"] = enumor.Numeric

	opt.RuleFields["enum_array"] = enumor.Boolean
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a boolean") {
		t.Errorf("validate bool type failed, err: %v", err)
		return
	}
	opt.RuleFields["enum_array"] = enumor.String

	opt.RuleFields["bool"] = enumor.Time
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "is not of time type") {
		t.Errorf("validate time type failed, err: %v", err)
		return
	}
	opt.RuleFields["bool"] = enumor.Boolean

	opt.RuleFields["time"] = enumor.Boolean
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a boolean") {
		t.Errorf("validate boolean type failed, err: %v", err)
		return
	}
	opt.RuleFields["time"] = enumor.Time

	opt.MaxRulesDepth = 2
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "expression rules depth exceeds maximum") {
		t.Errorf("validate rule depth failed, err: %v", err)
		return
	}
	opt.MaxRulesDepth = 3

	opt.MaxRulesLimit = 3
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "rules elements number exceeds limit: 3") {
		t.Errorf("validate rule limit failed, err: %v", err)
		return
	}
	opt.MaxRulesLimit = 4

	opt.MaxInLimit = 1
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "elements length 2 exceeds maximum 1") {
		t.Errorf("validate rule in limit failed, err: %v", err)
		return
	}
	opt.MaxInLimit = 2

	opt.MaxNotInLimit = 2
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "elements length 3 exceeds maximum 2") {
		t.Errorf("validate rule in limit failed, err: %v", err)
		return
	}
	opt.MaxNotInLimit = 3
}
