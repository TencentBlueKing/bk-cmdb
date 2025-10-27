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

package conv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
)

type arrayTestModel struct {
	ID            int                  `gorm:"primaryKey"`
	Name          string               `gorm:"column:name"`
	IntArray      types.Array[int64]   `gorm:"column:int_array"`
	StringArray   types.Array[string]  `gorm:"column:string_array"`
	FloatArray    types.Array[float64] `gorm:"column:float_array"`
	BoolArray     *types.Array[bool]   `gorm:"column:bool_array"`
	NullableArray *types.Array[[]byte] `gorm:"column:nullable_array"`
}

func (arrayTestModel) TableName() string {
	return "array_test"
}

type userIntType int

func Test_arrayRuleToClauseExpr(t *testing.T) {
	type args struct {
		rule *filter.AtomRule
	}
	tests := []struct {
		name         string
		args         args
		want         clause.Expression
		wantConvErr  string
		wantQueryErr string
		wantSQL      string
		wantVars     []any
		shouldFound  []arrayTestModel
	}{
		{
			name: "EmptySlice-Equal",
			args: args{
				rule: filter.RuleArrayEqual("int_array", []int64{}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{}`},
			shouldFound: []arrayTestModel{},
		},
		{
			name: "EmptySlice-Equal",
			args: args{
				rule: filter.RuleArrayEqual("bool_array", []bool{}),
			},
			want:        &SimpleExpression{Column: "bool_array", OP: "=", Value: `{}`},
			wantSQL:     `"bool_array" = $1`,
			wantVars:    []any{`{}`},
			shouldFound: []arrayTestModel{arrayInst2},
		},
		{
			name: "NilSlice-ArrayEqual",
			args: args{
				// will be converted to empty array, is want to check column is NULL, use RuleEqual
				rule: filter.RuleArrayEqual("int_array", []int64(nil)),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{}`},
			shouldFound: []arrayTestModel{},
		},
		{
			name: "EmptySlice-Contains",
			args: args{
				rule: filter.RuleArrayContains("int_array", []int64(nil)),
			},
			want:     &SimpleExpression{Column: "int_array", OP: "@>", Value: `{}`},
			wantSQL:  `"int_array" @> $1`,
			wantVars: []any{`{}`},
			// contains empty array will return all data, should be excluded by validate
			shouldFound: allArrayTestModels,
		},
		{
			name: "EmptySlice-Overlap",
			args: args{
				rule: filter.RuleArrayOverlap("int_array", []int64(nil)),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "&&", Value: `{}`},
			wantSQL:     `"int_array" && $1`,
			wantVars:    []any{`{}`},
			shouldFound: []arrayTestModel{},
		},
		{
			name: "Equal-int64",
			args: args{
				rule: filter.RuleArrayEqual("int_array", []int64{1, 2, 3}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{1,2,3}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{1,2,3}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			// expect array elem types, other type may also be supported for query, but not is not recommend
			name: "int16",
			args: args{
				rule: &filter.AtomRule{Field: "int_array", Op: filter.ArrayEqual.Factory(), Value: []int16{1, 2, 3}},
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{1,2,3}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{1,2,3}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "unsupportedArrayType-any",
			args: args{
				rule: &filter.AtomRule{Field: "int_array", Op: filter.ArrayEqual.Factory(), Value: []any{1, 2, 3}},
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{1,2,3}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{1,2,3}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "Equal-userIntType",
			args: args{
				rule: &filter.AtomRule{Field: "int_array", Op: filter.ArrayEqual.Factory(), Value: []userIntType{1, 2, 3}},
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{1,2,3}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{1,2,3}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "NotEqual-int64",
			args: args{
				rule: filter.RuleArrayNotEqual("int_array", []int64{1, 2, 3}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "<>", Value: `{1,2,3}`},
			wantSQL:     `"int_array" <> $1`,
			wantVars:    []any{`{1,2,3}`},
			shouldFound: []arrayTestModel{arrayInst2, arrayInst3},
		},
		{
			name: "Equal-bad-type-equal",
			args: args{
				rule: filter.RuleArrayEqual("int_array", []string{"1", "2", "3"}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{"1","2","3"}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{"1","2","3"}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "Equal-bad-type-no-match",
			args: args{
				rule: filter.RuleArrayEqual("int_array", []string{"a", "b", "c"}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "=", Value: `{"a","b","c"}`},
			wantSQL:     `"int_array" = $1`,
			wantVars:    []any{`{"a","b","c"}`},
			shouldFound: []arrayTestModel{},
		},
		{
			name: "Contains",
			args: args{
				rule: filter.RuleArrayContains("string_array", []string{"a", "b"}),
			},
			want:        &SimpleExpression{Column: "string_array", OP: "@>", Value: `{"a","b"}`},
			wantSQL:     `"string_array" @> $1`,
			wantVars:    []any{`{"a","b"}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "In",
			args: args{
				rule: filter.RuleArraySubset("string_array", []string{"a", "b", "c", "d"}),
			},
			want:        &SimpleExpression{Column: "string_array", OP: "<@", Value: `{"a","b","c","d"}`},
			wantSQL:     `"string_array" <@ $1`,
			wantVars:    []any{`{"a","b","c","d"}`},
			shouldFound: []arrayTestModel{arrayInst1},
		},
		{
			name: "In2",
			args: args{
				rule: filter.RuleArraySubset("string_array", []string{"a", "b", "c", "d", "e", "f"}),
			},
			want:        &SimpleExpression{Column: "string_array", OP: "<@", Value: `{"a","b","c","d","e","f"}`},
			wantSQL:     `"string_array" <@ $1`,
			wantVars:    []any{`{"a","b","c","d","e","f"}`},
			shouldFound: []arrayTestModel{arrayInst1, arrayInst2},
		},
		{
			name: "Overlap",
			args: args{
				rule: filter.RuleArrayOverlap("int_array", []int64{5}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "&&", Value: `{5}`},
			wantSQL:     `"int_array" && $1`,
			wantVars:    []any{`{5}`},
			shouldFound: []arrayTestModel{arrayInst2},
		},
		{
			name: "Overlap3",
			args: args{
				rule: filter.RuleArrayOverlap("int_array", []int64{1, 5, 9}),
			},
			want:        &SimpleExpression{Column: "int_array", OP: "&&", Value: `{1,5,9}`},
			wantSQL:     `"int_array" && $1`,
			wantVars:    []any{`{1,5,9}`},
			shouldFound: []arrayTestModel{arrayInst1, arrayInst2, arrayInst3},
		},
		{
			name: "OverlapNullableBytea",
			args: args{
				rule: filter.RuleArrayOverlap("nullable_array", [][]byte{{1}, {2}}),
			},
			want:        &SimpleExpression{Column: "nullable_array", OP: "&&", Value: `{"\\x01","\\x02"}`},
			wantSQL:     `"nullable_array" && $1`,
			wantVars:    []any{`{"\\x01","\\x02"}`},
			shouldFound: []arrayTestModel{arrayInst2, arrayInst3},
		},
		{
			name: "OverlapNullableBytea-single",
			args: args{
				rule: filter.RuleArrayOverlap("nullable_array", [][]byte{{1}}),
			},
			want:        &SimpleExpression{Column: "nullable_array", OP: "&&", Value: `{"\\x01"}`},
			wantSQL:     `"nullable_array" && $1`,
			wantVars:    []any{`{"\\x01"}`},
			shouldFound: []arrayTestModel{arrayInst2, arrayInst3},
		},
		{
			name: "isEmpty",
			args: args{
				rule: filter.RuleArrayIsEmpty("bool_array"),
			},
			want: &SimpleExpression{Column: "bool_array", OP: "IS NULL",
				ColumnWrapper: &functionWrapper{Name: "array_length", ExtraArgs: []any{1}},
				ValueWrapper:  &emptyWrapper{}},
			wantSQL:     `array_length("bool_array",$1) IS NULL `,
			wantVars:    []any{1},
			shouldFound: []arrayTestModel{arrayInst1, arrayInst2},
		},
		{
			name: "notEmpty",
			args: args{
				rule: filter.RuleArrayNotEmpty("bool_array"),
			},
			want: &SimpleExpression{Column: "bool_array", OP: "IS NOT NULL",
				ColumnWrapper: &functionWrapper{Name: "array_length", ExtraArgs: []any{1}},
				ValueWrapper:  &emptyWrapper{}},
			wantSQL:     `array_length("bool_array",$1) IS NOT NULL `,
			wantVars:    []any{1},
			shouldFound: []arrayTestModel{arrayInst3},
		},
	}

	g := prepareTestDB(t, arrayTestModel{}, allArrayTestModels)
	if !assert.NotNil(t, g, "failed to open gorm db") {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := arrayRuleToClauseExpr(tt.args.rule)
			if tt.wantConvErr != "" {
				assert.ErrorContains(t, err, tt.wantConvErr, "convert to clause failed got unexpected error")
				return
			}
			if !assert.Nil(t, err, "convert to clause failed") {
				return
			}
			if !assert.Equal(t, tt.want, got, "clause mismatch") {
				return
			}

			s := &gorm.Statement{DB: g.Session(&gorm.Session{DryRun: true})}
			got.Build(s)
			sql := s.SQL.String()
			t.Log("SQL:", sql)
			assert.Equal(t, tt.wantSQL, sql, "sql mismatch")
			assert.Equal(t, tt.wantVars, s.Vars, "bind vars mismatch")

			// try query
			found := make([]arrayTestModel, 0)
			err = g.Model(arrayTestModel{}).Clauses(got).Find(&found).Error
			if (err != nil) != (tt.wantConvErr != "") {
				assert.ErrorContainsf(t, err, tt.wantQueryErr, "query failed: %v", err)
			}
			assert.Equalf(t, tt.shouldFound, found, "found data mismatch")
		})
	}
}

func TestArrayFilter(t *testing.T) {
	type args struct {
		rule filter.RuleFactory
	}
	tests := []struct {
		name         string
		args         args
		want         clause.Expression
		wantConvErr  string
		wantQueryErr string
		wantSQL      string
		wantVars     []any
		shouldFound  []arrayTestModel
	}{
		{
			name:        "ContainsFloatRule",
			args:        args{rule: filter.RuleArrayContains("float_array", []float32{-0.2})},
			want:        &SimpleExpression{Column: "float_array", OP: "@>", Value: `{-0.2}`},
			wantSQL:     `"float_array" @> $1`,
			wantVars:    []any{`{-0.2}`},
			shouldFound: []arrayTestModel{arrayInst1, arrayInst2},
		},
		{
			name: "ContainsFloatExpression",
			args: args{rule: filter.ExpressionAnd(
				filter.RuleArrayContains("float_array", []float64{-0.4}))},
			want:        &SimpleExpression{Column: "float_array", OP: "@>", Value: `{-0.4}`},
			wantSQL:     `"float_array" @> $1`,
			wantVars:    []any{`{-0.4}`},
			shouldFound: []arrayTestModel{arrayInst2, arrayInst3},
		},
		{
			name: "arraySlice",
			args: args{rule: filter.ExpressionAnd(
				filter.RuleArrayContains("float_array[0:1]", []float64{-0.4}))},
			want:     &SimpleExpression{Column: "float_array[0:1]", OP: "@>", Value: `{-0.4}`},
			wantSQL:  `"float_array[0:1]" @> $1`,
			wantVars: []any{`{-0.4}`},
			// array slice is not supported currently
			wantQueryErr: "column \"float_array[0:1]\" does not exist",
			shouldFound:  []arrayTestModel{},
		},
	}

	g := prepareTestDB(t, arrayTestModel{}, allArrayTestModels)
	if !assert.NotNil(t, g, "failed to open gorm db") {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.args.rule)
			if tt.wantConvErr != "" {
				assert.ErrorContains(t, err, tt.wantConvErr, "convert to clause failed got unexpected error")
				return
			}
			if !assert.Nil(t, err, "convert to clause failed") {
				return
			}

			if !assert.Equal(t, tt.want, got, "clause mismatch") {
				return
			}

			s := &gorm.Statement{DB: g.Session(&gorm.Session{DryRun: true})}
			got.Build(s)
			sql := s.SQL.String()
			t.Log("SQL:", sql)
			assert.Equal(t, tt.wantSQL, sql, "sql mismatch")
			assert.Equal(t, tt.wantVars, s.Vars, "bind vars mismatch")

			// try query
			found := make([]arrayTestModel, 0)
			err = g.Model(arrayTestModel{}).Clauses(got).Find(&found).Error
			if (err != nil) != (tt.wantConvErr != "") {
				assert.ErrorContainsf(t, err, tt.wantQueryErr, "query failed: %v", err)
			}
			assert.Equalf(t, tt.shouldFound, found, "found data mismatch")
		})
	}
}

var allArrayTestModels = []arrayTestModel{arrayInst1, arrayInst2, arrayInst3}
var arrayInst1 = arrayTestModel{
	ID:            1,
	Name:          "array-1",
	IntArray:      types.Array[int64]{1, 2, 3},
	StringArray:   types.Array[string]{"a", "b", "c"},
	FloatArray:    types.Array[float64]{-0.1, -0.2, -0.3},
	NullableArray: nil,
	BoolArray:     nil,
}
var arrayInst2 = arrayTestModel{
	ID:            2,
	Name:          "array-2",
	IntArray:      types.Array[int64]{4, 5, 6},
	StringArray:   types.Array[string]{"d", "e", "f"},
	FloatArray:    types.Array[float64]{-0.2, -0.3, -0.4},
	NullableArray: &types.Array[[]byte]{{1}, {2}, {4}},
	BoolArray:     &types.Array[bool]{},
}
var arrayInst3 = arrayTestModel{
	ID:            3,
	Name:          "array-3",
	IntArray:      types.Array[int64]{7, 8, 9},
	StringArray:   types.Array[string]{"g", "h", "i"},
	FloatArray:    types.Array[float64]{-0.3, -0.4, -0.5},
	NullableArray: &types.Array[[]byte]{{1}, {2}, {3}},
	BoolArray:     &types.Array[bool]{true},
}

func Test_buildArraySQL(t *testing.T) {
	tests := []struct {
		name    string
		arg     any
		wantSql string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "string",
			arg:     []string{"a"},
			wantSql: `{"a"}`,
			wantErr: assert.NoError,
		},
		{
			name:    "bool",
			arg:     []bool{true, false},
			wantSql: `{t,f}`,
			wantErr: assert.NoError,
		},
		{
			name:    "bytea",
			arg:     [][]byte{{1, 2, 3}, []byte("abc")},
			wantSql: `{"\\x010203","\\x616263"}`,
			wantErr: assert.NoError,
		},
		{
			name:    "int32",
			arg:     []int32{1, 2},
			wantSql: `{1,2}`,
			wantErr: assert.NoError,
		},
		{
			name:    "int64",
			arg:     []int64{1, 2},
			wantSql: `{1,2}`,
			wantErr: assert.NoError,
		},
		{
			name:    "float32",
			arg:     []float32{1, 2},
			wantSql: `{1,2}`,
			wantErr: assert.NoError,
		},
		{
			name:    "float64",
			arg:     []float64{1, 2},
			wantSql: `{1,2}`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, err := buildArraySQL(tt.arg)
			if !tt.wantErr(t, err, fmt.Sprintf("buildArraySQL(%v)", tt.arg)) {
				return
			}
			assert.Equalf(t, tt.wantSql, gotSql, "buildArraySQL(%v)", tt.arg)
		})
	}
}
