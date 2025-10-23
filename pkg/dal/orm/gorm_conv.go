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

// Package orm ...
package orm

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/util"
)

// ConvFilter convert non-nil filter to gorm clause expression
func ConvFilter(flt filter.RuleFactory) (clause.Expression, error) {
	if flt == nil {
		return nil, errors.New("filter expression is nil")
	}

	switch typed := flt.(type) {
	case *filter.AtomRule:
		return atomRuleToGormClause(typed)
	case *filter.Expression:
		return expressionToGormClause(typed)
	default:
		return nil, fmt.Errorf("filter type is not supported: %T", flt)
	}
}

// expressionToGormClause convert *filter.Expression to gorm clause expression, return nil if flt is empty
func expressionToGormClause(flt *filter.Expression) (exp clause.Expression, err error) {
	if flt.IsEmpty() {
		return nil, errors.New("expression is empty")
	}
	if flt.Op == filter.All {
		return clause.Expr{
			SQL:                "1 = 1",
			Vars:               nil,
			WithoutParentheses: false,
		}, nil
	}

	var exps []clause.Expression
	var expr clause.Expression
	for _, sub := range flt.Rules {
		expr, err = ConvFilter(sub)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			exps = append(exps, expr)
		}
	}

	switch flt.Op {
	case filter.And:
		exp = clause.And(exps...)
	case filter.Or:
		exp = clause.Or(exps...)
	default:
		return nil, fmt.Errorf("expression op is not supported: %s", flt.Op)
	}
	return exp, nil
}

// atomRuleToGormClause convert *filter.AtomRule to gorm clause expression
func atomRuleToGormClause(rule *filter.AtomRule) (clause.Expression, error) {
	if rule == nil {
		return nil, errors.New("rule is nil")
	}
	op := filter.OpType(rule.Op)
	switch op {
	case filter.Equal:
		return clause.Eq{Column: rule.Field, Value: rule.Value}, nil
	case filter.NotEqual:
		return clause.Neq{Column: rule.Field, Value: rule.Value}, nil
	case filter.LessThan:
		return clause.Lt{Column: rule.Field, Value: rule.Value}, nil
	case filter.LessThanEqual:
		return clause.Lte{Column: rule.Field, Value: rule.Value}, nil
	case filter.GreaterThan:
		return clause.Gt{Column: rule.Field, Value: rule.Value}, nil
	case filter.GreaterThanEqual:
		return clause.Gte{Column: rule.Field, Value: rule.Value}, nil
	case filter.In:
		return buildIN(rule)
	case filter.NotIn:
		return buildNotIN(rule)
	case filter.ContainsInsensitive:
		return buildCIS(rule)
	case filter.ContainsSensitive:
		return buildCS(rule)
	default:
		// try other operator below
	}
	if filter.IsJSONOperator(op) {
		return jsonRuleToClauseExpr(rule)
	}
	if filter.IsArrayOperator(op) {
		return arrayRuleToClauseExpr(rule)
	}
	return nil, fmt.Errorf("rule op is not supported: %s", op)
}

func buildCS(rule *filter.AtomRule) (clause.Expression, error) {
	s, ok := util.GetString(rule.Value)
	if !ok {
		return nil, errors.New("cs operator's value should be an string")
	}

	if len(s) == 0 {
		return nil, errors.New("cs operator's value can not be a empty string")
	}

	likeExpr := clause.Like{
		Column: rule.Field,
		Value:  "%" + s + "%",
	}
	return likeExpr, nil
}

func buildCIS(rule *filter.AtomRule) (clause.Expression, error) {
	s, ok := util.GetString(rule.Value)
	if !ok {
		return nil, errors.New("cis operator's value should be an string")
	}

	if len(s) == 0 {
		return nil, errors.New("cis operator's value can not be a empty string")
	}
	likeExpr := clause.Like{
		Column: clause.Expr{
			SQL:                "LOWER(?)",
			Vars:               []any{clause.Column{Name: rule.Field}},
			WithoutParentheses: true,
		},
		Value: "%" + strings.ToLower(s) + "%",
	}
	return likeExpr, nil
}

func buildIN(rule *filter.AtomRule) (clause.Expression, error) {
	values, ok := rule.Value.([]any)
	if !ok {
		return nil, fmt.Errorf("filter value is not []any: %T", rule.Value)
	}
	return clause.IN{Column: rule.Field, Values: values}, nil
}

func buildNotIN(rule *filter.AtomRule) (clause.Expression, error) {
	in, err := buildIN(rule)
	if err != nil {
		return nil, err
	}
	return clause.Not(in), nil
}
