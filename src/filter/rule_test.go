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
	"reflect"
	"testing"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/json"
	"configcenter/src/common/util"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	exampleRule = &CombinedRule{
		Condition: And,
		Rules: []RuleFactory{
			&AtomRule{
				Field:    "test",
				Operator: Equal.Factory(),
				Value:    1,
			},
			&CombinedRule{
				Condition: Or,
				Rules: []RuleFactory{
					&AtomRule{
						Field:    "test1",
						Operator: Array.Factory(),
						Value: &AtomRule{
							Field:    ArrayElement,
							Operator: Object.Factory(),
							Value: &CombinedRule{
								Condition: And,
								Rules: []RuleFactory{
									&AtomRule{
										Field:    "test2",
										Operator: In.Factory(),
										Value:    []string{"b", "c"},
									},
								},
							},
						},
					},
					&AtomRule{
						Field:    "test3",
						Operator: DatetimeLess.Factory(),
						Value:    1,
					},
				},
			},
		},
	}
)

func TestJsonUnmarshalRule(t *testing.T) {
	ruleJson, err := json.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	rule := new(CombinedRule)
	err = json.Unmarshal(ruleJson, rule)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, rule)
}

func TestBsonUnmarshalRule(t *testing.T) {
	ruleBson, err := bson.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	rule := new(CombinedRule)
	err = bson.Unmarshal(ruleBson, rule)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, rule)
}

func TestRuleWithType(t *testing.T) {
	var rule RuleFactory

	rule = new(AtomRule)
	if rule.WithType() != AtomType {
		t.Errorf("rule type %s is invalid", rule.WithType())
		return
	}

	rule = new(CombinedRule)
	if rule.WithType() != CombinedType {
		t.Errorf("rule type %s is invalid", rule.WithType())
		return
	}
}

func TestRuleValidate(t *testing.T) {
	var rule RuleFactory

	// test atomic rule validation
	rule = &AtomRule{
		Field:    "test1",
		Operator: NotIn.Factory(),
		Value:    []string{"a", "b", "c"},
	}

	if err := rule.Validate(NewDefaultExprOpt(map[string]enumor.FieldType{"test1": enumor.String})); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	opt := &ExprOption{
		RuleFields: map[string]enumor.FieldType{
			"test1": enumor.String,
		},
		MaxNotInLimit: 3,
	}

	if err := rule.Validate(opt); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	// test invalid atomic rule scenario
	opt = &ExprOption{
		RuleFields: map[string]enumor.FieldType{
			"test2": enumor.String,
		},
	}

	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}

	opt = &ExprOption{
		MaxRulesLimit: 10,
		MaxNotInLimit: 2,
	}

	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}

	// test combined rule validation
	rule = exampleRule

	// TODO confirm how to deal with object & array
	opt = &ExprOption{
		RuleFields: map[string]enumor.FieldType{
			"test":                enumor.Numeric,
			"test1":               enumor.Array,
			"test1.element":       enumor.Object,
			"test1.element.test2": enumor.String,
			"test3":               enumor.Time,
		},
		MaxInLimit:    2,
		MaxRulesLimit: 2,
		MaxRulesDepth: 6,
	}

	if err := rule.Validate(opt); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	// test invalidate scenario
	opt.RuleFields["test"] = enumor.String
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.RuleFields["test"] = enumor.Numeric

	delete(opt.RuleFields, "test")
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.RuleFields["test"] = enumor.Numeric

	opt.MaxInLimit = 1
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxInLimit = 0

	opt.MaxRulesLimit = 1
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxRulesLimit = 0

	opt.MaxRulesDepth = 5
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxRulesDepth = 0
}

func TestRuleFields(t *testing.T) {
	var rule RuleFactory

	rule = &AtomRule{
		Field:    "test1",
		Operator: Equal.Factory(),
		Value:    1,
	}

	fields := rule.RuleFields()
	if !reflect.DeepEqual(fields, []string{"test1"}) {
		t.Errorf("rule fields %+v is invalid", fields)
		return
	}

	rule = exampleRule
	fields = rule.RuleFields()
	// TODO confirm how to deal with filter object & array
	if !reflect.DeepEqual(fields, []string{"test", "test1.element.test2", "test3"}) {
		t.Errorf("rule fields %+v is invalid", fields)
		return
	}
}

func TestRuleToMgo(t *testing.T) {
	var rule RuleFactory

	// test atomic rule to mongo
	rule = &AtomRule{
		Field:    "test1",
		Operator: NotIn.Factory(),
		Value:    []string{"a", "b", "c"},
	}

	mgo, err := rule.ToMgo(nil)
	if err != nil {
		t.Errorf("covert rule to mongo failed, err: %v", err)
		return
	}

	expectMgo := map[string]interface{}{
		"test1": map[string]interface{}{
			common.BKDBNIN: []string{"a", "b", "c"},
		},
	}

	if !reflect.DeepEqual(mgo, expectMgo) {
		t.Errorf("rule mongo condition %+v is invalid", mgo)
		return
	}

	// test combined rule to mongo
	rule = exampleRule
	mgo, err = rule.ToMgo(nil)
	if err != nil {
		t.Errorf("covert rule to mongo failed, err: %v", err)
		return
	}

	expectMgo = map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{{
			"test": map[string]interface{}{common.BKDBEQ: 1},
		}, {
			common.BKDBOR: []map[string]interface{}{{
				common.BKDBAND: []map[string]interface{}{{
					"test1.test2": map[string]interface{}{common.BKDBIN: []string{"b", "c"}},
				}},
			}, {
				"test3": map[string]interface{}{common.BKDBLT: time.Unix(1, 0)},
			}},
		}},
	}

	if !reflect.DeepEqual(mgo, expectMgo) {
		t.Errorf("rule mongo condition %+v is invalid", mgo)
		return
	}

	// test invalid combined rule to mongo scenario
	rule = &CombinedRule{
		Condition: "test",
		Rules:     []RuleFactory{exampleRule},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}

	rule = &CombinedRule{
		Condition: "",
		Rules:     []RuleFactory{exampleRule},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}

	rule = &CombinedRule{
		Condition: "test",
		Rules: []RuleFactory{
			&AtomRule{
				Field:    "test1",
				Operator: In.Factory(),
				Value:    []interface{}{"a", 1, "c"},
			},
		},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}
}

