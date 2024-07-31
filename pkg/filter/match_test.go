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
	"testing"
	"time"

	"configcenter/src/common"

	"github.com/stretchr/testify/assert"
)

func TestEqualMatch(t *testing.T) {
	op := Equal.Factory().Operator()

	// test equal int type
	matched, err := op.Match(int32(1), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(int32(2), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test equal string type
	matched, err = op.Match("a", "a")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test equal bool type
	matched, err = op.Match(false, false)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(true, false)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotEqualMatch(t *testing.T) {
	op := NotEqual.Factory().Operator()

	// test not equal int type
	matched, err := op.Match(int32(1), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(int32(2), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not equal string type
	matched, err = op.Match("a", "a")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not equal bool type
	matched, err = op.Match(false, false)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(true, false)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}

func TestInMatch(t *testing.T) {
	op := In.Factory().Operator()

	// test in int array type
	matched, err := op.Match(1.0, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(3, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test in string array type
	matched, err = op.Match("b", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("c", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test in bool array type
	matched, err = op.Match(false, []interface{}{true, false})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(false, []interface{}{true})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotInMatch(t *testing.T) {
	op := NotIn.Factory().Operator()

	// test not in int array type
	matched, err := op.Match(1.0, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(3, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not in string array type
	matched, err = op.Match("b", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("c", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not in bool array type
	matched, err = op.Match(true, []interface{}{false, true})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(false, []interface{}{true})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}

func TestLessMatch(t *testing.T) {
	op := Less.Factory().Operator()

	// test less int type
	matched, err := op.Match(0.01, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(3, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less than 0
	matched, err = op.Match(-1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1.1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less than a negative number
	matched, err = op.Match(-1.23, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-1, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestLessOrEqualMatch(t *testing.T) {
	op := LessOrEqual.Factory().Operator()

	// test less or equal int type
	matched, err := op.Match(0.01, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(3, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less or equal than 0
	matched, err = op.Match(-1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1.1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less or equal than a negative number
	matched, err = op.Match(-1.0, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-0.01, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestGreaterMatch(t *testing.T) {
	op := Greater.Factory().Operator()

	// test greater int type
	matched, err := op.Match(3, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(0.01, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test greater than 0
	matched, err = op.Match(1.1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test greater than a negative number
	matched, err = op.Match(-0.01, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-1, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestGreaterOrEqualMatch(t *testing.T) {
	op := GreaterOrEqual.Factory().Operator()

	// test greater or equal int type
	matched, err := op.Match(3, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	assert.NoError(t, err)
	matched, err = op.Match(0.01, 1)
	assert.Equal(t, false, matched)

	// test greater or equal than 0
	matched, err = op.Match(1.1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test greater or equal than a negative number
	matched, err = op.Match(-1.0, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-1.23, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestDatetimeLessMatch(t *testing.T) {
	op := DatetimeLess.Factory().Operator()

	// test datetime less time type
	now := time.Now()
	matched, err := op.Match(now.Unix()-1, now)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(time.Second).Format(common.TimeTransferModel), now)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime less timestamp type
	matched, err = op.Match(now.Add(-time.Second).Format(common.TimeTransferModel), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now, now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime less time string
	nowStr := now.Format(common.TimeTransferModel)
	matched, err = op.Match(now.Add(-time.Second), nowStr)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Unix()+1, nowStr)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestDatetimeLessOrEqualMatch(t *testing.T) {
	op := DatetimeLessOrEqual.Factory().Operator()

	// test datetime less or equal time type
	now := time.Now()
	matched, err := op.Match(now.Unix()-1, now)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(time.Second).Format(common.TimeTransferModel), now)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime less or equal timestamp type
	matched, err = op.Match(now.Format(common.TimeTransferModel), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(time.Second), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime less or equal time string
	nowStr := now.Format(common.TimeTransferModel)
	matched, err = op.Match(now.Add(-time.Second), nowStr)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Unix()+1, nowStr)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestDatetimeGreaterMatch(t *testing.T) {
	op := DatetimeGreater.Factory().Operator()

	// test datetime greater time type
	now := time.Now()
	matched, err := op.Match(now.Add(time.Second).Format(common.TimeTransferModel), now)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Unix()-1, now)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime greater timestamp type
	matched, err = op.Match(now.Add(time.Second), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Format(common.TimeTransferModel), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime greater time string
	nowStr := now.Format(common.TimeTransferModel)
	matched, err = op.Match(now.Unix()+1, nowStr)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(-time.Second), nowStr)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestDatetimeGreaterOrEqualMatch(t *testing.T) {
	op := DatetimeGreaterOrEqual.Factory().Operator()

	// test datetime greater or equal time type
	now := time.Now()
	matched, err := op.Match(now.Add(time.Second).Format(common.TimeTransferModel), now)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Unix()-1, now)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime greater or equal timestamp type
	matched, err = op.Match(now.Format(common.TimeTransferModel), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(-time.Second), now.Unix())
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test datetime greater or equal time string
	nowStr := now.Format(common.TimeTransferModel)
	matched, err = op.Match(now.Unix()+1, nowStr)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(now.Add(-time.Second).Format(common.TimeTransferModel), nowStr)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestBeginsWithMatch(t *testing.T) {
	op := BeginsWith.Factory().Operator()

	// test begins with string type
	matched, err := op.Match("abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestBeginsWithInsensitiveMatch(t *testing.T) {
	op := BeginsWithInsensitive.Factory().Operator()

	// test begins with insensitive string type
	matched, err := op.Match("aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("Abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotBeginsWithMatch(t *testing.T) {
	op := NotBeginsWith.Factory().Operator()

	// test not begins with string type
	matched, err := op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotBeginsWithInsensitiveMatch(t *testing.T) {
	op := NotBeginsWithInsensitive.Factory().Operator()

	// test not begins with insensitive string type
	matched, err := op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestContainsMatch(t *testing.T) {
	op := Contains.Factory().Operator()

	// test contains string type
	matched, err := op.Match("123aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestContainsSensitiveMatch(t *testing.T) {
	op := ContainsSensitive.Factory().Operator()

	// test contains string type
	matched, err := op.Match("123abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotContainsMatch(t *testing.T) {
	op := NotContains.Factory().Operator()

	// test not contains string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotContainsInsensitiveMatch(t *testing.T) {
	op := NotContainsInsensitive.Factory().Operator()

	// test not contains insensitive string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123Abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestEndsWithMatch(t *testing.T) {
	op := EndsWith.Factory().Operator()

	// test ends with string type
	matched, err := op.Match("123abcdef", "ef")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "eF")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("123abcdef", "df")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestEndsWithInsensitiveMatch(t *testing.T) {
	op := EndsWithInsensitive.Factory().Operator()

	// test ends with insensitive string type
	matched, err := op.Match("123abcDef", "dEf")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123Abcdef", "abc")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotEndsWithMatch(t *testing.T) {
	op := NotEndsWith.Factory().Operator()

	// test not ends with string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ef")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotEndsWithInsensitiveMatch(t *testing.T) {
	op := NotEndsWithInsensitive.Factory().Operator()

	// test not ends with insensitive string type
	matched, err := op.Match("123Abcdef", "abc")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcDef", "dEf")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsEmptyMatch(t *testing.T) {
	op := IsEmpty.Factory().Operator()

	// test is empty matched
	matched, err := op.Match([]int64{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]float64{1, 2}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{""}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{true}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsNotEmptyMatch(t *testing.T) {
	op := IsNotEmpty.Factory().Operator()

	// test is not empty matched
	matched, err := op.Match([]int64{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]float64{1, 2}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{""}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{false}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}

func TestSizeMatch(t *testing.T) {
	op := Size.Factory().Operator()

	// test size matched
	matched, err := op.Match([]int64{1, 2, 3}, 3)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]float64{1, 2}, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{"1", "2"}, 2)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{""}, 2)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{}, 0)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{true}, 2)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsNullMatch(t *testing.T) {
	op := IsNull.Factory().Operator()

	// test is null matched
	matched, err := op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsNotNullMatch(t *testing.T) {
	op := IsNotNull.Factory().Operator()

	// test is not null matched
	matched, err := op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestExistMatch(t *testing.T) {
	op := Exist.Factory().Operator()

	// test exist matched
	matched, err := op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotExistMatch(t *testing.T) {
	op := NotExist.Factory().Operator()

	// test not exist matched
	matched, err := op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestObjectMatch(t *testing.T) {
	op := Object.Factory().Operator()

	// test filter object normal type
	matched, err := op.Match(`{"test":1,"test1":[{"test2":"d"}],"test3":0}`, exampleRule)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}

func TestArrayMatch(t *testing.T) {
	op := Array.Factory().Operator()

	// test filter array matched
	matched, err := op.Match([]string{"a", "cc"}, &CombinedRule{
		Condition: And,
		Rules: []RuleFactory{
			&AtomRule{
				Field:    ArrayElement,
				Operator: Contains.Factory(),
				Value:    "c",
			},
			&CombinedRule{
				Condition: Or,
				Rules: []RuleFactory{
					&AtomRule{
						Field:    ArrayElement,
						Operator: NotEqual.Factory(),
						Value:    "a",
					},
					&AtomRule{
						Field:    ArrayElement,
						Operator: In.Factory(),
						Value:    []string{"bb", "cc"},
					},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}

func TestRuleMatch(t *testing.T) {
	matched, err := exampleRule.Match(JsonString(`{"test":1,"test1":[{"test2":"b"}],"test3":111}`))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = exampleRule.Match(JsonString(`{"test":1,"test1":[{"test2":"a"}],"test3":111}`))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
