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
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

type testModel struct {
	ID           int       `gorm:"column:id;primaryKey"`
	Str          string    `gorm:"column:str"`
	Int          int       `gorm:"column:int"`
	Float64      float64   `gorm:"column:float64"`
	Time         time.Time `gorm:"column:time"`
	NullableBool *bool     `gorm:"column:nullable_bool"`
}

func (m testModel) TableName() string {
	return "test_model"
}

func Test_atomRuleToGormClause(t *testing.T) {
	type args struct {
		rule *filter.AtomRule
	}

	tests := []struct {
		name         string
		args         args
		want         clause.Expression
		wantConvErr  string
		wantSQL      string
		wantVars     []any
		wantQueryErr string
		shouldFound  []testModel
	}{
		{
			name: "TestStrEqual",
			args: args{
				rule: filter.RuleEqual("str", "abc"),
			},
			want:        clause.Eq{Column: "str", Value: "abc"},
			wantSQL:     `"str" = $1`,
			wantVars:    []any{"abc"},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name: "TestStrNotEqual",
			args: args{
				rule: filter.RuleNotEqual("str", "abc"),
			},
			want:        clause.Neq{Column: "str", Value: "abc"},
			wantSQL:     `"str" <> $1`,
			wantVars:    []any{"abc"},
			shouldFound: []testModel{defTestModel},
		},
		{
			name: "TestStrCis",
			args: args{
				rule: filter.RuleCis("str", "b"),
			},
			want: clause.Like{
				Column: clause.Expr{SQL: "LOWER(?)", WithoutParentheses: true, Vars: []any{clause.Column{Name: "str"}}},
				Value:  "%b%"},
			wantSQL:     `LOWER("str") LIKE $1`,
			wantVars:    []any{"%b%"},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name: "TestStrIn",
			args: args{
				rule: filter.RuleIn("str", []string{"abc", "123"}),
			},
			want:        clause.IN{Column: "str", Values: []any{"abc", "123"}},
			wantSQL:     `"str" IN ($1,$2)`,
			wantVars:    []any{"abc", "123"},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name: "TestStrInEmpty",
			args: args{
				rule: filter.RuleIn("str", []string{}),
			},
			want:        clause.IN{Column: "str", Values: []any{}},
			wantSQL:     `"str" IN (NULL)`,
			wantVars:    []any(nil),
			shouldFound: []testModel{},
		},
		{
			name: "TestStrNotIn",
			args: args{
				rule: filter.RuleNotIn("str", []string{"abc", "123"}),
			},
			want:        clause.Not(clause.IN{Column: "str", Values: []any{"abc", "123"}}),
			wantSQL:     `"str" NOT IN ($1,$2)`,
			wantVars:    []any{"abc", "123"},
			shouldFound: []testModel{defTestModel},
		},
		{
			name: "TestIntGt",
			args: args{
				rule: filter.RuleGreaterThan("int", 1),
			},
			want:        clause.Gt{Column: "int", Value: 1},
			wantSQL:     `"int" > $1`,
			wantVars:    []any{1},
			shouldFound: []testModel{defTestModel},
		},
		{
			name: "TestIntGte",
			args: args{
				rule: filter.RuleGreaterThanEqual("int", 1),
			},
			want:        clause.Gte{Column: "int", Value: 1},
			wantSQL:     `"int" >= $1`,
			wantVars:    []any{1},
			shouldFound: []testModel{abcTestModel, defTestModel},
		},
		{
			name: "TestFloat64Lt",
			args: args{
				rule: filter.RuleLessThan("float64", 1),
			},
			want:        clause.Lt{Column: "float64", Value: 1},
			wantSQL:     `"float64" < $1`,
			wantVars:    []any{1},
			shouldFound: []testModel{defTestModel},
		},
		{
			name: "TestFloat64Lte",
			args: args{
				rule: filter.RuleLessThanEqual("float64", 1.1),
			},
			want:        clause.Lte{Column: "float64", Value: 1.1},
			wantSQL:     `"float64" <= $1`,
			wantVars:    []any{1.1},
			shouldFound: []testModel{abcTestModel, defTestModel},
		},
		{
			name: "TestBoolEqualNil",
			args: args{
				rule: filter.RuleEqual("nullable_bool", nil),
			},
			want:        clause.Eq{Column: "nullable_bool", Value: nil},
			wantSQL:     `"nullable_bool" IS NULL`,
			wantVars:    nil,
			shouldFound: []testModel{defTestModel},
		},
	}
	g := prepareTestDB(t, &testModel{}, testModels)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := atomRuleToGormClause(tt.args.rule)
			if tt.wantConvErr != "" {
				if !assert.ErrorContains(t, err, tt.wantConvErr, "convert to clause failed got unexpected error") {
					return
				}
				return
			} else {
				assert.Nil(t, err, "convert to clause failed")
			}
			if !assert.Equal(t, tt.want, got, "got unexpected clause") {
				return
			}
			s := &gorm.Statement{DB: g.Session(&gorm.Session{DryRun: true})}
			got.Build(s)
			sql := s.SQL.String()
			t.Log("SQL:", sql)
			assert.Equal(t, tt.wantSQL, sql, "sql mismatch")
			assert.Equal(t, tt.wantVars, s.Vars, "bind vars mismatch")
			// try query
			found := make([]testModel, 0)
			err = g.Model(testModel{}).Clauses(got).Find(&found).Error
			if (err != nil) != (tt.wantConvErr != "") {
				assert.ErrorContainsf(t, err, tt.wantQueryErr, "query failed: %v", err)
			}
			assert.Equalf(t, tt.shouldFound, found, "found data mismatch")
		})
	}
}