func testExampleRule(t *testing.T, r RuleFactory) {
	if r == nil {
		t.Errorf("rule is nil")
		return
	}

	rule, ok := r.(*CombinedRule)
	if !ok {
		t.Errorf("rule %+v is not combined type", r)
		return
	}

	if rule.Condition != And {
		t.Errorf("rule condition %s is not and", rule.Condition)
		return
	}

	if len(rule.Rules) != 2 {
		t.Errorf("rules length %d is not 2", len(rule.Rules))
		return
	}

	subAtomRule, ok := rule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("first sub rule %+v is not atom type", rule.Rules[0])
		return
	}

	if subAtomRule.Field != "test" {
		t.Errorf("first sub rule field %s is not test", subAtomRule.Field)
		return
	}

	if subAtomRule.Operator != Equal.Factory() {
		t.Errorf("first sub rule op %s is not equal", subAtomRule.Operator)
		return
	}

	intVal, err := util.GetInt64ByInterface(subAtomRule.Value)
	if err != nil {
		t.Errorf("first sub rule value %v is invalid", subAtomRule.Value)
		return
	}

	if intVal != 1 {
		t.Errorf("first sub rule value %v is not 1", subAtomRule.Value)
		return
	}

	subCombinedRule, ok := rule.Rules[1].(*CombinedRule)
	if !ok {
		t.Errorf("second sub rule %+v is not combined type", rule.Rules[1])
		return
	}

	if subCombinedRule.Condition != Or {
		t.Errorf("second sub rule condition %s is not combined type", rule.Condition)
		return
	}

	if len(subCombinedRule.Rules) != 2 {
		t.Errorf("second sub rules length %d is not 2", len(rule.Rules))
		return
	}

	subAtomRule1, ok := subCombinedRule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("first sub sub rule %+v is not atom type", subCombinedRule.Rules[0])
		return
	}

	if subAtomRule1.Field != "test1" {
		t.Errorf("first sub sub rule field %s is not test1", subAtomRule1.Field)
		return
	}

	if subAtomRule1.Operator != Array.Factory() {
		t.Errorf("first sub sub rule op %s is not ne", subAtomRule1.Operator)
		return
	}

	filterArrVal, ok := subAtomRule1.Value.(*AtomRule)
	if !ok {
		t.Errorf("first sub sub rule value %v is invalid", subAtomRule1.Value)
		return
	}

	if filterArrVal.Field != ArrayElement {
		t.Errorf("filter array rule field %s is not %s", subAtomRule.Field, ArrayElement)
		return
	}

	if filterArrVal.Operator != Object.Factory() {
		t.Errorf("filter array rule op %s is not filter object", subAtomRule.Operator)
		return
	}

	filterArrValRule, ok := filterArrVal.Value.(*CombinedRule)
	if !ok {
		t.Errorf("filter array rule value %v is invalid", filterArrVal.Value)
		return
	}

	if filterArrValRule.Condition != And {
		t.Errorf("filter array sub condition %s is not and", filterArrValRule.Condition)
		return
	}

	if len(filterArrValRule.Rules) != 1 {
		t.Errorf("filter array sub rules length %d is not 1", len(rule.Rules))
		return
	}

	filterObjRule, ok := filterArrValRule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("filter object rule %+v is not atom type", filterArrValRule.Rules[0])
		return
	}

	if filterObjRule.Field != "test2" {
		t.Errorf("filter object rule field %s is not test2", filterObjRule.Field)
		return
	}

	if filterObjRule.Operator != In.Factory() {
		t.Errorf("filter object rule op %s is not in", filterObjRule.Operator)
		return
	}

	arrVal, ok := filterObjRule.Value.([]interface{})
	if !ok {
		t.Errorf("filter object rule value %v is invalid", filterObjRule.Value)
		return
	}

	if len(arrVal) != 2 {
		t.Errorf("array value length %d is not 2", len(arrVal))
		return
	}

	strVal1, ok := arrVal[0].(string)
	if !ok {
		t.Errorf("first array value %v is invalid", arrVal[0])
		return
	}

	if strVal1 != "b" {
		t.Errorf("first array value %v is not b", arrVal[0])
		return
	}

	strVal2, ok := arrVal[1].(string)
	if !ok {
		t.Errorf("second array value %v is invalid", arrVal[1])
		return
	}

	if strVal2 != "c" {
		t.Errorf("second array value %v is not c", arrVal[1])
		return
	}

	subAtomRule2, ok := subCombinedRule.Rules[1].(*AtomRule)
	if !ok {
		t.Errorf("second sub sub rule %+v is not atom type", subCombinedRule.Rules[1])
		return
	}

	if subAtomRule2.Field != "test3" {
		t.Errorf("second sub sub rule field %s is not test3", subAtomRule2.Field)
		return
	}

	if subAtomRule2.Operator != DatetimeLess.Factory() {
		t.Errorf("second sub sub rule op %s is not datetime less", subAtomRule2.Operator)
		return
	}

	if !util.IsNumeric(subAtomRule2.Value) {
		t.Errorf("second sub rule value %v is invalid", subAtomRule2.Value)
		return
	}
}
