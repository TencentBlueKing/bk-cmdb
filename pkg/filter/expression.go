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
	"errors"
	"fmt"

	"configcenter/src/common/criteria/enumor"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// MaxRulesDepth defines the maximum number of rules depth
	MaxRulesDepth = uint(3)
)

// ExprOption defines how to validate an expression.
type ExprOption struct {
	// RuleFields:
	// 1. used to test if all the expression rule's field
	//    is in the RuleFields' key restricts.
	// 2. all the expression's rule field should be a sub-set
	//    of the RuleFields' key.
	RuleFields map[string]enumor.FieldType
	// IgnoreRuleFields defines if expression rule field validation needs to be ignored.
	IgnoreRuleFields bool
	// MaxInLimit defines the maximum element of the in operator
	MaxInLimit uint
	// MaxNotInLimit defines the maximum element of the nin operator
	MaxNotInLimit uint
	// MaxRulesLimit defines the maximum number of rules an expression allows.
	MaxRulesLimit uint
	// MaxRulesDepth defines the maximum depth of rules an expression allows.
	MaxRulesDepth uint
}

// NewDefaultExprOpt init an expression option with default limit option.
func NewDefaultExprOpt(ruleFields map[string]enumor.FieldType) *ExprOption {
	return &ExprOption{
		RuleFields:    ruleFields,
		MaxInLimit:    200,
		MaxNotInLimit: 200,
		MaxRulesLimit: 50,
		MaxRulesDepth: MaxRulesDepth,
	}
}

// NewDefaultExprOpt init an expression option with default limit option.
func cloneExprOption(opt *ExprOption) *ExprOption {
	return &ExprOption{
		RuleFields:       opt.RuleFields,
		IgnoreRuleFields: opt.IgnoreRuleFields,
		MaxInLimit:       opt.MaxInLimit,
		MaxNotInLimit:    opt.MaxNotInLimit,
		MaxRulesLimit:    opt.MaxRulesLimit,
		MaxRulesDepth:    opt.MaxRulesDepth,
	}
}

// Expression is to build a query expression
type Expression struct {
	RuleFactory
}

// Validate if the expression is valid or not.
func (exp Expression) Validate(opt *ExprOption) error {
	if opt == nil {
		return errors.New("expression's validate option must be set")
	}

	if exp.RuleFactory == nil {
		return errors.New("expression should not be nil")
	}

	if opt.MaxRulesDepth > opt.MaxRulesDepth {
		return fmt.Errorf("max rule depth exceeds maximum limit: %d", opt.MaxRulesDepth)
	}

	return exp.RuleFactory.Validate(opt)
}

// MarshalJSON marshal Expression into json value
func (exp Expression) MarshalJSON() ([]byte, error) {
	if exp.RuleFactory != nil {
		return json.Marshal(exp.RuleFactory)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON unmarshal Expression from json value
func (exp *Expression) UnmarshalJSON(raw []byte) error {
	rule, err := parseJsonRule(raw)
	if err != nil {
		return fmt.Errorf("parse rule(%s) failed, err: %v", string(raw), err)
	}

	exp.RuleFactory = rule
	return nil
}

// MarshalBSON marshal Expression into bson value
func (exp *Expression) MarshalBSON() ([]byte, error) {
	// right now bson will panic if MarshalBSON is defined using a value receiver and called by a nil pointer
	// TODO this is compatible for nil pointer, but struct marshalling is not supported, find a way to support both
	if exp == nil {
		return bson.Marshal((map[string]interface{})(nil))
	}

	if exp.RuleFactory != nil {
		return bson.Marshal(exp.RuleFactory)
	}

	return bson.Marshal((map[string]interface{})(nil))
}

// UnmarshalBSON unmarshal Expression from bson value
func (exp *Expression) UnmarshalBSON(raw []byte) error {
	rule, err := parseBsonRule(raw)
	if err != nil {
		return fmt.Errorf("parse rule failed, err: %v", err)
	}

	exp.RuleFactory = rule
	return nil
}

func parseJsonRule(raw []byte) (RuleFactory, error) {
	// rule with 'condition' key means that it is a combined rule
	if gjson.GetBytes(raw, "condition").Exists() {
		rule := new(CombinedRule)
		err := json.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
		}
		return rule, nil
	}

	// rule with 'operator' key means that it is an atomic rule
	if gjson.GetBytes(raw, "operator").Exists() {
		rule := new(AtomRule)
		err := json.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into atomic rule failed, err: %v", err)
		}
		return rule, nil
	}

	return nil, errors.New("no rule is found")
}

func parseBsonRule(raw []byte) (RuleFactory, error) {
	// rule with 'condition' key means that it is a combined rule
	if _, ok := bson.Raw(raw).Lookup("condition").StringValueOK(); ok {
		rule := new(CombinedRule)
		err := bson.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
		}
		return rule, nil
	}

	// rule with 'operator' key means that it is an atomic rule
	if _, ok := bson.Raw(raw).Lookup("operator").StringValueOK(); ok {
		rule := new(AtomRule)
		err := bson.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into atomic rule failed, err: %v", err)
		}
		return rule, nil
	}

	return nil, errors.New("no rule is found")
}