func TestConvToGormClause(t *testing.T) {
	type args struct {
		flt *filter.Expression
	}
	tests := []struct {
		name         string
		args         args
		want         clause.Expression
		wantConvErr  string
		wantSQL      string
		wantVars     []any
		wantQueryErr string
		shouldFound  []testModel
	}{
		{
			name:        "TestAllExpression",
			args:        args{filter.AllExpression()},
			want:        clause.Expr{SQL: "1 = 1", WithoutParentheses: false},
			wantSQL:     `1 = 1`,
			wantVars:    nil,
			shouldFound: []testModel{abcTestModel, defTestModel},
		},
		{
			name: "TestAllExpressionWithOtherRule",
			args: args{filter.ExpressionAnd(
				filter.RuleIn("str", []any{"abc", "123"}),
				filter.AllExpression())},
			want: clause.And(
				clause.IN{Column: "str", Values: []any{"abc", "123"}},
				clause.Expr{SQL: "1 = 1", WithoutParentheses: false}),
			wantSQL:     `("str" IN ($1,$2) AND 1 = 1)`,
			wantVars:    []any{"abc", "123"},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name:        "TestSimpleContainsExpression",
			args:        args{filter.ContainersExpression("str", []any{"abc", "123"})},
			want:        clause.IN{Column: "str", Values: []any{"abc", "123"}},
			wantSQL:     `"str" IN ($1,$2)`,
			wantVars:    []any{"abc", "123"},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name: "TestMergeExpression",
			args: args{filter.ExpressionAnd(
				filter.RuleIn("str", []any{"abc", "123"}),
				filter.RuleGreaterThan("float64", 0.5),
			)},
			want:        clause.And(clause.IN{Column: "str", Values: []any{"abc", "123"}}, clause.Gt{Column: "float64", Value: 0.5}),
			wantSQL:     `("str" IN ($1,$2) AND "float64" > $3)`,
			wantVars:    []any{"abc", "123", 0.5},
			shouldFound: []testModel{abcTestModel},
		},
		{
			name: "TestComplexExpression",
			args: args{&filter.Expression{
				Op: filter.Or,
				Rules: []filter.RuleFactory{
					filter.ExpressionAnd(
						filter.RuleIn("str", []any{"abc", "123"}),
						filter.RuleGreaterThan("float64", 0.5),
					),
					filter.RuleEqual("str", "def"),
				},
			}},
			want:        clause.Or(clause.And(clause.IN{Column: "str", Values: []any{"abc", "123"}}, clause.Gt{Column: "float64", Value: 0.5}), clause.Eq{Column: "str", Value: "def"}),
			wantSQL:     `(("str" IN ($1,$2) AND "float64" > $3) OR "str" = $4)`,
			wantVars:    []any{"abc", "123", 0.5, "def"},
			shouldFound: []testModel{abcTestModel, defTestModel},
		},
	}
	g := prepareTestDB(t, &testModel{}, testModels)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expressionToGormClause(tt.args.flt)
			if err != nil {
				if tt.wantQueryErr == "" {
					t.Errorf("test %s should not have error: %v", tt.name, err)
					return
				} else {
					assert.ErrorContains(t, err, tt.wantConvErr, "convert to clause failed")
				}
			}
			assert.Equal(t, tt.want, got)
			s := &gorm.Statement{DB: g.Session(&gorm.Session{DryRun: true})}
			got.Build(s)
			sql := s.SQL.String()
			t.Log("SQL:", sql)
			assert.Equal(t, tt.wantSQL, sql, "sql mismatch")
			assert.Equal(t, tt.wantVars, s.Vars, "bind vars mismatch")
			// try query
			found := make([]testModel, 0)
			err = g.Model(testModel{}).Clauses(got).Find(&found).Error
			if (err != nil) != (tt.wantConvErr != "") {
				assert.ErrorContainsf(t, err, tt.wantQueryErr, "query failed: %v", err)
			}
			assert.Equalf(t, tt.shouldFound, found, "found data mismatch")
		})
	}
}

