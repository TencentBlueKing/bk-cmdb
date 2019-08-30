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

package querybuilder_test

import (
	"testing"

	"configcenter/src/common/querybuilder"

	"github.com/stretchr/testify/assert"
)

func TestNormalAtomRule(t *testing.T) {
	rules := []querybuilder.AtomRule{
		{
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field1.field2",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    1.0,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    uint(1),
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    byte(1),
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    true,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorNotEqual,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    byte(1),
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    []int64{1, 2, 3},
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    []float64{1.0},
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    []string{"test"},
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    []bool{true, false, true},
		}, {
			Operator: querybuilder.OperatorNotIn,
			Field:    "field",
			Value:    []int64{1, 2, 3},
		}, {
			Operator: querybuilder.OperatorLess,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorLess,
			Field:    "field",
			Value:    1.0,
		}, {
			Operator: querybuilder.OperatorLessOrEqual,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorGreater,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorGreaterOrEqual,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorDatetimeLess,
			Field:    "field",
			Value:    "2019-08-04T14:08:00.00Z",
		}, {
			Operator: querybuilder.OperatorDatetimeLessOrEqual,
			Field:    "field",
			Value:    "2019-08-04T14:08:00.00Z",
		}, {
			Operator: querybuilder.OperatorDatetimeGreater,
			Field:    "field",
			Value:    "2019-08-04T14:08:00.00Z",
		}, {
			Operator: querybuilder.OperatorDatetimeGreaterOrEqual,
			Field:    "field",
			Value:    "2019-08-04T14:08:00.00Z",
		}, {
			Operator: querybuilder.OperatorBeginsWith,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorNotBeginsWith,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorContains,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorNotContains,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorsEndsWith,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorNotEndsWith,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorIsEmpty,
			Field:    "field",
			Value:    nil,
		}, {
			Operator: querybuilder.OperatorIsNotEmpty,
			Field:    "field",
			Value:    nil,
		}, {
			Operator: querybuilder.OperatorIsNull,
			Field:    "field",
			Value:    nil,
		}, {
			Operator: querybuilder.OperatorIsNotNull,
			Field:    "field",
			Value:    nil,
		}, {
			Operator: querybuilder.OperatorExist,
			Field:    "field",
			Value:    nil,
		}, {
			Operator: querybuilder.OperatorNotExist,
			Field:    "field",
			Value:    nil,
		},
	}
	for idx, rule := range rules {
		t.Logf("running case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.Nil(t, err)
		assert.Empty(t, errKey)
		assert.NotNil(t, filter)

		deep := rule.GetDeep()
		assert.Equal(t, deep, 1)

		errKey, err = rule.Validate()
		assert.Empty(t, errKey)
		assert.Nil(t, err)
	}
}

func TestInvalidateFieldAtomRule(t *testing.T) {
	rules := []querybuilder.AtomRule{
		{
			Operator: querybuilder.OperatorEqual,
			Field:    "1field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    ".field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "-field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "_field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field?",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "?field",
			Value:    1,
		},
	}

	for idx, rule := range rules {
		t.Logf("running invalid field case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		assert.Nil(t, filter)

		deep := rule.GetDeep()
		assert.Equal(t, deep, 1)

		errKey, err = rule.Validate()
		assert.NotEmpty(t, errKey)
		assert.NotNil(t, err)
	}
}

func TestInvalidateValueAtomRule(t *testing.T) {
	rules := []querybuilder.AtomRule{
		{
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    map[string]interface{}{"key": 1},
		}, {
			Operator: querybuilder.OperatorEqual,
			Field:    "field",
			Value:    []int{1, 2, 3},
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorIn,
			Field:    "field",
			Value:    map[string]interface{}{"key": 1},
		}, {
			Operator: querybuilder.OperatorLess,
			Field:    "field",
			Value:    true,
		}, {
			Operator: querybuilder.OperatorLess,
			Field:    "field",
			Value:    []int{1, 2},
		}, {
			Operator: querybuilder.OperatorDatetimeLess,
			Field:    "field",
			Value:    "test",
		}, {
			Operator: querybuilder.OperatorDatetimeLess,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorDatetimeLess,
			Field:    "field",
			Value:    []string{"2019-08-04T14:08:00.00Z"},
		}, {
			Operator: querybuilder.OperatorBeginsWith,
			Field:    "field",
			Value:    1,
		}, {
			Operator: querybuilder.OperatorBeginsWith,
			Field:    "field",
			Value:    []string{"test"},
		},
	}
	for idx, rule := range rules {
		t.Logf("running invalid value case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		assert.Nil(t, filter)

		deep := rule.GetDeep()
		assert.Equal(t, deep, 1)

		errKey, err = rule.Validate()
		assert.NotEmpty(t, errKey)
		assert.NotNil(t, err)
	}
}

func TestInvalidateOperatorAtomRule(t *testing.T) {
	rules := []querybuilder.AtomRule{
		{
			Operator: querybuilder.Operator("unknown"),
			Field:    "field",
			Value:    map[string]interface{}{"key": 1},
		},
	}
	for idx, rule := range rules {
		t.Logf("running invalid operator case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		assert.Nil(t, filter)

		deep := rule.GetDeep()
		assert.Equal(t, deep, 1)

		errKey, err = rule.Validate()
		assert.NotEmpty(t, errKey)
		assert.NotNil(t, err)
	}
}

func TestNormalCombinedRule(t *testing.T) {

	rules := []querybuilder.CombinedRule{
		{
			Condition: querybuilder.ConditionAnd,
			Rules: []querybuilder.Rule{
				querybuilder.AtomRule{
					Operator: querybuilder.OperatorEqual,
					Field:    "field",
					Value:    1,
				},
			},
		}, {
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				querybuilder.AtomRule{
					Operator: querybuilder.OperatorEqual,
					Field:    "field",
					Value:    1,
				},
			},
		},
	}
	for idx, rule := range rules {
		t.Logf("running normal combined rule case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.Nil(t, err)
		assert.Empty(t, errKey)
		assert.NotNil(t, filter)

		deep := rule.GetDeep()
		assert.Equal(t, deep, 2)

		errKey, err = rule.Validate()
		assert.Empty(t, errKey)
		assert.Nil(t, err)
	}
}

func TestInvalidCombinedRule(t *testing.T) {
	rules := []querybuilder.CombinedRule{
		{
			Condition: querybuilder.Condition("unknown"),
			Rules: []querybuilder.Rule{
				querybuilder.AtomRule{
					Operator: querybuilder.OperatorEqual,
					Field:    "field",
					Value:    1,
				},
			},
		}, {
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				querybuilder.AtomRule{
					Operator: querybuilder.OperatorEqual,
					Field:    "-field",
					Value:    1,
				},
			},
		}, {
			Condition: querybuilder.ConditionOr,
		},
	}
	for idx, rule := range rules {
		t.Logf("running invalid combined rule case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		assert.Nil(t, filter)

		errKey, err = rule.Validate()
		assert.NotEmpty(t, errKey)
		assert.NotNil(t, err)
	}
}

func TestExceedMaxDeep(t *testing.T) {
	rules := []querybuilder.CombinedRule{
		{
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				querybuilder.CombinedRule{
					Condition: querybuilder.ConditionOr,
					Rules: []querybuilder.Rule{
						querybuilder.AtomRule{
							Operator: querybuilder.OperatorEqual,
							Field:    "field",
							Value:    1,
						},
					},
				},
			},
		},
	}
	for idx, rule := range rules {
		t.Logf("running exceed max deep case %d, rule: %+v", idx, rule)
		filter, errKey, err := rule.ToMgo()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		assert.Nil(t, filter)

		errKey, err = rule.Validate()
		assert.NotEmpty(t, errKey)
		assert.NotNil(t, err)
	}
}