func prepareTestDB(t *testing.T, table interface{ TableName() string }, preInserts any) *gorm.DB {
	g, err := tests.GetTestGORM(t)
	if err != nil {
		t.Errorf("failed to open gorm db: %v", err)
		return nil
	}
	if g.Migrator().HasTable(table.TableName()) {
		err := g.Migrator().DropTable(table.TableName())
		assert.Nil(t, err, "failed to drop table")
	}
	err = g.Migrator().AutoMigrate(table)
	assert.Nil(t, err, "failed to migrate table")

	if preInserts != nil {
		err = g.Create(preInserts).Error
		assert.Nil(t, err, "failed to create data")
	}
	t.Cleanup(func() {
		if g.Migrator().HasTable(table.TableName()) {
			err := g.Migrator().DropTable(table.TableName())
			assert.Nil(t, err, "failed to drop table on cleanup")
		}
	})
	return g
}

var abcTestModel = testModel{
	ID:           1,
	Str:          "abc",
	Int:          1,
	Float64:      1.1,
	Time:         time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
	NullableBool: new(bool),
}
var defTestModel = testModel{
	ID:      2,
	Str:     "def",
	Int:     2,
	Float64: -2.2,
	Time:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
}
var testModels = []testModel{
	abcTestModel,
	defTestModel,
}

func TestConvFilter(t *testing.T) {
	tests := []struct {
		name    string
		flt     filter.RuleFactory
		want    clause.Expression
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "RuleEqual",
			flt:     filter.RuleEqual("str", "abc"),
			want:    clause.Eq{Column: "str", Value: "abc"},
			wantErr: assert.NoError,
		},
		{
			name:    "ExpressionAnd",
			flt:     filter.ExpressionAnd(filter.RuleEqual("str", "abc"), filter.RuleEqual("int", 1)),
			want:    clause.And(clause.Eq{Column: "str", Value: "abc"}, clause.Eq{Column: "int", Value: 1}),
			wantErr: assert.NoError,
		},
		{
			name: "nil",
			flt:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "filter expression is nil")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.flt)
			if !tt.wantErr(t, err, fmt.Sprintf("Filter(%v)", tt.flt)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Filter(%v)", tt.flt)
		})
	}
}
